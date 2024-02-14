package user

import (
	"net/http"

	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"github.com/Manish-Mehta/tigerhall/pkg/interceptor"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Manish-Mehta/tigerhall/api/dto"
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

// UserSignup godoc
//
//	@Summary		User Signup
//	@Description	creates a new user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user			body		dto.SignupRequest true	"User Details"
//	@Success		201
//	@Failure		409	{object}	errorhandler.Error
//	@Failure		400	{object}	errorhandler.Error
//	@Failure		500	{object}	errorhandler.Error
//	@Router			/user [post]
func (uc userController) Signup(c *gin.Context) {
	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while creating user")

	request := &dto.SignupRequest{}
	if err := c.ShouldBind(request); err != nil {
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

// UserLogin godoc
//
//	@Summary		User Login
//	@Description	log the user in
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user			body		dto.LoginRequest true	"User Creds"
//	@Success		200 {object}	dto.LoginResponse
//	@Failure		400	{object}	errorhandler.Error
//	@Failure		500	{object}	errorhandler.Error
//	@Router			/user/login [post]
func (uc userController) Login(c *gin.Context) {
	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while logging in")

	request := &dto.LoginRequest{}
	if err := c.ShouldBind(request); err != nil {
		interceptor.SendErrRes(c, "Invalid request body", "Check your request body data with proper validations", http.StatusBadRequest)
		return
	}
	token, err := uc.service.Login(request)
	if err != nil {
		interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
		return
	}
	interceptor.SendSuccessRes(c, dto.LoginResponse{Token: token}, http.StatusOK)
}

// UserTokenRefresh godoc
//
//	@Summary		User Token Refresh
//	@Description	Refreshes the user access token
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200 {object}	dto.LoginResponse
//	@Failure		400	{object}	errorhandler.Error
//	@Failure		500	{object}	errorhandler.Error
//	@Router			/user/refresh [get]
//
// @Security ApiKeyAuth
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
	interceptor.SendSuccessRes(c, dto.LoginResponse{Token: token}, http.StatusOK)
}

func NewUserController(userService us.UserService) UserController {
	return &userController{
		service: userService,
	}
}
