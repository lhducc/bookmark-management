package main

import (
	"github.com/lhducc/bookmark-management/internal/api"
	redisPkg "github.com/lhducc/bookmark-management/pkg/redis"
)

// @title Bookmark Management API
// @version 1.0
// @description API documentation for bookmark service.
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}

	app := api.New(cfg, redisClient)
	app.Start()
}
