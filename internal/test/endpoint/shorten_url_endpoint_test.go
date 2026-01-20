package endpoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lhducc/bookmark-management/internal/api"
	redisPkg "github.com/lhducc/bookmark-management/pkg/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUrlShortenEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatus  int
		expectedCodeLen int
		expectedMessage string
	}{
		{
			name: "success",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"url": "https://google.com",
					"exp": 604800,
				}
				jsonBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:  http.StatusOK,
			expectedCodeLen: 7,
			expectedMessage: "Shorten URL generated successfully!",
		},
		{
			name: "wrong input - empty url",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"url": "",
					"exp": 10,
				}
				jsonBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:  http.StatusBadRequest,
			expectedCodeLen: 0,
			expectedMessage: "Invalid request",
		},
	}

	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := api.New(cfg, redisPkg.InitMockRedis(t))
			rec := tc.setupTestHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp map[string]string
			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedMessage, resp["message"])
			if tc.expectedCodeLen > 0 {
				assert.Equal(t, tc.expectedCodeLen, len(resp["code"]))
			}

		})
	}
}

func TestUrlRedirectEndpoint(t *testing.T) {
	t.Parallel()

	type shortenResp struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}

	testCases := []struct {
		name string

		setupTestHTTP func(app api.Engine) *httptest.ResponseRecorder

		expectedStatus   int
		expectedLocation string
		expectedMessage  string
	}{
		{
			name: "success redirect",

			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"url": "https://google.com",
					"exp": 604800,
				}
				jsonBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsonBody))
				req.Header.Set("Content-Type", "application/json")

				shortenRec := httptest.NewRecorder()
				app.ServeHTTP(shortenRec, req)
				require.Equal(t, http.StatusOK, shortenRec.Code)

				var sr shortenResp
				err := json.Unmarshal(shortenRec.Body.Bytes(), &sr)
				require.NoError(t, err)
				require.NotEmpty(t, sr.Code)

				redirectReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/links/redirect/%s", sr.Code), nil)
				redirectRec := httptest.NewRecorder()
				app.ServeHTTP(redirectRec, redirectReq)
				return redirectRec
			},

			expectedStatus:   http.StatusFound,
			expectedLocation: "https://google.com",
		},
		{
			name: "not found",

			setupTestHTTP: func(app api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/notfound", nil)
				rec := httptest.NewRecorder()
				app.ServeHTTP(rec, req)
				return rec
			},

			expectedStatus:  http.StatusNotFound,
			expectedMessage: "url not found",
		},
	}

	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := api.New(cfg, redisPkg.InitMockRedis(t))
			rec := tc.setupTestHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusFound {
				assert.Equal(t, tc.expectedLocation, rec.Header().Get("Location"))
				return
			}

			// JSON response cases
			if tc.expectedMessage != "" {
				var resp map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedMessage, resp["message"])
			}
		})
	}

}
