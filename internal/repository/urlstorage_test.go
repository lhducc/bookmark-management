package repository

import (
	"context"
	redisPkg "github.com/lhducc/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUrlStorage_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func() *redis.Client

		expectErr  error
		verifyFunc func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",

			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},

			expectErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				url, err := r.Get(ctx, "123").Result()
				assert.Nil(t, err)
				assert.Equal(t, url, "https://google.com")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMock := tc.setupMock()
			testRepo := NewUrlStorage(redisMock)

			err := testRepo.StoreURL(ctx, "123", "https://google.com")
			assert.Equal(t, tc.expectErr, err)
			if err == nil {
				tc.verifyFunc(ctx, redisMock)
			}
		})
	}
}

func TestUrlStorage_StoreURLIfNotExists(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		code string
		url  string
		exp  int

		setupMock func() *redis.Client

		expectOK  bool
		expectErr error
	}{
		{
			name: "normal case",

			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},

			code: "123",
			url:  "https://google.com",
			exp:  10,

			expectOK:  true,
			expectErr: nil,
		},
		{
			name: "key already exists",

			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				err := mock.Set(context.Background(), "123", "old-url", time.Hour).Err()
				require.NoError(t, err)
				return mock
			},

			code: "123",
			url:  "https://google.com",
			exp:  10,

			expectOK:  false,
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMock := tc.setupMock()
			testRepo := NewUrlStorage(redisMock)

			ok, err := testRepo.StoreURLIfNotExists(ctx, tc.code, tc.url, tc.exp)

			assert.Equal(t, tc.expectErr, err)
			assert.Equal(t, tc.expectOK, ok)

		})
	}
}
