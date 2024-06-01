package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
)

func JWTAuthWithRole(role []string) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := CheckRole(context, role)
		if err != nil {
			if err.Error() == "token has invalid claims: token is expired" || err.Error() == "need to renew token" {
				utils.WrapperResponse(context, http.StatusUnauthorized, "")
				context.Abort()
				return
			}
			context.Abort()
			return
		}
		context.Next()
	}
}
