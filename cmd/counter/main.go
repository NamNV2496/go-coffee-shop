package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/app"
	"google.golang.org/grpc"
)

func main() {

	server := grpc.NewServer()
	defer server.GracefulStop()

	app, cleanup, err := counter.Initialize(server, "")
	if err != nil {
		panic("Error start app")
	}
	defer cleanup()
	r := SetupGin()

	ctx := context.Background()
	app.Start(ctx)

	routing(app, r, ctx)
	r.Run(":8081")
}

func SetupGin() *gin.Engine {

	r := gin.Default()
	return r
}

func routing(app *app.App, r *gin.Engine, ctx context.Context) {
	// create order and save to cache
	r.POST("/api/v1/createOrder", func(req *gin.Context) { app.CreateOrder(ctx, req) })
	// submit order to get data in cache and save it to DB
	r.POST("/api/v1/submitOrder", func(req *gin.Context) { app.SubmitOrder(ctx, req) })
	/// update status of order (only for cancel)
	r.PUT("/api/v1/updateOrderStatus", func(req *gin.Context) { app.UpdateOrderStatus(ctx, req) })
	// get item by id or name
	r.GET("/api/v1/getItems", func(req *gin.Context) { app.GetItem(ctx, req) })
	// view order by orderId or customerId
	r.GET("/api/v1/getOrders", func(req *gin.Context) { app.GetOrders(ctx, req) })
}
