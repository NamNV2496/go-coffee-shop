package app

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/handler/jobs"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
)

type AppInterface interface {
	Start() error
	TriggerJob(ctx context.Context, req *gin.Context)
}

type App struct {
	jobs       jobs.ClearAllOrderEOD
	cronConfig configs.Cron
}

func NewApp(
	jobs jobs.ClearAllOrderEOD,
	cronConfig configs.Cron,
) *App {
	return &App{
		jobs:       jobs,
		cronConfig: cronConfig,
	}
}

func (app App) Start() error {

	// app.ConsumerHandler.StartConsumerUp(ctx)
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
	// add generate pdf file
}
