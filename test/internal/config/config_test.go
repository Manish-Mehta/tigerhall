package config

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Manish-Mehta/tigerhall/internal/config"
)

func TestSetConfig(t *testing.T) {

	var Getenv = func(key string) string {
		val := os.Getenv(key)
		return strings.Trim(val, " ")
	}

	tests := []struct {
		name           string
		envVars        map[string]string
		expectedRegion string
		expectedKey    string
		expectedId     string
		failureCase    bool
	}{
		{
			name:           "ValidEnvironmentVariables",
			envVars:        map[string]string{"AWS_REGION": "test_region", "AWS_SECRET_ACCESS_KEY": "test_key", "AWS_ACCESS_KEY_ID": "test_id"},
			expectedRegion: "test_region",
			expectedKey:    "test_key",
			expectedId:     "test_id",
			failureCase:    false,
		},
		{
			name:           "MissingInvalidEnvironmentVariable",
			envVars:        map[string]string{"AWS_SECRET_ACCESS_KEY": "test_key", "AWS_ACCESS_KEY_ID": "test_id"},
			expectedRegion: "",
			expectedKey:    "test",
			expectedId:     "test",
			failureCase:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envVars {
				os.Setenv(k, v)
			}

			config.SetConfig(Getenv)
			// Assert results
			if tc.failureCase {
				assert.NotEqual(t, tc.expectedRegion, os.Getenv("AWS_REGION"))
				assert.NotEqual(t, tc.expectedKey, os.Getenv("AWS_SECRET_ACCESS_KEY"))
				assert.NotEqual(t, tc.expectedId, os.Getenv("AWS_ACCESS_KEY_ID"))
			} else {
				assert.Equal(t, tc.expectedRegion, os.Getenv("AWS_REGION"))
				assert.Equal(t, tc.expectedKey, os.Getenv("AWS_SECRET_ACCESS_KEY"))
				assert.Equal(t, tc.expectedId, os.Getenv("AWS_ACCESS_KEY_ID"))
			}
		})
	}
}

func TestInitServer(t *testing.T) {
	server := &config.Server{}
	server.InitServer()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	server.GinInstance.ServeHTTP(w, req)

	t.Run("ServerInit", func(t *testing.T) {
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
