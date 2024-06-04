package repo

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/domain"
)

type CustomerRepo interface {
	GetCustomer(ctx context.Context, customerId int32) (domain.Customer, error)
	GetCustomers(ctx context.Context) ([]domain.Customer, error)
}

type customerRepo struct {
	database *goqu.Database
}

func NewCustomerRepo(
	database *goqu.Database,
) CustomerRepo {

	return &customerRepo{
		database: database,
	}
}

func (customerrepo customerRepo) GetCustomer(ctx context.Context, customerId int32) (domain.Customer, error) {

	query := customerrepo.database.
		From(domain.TabNameCustomer).
		Where(goqu.C(domain.ColId).Eq(customerId))

	rows, err := query.Executor().QueryContext(context.Background())
	if err != nil {
		fmt.Println("Error executing query:", err)
	}

	var customer domain.Customer
	if rows.Next() {
		err = rows.Scan(&customer.Id, &customer.Name, &customer.Age, &customer.Loyalty_point)
		if err != nil {
			return customer, fmt.Errorf("error scanning row: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		return customer, fmt.Errorf("error during row iteration: %v", err)
	}
	return customer, nil
}

func (customerrepo customerRepo) GetCustomers(ctx context.Context) ([]domain.Customer, error) {

	query := customerrepo.database.
		From(domain.TabNameCustomer)

	var customers []domain.Customer
	if err := query.ScanStructs(&customers); err != nil {
		return nil, err
	}
	return customers, nil
}
