package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"testing"
)

// InitMockRedis initializes a mock Redis instance for testing purposes.
// It returns a pointer to a Redis client connected to the mock instance.
// The mock instance is automatically shut down when the test function returns.
// The Addr field of the Redis options is set to the address of the mock instance.
// The function should be used as a helper when writing unit tests for code that interacts with Redis.
func InitMockRedis(t *testing.T) *redis.Client {
	mock := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{
		Addr: mock.Addr(),
	})
}
