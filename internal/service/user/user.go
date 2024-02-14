package user

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Manish-Mehta/tigerhall/api/dto"
	"github.com/Manish-Mehta/tigerhall/internal/config"
	"github.com/Manish-Mehta/tigerhall/model/datastore"
	"github.com/Manish-Mehta/tigerhall/model/entities"
	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
)

type UserService interface {
	Signup(request *dto.SignupRequest) *errorHandler.Error
	Login(request *dto.LoginRequest) (string, *errorHandler.Error)
	Refresh(string, time.Time) (string, *errorHandler.Error)
}

type userService struct {
	dataStore datastore.UserStore
}

func (service *userService) Signup(request *dto.SignupRequest) *errorHandler.Error {
	// Only doing basic validations
	exists, err := service.dataStore.EmailExists(request.Email)
	log.Println(err)
	if err != nil {
		return &errorHandler.Error{
			Err:        "User check failed",
			ErrMsg:     "Error while creating user in the system",
			StatusCode: http.StatusInternalServerError,
		}
	}
	if exists {
		return &errorHandler.Error{
			Err:        "User email exists",
			ErrMsg:     "User already exists",
			StatusCode: http.StatusConflict,
		}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	userEntity := &entities.User{
		UserName: request.UserName,
		Email:    request.Email,
		Password: string(hashedPassword),
	}
	err = service.dataStore.Create(userEntity)
	if err != nil {
		return &errorHandler.Error{
			Err:        "User creation failed",
			ErrMsg:     "Error while creating user in the system",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (service *userService) Login(request *dto.LoginRequest) (string, *errorHandler.Error) {
	// Only doing basic validations

	userEntity := &entities.User{}
	err := service.dataStore.Get(userEntity, &entities.User{Email: request.Email}, []string{"id", "email", "password"})
	if err != nil {
		return "", &errorHandler.Error{
			Err:        "User fetch failed",
			ErrMsg:     "Error while getting user data",
			StatusCode: http.StatusInternalServerError,
		}
	}
	if userEntity.Email == "" {
		return "", &errorHandler.Error{
			Err:        "User not found",
			ErrMsg:     "User doesn't exists",
			StatusCode: http.StatusBadRequest,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(request.Password))
	if err != nil {
		return "", &errorHandler.Error{
			Err:        "Password verification failed",
			ErrMsg:     "Enter correct password",
			StatusCode: http.StatusUnauthorized,
		}
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userEntity.Email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 1)),
		Issuer:    strconv.FormatUint(uint64(userEntity.ID), 10),
	})
	token, err := claims.SignedString([]byte(config.TOKEN_SECRET))
	if err != nil {
		return "", &errorHandler.Error{
			Err:        "Token generation failed",
			ErrMsg:     "Error in login process",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return token, nil
}

func (service *userService) Refresh(email string, expiry time.Time) (string, *errorHandler.Error) {

	// New token will only be issued if the old token is within 1 hour of expiry.
	if time.Until(expiry) > time.Hour {
		return "", &errorHandler.Error{
			Err:        "High token expiry",
			ErrMsg:     "Token refresh only happens within 1 hour of expiry",
			StatusCode: http.StatusBadRequest,
		}
	}

	// Move the below code in a common lib function as to be used by login also.
	userEntity := &entities.User{}
	err := service.dataStore.Get(userEntity, &entities.User{Email: email}, []string{"id", "email"})
	if err != nil {
		return "", &errorHandler.Error{
			Err:        "User fetch failed",
			ErrMsg:     "Error while getting user data",
			StatusCode: http.StatusInternalServerError,
		}
	}
	if userEntity.Email == "" {
		return "", &errorHandler.Error{
			Err:        "User not found",
			ErrMsg:     "User doesn't exists",
			StatusCode: http.StatusBadRequest,
		}
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userEntity.Email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 1)),
		Issuer:    strconv.FormatUint(uint64(userEntity.ID), 10),
	})
	token, err := claims.SignedString([]byte(config.TOKEN_SECRET))
	if err != nil {
		return "", &errorHandler.Error{
			Err:        "Token generation failed",
			ErrMsg:     "Error in login process",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return token, nil
}

func NewUserService(ds datastore.UserStore) UserService {
	return &userService{
		dataStore: ds,
	}
}
