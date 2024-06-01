package main

import (
	"context"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/app"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/ratelimit"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/security"
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
		// add ratelimit
		// 100: the maximum frequency of some events. 100 request per second
		// 500: The limiter can handle a burst of up to 500 requests at once. This means if there has been no traffic for a while, the limiter can accumulate up to 500 "tokens" and allow a burst of up to 50 requests to be processed immediately. After the burst, the limiter will revert to allowing 1 request per second.
		rl := ratelimit.NewIPRateLimiter(100, 500)
		http.ListenAndServe(":8081", ratelimit.LimitMiddleware(r, rl))
		// r.Run(":8081")
	}()
	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
}

func SetupGin() *gin.Engine {

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

func routing(app *app.App, r *gin.Engine, ctx context.Context) {

	// create order and save to cache
	counterRoutes := r.Group("/api/v1")
	counterRoutes.Use(security.JWTAuthWithRole([]string{"counter", "admin"}))
	counterRoutes.POST("/createOrder", func(req *gin.Context) { app.CreateOrder(ctx, req) })
	// get item by id or name
	counterRoutes.GET("/getItems", func(req *gin.Context) { app.GetItem(ctx, req) })
	// submit order to get data in cache and save it to DB
	counterRoutes.POST("/submitOrder", func(req *gin.Context) { app.SubmitOrder(ctx, req) })
	// view order by orderId or customerId
	memberRoutes := r.Group("/api/v1/getSuccessOrders")
	memberRoutes.Use(security.JWTAuthWithRole([]string{"counter", "admin"}))
	memberRoutes.GET("", func(req *gin.Context) { app.GetOrders(ctx, req) })

	// update status of order (only for cancel)
	// publicRoutes.PUT("/updateOrderStatus", func(req *gin.Context) { app.UpdateOrderStatus(ctx, req) })
}
