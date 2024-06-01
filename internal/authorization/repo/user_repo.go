package repo

import (
	"context"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/go-coffee-shop-demo/internal/authorization/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	GetUser(ctx context.Context, userId string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
	InactiveUser(ctx context.Context, userId string) error
}

type userRepo struct {
	database *goqu.Database
}

func NewUserRepo(
	database *goqu.Database,
) UserRepo {
	return &userRepo{
		database: database,
	}
}

func (u userRepo) CreateUser(ctx context.Context, user domain.User) (int, error) {

	user.CreatedDate = time.Now()
	query := u.database.
		Insert(domain.TabNameUser).
		Rows(user)
	// fmt.Println(query.ToSQL())
	result, err := query.Executor().ExecContext(ctx)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (u userRepo) GetUser(ctx context.Context, userId string) (domain.User, error) {

	var result domain.User
	query := u.database.
		Select().
		From(domain.TabNameUser).
		Where(
			goqu.C(domain.ColUserId).Eq(userId),
		)
	// fmt.Println(query.ToSQL())
	found, err := query.Executor().ScanStructContext(ctx, &result)
	if err != nil {
		return domain.User{}, status.Error(codes.Internal, "userId is not found")
	}
	if !found {
		return domain.User{}, status.Error(codes.Internal, "userId is not found")
	}
	return result, nil
}

func (u userRepo) UpdateUser(ctx context.Context, user domain.User) error {

	query := u.database.
		Update(domain.TabNameUser).
		Set(user).
		Where(
			goqu.C(domain.ColUserId).Eq(user.UserId),
		)

	// fmt.Println(query.ToSQL())
	result, err := query.Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	if id != int64(user.Id) {
		return errors.New("failed to update")
	}
	return nil
}

func (u userRepo) InactiveUser(ctx context.Context, userId string) error {

	user, err := u.GetUser(ctx, userId)
	if err != nil {
		return err
	}
	user.IsActive = false
	query := u.database.
		Update(domain.TabNameUser).
		Set(user).
		Where(
			goqu.C(domain.ColUserId).Eq(user.UserId),
		)
	// fmt.Println(query.ToSQL())
	result, err := query.Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	if id != int64(user.Id) {
		return errors.New("failed to update")
	}
	return nil
}
