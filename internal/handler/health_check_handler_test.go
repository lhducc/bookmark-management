package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lhducc/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
				mockSvc.On("Check").Return("OK", "bookmark-management", "2025")
				return mockSvc
			},

			expectedResponseCode: http.StatusOK,
			expectedResponseBody: `{"message":"OK","serviceName":"bookmark-management","instanceID":"2025"}`,
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
