package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lhducc/bookmark-management/internal/service"
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
					"exp": 604800,
				}
				jsonBody, _ := json.Marshal(body)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("ShortenUrl",
					ctx,
					"https://example.com",
					604800).Return("123", nil)
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
					"exp": 604800,
				}
				jsonBody, _ := json.Marshal(body)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
			},
			setupMockSvc: func(ctx context.Context) *mocks.ShortenUrl {
				svcMock := mocks.NewShortenUrl(t)
				svcMock.On("ShortenUrl",
					ctx,
					"https://example.com",
					604800).Return("", assert.AnError)
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
					"exp": 604800,
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

func TestUrlShortenHandler_GetUrl(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockSvc func(t *testing.T, ctx context.Context) *mocks.ShortenUrl

		expectedResponseCode int
		expectedResponseBody string
		expectedLocation     string
	}{
		{
			name: "empty code -> 400",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/", nil)
				ctx.Params = gin.Params{{Key: "code", Value: ""}}
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenUrl {
				return mocks.NewShortenUrl(t)
			},

			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"wrong format"}`,
		},
		{
			name: "code not found -> 400",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/notfound", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "notfound"}}
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenUrl {
				mockSvc := mocks.NewShortenUrl(t)
				mockSvc.On("GetUrl", ctx, "notfound").
					Return("", service.ErrCodeNotFound).
					Once()
				return mockSvc
			},

			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"url not found"}`,
		},
		{
			name: "service returns other error -> 500",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/boom", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "boom"}}
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenUrl {
				mockSvc := mocks.NewShortenUrl(t)
				mockSvc.On("GetUrl", ctx, "boom").
					Return("", errors.New("some error")).
					Once()
				return mockSvc
			},

			expectedResponseCode: http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal server error"}`,
		},
		{
			name: "success -> 302 redirect",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/abc1234", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "abc1234"}}
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenUrl {
				mockSvc := mocks.NewShortenUrl(t)
				mockSvc.On("GetUrl", ctx, "abc1234").
					Return("https://google.com", nil).
					Once()
				return mockSvc
			},

			expectedResponseCode: http.StatusFound,
			expectedLocation:     "https://google.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)

			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc(t, gc)

			testHandler := NewUrlShortenHandler(mockSvc)
			testHandler.GetUrl(gc)

			assert.Equal(t, tc.expectedResponseCode, rec.Code)

			if tc.expectedResponseCode == http.StatusFound {
				assert.Equal(t, tc.expectedLocation, rec.Header().Get("Location"))
				return
			}

			assert.Equal(t, tc.expectedResponseBody, rec.Body.String())
		})
	}
}
