package repo

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/data_access/enums"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, orders []domain.OrderItem, customerId int32) (int32, error)
	UpdateStatusOrder(ctx context.Context, orderId int32, status int32) error
	GetOrderById(ctx context.Context, orderId int32) (domain.Order, error)
	GetOrderByCustomerId(ctx context.Context, customerId int32) ([]domain.Order, error)
	GetOrders(ctx context.Context) ([]domain.Order, error)
	GetOrderItem(ctx context.Context, orderIds []int32) ([]domain.OrderItem, error)
}

type orderRepo struct {
	database *goqu.Database
}

func NewOrderRepo(
	database *goqu.Database,
) OrderRepo {

	return &orderRepo{
		database: database,
	}
}

func (order orderRepo) CreateOrder(ctx context.Context, orders []domain.OrderItem, customerId int32) (int32, error) {
	totalAmount := 0
	for _, order := range orders {
		totalAmount += int(order.Quantity) * int(order.Price)
	}

	query := order.database.
		Insert(domain.TabNameOrder).
		Cols(domain.ColCustomerId, domain.ColTotalAmount, domain.ColStatus).
		Vals(
			goqu.Vals{customerId, totalAmount, enums.Ordering.Code},
		)

	result, err := query.Executor().ExecContext(ctx)
	if err != nil {
		return 0, status.Error(codes.Internal, "failed to create order")
	}
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, status.Error(codes.Internal, "failed to create order")
	}
	addOrderItems(order.database, orders, int32(lastInsertedID))

	return int32(lastInsertedID), nil
}

func addOrderItems(database *goqu.Database, orders []domain.OrderItem, id int32) error {

	var vals [][]interface{}
	for _, order := range orders {
		vals = append(vals, []interface{}{id, order.ItemId, order.Quantity, order.Price})
	}
	query := database.
		Insert(domain.TabNameOrderItem).
		Cols(domain.ColOrderId, domain.ColItemId, domain.ColQuantity, domain.ColPrice).
		Vals(vals...)

	// query := order.database.
	// 	Insert(domain.TabNameOrder).
	// 	Rows(goqu.Record{
	// 		domain.ColId:          1,
	// 		domain.ColTotalAmount: totalAmount,
	// 	})

	_, err := query.Executor().ExecContext(context.Background())
	if err != nil {
		return status.Error(codes.Internal, "failed to create order")
	}
	return nil
}

func (order orderRepo) UpdateStatusOrder(ctx context.Context, orderId int32, status int32) error {
	newStatus, err := enums.FindStatus(status)
	if err != nil {
		fmt.Println("Status if invalid")
	}
	query := order.database.
		From(domain.TabNameOrder).
		Where(goqu.C(domain.ColId).Eq(orderId))

	var orderUpdate domain.Order
	found, err := query.ScanStruct(&orderUpdate)
	if err != nil {
		fmt.Println("Failed to get order by id:", err)
	}
	if !found {
		fmt.Println("No order found")
	}
	orderUpdate.Status = newStatus.Code

	updateQuery := order.database.
		Update(domain.TabNameOrder).
		Set(goqu.Record{domain.ColStatus: orderUpdate.Status}).
		Where(goqu.C(domain.ColId).Eq(orderId))

	_, err = updateQuery.Executor().Exec()
	if err != nil {
		fmt.Println("Failed to update order status:", err)
	}
	return nil
}

func (order orderRepo) GetOrderById(ctx context.Context, orderId int32) (domain.Order, error) {

	/*
		query := order.database.
			From(domain.TabNameOrder).
			Where(goqu.C(domain.ColId).Eq(orderId))
		result, err := query.Executor().Scanner()
		if err != nil {
			fmt.Println("Fail to get order by id")
		}

		var res domain.Order
		defer result.Close()

		for result.Next() {
			err = result.ScanStruct(&res)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		return res, nil
	*/

	query := order.database.
		From(domain.TabNameOrder).
		Where(
			goqu.C(domain.ColId).Eq(orderId),
		)
	var orders []domain.Order
	query.ScanStructs(&orders)

	return orders[0], nil
}

func (order orderRepo) GetOrderByCustomerId(ctx context.Context, customerId int32) ([]domain.Order, error) {
	query := order.database.
		From(domain.TabNameOrder).
		Where(
			goqu.C(domain.ColCustomerId).Eq(customerId),
		)
	var orders []domain.Order
	query.ScanStructs(&orders)

	return orders, nil
}

func (order orderRepo) GetOrderItem(ctx context.Context, orderIds []int32) ([]domain.OrderItem, error) {

	query := order.database.
		From(domain.TabNameOrderItem).
		Where(
			goqu.C(domain.ColOrderId).In(orderIds),
		)
	var orders []domain.OrderItem
	query.ScanStructs(&orders)

	return orders, nil
}

func (order orderRepo) GetOrders(ctx context.Context) ([]domain.Order, error) {

	query := order.database.
		From(domain.TabNameOrder)

	var orders []domain.Order
	query.ScanStructs(&orders)
	return orders, nil
}
