package mq

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-coffee-shop-demo/internal/mq/consumer"
	"github.com/namnv2496/go-coffee-shop-demo/internal/mq/producer"
)

var MQWireSet = wire.NewSet(
	consumer.NewConsumer,
	producer.NewClient,
)
