package repo

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ItemRepo interface {
	GetAll(ctx context.Context, offset int32, limit int32) ([]domain.Item, error)
	GetByIdOrName(ctx context.Context, id int32, name string, offset int32, limit int32) ([]domain.Item, error)
	AddNewProduct(ctx context.Context, item domain.Item, img string) (int32, error)
}

type itemRepo struct {
	database *goqu.Database
}

func NewItemRepo(
	database *goqu.Database,
) ItemRepo {
	return &itemRepo{
		database: database,
	}
}

func (itemRepo *itemRepo) GetAll(
	ctx context.Context,
	offset int32,
	limit int32,
) ([]domain.Item, error) {

	itemList := make([]domain.Item, 0)
	err := itemRepo.database.
		Select().
		From(domain.TabNameItem).
		Where().
		Offset(uint(offset)*uint(limit)).
		Limit(uint(limit)).
		Executor().
		ScanStructsContext(ctx, &itemList)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create account")
	}
	return itemList, nil
}

func (itemRepo *itemRepo) GetByIdOrName(
	ctx context.Context,
	id int32,
	name string,
	offset int32,
	limit int32,
) ([]domain.Item, error) {

	// Build the conditions dynamically
	var conditions []goqu.Expression
	if id != 0 {
		conditions = append(conditions, goqu.C(domain.ColId).Eq(id))
	}
	if name != "" {
		conditions = append(conditions, goqu.C(domain.ColName).Like("%"+name+"%"))
	}
	query := itemRepo.database.
		From(domain.TabNameItem).
		Where(goqu.Or(conditions...)).
		Offset(uint(offset) * uint(limit)).
		Limit(uint(limit))

	itemList := make([]domain.Item, 0)
	err := query.Executor().ScanStructsContext(ctx, &itemList)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return itemList, nil
}

func (itemRepo *itemRepo) AddNewProduct(
	ctx context.Context,
	item domain.Item,
	img string,
) (int32, error) {
	query := itemRepo.database.
		Insert(domain.TabNameItem).
		Cols(domain.ColName, domain.ColPrice, domain.ColType, domain.ColImg).
		Vals(
			goqu.Vals{item.Name, item.Price, item.Type, img},
		)

	result, err := query.Executor().ExecContext(ctx)
	if err != nil {
		return 0, status.Error(codes.Internal, "failed to create order")
	}
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, status.Error(codes.Internal, "failed to create order")
	}
	return int32(lastInsertedID), nil
}
