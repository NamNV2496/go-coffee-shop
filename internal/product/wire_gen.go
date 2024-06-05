// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package product

import (
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/app"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/handler/router"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/repo"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/data_access"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/s3"
	"google.golang.org/grpc"
)

// Injectors from wire.go:

func Initialize(grpcServer *grpc.Server, filePath configs.ConfigFilePath) (*app.App, func(), error) {
	config, err := configs.GetConfigFromYaml(filePath)
	if err != nil {
		return nil, nil, err
	}
	db, cleanup, err := database.InitializeAndMigrateUpDB(config)
	if err != nil {
		return nil, nil, err
	}
	goquDatabase := database.InitializeGoquDB(db)
	itemRepo := repo.NewItemRepo(goquDatabase)
	configsS3 := config.S3
	s3Client := s3.NewS3Client(configsS3)
	productService := service.NewProductService(itemRepo, s3Client)
	productServiceServer := router.NewHandler(productService)
	productServer := router.NewGrpcRouterServer(config, grpcServer, productServiceServer)
	appApp := app.NewApp(productServer, productService)
	return appApp, func() {
		cleanup()
	}, nil
}
