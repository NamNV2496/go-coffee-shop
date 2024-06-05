//go:build wireinject

package product

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/app"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/handler/router"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/repo"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	data_access "github.com/namnv2496/go-coffee-shop-demo/pkg/data_access"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/s3"
	"google.golang.org/grpc"
)

func Initialize(
	grpcServer *grpc.Server,
	filePath configs.ConfigFilePath,
) (*app.App, func(), error) {
	panic(wire.Build(
		configs.ConfigWireSet,
		data_access.DataWireSet,
		s3.FileWireSet,
		router.GrpcWireSet,
		repo.RepoWireSet,
		service.ServiceWireSet,
		app.NewApp,
	))

	return nil, nil, nil
}
