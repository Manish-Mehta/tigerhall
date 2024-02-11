package user

import (
	"net/http"

	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"github.com/Manish-Mehta/tigerhall/pkg/interceptor"
	"github.com/gin-gonic/gin"
)

var CreateUser = func() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while creating user")

		interceptor.SendSuccessRes(c, map[string]string{"message": "user created"}, http.StatusCreated)
	}
}
