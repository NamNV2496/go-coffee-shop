package main

import (
	"context"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/app"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
)

func main() {

	app, cleanup, err := batch.Initialize("")
	if err != nil {
		panic("Error start app")
	}
	defer cleanup()
	r := SetupGin()

	app.Start()

	go func() {
		routing(app, r, context.Background())
		r.Run(":8081")
	}()
	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
}

func SetupGin() *gin.Engine {

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Respond with 200 OK status
		c.Status(http.StatusOK)
	})
	return r
}

func routing(app *app.App, r *gin.Engine, ctx context.Context) {

	r.POST("/api/v1/triggerJob", func(req *gin.Context) { app.TriggerJob(ctx, req) })
}
