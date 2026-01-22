package service

import (
	"context"
	"errors"
	"github.com/lhducc/bookmark-management/internal/repository"
	"github.com/lhducc/bookmark-management/internal/repository/mocks"
	mockKeyGen "github.com/lhducc/bookmark-management/pkg/stringutils/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var testError = errors.New("test error")

func TestShortenUrl_ShortenUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		url string
		exp int

		setupMockRepo   func(t *testing.T, ctx context.Context, url string, exp int) repository.UrlStorage
		setupMockKeyGen func() *mockKeyGen.KeyGen

		expectedCode string
		expectErr    error
		expectedLen  int
	}{
		{
			name: "normal case",

			url: "https://www.google.com",
			exp: 10,

			setupMockRepo: func(t *testing.T, ctx context.Context, url string, exp int) repository.UrlStorage {
				repoMock := mocks.NewUrlStorage(t)
				repoMock.On(
					"StoreURLIfNotExists",
					ctx,
					mock.AnythingOfType("string"),
					url,
					exp,
				).Return(true, nil)
				return repoMock
			},
			setupMockKeyGen: func() *mockKeyGen.KeyGen {
				keyGenMock := mockKeyGen.NewKeyGen(t)
				keyGenMock.On("GenerateCode", urlCodeLength).Return("abc1237", nil)
				return keyGenMock
			},

			expectedCode: "abc1237",
			expectedLen:  7,
			expectErr:    nil,
		},

		{
			name: "key gen error",

			url: "https://www.google.com",
			exp: 10,

			setupMockRepo: func(t *testing.T, ctx context.Context, url string, exp int) repository.UrlStorage {
				repoMock := mocks.NewUrlStorage(t)
				return repoMock
			},
			setupMockKeyGen: func() *mockKeyGen.KeyGen {
				keyGenMock := mockKeyGen.NewKeyGen(t)
				keyGenMock.On("GenerateCode", urlCodeLength).Return("", testError)
				return keyGenMock
			},

			expectedCode: "",
			expectErr:    testError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cxt := context.Background()

			urlStorageMock := tc.setupMockRepo(t, cxt, tc.url, tc.exp)
			mockKeyGen := tc.setupMockKeyGen()
			testSvc := NewShortenUrl(urlStorageMock, mockKeyGen)

			urlCode, err := testSvc.ShortenUrl(cxt, tc.url, tc.exp)

			assert.Equal(t, tc.expectedLen, len(urlCode))
			assert.Equal(t, tc.expectErr, err)
			if err == nil {
				assert.Equal(t, urlSafeRegex.MatchString(urlCode), true)
			}

			assert.Equal(t, tc.expectedCode, urlCode)
		})
	}
}

func TestShortenUrl_GetUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		code string

		setupMock func(t *testing.T) *mocks.UrlStorage

		expURL    string
		expectErr error
	}{
		{
			name: "normal case",

			code: "abc1234",

			setupMock: func(t *testing.T) *mocks.UrlStorage {
				repo := mocks.NewUrlStorage(t)
				repo.
					On("GetURL", mock.Anything, "abc1234").
					Return("https://google.com", nil).
					Once()
				return repo
			},

			expURL:    "https://google.com",
			expectErr: nil,
		},
		{
			name: "code not found -> map redis.Nil to ErrCodeNotFound",

			code: "notfound",

			setupMock: func(t *testing.T) *mocks.UrlStorage {
				repo := mocks.NewUrlStorage(t)
				repo.
					On("GetURL", mock.Anything, "notfound").
					Return("", redis.Nil).
					Once()
				return repo
			},

			expURL:    "",
			expectErr: ErrCodeNotFound,
		},
		{
			name: "repo returns other error -> passthrough",

			code: "errcode",

			setupMock: func(t *testing.T) *mocks.UrlStorage {
				repo := mocks.NewUrlStorage(t)
				repo.
					On("GetURL", mock.Anything, "errcode").
					Return("", redis.ErrClosed).
					Once()
				return repo
			},

			expURL:    "",
			expectErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			repoMock := tc.setupMock(t)

			svc := NewShortenUrl(repoMock, nil)

			url, err := svc.GetUrl(ctx, tc.code)

			if tc.expectErr != nil {
				assert.True(t, errors.Is(err, tc.expectErr), "expected error to match")
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expURL, url)
		})
	}
}
