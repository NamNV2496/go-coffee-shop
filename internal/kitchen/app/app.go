package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/handler/consumers"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/service"
	"github.com/namnv2496/go-coffee-shop-demo/internal/mq/producer"
	"google.golang.org/grpc"
)

type AppInterface interface {
	Start(ctx context.Context)
	UpdateCoockedDoneStatus(ctx context.Context, req *gin.Context)
}

type App struct {
	grpc            *grpc.Server
	Producer        producer.Client
	ConsumerHandler consumers.ConsumerHandler
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
		Producer:        producer,
		ConsumerHandler: consumerHandler,
		KitchenService:  kitchenService,
	}
}

func (app App) Start(ctx context.Context) {
	app.ConsumerHandler.StartConsumerUp(ctx)
}

func (app App) UpdateCoockedDoneStatus(ctx context.Context, req *gin.Context) {

	req.JSON(http.StatusCreated, gin.H{"message": ""})
}
