package main

import (
	"context"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/app"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
	"google.golang.org/grpc"
)

func main() {

	grpc := grpc.NewServer()
	defer grpc.GracefulStop()

	app, cleanup, err := counter.Initialize(grpc, "")
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
	// create order and save to cache
	r.POST("/api/v1/createOrder", func(req *gin.Context) { app.CreateOrder(ctx, req) })
	// submit order to get data in cache and save it to DB
	r.POST("/api/v1/submitOrder", func(req *gin.Context) { app.SubmitOrder(ctx, req) })
	/// update status of order (only for cancel)
	// r.PUT("/api/v1/updateOrderStatus", func(req *gin.Context) { app.UpdateOrderStatus(ctx, req) })
	// get item by id or name
	r.GET("/api/v1/getItems", func(req *gin.Context) { app.GetItem(ctx, req) })
	// view order by orderId or customerId
	r.GET("/api/v1/getOrders", func(req *gin.Context) { app.GetOrders(ctx, req) })
}
