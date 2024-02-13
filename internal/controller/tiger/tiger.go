package tiger

import (
	"net/http"

	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"github.com/Manish-Mehta/tigerhall/pkg/interceptor"
	"github.com/gin-gonic/gin"

	"github.com/Manish-Mehta/tigerhall/dto"
	ts "github.com/Manish-Mehta/tigerhall/internal/service/tiger"
)

type TigerController interface {
	Create(c *gin.Context)
}
type tigerController struct {
	service ts.TigerService
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

func NewTigerController(tigerService ts.TigerService) TigerController {
	return &tigerController{
		service: tigerService,
	}
}
