package tiger

import (
	"net/http"
	"time"

	"github.com/Manish-Mehta/tigerhall/dto"
	"github.com/Manish-Mehta/tigerhall/model/datastore"
	"github.com/Manish-Mehta/tigerhall/model/entities"
	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"gorm.io/datatypes"
)

type TigerService interface {
	Create(request *dto.TigerCreateRequest) *errorHandler.Error
	// Login(request *dto.LoginRequest) (string, *errorHandler.Error)
	// Refresh(string, time.Time) (string, *errorHandler.Error)
}

type tigerService struct {
	dataStore datastore.TigerStore
}

func (service *tigerService) Create(request *dto.TigerCreateRequest) *errorHandler.Error {
	format := "2006-01-02"
	dob, err := time.Parse(format, request.DOB)
	if err != nil {
		return &errorHandler.Error{
			Err:        "Invalid DOB",
			ErrMsg:     "Provide valid DOB date in yyyy-mm-dd format",
			StatusCode: http.StatusBadRequest,
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
		Name: request.Name,
		Dob:  datatypes.Date(dob),
		// LastSeen: request.LastSeen,
		// Lat:      request.Coordinate.Lat,
		// Lon:      request.Coordinate.Lon,
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

// func (service *tigerService) Login(request *dto.LoginRequest) (string, *errorHandler.Error) {
// 	// Only doing basic validations

// 	tigerEntity := &entities.Tiger{}
// 	err := service.dataStore.Get(tigerEntity, &entities.Tiger{Email: request.Email}, []string{"email", "password"})
// 	if err != nil {
// 		return "", &errorHandler.Error{
// 			Err:        "Tiger fetch failed",
// 			ErrMsg:     "Error while getting user data",
// 			StatusCode: http.StatusInternalServerError,
// 		}
// 	}
// 	if tigerEntity.Email == "" {
// 		return "", &errorHandler.Error{
// 			Err:        "Tiger not found",
// 			ErrMsg:     "Tiger doesn't exists",
// 			StatusCode: http.StatusBadRequest,
// 		}
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(tigerEntity.Password), []byte(request.Password))
// 	if err != nil {
// 		return "", &errorHandler.Error{
// 			Err:        "Password verification failed",
// 			ErrMsg:     "Enter correct password",
// 			StatusCode: http.StatusUnauthorized,
// 		}
// 	}

// 	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
// 		Subject:   tigerEntity.Email,
// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 1)),
// 	})
// 	token, err := claims.SignedString([]byte(config.TOKEN_SECRET))
// 	if err != nil {
// 		return "", &errorHandler.Error{
// 			Err:        "Token generation failed",
// 			ErrMsg:     "Error in login process",
// 			StatusCode: http.StatusInternalServerError,
// 		}
// 	}

// 	return token, nil
// }

func NewTigerService(ds datastore.TigerStore) TigerService {
	return &tigerService{
		dataStore: ds,
	}
}
