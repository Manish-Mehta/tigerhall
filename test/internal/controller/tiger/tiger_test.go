package tiger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Manish-Mehta/tigerhall/api/dto"
	uc "github.com/Manish-Mehta/tigerhall/internal/controller/tiger"
	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockTigerService struct {
}

func (m *MockTigerService) Create(request *dto.TigerCreateRequest) *errorHandler.Error {
	return nil
}

func (m *MockTigerService) List(int, int) (*[]*dto.ListTigerResponse, *errorHandler.Error) {
	return &[]*dto.ListTigerResponse{
		{
			ID:   1,
			Name: "Tiger 1",
		},
		{
			ID:   2,
			Name: "Tiger 2",
		},
	}, nil
}

func TestCreateController(t *testing.T) {

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
			name:         "SuccessCreate",
			reqBody:      `{"name": "testuser", "dob": "2020-07-17", "lastSeen": "2023-02-12T14:58:46Z", "coordinate": { "lat": 12.964597734871699, "lon": 77.63829741684891 }}`,
			expectedCode: 201,
		},
	}
	mockTigerService := &MockTigerService{}
	uController := uc.NewTigerController(mockTigerService)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/create", bytes.NewBufferString(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			uController.Create(c)

			assert.Equal(t, tc.expectedCode, c.Writer.Status())
		})
	}
}

func TestListController(t *testing.T) {

	tests := []struct {
		name         string
		reqBody      string
		expectedCode int
	}{
		{
			name:         "NonQueryParams",
			reqBody:      `{"user": "testuser"}`,
			expectedCode: 200,
		},
		{
			name:         "QueryParams",
			expectedCode: 200,
		},
	}
	mockTigerService := &MockTigerService{}
	uController := uc.NewTigerController(mockTigerService)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/?page=1&limit=2", nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			uController.List(c)

			assert.Equal(t, tc.expectedCode, c.Writer.Status())
		})
	}
}
