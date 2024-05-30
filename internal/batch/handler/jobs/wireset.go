package jobs

import "github.com/google/wire"

var JobWireSet = wire.NewSet(
	NewExecuteClearAllOrderEOD,
)
