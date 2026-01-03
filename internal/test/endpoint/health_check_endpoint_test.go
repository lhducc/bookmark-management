package endpoint

import (
	"encoding/json"
	"github.com/lhducc/bookmark-management/internal/api"
	redisPkg "github.com/lhducc/bookmark-management/pkg/redis"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthCheckEndpoint tests the healthCheckEndpoint function.
// It tests that the function returns a HTTP 200 OK response with the correct JSON body.
// The JSON body is expected to have the following structure:
//
//	{
//	  "message": string,
//	  "serviceName": string,
//	  "instanceID": string
//	}
//
// The character set used for generating the JSON body is constant and does not change across different implementations of the interface. The length of the generated JSON body is constant and does not change across different implementations of the interface.
// If an error occurs while generating the JSON body, the error is returned immediately and the generated JSON body is an empty string.
func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"OK","serviceName":"bookmark-management","instanceID":"2025"}`,
		},
	}

	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := api.New(cfg, redisClient)
			rec := tc.setupTestHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp map[string]string

			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, "OK", resp["message"])
			assert.Equal(t, "bookmark-management", resp["serviceName"])
			assert.NotEmpty(t, resp["instanceID"])

		})
	}
}
