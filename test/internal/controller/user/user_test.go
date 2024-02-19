package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Manish-Mehta/tigerhall/api/dto"
	uc "github.com/Manish-Mehta/tigerhall/internal/controller/user"
	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type MockUserService struct {
}

func (m *MockUserService) Signup(request *dto.SignupRequest) *errorHandler.Error {
	return nil
}

func (m *MockUserService) Login(request *dto.LoginRequest) (string, *errorHandler.Error) {
	return "token", nil
}

func (m *MockUserService) Refresh(string, time.Time) (string, *errorHandler.Error) {
	return "token", nil
}

func TestSignupController(t *testing.T) {

	tests := []struct {
		name         string
		reqBody      string
		expectedCode int
	}{
		{
			name:         "BadRequest",
			reqBody:      `{"user": "testuser"}`,
			expectedCode: 400,
		},
		{
			name:         "SuccessSignup",
			reqBody:      `{"userName": "testuser", "password": "password123", "email": "test@example.com"}`,
			expectedCode: 201,
		},
	}
	mockUserService := &MockUserService{}
	uController := uc.NewUserController(mockUserService)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/signup", bytes.NewBufferString(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			uController.Signup(c)

			assert.Equal(t, tc.expectedCode, c.Writer.Status())
		})
	}
}

func TestLoginController(t *testing.T) {

	tests := []struct {
		name         string
		reqBody      string
		expectedCode int
	}{
		{
			name:         "BadRequest",
			reqBody:      `{"user": "testuser"}`,
			expectedCode: 400,
		},
		{
			name:         "SuccessLogin",
			reqBody:      `{"password": "password123", "email": "test@example.com"}`,
			expectedCode: 200,
		},
	}
	mockUserService := &MockUserService{}
	uController := uc.NewUserController(mockUserService)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			uController.Login(c)

			assert.Equal(t, tc.expectedCode, c.Writer.Status())
		})
	}
}

func TestRefreshController(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		exp          *jwt.NumericDate
		expectedCode int
	}{
		{
			name:         "BadRequest",
			exp:          &jwt.NumericDate{time.Now()},
			expectedCode: 400,
		},
		{
			name:         "SuccessRefresh",
			email:        "test@example.com",
			exp:          &jwt.NumericDate{time.Now()},
			expectedCode: 200,
		},
	}
	mockUserService := &MockUserService{}
	uController := uc.NewUserController(mockUserService)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/refresh", nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tc.email != "" {
				c.Set("Email", tc.email)
			}
			c.Set("TokenExpiry", tc.exp)
			c.Request = req

			uController.Refresh(c)

			assert.Equal(t, tc.expectedCode, c.Writer.Status())
		})
	}
}
