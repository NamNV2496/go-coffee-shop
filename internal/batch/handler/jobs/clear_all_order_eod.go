package jobs

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-co-op/gocron/v2"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
)

type ClearAllOrderEOD interface {
	Run(context.Context) error
	StartClearJobEOD(ctx context.Context) error
}

type clearAllOrderEOD struct {
	batchService service.BatchService
	cronConfig   configs.Cron
}

func NewExecuteClearAllOrderEOD(
	batchService service.BatchService,
	cronConfig configs.Cron,
) ClearAllOrderEOD {
	return &clearAllOrderEOD{
		batchService: batchService,
		cronConfig:   cronConfig,
	}
}

func (j clearAllOrderEOD) Run(ctx context.Context) error {
	fmt.Println("Trigger clearAllOrderEOD")
	if err := j.batchService.GenerateReport(ctx); err != nil {
		return err
	}
	return j.batchService.ClearAllOrderEOD(ctx)
}

func (j clearAllOrderEOD) StartClearJobEOD(ctx context.Context) error {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return errors.New("failed to initialize scheduler")
	}

	hour, err := strconv.Atoi(j.cronConfig.ClearAllOrder.Hour)
	if err != nil {
		return err
	}
	minute, err := strconv.Atoi(j.cronConfig.ClearAllOrder.Minute)
	if err != nil {
		return err
	}

	job, err := scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(uint(hour), uint(minute), 0),
			),
		),
		gocron.NewTask(func() {
			if err := j.Run(context.Background()); err != nil {
				fmt.Println("failed to run execute all pending download task job")
			}
			j.exampleFpdf_CellFormat_tables(context.Background())
		}),
	)
	if err != nil {
		fmt.Println("failed to schedule execute all pending download task job")
		return err
	}
	fmt.Println("Daily scheduler job id: ", job.ID())

	scheduler.Start()
	return nil
}

func (j clearAllOrderEOD) exampleFpdf_CellFormat_tables(ctx context.Context) {
	if err := j.batchService.GenerateReport(ctx); err != nil {
		return
	}
}
