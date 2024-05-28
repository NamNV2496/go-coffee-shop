package consumers

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/service"
// 	"github.com/namnv2496/go-coffee-shop-demo/internal/mq"
// 	"github.com/namnv2496/go-coffee-shop-demo/internal/mq/consumer"
// )

// type ConsumerHandler interface {
// 	StartConsumerUp(ctx context.Context) error
// }

// type consumerHandler struct {
// 	OrderService service.OrderService
// 	Consumer     consumer.Consumer
// }

// func NewHandler(
// 	orderService service.OrderService,
// 	consumer consumer.Consumer,
// ) ConsumerHandler {
// 	return &consumerHandler{
// 		OrderService: orderService,
// 		Consumer:     consumer,
// 	}
// }

// func (c consumerHandler) StartConsumerUp(ctx context.Context) error {
// 	fmt.Println("Add consumer for topic: ", mq.TOPIC_ORDER_STATUS_UPDATE)
// 	c.Consumer.RegisterHandler(
// 		mq.TOPIC_ORDER_STATUS_UPDATE,
// 		func(ctx context.Context, queueName string, payload []byte) error {
// 			fmt.Println("listen from queue: " + queueName + ". Data: " + string(payload))
// 			var event mq.OrderRequestToKitchen
// 			if err := json.Unmarshal(payload, &event); err != nil {
// 				return err
// 			}
// 			return nil
// 		},
// 	)
// 	return nil
// }
