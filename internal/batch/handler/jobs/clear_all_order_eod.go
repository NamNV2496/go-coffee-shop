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
	batchService service.BatchService
}

func NewExecuteClearAllOrderEOD(
	batchService service.BatchService,
) ClearAllOrderEOD {
	return &clearAllOrderEOD{
		batchService: batchService,
	}
}

func (j clearAllOrderEOD) Run(ctx context.Context) error {
	fmt.Println("Trigger clearAllOrderEOD")
	return j.batchService.ClearAllOrderEOD(ctx)
}
