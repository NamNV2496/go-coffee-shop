package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/app"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	app, cleanup, err := product.Initialize(server, "")

	ctx := context.Background()
	r := ginSetup(ctx)
	rounting(ctx, app, r)
	go func() {
		app.Start(ctx)
	}()
	r.Run(":8080")
	defer cleanup()
	if err != nil {
		panic("Error when start")
	}
}

func ginSetup(ctx context.Context) *gin.Engine {
	r := gin.Default()
	return r
}

func rounting(ctx context.Context, app *app.App, r *gin.Engine) {
	r.POST("/api/v1/product", func(req *gin.Context) { app.AddNewProduct(ctx, req) })
	r.GET("/api/v1/product", func(req *gin.Context) { app.GetImageInMinio(ctx, req) })
}
