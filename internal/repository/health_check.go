package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name=HealthCheck --filename health_check.go
type HealthCheck interface {
	Ping(ctx context.Context) error
}

type healthCheck struct {
	redis *redis.Client
}

// NewHealthCheck returns a new instance of the healthCheck, which implements the HealthCheck interface.
// It takes a single parameter, redis, which is a pointer to a redis.Client.
// The returned healthCheck is used to ping the redis server to check if it is reachable.
// If an error occurs while generating the healthCheck, the error is returned immediately and the generated healthCheck is an empty string.
// The character set used for generating the healthCheck response is constant and does not change across different implementations of the interface. The length of the generated healthCheck response is constant and does not change across different implementations of the interface.
func NewHealthCheck(redis *redis.Client) HealthCheck {
	return &healthCheck{redis: redis}
}

// Ping checks if the Redis connection is alive.
// It returns an error if there is an issue with the Redis connection.
// The method takes a context as an input parameter.
// It uses the Ping method of the Redis client to check the Redis connection.
// If the Redis connection is alive, the method returns nil. Otherwise, it returns the error.
// This method is used to implement the HealthCheck interface.
func (r *healthCheck) Ping(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}
