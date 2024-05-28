package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/app"
	"github.com/namnv2496/go-coffee-shop-demo/internal/mq"
	"google.golang.org/grpc"
)

func main() {
	grpc := grpc.NewServer()
	app, cleanup, err := kitchen.Initialize(grpc, "")
	if err != nil {
		panic("fail to start kitchen")
	}
	ctx := context.Background()
	go func() {
		app.Start(ctx)
	}()
	defer cleanup()
	r := SetupGin()
	rounting(ctx, app, r)
	r.Run(":8082")
}

func SetupGin() *gin.Engine {
	r := gin.Default()
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
		finishNumber, err := mq.FindCockStatus(int32(finishedConv))
		if err != nil {
			panic("Invalid input")
		}
		app.KitchenService.UpdateStatusOrderToRedis(ctx, int32(idConv), int32(itemConv), int32(finishNumber))
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
