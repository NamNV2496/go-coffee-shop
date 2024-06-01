package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/authorization/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/authorization/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/security"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
)

type AppInterface interface {
	Start() error
	CreateUser(ctx context.Context, req *gin.Context)
	Login(ctx context.Context, req *gin.Context)
	RenewToken(ctx context.Context, req *gin.Context)
	Ping(ctx context.Context, req *gin.Context)
	GetUser(ctx context.Context, req *gin.Context) (domain.User, error)
	UpdateUser(ctx context.Context, req *gin.Context)
	InactiveUser(ctx context.Context, req *gin.Context)
}

type App struct {
	userService service.UserService
}

func NewApp(
	userService service.UserService,
) *App {
	return &App{
		userService: userService,
	}
}

func (app App) Start() error {
	security.InitJWT("./pkg/security/.env")
	return nil
}

func (app App) CreateUser(ctx context.Context, req *gin.Context) {
	var user domain.User
	if err := req.BindJSON(&user); err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	id, err := app.userService.CreateUser(ctx, user)
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	utils.WrapperResponse(req, http.StatusCreated, id)
}

func (app App) Login(ctx context.Context, req *gin.Context) {
	var user domain.User
	if err := req.BindJSON(&user); err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	if token, err := app.userService.Login(ctx, user); err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	} else {
		utils.WrapperResponse(req, http.StatusOK, gin.H{
			"token":        token[0],
			"refreshToken": token[1],
		})
	}
}

func (app App) RenewToken(ctx context.Context, req *gin.Context) {
	newToken, err := security.RenewToken(req)
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	utils.WrapperResponse(req, http.StatusOK, gin.H{
		"token":        newToken[0],
		"refreshToken": newToken[1],
	})
}

func (app App) Ping(ctx context.Context, req *gin.Context) {

}

func (app App) GetUser(ctx context.Context, req *gin.Context) (domain.User, error) {
	userId := req.Query("userId")
	return app.userService.GetUser(ctx, userId)
}

func (app App) UpdateUser(ctx context.Context, req *gin.Context) {
	var user domain.User
	if err := req.BindJSON(&user); err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	if err := app.userService.UpdateUser(ctx, user); err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
	}
	utils.WrapperResponse(req, http.StatusOK, "")
}

func (app App) InactiveUser(ctx context.Context, req *gin.Context) {
	userId := req.Query("userId")
	err := app.userService.InactiveUser(ctx, userId)
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	utils.WrapperResponse(req, http.StatusOK, "")

}
