package consumers

import (
	"context"
	"fmt"

	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq/consumer"
)

type ConsumerHandler interface {
	StartConsumerUp(ctx context.Context) error
}

type consumerHandler struct {
	consumer consumer.Consumer
}

func NewKafkaHandler(
	consumer consumer.Consumer,
) ConsumerHandler {
	return &consumerHandler{
		consumer: consumer,
	}
}

func (c consumerHandler) StartConsumerUp(ctx context.Context) error {
	fmt.Println("Add consumer for topic: ", mq.TOPIC_PROCESS_COOK)
	c.consumer.RegisterHandler(
		mq.TOPIC_PROCESS_COOK,
		func(ctx context.Context, queueName string, payload []byte) error {
			fmt.Println("listen from queue: " + queueName + ". Data: " + string(payload))
			fmt.Println("have new order. Please check and do it!")
			return nil
		},
	)
	return nil
}
