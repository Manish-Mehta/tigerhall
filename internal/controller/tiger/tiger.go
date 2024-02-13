package tiger

import (
	"log"
	"net/http"

	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"github.com/Manish-Mehta/tigerhall/pkg/interceptor"
	"github.com/gin-gonic/gin"

	"github.com/Manish-Mehta/tigerhall/dto"
	tg "github.com/Manish-Mehta/tigerhall/internal/service/tiger"
)

type TigerController interface {
	Create(c *gin.Context)
	CreateSighting(c *gin.Context)
	// Refresh(c *gin.Context)
}
type tigerController struct {
	service tg.TigerService
}

func (tc tigerController) Create(c *gin.Context) {
	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while creating tiger")

	request := &dto.TigerCreateRequest{}
	if err := c.ShouldBind(request); err != nil {
		interceptor.SendErrRes(c, "Invalid request body", "Check your request body data with proper validations", http.StatusBadRequest)
		return
	}

	err := tc.service.Create(request)
	if err != nil {
		interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
		return
	}
	interceptor.SendSuccessRes(c, map[string]string{"message": "Tiger created"}, http.StatusCreated)
}

func (tc tigerController) CreateSighting(c *gin.Context) {
	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while creating sighting")

	userId, exists := c.Get("Id")
	if !exists {
		interceptor.SendErrRes(c, "Access token malformed", "Invalid access token", http.StatusBadRequest)
		return
	}

	request := &dto.TigerCreateSightingRequest{}
	request.UserID = uint(userId.(uint64))

	if err := c.ShouldBind(request); err != nil {
		log.Println(err)
		interceptor.SendErrRes(c, "Invalid request body", "Check your request body data with proper validations", http.StatusBadRequest)
		return
	}

	err := tc.service.CreateSighting(request)
	if err != nil {
		interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
		return
	}
	interceptor.SendSuccessRes(c, map[string]string{"message": "Tiger created"}, http.StatusCreated)
}

// func (tc tigerController) Login(c *gin.Context) {
// 	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while logging in")

// 	request := &dto.LoginRequest{}
// 	if err := c.ShouldBind(request); err != nil {
// 		interceptor.SendErrRes(c, "Invalid request body", "Check your request body data with proper validations", http.StatusBadRequest)
// 		return
// 	}
// 	token, err := tc.service.Login(request)
// 	if err != nil {
// 		interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
// 		return
// 	}
// 	interceptor.SendSuccessRes(c, map[string]string{"access_token": token}, http.StatusOK)
// }

// func (tc tigerController) Refresh(c *gin.Context) {
// 	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while refreshing token")

// 	email, exists1 := c.Get("Email")
// 	expiry, exists2 := c.Get("TokenExpiry")
// 	if !exists1 || !exists2 {
// 		interceptor.SendErrRes(c, "Access token malformed", "Invalid access token", http.StatusBadRequest)
// 		return
// 	}

// 	token, err := tc.service.Refresh(email.(string), expiry.(*jwt.NumericDate).Time)
// 	if err != nil {
// 		interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
// 		return
// 	}
// 	interceptor.SendSuccessRes(c, map[string]interface{}{"access_token": token}, http.StatusOK)
// }

func NewTigerController(tigerService tg.TigerService) TigerController {
	return &tigerController{
		service: tigerService,
	}
}
