package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lhducc/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUrlStorageHandler_ShortenUrl(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(ctx context.Context) *mocks.ShortenUrl

		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"url": "https://example.com",
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("ShortenUrl",
					ctx,
					"https://example.com",
					3600).Return("123", nil)
				return svcMock
			},

			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"message": "Shorten URL generated successfully!",
				"code":    "123",
			},
		},
		{
			name: "serivce error",

			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"url": "https://example.com",
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("ShortenUrl",
					ctx,
					"https://example.com",
					3600).Return("", assert.AnError)
				return svcMock
			},

			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]any{
				"message": "internal server error",
			},
		},
		{
			name: "wrong input",

			setupRequest: func(ctx *gin.Context) {
				body := map[string]any{
					"url": "",
					"exp": 3600,
				}
				jsonBody, _ := json.Marshal(body)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				return mocks.NewShortenUrl(t)
			},

			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]any{
				"message": "Invalid request",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc(gc)
			testHandler := NewUrlShortenHandler(mockSvc)

			testHandler.ShortenUrl(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != nil {
				var actualBody map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &actualBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, actualBody)
			}
		})
	}
}
