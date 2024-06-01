package main

import (
	"context"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	author "github.com/namnv2496/go-coffee-shop-demo/internal/authorization"
	"github.com/namnv2496/go-coffee-shop-demo/internal/authorization/app"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/security"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
)

func main() {

	app, cleanup, err := author.Initialize("")
	if err != nil {
		panic("Error start app")
	}
	defer cleanup()
	r := SetupGin()
	app.Start()

	go func() {

		routing(app, r, context.Background())
		r.Run(":8083")
	}()
	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
}

func SetupGin() *gin.Engine {

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Respond with 200 OK status
		c.Status(http.StatusOK)
	})
	return r
}

func routing(app *app.App, r *gin.Engine, ctx context.Context) {

	// create account
	r.POST("/api/v1/register", func(req *gin.Context) { app.CreateUser(ctx, req) })
	r.POST("/api/v1/login", func(req *gin.Context) { app.Login(ctx, req) })
	//get uset by userId
	r.POST("/api/v1/renewToken", func(req *gin.Context) { app.RenewToken(ctx, req) })

	adminMemberRoutes := r.Group("/api/v1/ping")
	adminMemberRoutes.Use(security.JWTAuthWithRole([]string{"admin", "counter", "kitchen"}))
	adminMemberRoutes.GET("", func(req *gin.Context) { app.Ping(ctx, req) })

	adminMemberRoutes = r.Group("/api/v1/getUser")
	adminMemberRoutes.Use(security.JWTAuthWithRole([]string{"admin"}))

	/* Check role:
	if role is "admin" then return information of all members
	if role is member then only return his informatin
	*/
	adminMemberRoutes.GET("", func(req *gin.Context) {
		user, err := app.GetUser(ctx, req)
		if err != nil {
			utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
			return
		}
		roles, err := security.GetRole(req)
		if err != nil {
			utils.WrapperResponse(req, http.StatusForbidden, err.Error())
			return
		}
		for _, role := range roles {
			if role == "admin" {
				utils.WrapperResponse(req, http.StatusOK, user)
				return
			} else if role == "member" {
				userId, err := security.GetUserId(req)
				if err != nil {
					utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
					return
				}
				if userId == user.UserId {
					utils.WrapperResponse(req, http.StatusOK, user)

					return
				}
			}
		}
		utils.WrapperResponse(req, http.StatusNotFound, nil)
	})

	// update user
	adminRoutes := r.Group("/api/v1")
	adminRoutes.Use(security.JWTAuthWithRole([]string{"admin"}))
	adminRoutes.PUT("/updateUser", func(req *gin.Context) { app.UpdateUser(ctx, req) })
	// active/inactive user
	adminRoutes.Use(security.JWTAuthWithRole([]string{"admin"}))
	adminRoutes.PUT("/activeUser", func(req *gin.Context) { app.InactiveUser(ctx, req) })
}
