package service

import (
	"github.com/google/wire"
)

var ServiceWireSet = wire.NewSet(
	NewUserService,
)
