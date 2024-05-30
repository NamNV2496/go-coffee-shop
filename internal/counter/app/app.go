package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/handler/jobs"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/handler/router"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq/producer"
)

type AppInterface interface {
	Start() error
	CreateOrder(ctx context.Context, req *gin.Context)
	SubmitOrder(ctx context.Context, req *gin.Context)
	// UpdateOrderStatus(ctx context.Context, req *gin.Context)
	GetItem(ctx context.Context, req *gin.Context)
	GetOrders(ctx context.Context, req *gin.Context)
}

type App struct {
	orderService service.OrderService
	grpcClient   router.ProductGRPCClient
	producer     producer.Client
	// ConsumerHandler consumers.ConsumerHandler
	jobs       jobs.ClearAllOrderEOD
	cronConfig configs.Cron
}

func NewApp(
	orderService service.OrderService,
	grpcClient router.ProductGRPCClient,
	producer producer.Client,
	// consumerHandler consumers.ConsumerHandler,
	jobs jobs.ClearAllOrderEOD,
	cronConfig configs.Cron,
) *App {
	return &App{
		orderService: orderService,
		grpcClient:   grpcClient,
		producer:     producer,
		// ConsumerHandler: consumerHandler,
		jobs:       jobs,
		cronConfig: cronConfig,
	}
}

func (app App) Start() error {

	// app.ConsumerHandler.StartConsumerUp(ctx)
	go func() {
		err := app.startScheduler(context.Background())
		if err != nil {
			panic("failed to start scheduler")
		}
	}()
	return nil
}

func (app App) startScheduler(ctx context.Context) error {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return errors.New("failed to initialize scheduler")
	}

	hour, err := strconv.Atoi(app.cronConfig.ClearAllOrder.Hour)
	if err != nil {
		return err
	}
	minute, err := strconv.Atoi(app.cronConfig.ClearAllOrder.Minute)
	if err != nil {
		return err
	}

	j, err := scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(uint(hour), uint(minute), 0),
			),
		),
		gocron.NewTask(func() {
			if err := app.jobs.Run(context.Background()); err != nil {
				fmt.Println("failed to run execute all pending download task job")
			}
		}),
	)
	if err != nil {
		fmt.Println("failed to schedule execute all pending download task job")
		return err
	}
	fmt.Println("Daily scheduler job id: ", j.ID())

	scheduler.Start()
	return nil
}

// func (app App) UpdateOrderStatus(ctx context.Context, req *gin.Context) {
// 	orderId, err := strconv.Atoi(req.Query("orderId"))
// 	if err != nil {
// 		orderId = 0
// 	}
// 	status, err := strconv.Atoi(req.Query("status"))
// 	if err != nil || status != int(mq.Canceled) {
// 		panic("Status is invalide")
// 	}
// 	app.OrderService.UpdateStatusOrder(context.Background(), int32(orderId), int32(status))
// 	req.JSON(http.StatusCreated, gin.H{"message": "Thành công"})
// }

func (app App) CreateOrder(ctx context.Context, req *gin.Context) {

	var orderList domain.OrderItemListDto

	if err := req.BindJSON(&orderList); err != nil {
		req.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("Failed: ", err)
		return
	}

	err := app.orderService.CreateOrder(context.Background(), orderList.OrderItems, orderList.CustomerId)
	if err != nil {
		panic(fmt.Sprintf("Insert fail: %v", err))
	}
	// trigger to kitchen
	data, err := json.Marshal(orderList)
	if err != nil {
		fmt.Println("Error marshall data to send to kitchen")
	}
	go func() {
		fmt.Println("Call trigger to kitchen: ", orderList)
		app.producer.Produce(context.Background(), mq.TOPIC_PROCESS_COOK, data)
	}()
	req.JSON(http.StatusCreated, gin.H{"message": err})
}

func (app App) SubmitOrder(ctx context.Context, req *gin.Context) {
	queryId := req.Query("customerId")
	customerId, err := strconv.Atoi(queryId)
	if err != nil {
		panic(fmt.Sprintf("Insert fail: %v", err))
	}
	id, err := app.orderService.SubmitOrder(context.Background(), int32(customerId))
	if err != nil {
		panic(fmt.Sprintf("Insert fail: %v", err))
	}
	req.JSON(http.StatusCreated, gin.H{"message": id})
}

func (app App) GetItem(ctx context.Context, req *gin.Context) {
	// get from parameters
	id := req.Query("id")
	name := req.Query("name")
	num, _ := strconv.Atoi(id)
	page := 0
	size := 50
	pageQuery := req.Query("page")
	page, _ = strconv.Atoi(pageQuery)
	sizeQuery := req.Query("size")
	size, _ = strconv.Atoi(sizeQuery)
	items, err := app.grpcClient.GetProductByIdOrName(int32(num), name, int32(page), int32(size))
	if err != nil {
		panic("get items fail")
	}
	req.JSON(http.StatusCreated, gin.H{"message": items})
}

func (app App) GetOrders(ctx context.Context, req *gin.Context) {

	orderId, err := strconv.Atoi(req.Query("orderId"))
	if err != nil {
		orderId = 0
	}
	customerId, err := strconv.Atoi(req.Query("customerId"))
	if err != nil {
		customerId = 0
	}

	orderDtoRes, err := app.orderService.GetOrder(context.Background(), int32(orderId), int32(customerId))
	if err != nil {
		fmt.Println("Error")
	}
	req.JSON(http.StatusCreated, gin.H{"message": orderDtoRes})
}
