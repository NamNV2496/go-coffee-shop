package main

import (
	"context"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/app"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/security"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
)

func main() {
	app, cleanup, err := kitchen.Initialize("")
	if err != nil {
		panic("fail to start kitchen")
	}
	defer cleanup()

	err = app.Start()
	if err != nil {
		panic("Error when start")
	}
	go func() {
		r := SetupGin()
		rounting(context.Background(), app, r)
		// Run Gin server
		if err := http.ListenAndServe(":8082", r); err != nil {
			return
		}
		// r.Run(":8082")
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

func rounting(ctx context.Context, app *app.App, r *gin.Engine) {
	// Update order status by customerId and itemId

	kitchenRoutes := r.Group("/api/v1")
	kitchenRoutes.Use(security.JWTAuthWithRole([]string{"kitchen", "admin"}))
	kitchenRoutes.PUT("/updateOrderStatus", func(req *gin.Context) { app.UpdateOrderStatus(ctx, req) })

	memberRoutes := r.Group("/api/v1")
	memberRoutes.Use(security.JWTAuthWithRole([]string{"kitchen", "admin", "counter"}))
	memberRoutes.GET("/getOrdersByCustomerId", func(req *gin.Context) { app.GetOrdersByCustomerId(ctx, req) })
	memberRoutes.GET("/getOrders", func(req *gin.Context) { app.GetOrders(ctx, req) })
}
