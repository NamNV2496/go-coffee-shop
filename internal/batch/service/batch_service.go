package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/repo"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/cache"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/s3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BatchService interface {
	ClearAllOrderEOD(ctx context.Context) error
	GenerateReport(ctx context.Context) error
}

type batchService struct {
	s3Client  s3.S3Client
	orderRepo repo.OrderRepo
	cache     cache.Client
}

func NewBatchService(
	s3Client s3.S3Client,
	orderRepo repo.OrderRepo,
	cache cache.Client,
) BatchService {
	return &batchService{
		s3Client:  s3Client,
		orderRepo: orderRepo,
		cache:     cache,
	}
}

func (s batchService) ClearAllOrderEOD(
	ctx context.Context,
) error {

	var jsonData = make([]mq.RedisOrderDTO, 0)
	json, ok := json.Marshal(jsonData)
	if ok != nil {
		return status.Error(codes.Internal, "failed to marshall data into cache")
	}
	if err := s.cache.Set(ctx, mq.REDIS_KEY_ORDER, json); err != nil {
		return err
	}
	return nil
}

func (s batchService) getOrderOfToday(ctx context.Context) ([]domain.OrderItem, error) {

	orderItems, err := s.orderRepo.GetOrderItems(ctx)
	if err != nil {
		return []domain.OrderItem{}, nil
	}
	return orderItems, nil
}

func (s batchService) GenerateReport(ctx context.Context) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	orderItemList := make([]domain.OrderItem, 0, 8)
	header := []string{"ItemId", "Quantity", "Price"}

	// Load data function
	loadData := func() {
		var err error
		orderItemList, err = s.getOrderOfToday(context.Background()) // Adjust context as necessary
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
		return nil
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
		return err
	}

	// Open the file for reading
	file, err := os.Open(tempFile)
	if err != nil {
		fmt.Println("Error opening PDF file:", err)
		return err
	}
	defer file.Close()

	// Get the file info (for size)
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return err
	}
	fileSize := fileInfo.Size()

	// Upload to MinIO (assuming app.s3Client is properly initialized)
	_, err = s.s3Client.Write(ctx, fileStr, s3.BUCKETNAME, file, fileSize, "application/pdf")
	if err != nil {
		fmt.Println("Error uploading file to MinIO:", err)
		return err
	}

	fmt.Println("PDF successfully created and uploaded to MinIO:", fileStr)

	// Optionally, remove the temporary file
	// err = os.Remove(tempFile)
	// if err != nil {
	// 	fmt.Println("Error removing temporary file:", err)
	// }
	return nil
}
