package redis

import "github.com/redis/go-redis/v9"

// NewClient returns a new instance of the redis.Client, which is used to interact with the Redis server.
// It takes an environment prefix string as an argument, which is used to load the configuration for the Redis client from the environment variables.
// The configuration is loaded using the newConfig function, which returns an error if there was an issue loading the configuration.
// If an error occurs while loading the configuration, the error is returned immediately and the returned redis.Client is nil.
// The returned redis.Client is ready to be used and does not require any additional setup before interacting with the Redis server.
// It is created with the configuration loaded from the environment variables, and is used to connect to the Redis server at the specified address, using the specified password and database.
// The returned redis.Client is an instance of the github.com/redis/go-redis/v9.Client struct, which is used to interact with the Redis server.
func NewClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return redisClient, nil
}
