package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/handler/jobs"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
)

type AppInterface interface {
	Start() error
	TriggerJob(ctx context.Context, req *gin.Context)
}

type App struct {
	batchService service.BatchService
	jobs         jobs.ClearAllOrderEOD
}

func NewApp(
	batchService service.BatchService,
	jobs jobs.ClearAllOrderEOD,
) *App {
	return &App{
		batchService: batchService,
		jobs:         jobs,
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
	if err := app.jobs.StartClearJobEOD(ctx); err != nil {
		return err
	}
	return nil
}

func (app App) TriggerJob(ctx context.Context, req *gin.Context) {
	fmt.Println("Trigger job report")
	if err := app.jobs.Run(ctx); err != nil {
		return
	}
	utils.WrapperResponse(req, http.StatusOK, "")
}
