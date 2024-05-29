package main

import (
	"context"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/app"
	"github.com/namnv2496/go-coffee-shop-demo/internal/utils"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	app, cleanup, err := product.Initialize(server, "")
	defer cleanup()

	go func() {
		ctx := context.Background()
		r := ginSetup()
		rounting(ctx, app, r)
		r.Run(":8080")

	}()
	err = app.Start()
	if err != nil {
		panic("Error when start")
	}
	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
}

func ginSetup() *gin.Engine {
	r := gin.Default()
	return r
}

func rounting(ctx context.Context, app *app.App, r *gin.Engine) {
	r.POST("/api/v1/product", func(req *gin.Context) { app.AddNewProduct(ctx, req) })
	r.GET("/api/v1/product", func(req *gin.Context) { app.GetImageInMinio(ctx, req) })
}
