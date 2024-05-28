package domain

import "github.com/doug-martin/goqu/v9"

var (
	TabNameOrder     = goqu.T("orders")
	TabNameOrderItem = goqu.T("order_item")
)

type Order struct {
	Id          int32 `db:"id"`
	Customer_id int32 `db:"customer_id"`
	TotalAmount int32 `db:"total_amount"`
	Status      int32 `db:"status"`
}

type OrderItemListDto struct {
	OrderItems []OrderItem
	CustomerId int32
}

type OrderItem struct {
	Id       int32 `db:"id"`
	OrderId  int32 `db:"order_id"`
	ItemId   int32 `db:"item_id"`
	Quantity int32 `db:"quantity"`
	Price    int32 `db:"price"`
}

type OrderDto struct {
	Order      Order
	Customer   Customer
	OrderItems []OrderItem
}

type OrderDtoRes struct {
	Orders []OrderDto
}
