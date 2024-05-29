//go:build wireinject

package counter

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/app"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/handler"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/handler/jobs"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/repo"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/cache"
	configs "github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	data_access "github.com/namnv2496/go-coffee-shop-demo/pkg/data_access"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq"
	"google.golang.org/grpc"
)

func Initialize(grpc *grpc.Server, filePath configs.ConfigFilePath) (*app.App, func(), error) {

	wire.Build(
		configs.ConfigWireSet,
		data_access.DataWireSet,
		mq.MQWireSet,
		cache.CacheWireSet,

		repo.RepoWireSet,
		service.ServiceWireSet,
		handler.HandlerWireSet,
		app.NewApp,
		jobs.JobWireSet,
	)
	return nil, nil, nil
}
