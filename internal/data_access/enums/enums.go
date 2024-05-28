package enums

import "errors"

type OrderStatus struct {
	Code int32
	Msg  string
}

var (
	Ordering   = OrderStatus{0, "ordered"}    // order dang order
	Paid       = OrderStatus{1, "paid"}       // da thanh toan -> chuyen den bep de nau
	Processing = OrderStatus{2, "processing"} // bep dang nau
	Shipping   = OrderStatus{3, "Shipping"}   // dang ship
	Done       = OrderStatus{4, "done"}       // ship xong
	Canceled   = OrderStatus{5, "canceled"}   // huy bo
)

var AllStatuses = []OrderStatus{Ordering, Paid, Processing, Shipping, Done, Canceled}

func FindStatus(code int32) (OrderStatus, error) {
	for _, status := range AllStatuses {
		if status.Code == code {
			return status, nil
		}
	}
	return OrderStatus{}, errors.New("order status not found")
}
