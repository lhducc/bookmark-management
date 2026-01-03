package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lhducc/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testConnectError = errors.New("can't connect to redis")

func TestHealthCheckHandler_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(t *testing.T) *mocks.HealthCheck

		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name: "normal case",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health_check", nil)
			},
			setupMockSvc: func(t *testing.T) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check", mock.Anything).Return("OK", "bookmark-management", "2025", nil)
				return mockSvc
			},

			expectedResponseCode: http.StatusOK,
			expectedResponseBody: `{"message":"OK","serviceName":"bookmark-management","instanceID":"2025"}`,
		},
		{
			name: "error case",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health_check", nil)
			},
			setupMockSvc: func(t *testing.T) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check", mock.Anything).Return("NOT OK", "bookmark-management", "2025", testConnectError)
				return mockSvc
			},

			expectedResponseCode: http.StatusServiceUnavailable,
			expectedResponseBody: `{"error":"Internal Server Error","message":"NOT OK","service_name":"bookmark-management","instance_id":"2025"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc(t)
			testHandler := NewHealthCheckHandler(mockSvc)
			testHandler.Check(gc)

			assert.Equal(t, tc.expectedResponseCode, rec.Code)
			assert.Equal(t, tc.expectedResponseBody, rec.Body.String())
		})
	}
}
