package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/jung-kurt/gofpdf"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/handler/jobs"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/s3"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
)

type AppInterface interface {
	Start() error
	TriggerJob(ctx context.Context, req *gin.Context)
}

type App struct {
	s3Client     s3.S3Client
	batchService service.BatchService
	jobs         jobs.ClearAllOrderEOD
	cronConfig   configs.Cron
}

func NewApp(
	s3Client s3.S3Client,
	batchService service.BatchService,
	jobs jobs.ClearAllOrderEOD,
	cronConfig configs.Cron,
) *App {
	return &App{
		s3Client:     s3Client,
		batchService: batchService,
		jobs:         jobs,
		cronConfig:   cronConfig,
	}
}

func (app App) Start() error {

	// app.ExampleFpdf_CellFormat_tables(context.Background())
	go func() {
		err := app.startScheduler(context.Background())
		if err != nil {
			panic("failed to start scheduler")
		}
	}()
	return nil
}

func (app App) startScheduler(ctx context.Context) error {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return errors.New("failed to initialize scheduler")
	}

	hour, err := strconv.Atoi(app.cronConfig.ClearAllOrder.Hour)
	if err != nil {
		return err
	}
	minute, err := strconv.Atoi(app.cronConfig.ClearAllOrder.Minute)
	if err != nil {
		return err
	}

	j, err := scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(uint(hour), uint(minute), 0),
			),
		),
		gocron.NewTask(func() {
			if err := app.jobs.Run(context.Background()); err != nil {
				fmt.Println("failed to run execute all pending download task job")
			}
			app.ExampleFpdf_CellFormat_tables(context.Background())
		}),
	)
	if err != nil {
		fmt.Println("failed to schedule execute all pending download task job")
		return err
	}
	fmt.Println("Daily scheduler job id: ", j.ID())

	scheduler.Start()
	return nil
}

func (app App) TriggerJob(ctx context.Context, req *gin.Context) {
	fmt.Println("Trigger job report")
	app.jobs.Run(ctx)
	utils.WrapperResponse(req, http.StatusOK, "")
	// add generate pdf file
}

func (app App) ExampleFpdf_CellFormat_tables(ctx context.Context) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	orderItemList := make([]domain.OrderItem, 0, 8)
	header := []string{"ItemId", "Quantity", "Price"}

	// Load data function
	loadData := func() {
		var err error
		orderItemList, err = app.batchService.GetOrderOfToday(context.Background()) // Adjust context as necessary
		if err != nil {
			fmt.Println("Error loading data:", err)
			return
		}
	}

	// Simple table function
	basicTable := func() {
		left := (210.0 - 4*40) / 2
		pdf.SetX(left)
		for _, str := range header {
			pdf.CellFormat(40, 7, str, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
		for _, c := range orderItemList {
			pdf.SetX(left)
			pdf.CellFormat(40, 6, strconv.Itoa(int(c.ItemId)), "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 6, strconv.Itoa(int(c.Quantity)), "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 6, strconv.Itoa(int(c.Price)), "1", 0, "", false, 0, "")
			pdf.Ln(-1)
		}
	}

	// Load data
	loadData()

	// Check if data is loaded
	if len(orderItemList) == 0 {
		fmt.Println("No data loaded, exiting.")
		return
	}

	// Setup PDF
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()
	basicTable()

	// Output PDF file
	fileStr := "report_" + time.Now().Format("2006-01-02") + ".pdf"

	// push file to s3

	// Generate file name
	// Save PDF to a temporary file
	tempFile := "./" + fileStr
	err := pdf.OutputFileAndClose(tempFile)
	if err != nil {
		fmt.Println("Error creating PDF file:", err)
		return
	}

	// Open the file for reading
	file, err := os.Open(tempFile)
	if err != nil {
		fmt.Println("Error opening PDF file:", err)
		return
	}
	defer file.Close()

	// Get the file info (for size)
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	fileSize := fileInfo.Size()

	// Upload to MinIO (assuming app.s3Client is properly initialized)
	_, err = app.s3Client.Write(ctx, fileStr, s3.BUCKETNAME, file, fileSize, "application/pdf")
	if err != nil {
		fmt.Println("Error uploading file to MinIO:", err)
		return
	}

	fmt.Println("PDF successfully created and uploaded to MinIO:", fileStr)

	// Optionally, remove the temporary file
	// err = os.Remove(tempFile)
	// if err != nil {
	// 	fmt.Println("Error removing temporary file:", err)
	// }
}
