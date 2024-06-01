package app

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/handler/consumers"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq/producer"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
)

type AppInterface interface {
	Start() error
	UpdateOrderStatus(ctx context.Context, req *gin.Context)
	GetOrdersByCustomerId(ctx context.Context, req *gin.Context)
	GetOrders(ctx context.Context, req *gin.Context)
}

type App struct {
	producer       producer.Client
	kafkaHander    consumers.ConsumerHandler
	KitchenService service.KitchenService
}

func NewApp(
	producer producer.Client,
	kafkaHander consumers.ConsumerHandler,
	kitchenService service.KitchenService,
) *App {
	return &App{
		producer:       producer,
		kafkaHander:    kafkaHander,
		KitchenService: kitchenService,
	}
}

func (app App) Start() error {

	go func() {
		app.kafkaHander.StartConsumerUp(context.Background())
	}()
	return nil
}

func (app App) UpdateOrderStatus(ctx context.Context, req *gin.Context) {
	id := req.Query("customerId")
	idConv, err := strconv.Atoi(id)
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, "Invalid input")
		return
	}
	itemId := req.Query("itemId")
	itemConv, err := strconv.Atoi(itemId)
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, "Invalid input")
		return
	}
	cockStatus := req.Query("finished")
	finishedConv, err := strconv.Atoi(cockStatus)
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, "Invalid input")
		return
	}
	if err := app.KitchenService.UpdateStatusOrderToRedis(ctx, int32(idConv), int32(itemConv), int32(finishedConv)); err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	utils.WrapperResponse(req, http.StatusOK, "")
}

func (app App) GetOrdersByCustomerId(ctx context.Context, req *gin.Context) {
	id := req.Query("customerId")
	idConv, err := strconv.Atoi(id)
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, "invalid input")
		return
	}
	data, err := app.KitchenService.GetCustomerOrderInRedis(ctx, int32(idConv))
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	utils.WrapperResponse(req, http.StatusOK, data)
}

func (app App) GetOrders(ctx context.Context, req *gin.Context) {

	data, err := app.KitchenService.GetOrderInRedis(ctx)
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	utils.WrapperResponse(req, http.StatusOK, data)
}
