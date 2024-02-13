package tiger

import (
	"log"
	"net/http"

	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"github.com/Manish-Mehta/tigerhall/pkg/interceptor"
	"github.com/gin-gonic/gin"

	"github.com/Manish-Mehta/tigerhall/dto"
	ss "github.com/Manish-Mehta/tigerhall/internal/service/sight"
)

type SightController interface {
	Create(c *gin.Context)
}
type sightController struct {
	service ss.SightService
}

func (sc sightController) Create(c *gin.Context) {
	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while creating sight")

	userId, exists := c.Get("Id")
	if !exists {
		interceptor.SendErrRes(c, "Access token malformed", "Invalid access token", http.StatusBadRequest)
		return
	}

	request := &dto.CreateSightingRequest{}
	request.UserID = uint(userId.(uint64))

	if err := c.ShouldBind(request); err != nil {
		log.Println(err)
		interceptor.SendErrRes(c, "Invalid request body", "Check your request body data with proper validations", http.StatusBadRequest)
		return
	}

	err := sc.service.Create(request)
	if err != nil {
		interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
		return
	}
	interceptor.SendSuccessRes(c, map[string]string{"message": "Tiger sighting created"}, http.StatusCreated)
}

func NewSightController(sightService ss.SightService) SightController {
	return &sightController{
		service: sightService,
	}
}
