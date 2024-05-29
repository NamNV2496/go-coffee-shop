package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/handler/consumers"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq/producer"
	"google.golang.org/grpc"
)

type AppInterface interface {
	Start() error
	UpdateCoockedDoneStatus(ctx context.Context, req *gin.Context)
}

type App struct {
	grpc            *grpc.Server
	producer        producer.Client
	consumerHandler consumers.ConsumerHandler
	KitchenService  service.KitchenService
}

func NewApp(
	grpc *grpc.Server,
	producer producer.Client,
	consumerHandler consumers.ConsumerHandler,
	kitchenService service.KitchenService,
) *App {
	return &App{
		grpc:            grpc,
		producer:        producer,
		consumerHandler: consumerHandler,
		KitchenService:  kitchenService,
	}
}

func (app App) Start() error {

	go func() {
		app.consumerHandler.StartConsumerUp(context.Background())
	}()
	return nil
}

func (app App) UpdateCoockedDoneStatus(ctx context.Context, req *gin.Context) {

	req.JSON(http.StatusCreated, gin.H{"message": ""})
}
