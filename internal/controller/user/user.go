package user

import (
	"net/http"

	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"github.com/Manish-Mehta/tigerhall/pkg/interceptor"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Manish-Mehta/tigerhall/dto"
	us "github.com/Manish-Mehta/tigerhall/internal/service/user"
)

type UserController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}
type userController struct {
	service us.UserService
}

func (uc userController) Signup(c *gin.Context) {
	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while creating user")

	request := &dto.SignupRequest{}
	if err := c.ShouldBind(&request); err != nil {
		interceptor.SendErrRes(c, "Invalid request body", "Check your request body data with proper validations", http.StatusBadRequest)
		return
	}
	err := uc.service.Signup(request)
	if err != nil {
		interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
		return
	}
	interceptor.SendSuccessRes(c, map[string]string{"message": "user created"}, http.StatusCreated)
}

func (uc userController) Login(c *gin.Context) {
	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while logging in")

	request := &dto.LoginRequest{}
	if err := c.ShouldBind(&request); err != nil {
		interceptor.SendErrRes(c, "Invalid request body", "Check your request body data with proper validations", http.StatusBadRequest)
		return
	}
	token, err := uc.service.Login(request)
	if err != nil {
		interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
		return
	}
	interceptor.SendSuccessRes(c, map[string]string{"access_token": token}, http.StatusOK)
}

func (uc userController) Refresh(c *gin.Context) {
	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while refreshing token")

	email, exists1 := c.Get("Email")
	expiry, exists2 := c.Get("TokenExpiry")
	if !exists1 || !exists2 {
		interceptor.SendErrRes(c, "Access token malformed", "Invalid access token", http.StatusBadRequest)
		return
	}

	token, err := uc.service.Refresh(email.(string), expiry.(*jwt.NumericDate).Time)
	if err != nil {
		interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
		return
	}
	interceptor.SendSuccessRes(c, map[string]interface{}{"access_token": token}, http.StatusOK)
}

func NewUserController(userService us.UserService) UserController {
	return &userController{
		service: userService,
	}
}
