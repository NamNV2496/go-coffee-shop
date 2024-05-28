package consumers

import (
	"context"
	"fmt"

	"github.com/namnv2496/go-coffee-shop-demo/internal/mq"
	"github.com/namnv2496/go-coffee-shop-demo/internal/mq/consumer"
)

type ConsumerHandler interface {
	StartConsumerUp(ctx context.Context) error
}

type consumerHandler struct {
	Consumer consumer.Consumer
}

func NewHandler(
	consumer consumer.Consumer,
) ConsumerHandler {
	return &consumerHandler{
		Consumer: consumer,
	}
}

func (c consumerHandler) StartConsumerUp(ctx context.Context) error {
	fmt.Println("Add consumer for topic: ", mq.TOPIC_PROCESS_COOK)
	c.Consumer.RegisterHandler(
		mq.TOPIC_ORDER_STATUS_UPDATE,
		func(ctx context.Context, queueName string, payload []byte) error {
			fmt.Println("listen from queue: " + queueName + ". Data: " + string(payload))
			fmt.Println("have new order. Please check and do it!")
			return nil
		},
	)
	return nil
}
