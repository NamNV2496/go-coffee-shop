package router

import (
	"github.com/google/wire"
)

var GrpcWireSet = wire.NewSet(
	NewHandler,
	NewGrpcRouterServer,
)
