package main

import (
	"context"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/app"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/security"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	app, cleanup, err := product.Initialize(server, "")
	if err != nil {
		return
	}
	defer cleanup()

	go func() {
		ctx := context.Background()
		r := ginSetup()
		rounting(ctx, app, r)
		if err := http.ListenAndServe(":8080", r); err != nil {
			return
		}
		// r.Run(":8080")

	}()
	err = app.Start()
	if err != nil {
		panic("Error when start")
	}
	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
}

func ginSetup() *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Respond with 200 OK status
		c.Status(http.StatusOK)
	})
	return r
}

func rounting(ctx context.Context, app *app.App, r *gin.Engine) {
	publicRoutes := r.Group("/api/v1/product")
	publicRoutes.Use(security.JWTAuthWithRole([]string{"admin"}))
	publicRoutes.POST("", func(req *gin.Context) { app.AddNewProduct(ctx, req) })
	publicRoutes.GET("", func(req *gin.Context) { app.GetImageInMinio(ctx, req) })
}
