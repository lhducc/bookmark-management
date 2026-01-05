package endpoint

import (
	"bytes"
	"encoding/json"
	"github.com/lhducc/bookmark-management/internal/api"
	redisPkg "github.com/lhducc/bookmark-management/pkg/redis"
	"github.com/stretchr/testify/assert"
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
	}{
		{
			name: "success",

			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"url": "https://google.com",
					"exp": 10,
				}
				jsonBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/shorten-url", bytes.NewReader(jsonBody))
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:  http.StatusOK,
			expectedCodeLen: 7,
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

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, "Shorten URL generated successfully!", resp["message"])
			assert.Equal(t, tc.expectedCodeLen, len(resp["code"]))

		})
	}
}
