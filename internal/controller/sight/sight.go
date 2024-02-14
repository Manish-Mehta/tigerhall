package tiger

import (
	"log"
	"net/http"
	"strconv"

	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"github.com/Manish-Mehta/tigerhall/pkg/interceptor"
	"github.com/gin-gonic/gin"

	"github.com/Manish-Mehta/tigerhall/api/dto"
	ss "github.com/Manish-Mehta/tigerhall/internal/service/sight"
)

type SightController interface {
	Create(c *gin.Context)
	List(c *gin.Context)
}
type sightController struct {
	service ss.SightService
}

const LIMIT = 5

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

func (sc sightController) List(c *gin.Context) {
	defer errorHandler.RecoverAndSendErrRes(c, "Something went wrong while creating tiger")

	var page int64 = 1
	var limit int64 = LIMIT
	o := c.DefaultQuery("page", "1")
	if v, err := strconv.ParseInt(o, 10, 32); err == nil {
		if v < 0 {
			page = 0
		} else {
			page = v
		}
	}

	l := c.DefaultQuery("limit", strconv.Itoa(LIMIT))
	if v, err := strconv.ParseInt(l, 10, 32); err == nil {
		if v <= 0 || v > LIMIT {
			limit = LIMIT
		} else {
			limit = v
		}
	}

	sightings, err := sc.service.List(int(page), int(limit))
	if err != nil {
		interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
		return
	}
	interceptor.SendSuccessRes(c, sightings, http.StatusCreated)
}

func NewSightController(sightService ss.SightService) SightController {
	return &sightController{
		service: sightService,
	}
}
