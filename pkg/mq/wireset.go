package mq

import (
	"github.com/google/wire"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq/consumer"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq/producer"
)

var MQWireSet = wire.NewSet(
	consumer.NewConsumer,
	producer.NewClient,
)
