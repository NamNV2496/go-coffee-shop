package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/app"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	grpc := grpc.NewServer()
	defer grpc.GracefulStop()
	app, cleanup, err := kitchen.Initialize(grpc, "")
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
		r.Run(":8082")
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

func rounting(ctx context.Context, app *app.App, r *gin.Engine) {
	// Update order status by customerId and itemId
	r.PUT("/api/v1/updateOrderStatus", func(req *gin.Context) {
		id := req.Query("customerId")
		idConv, err := strconv.Atoi(id)
		if err != nil {
			panic("Invalid input")
		}
		itemId := req.Query("itemId")
		itemConv, err := strconv.Atoi(itemId)
		if err != nil {
			panic("Invalid input")
		}
		cockStatus := req.Query("finished")
		finishedConv, err := strconv.Atoi(cockStatus)
		if err != nil {
			panic("Invalid input")
		}
		app.KitchenService.UpdateStatusOrderToRedis(ctx, int32(idConv), int32(itemConv), int32(finishedConv))
	})
	r.GET("/api/v1/getOrdersByCustomerId", func(req *gin.Context) {
		id := req.Query("customerId")
		idConv, err := strconv.Atoi(id)
		if err != nil {
			panic("Invalid input")
		}
		data, err := app.KitchenService.GetCustomerOrderInRedis(ctx, int32(idConv))
		if err != nil {
			panic("Error when get order")
		}
		fmt.Println(data)
		req.JSON(http.StatusCreated, gin.H{"message": data})
	})
	r.GET("/api/v1/getOrders", func(req *gin.Context) {

		data, err := app.KitchenService.GetOrderInRedis(ctx)
		if err != nil {
			panic("Error when get order")
		}
		req.JSON(http.StatusCreated, gin.H{"message": data})
	})
}
