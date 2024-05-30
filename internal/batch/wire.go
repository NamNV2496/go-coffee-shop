//go:build wireinject

package batch

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/app"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/handler"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/repo"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/cache"
	configs "github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	data_access "github.com/namnv2496/go-coffee-shop-demo/pkg/data_access"
)

func Initialize(filePath configs.ConfigFilePath) (*app.App, func(), error) {

	wire.Build(
		configs.ConfigWireSet,
		data_access.DataWireSet,
		cache.CacheWireSet,
		handler.HandlerWireSet,
		repo.RepoWireSet,
		service.ServiceWireSet,
		app.NewApp,
	)
	return nil, nil, nil
}
