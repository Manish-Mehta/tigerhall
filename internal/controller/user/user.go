package user

import (
	"net/http"

	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"github.com/Manish-Mehta/tigerhall/pkg/interceptor"
	"github.com/gin-gonic/gin"

	"github.com/Manish-Mehta/tigerhall/dto"
	us "github.com/Manish-Mehta/tigerhall/internal/service/user"
)

type UserController interface {
	Signup(c *gin.Context)
}
type userController struct {
	service us.UserService
}

func (uc userController) Signup(c *gin.Context) {
	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while creating user")
	request := &dto.CreateUserRequest{}
	if err := c.ShouldBind(&request); err != nil {
		interceptor.SendErrRes(c, "Invalid request body", "Check your request body data with proper validations", http.StatusBadRequest)
		return
	}
	err := uc.service.Signup(request)
	if err != nil {
		if err != nil {
			interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
			return
		}
	}
	interceptor.SendSuccessRes(c, map[string]string{"message": "user created"}, http.StatusCreated)
}

func NewUserController(userService us.UserService) UserController {
	return &userController{
		service: userService,
	}
}
