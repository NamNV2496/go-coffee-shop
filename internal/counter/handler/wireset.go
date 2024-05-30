package handler

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/handler/router"
)

var HandlerWireSet = wire.NewSet(
	router.NewGRPCProductClient,
)
