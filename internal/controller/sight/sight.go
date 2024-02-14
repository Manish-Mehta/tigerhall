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

// CreateSighting godoc
//
//	@Summary		Create a new sighting of a tiger
//	@Description	Records last sighting of a tiger
//	@Description	NOTE: Access Token needed in Authorization header
//	@Tags			sighting
//	@Accept			json
//	@Produce		json
//	@Param			sight			body		dto.CreateSightingRequest true	"Sight Details"
//	@Success		201
//	@Failure		409	{object}	errorhandler.Error
//	@Failure		400	{object}	errorhandler.Error
//	@Failure		500	{object}	errorhandler.Error
//	@Router			/sight [post]
//
// @Security ApiKeyAuth
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

// ListTigers godoc
//
//	@Summary		List All Sighting of Tigers
//	@Description	Sorted by the last time the tigers were seen.
//	@Tags			sighting
//	@Accept			json
//	@Produce		json
//	@Param   		page         	query    	int        false  "Page number to be fetched"         		 minimum(1)
//	@Param   		limit         	query    	int        false  "Number of records to be fetched"          minimum(1)
//	@Success		200 {array}		dto.ListSightingResponse
//	@Failure		400	{object}	errorhandler.Error
//	@Failure		500	{object}	errorhandler.Error
//	@Router			/sight [get]
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
