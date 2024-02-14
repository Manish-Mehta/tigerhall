package tiger

import (
	"net/http"
	"strconv"

	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"github.com/Manish-Mehta/tigerhall/pkg/interceptor"
	"github.com/gin-gonic/gin"

	"github.com/Manish-Mehta/tigerhall/api/dto"
	ts "github.com/Manish-Mehta/tigerhall/internal/service/tiger"
)

type TigerController interface {
	Create(c *gin.Context)
	List(c *gin.Context)
}
type tigerController struct {
	service ts.TigerService
}

const LIMIT = 5

// CreateTiger godoc
//
//	@Summary		Create Tiger
//	@Description	Creates a new tiger
//	@Description	NOTE: Access Token needed in Authorization header
//	@Tags			tiger
//	@Accept			json
//	@Produce		json
//	@Param			tiger			body		dto.TigerCreateRequest true	"Tiger Details"
//	@Success		201
//	@Failure		409	{object}	errorhandler.Error
//	@Failure		400	{object}	errorhandler.Error
//	@Failure		500	{object}	errorhandler.Error
//	@Router			/tiger [post]
//
// @Security ApiKeyAuth
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

// ListTigers godoc
//
//	@Summary		List All Tigers
//	@Description	Sorted by the last time the tigers were seen.
//	@Tags			tiger
//	@Accept			json
//	@Produce		json
//	@Param   		page         	query    	int        false  "Page number to be fetched"         		 minimum(1)
//	@Param   		limit         	query    	int        false  "Number of records to be fetched"          minimum(1)
//	@Success		200 {array}		dto.ListTigerResponse
//	@Failure		400	{object}	errorhandler.Error
//	@Failure		500	{object}	errorhandler.Error
//	@Router			/tiger [get]
func (tc tigerController) List(c *gin.Context) {
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

	tigers, err := tc.service.List(int(page), int(limit))
	if err != nil {
		interceptor.SendErrRes(c, err.Err, err.ErrMsg, err.StatusCode)
		return
	}
	interceptor.SendSuccessRes(c, tigers, http.StatusOK)
}

func NewTigerController(tigerService ts.TigerService) TigerController {
	return &tigerController{
		service: tigerService,
	}
}
