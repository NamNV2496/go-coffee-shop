package jobs

import (
	"context"
	"fmt"

	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/service"
)

type ClearAllOrderEOD interface {
	Run(context.Context) error
}

type clearAllOrderEOD struct {
	orderService service.OrderService
}

func NewExecuteClearAllOrderEOD(
	orderService service.OrderService,
) ClearAllOrderEOD {
	return &clearAllOrderEOD{
		orderService: orderService,
	}
}

func (j clearAllOrderEOD) Run(ctx context.Context) error {
	fmt.Println("Trigger clearAllOrderEOD")
	return j.orderService.ClearAllOrderEOD(ctx)
}
