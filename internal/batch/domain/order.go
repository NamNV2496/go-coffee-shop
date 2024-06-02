package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameOrderItem   = goqu.T("order_item")
	TabNameCreatedDate = "created_date"
)

type OrderItem struct {
	Id          int32     `db:"id" goqu:"omitnil"`
	OrderId     int32     `db:"order_id" goqu:"omitnil"`
	ItemId      int32     `db:"item_id" goqu:"omitnil"`
	Quantity    int32     `db:"quantity" goqu:"omitnil"`
	Price       int32     `db:"price" goqu:"omitnil"`
	CreatedDate time.Time `db:"created_date" goqu:"omitnil"`
}
