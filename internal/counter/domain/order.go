package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameOrder     = goqu.T("orders")
	TabNameOrderItem = goqu.T("order_item")
)

type Order struct {
	Id          int32     `db:"id"`
	Customer_id int32     `db:"customer_id"`
	TotalAmount int32     `db:"total_amount"`
	Status      int32     `db:"status"`
	CreatedDate time.Time `db:"created_date" goqu:"omitnil"`
}

type OrderItemListDto struct {
	OrderItems []OrderItem
	CustomerId int32
}

type OrderItem struct {
	Id          int32     `db:"id" goqu:"omitnil"`
	OrderId     int32     `db:"order_id" goqu:"omitnil"`
	ItemId      int32     `db:"item_id" goqu:"omitnil"`
	Quantity    int32     `db:"quantity" goqu:"omitnil"`
	Price       int32     `db:"price" goqu:"omitnil"`
	CreatedDate time.Time `db:"created_date" goqu:"omitnil"`
}

type OrderDto struct {
	Order      Order
	Customer   Customer
	OrderItems []OrderItem
}

type OrderDtoRes struct {
	Orders []OrderDto
}
