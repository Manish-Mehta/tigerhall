package user

import (
	"log"
	"net/http"

	"github.com/Manish-Mehta/tigerhall/dto"
	"github.com/Manish-Mehta/tigerhall/model/datastore"
	"github.com/Manish-Mehta/tigerhall/model/entities"
	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Signup(request *dto.CreateUserRequest) *errorHandler.Error
}

type userService struct {
	dataStore datastore.UserStore
}

func (service *userService) Signup(request *dto.CreateUserRequest) *errorHandler.Error {
	// Only doing basic validations
	exists, err := service.dataStore.EmailExists(request.Email)
	log.Println(err)
	if err != nil {
		return &errorHandler.Error{
			Err:        "Not able to reach db",
			ErrMsg:     "Error while creating user in the system",
			StatusCode: http.StatusInternalServerError,
		}
	}
	if exists {
		return &errorHandler.Error{
			Err:        "User email exists",
			ErrMsg:     "User already exists",
			StatusCode: http.StatusInternalServerError,
		}
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 10)

	userEntity := &entities.User{
		UserName: request.UserName,
		Email:    request.Email,
		Password: string(hashedPassword),
	}
	err = service.dataStore.Create(userEntity)
	log.Println(err)
	if err != nil {
		return &errorHandler.Error{
			Err:        "db creation failed",
			ErrMsg:     "Error while creating user in the system",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

// func validatePassword(hashedPassword string, enteredPassword string) (bool, error) {
// 	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(enteredPassword))
// 	return err == nil, err
//   }

func NewUserService(ds datastore.UserStore) UserService {
	return &userService{
		dataStore: ds,
	}
}
