package mq

import (
	"errors"
	"time"
)

type OrderRequestToKitchen struct {
	OrderId int32 `json:"orderId"`
}

var (
	REDIS_KEY_ORDER = "kitchen"
	RedisTTL        = 86400 * time.Second
)

type CookStatus int32

const (
	Processing CookStatus = 0
	Canceled   CookStatus = 1
	Done       CookStatus = 2
)

var AllStatuses = []CookStatus{Processing, Canceled, Done}

func FindCockStatus(code int32) (CookStatus, error) {
	for _, status := range AllStatuses {
		if int32(status) == code {
			return status, nil
		}
	}
	return CookStatus(0), errors.New("order status not found")
}

type RedisOrder struct {
	ItemId   int32      `json:"itemId"`
	Quantity int32      `json:"quantity"`
	Price    int32      `json:"price"`
	Status   CookStatus `json:"status"`
	Finished int32      `json:"finshed"`
}

type RedisOrderDTO struct {
	CustomerId  int32 `json:"customerId"`
	RedisOrders []RedisOrder
}
