package handler

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/handler/jobs"
)

var HandlerWireSet = wire.NewSet(
	jobs.JobWireSet,
)
