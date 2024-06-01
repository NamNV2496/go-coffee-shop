//go:build wireinject

package kitchen

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/app"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/handler"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/cache"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq"
)

func Initialize(filePath configs.ConfigFilePath) (*app.App, func(), error) {

	wire.Build(
		configs.ConfigWireSet,
		mq.MQWireSet,
		cache.CacheWireSet,

		handler.HandlerWireSet,
		service.ServiceWireSet,
		app.NewApp,
	)
	return nil, nil, nil
}
