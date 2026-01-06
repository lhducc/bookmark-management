package service

import (
	"context"
	"github.com/lhducc/bookmark-management/internal/repository"
	"github.com/lhducc/bookmark-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestShortenUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		url string
		exp int

		setupMock func(t *testing.T, ctx context.Context, url string, exp int) repository.UrlStorage

		expectErr   error
		expectedLen int
	}{
		{
			name: "normal case",

			url: "https://www.google.com",
			exp: 10,

			setupMock: func(t *testing.T, ctx context.Context, url string, exp int) repository.UrlStorage {
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

			expectedLen: 7,
			expectErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cxt := context.Background()

			urlStorageMock := tc.setupMock(t, cxt, tc.url, tc.exp)

			testSvc := NewShortenUrl(urlStorageMock)

			urlCode, err := testSvc.ShortenUrl(cxt, tc.url, tc.exp)

			assert.Equal(t, tc.expectedLen, len(urlCode))
			assert.Equal(t, tc.expectErr, err)
			assert.Equal(t, urlSafeRegex.MatchString(urlCode), true)
		})
	}
}
