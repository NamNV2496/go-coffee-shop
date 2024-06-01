package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WrapperResponse(req *gin.Context, code int, messsage any) {

	req.JSON(http.StatusOK, gin.H{
		"code":    code,
		"status":  http.StatusText(code),
		"message": messsage,
	})
}
