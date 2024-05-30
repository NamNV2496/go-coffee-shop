package handler

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-coffee-shop-demo/internal/kitchen/handler/consumers"
)

var HandlerWireSet = wire.NewSet(
	consumers.NewHandler,
)
