package service

import (
	"context"
	"errors"
	"html"
	"strings"

	"github.com/namnv2496/go-coffee-shop-demo/internal/authorization/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/authorization/repo"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/security"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	Login(ctx context.Context, user domain.User) ([]string, error)
	GetUser(ctx context.Context, userId string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
	InactiveUser(ctx context.Context, userId string) error
}

type userService struct {
	userRepo repo.UserRepo
}

func NewUserService(
	userRepo repo.UserRepo,
) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s userService) CreateUser(ctx context.Context, user domain.User) (int, error) {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(passwordHash)
	user.UserId = html.EscapeString(strings.TrimSpace(user.UserId))
	id, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (s userService) Login(ctx context.Context, user domain.User) ([]string, error) {

	userDB, err := s.userRepo.GetUser(ctx, user.UserId)
	if err != nil {
		return []string{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if err == nil && userDB.IsActive {
		return security.GenerateJWTToken(userDB.UserId, strings.Split(userDB.Role, ","))
	}
	return []string{}, errors.New("wrong password")
}

func (s userService) GetUser(ctx context.Context, userId string) (domain.User, error) {
	user, err := s.userRepo.GetUser(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = ""
	return user, nil
}

func (s userService) UpdateUser(ctx context.Context, user domain.User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.UserId = html.EscapeString(strings.TrimSpace(user.UserId))
	return s.userRepo.UpdateUser(ctx, user)
}

func (s userService) InactiveUser(ctx context.Context, userId string) error {
	return s.userRepo.InactiveUser(ctx, userId)
}
