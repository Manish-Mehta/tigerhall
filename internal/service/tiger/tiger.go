package tiger

import (
	"log"
	"net/http"

	"gorm.io/datatypes"

	"github.com/Manish-Mehta/tigerhall/api/dto"
	"github.com/Manish-Mehta/tigerhall/internal/util"
	"github.com/Manish-Mehta/tigerhall/model/datastore"
	"github.com/Manish-Mehta/tigerhall/model/entities"
	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
)

type TigerService interface {
	Create(request *dto.TigerCreateRequest) *errorHandler.Error
	List(int, int) (*[]*entities.Tiger, *errorHandler.Error)
}

type tigerService struct {
	dataStore datastore.TigerStore
}

func (service *tigerService) Create(request *dto.TigerCreateRequest) *errorHandler.Error {
	validator := util.NewValidator()

	dob, err := validator.ValDateFormat(request.DOB)
	if err != nil {
		return &errorHandler.Error{
			Err:        "Invalid DOB",
			ErrMsg:     "Provide valid DOB date in yyyy-mm-dd format",
			StatusCode: http.StatusBadRequest,
		}
	}

	err = validator.ValCoord(request.Coordinate)
	if err != nil {
		return &errorHandler.Error{
			Err:        "Invalid coordinates",
			ErrMsg:     "Provide valid coordinates",
			StatusCode: http.StatusInternalServerError,
		}
	}

	exists, err := service.dataStore.NameExists(request.Name)
	if err != nil {
		return &errorHandler.Error{
			Err:        "Tiger check failed",
			ErrMsg:     "Error while creating tiger in the system",
			StatusCode: http.StatusInternalServerError,
		}
	}
	if exists {
		return &errorHandler.Error{
			Err:        "Tiger name exists",
			ErrMsg:     "Tiger already exists",
			StatusCode: http.StatusConflict,
		}
	}

	tigerEntity := &entities.Tiger{
		Name:     request.Name,
		Dob:      datatypes.Date(dob),
		LastSeen: request.LastSeen,
		Lat:      request.Coordinate.Lat,
		Lon:      request.Coordinate.Lon,
	}
	err = service.dataStore.Create(tigerEntity)
	if err != nil {
		return &errorHandler.Error{
			Err:        "Tiger creation failed",
			ErrMsg:     "Error while creating user in the system",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (service *tigerService) List(page int, limit int) (*[]*entities.Tiger, *errorHandler.Error) {

	var tigers []*entities.Tiger
	err := service.dataStore.List(&tigers, page, limit)
	if err != nil {
		log.Println(err)
		return nil, &errorHandler.Error{
			Err:        "Tiger fetch failed",
			ErrMsg:     "Error while getting tiger",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return &tigers, nil
}

func NewTigerService(ds datastore.TigerStore) TigerService {
	return &tigerService{
		dataStore: ds,
	}
}
