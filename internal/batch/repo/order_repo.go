package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/domain"
)

type OrderRepo interface {
	GetOrderItems(ctx context.Context) ([]domain.OrderItem, error)
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

func (order orderRepo) GetOrderItems(ctx context.Context) ([]domain.OrderItem, error) {

	year, month, day := time.Now().Date()
	query := order.database.
		From(domain.TabNameOrderItem).
		Where(
			goqu.C(domain.TabNameCreatedDate).Gte(time.Date(year, month, day, 0, 0, 0, 0, time.Now().UTC().Location())),
		)
	fmt.Println(query.ToSQL())
	var orderItems []domain.OrderItem
	if err := query.Executor().ScanStructsContext(ctx, &orderItems); err != nil {
		return nil, err
	}
	orderMap := make(map[int32]domain.OrderItem)

	for _, item := range orderItems {
		if v, ok := orderMap[item.ItemId]; ok {
			itemUpdate := v
			itemUpdate.Quantity += item.Quantity
			orderMap[item.ItemId] = itemUpdate
		} else {
			orderMap[item.ItemId] = item
		}
	}
	var finalOrderItems []domain.OrderItem
	for _, v := range orderMap {
		finalOrderItems = append(finalOrderItems, v)
	}
	return finalOrderItems, nil
}
