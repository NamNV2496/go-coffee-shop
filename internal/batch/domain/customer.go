package domain

import "github.com/doug-martin/goqu/v9"

var (
	TabNameCustomer = goqu.T("customer")
)

type Customer struct {
	Id            int32  `db:"id"`
	Name          string `db:"name"`
	Age           int32  `db:"age"`
	Loyalty_point int32  `db:"loyalty_point"`
}
