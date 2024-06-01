//go:build wireinject

package author

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-coffee-shop-demo/internal/authorization/app"
	"github.com/namnv2496/go-coffee-shop-demo/internal/authorization/repo"
	"github.com/namnv2496/go-coffee-shop-demo/internal/authorization/service"
	configs "github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	data_access "github.com/namnv2496/go-coffee-shop-demo/pkg/data_access"
)

func Initialize(filePath configs.ConfigFilePath) (*app.App, func(), error) {

	wire.Build(
		configs.ConfigWireSet,
		data_access.DataWireSet,

		repo.RepoWireSet,
		service.ServiceWireSet,
		app.NewApp,
	)
	return nil, nil, nil
}
