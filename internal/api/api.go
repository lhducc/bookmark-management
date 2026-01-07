package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lhducc/bookmark-management/docs"
	"github.com/lhducc/bookmark-management/internal/handler"
	"github.com/lhducc/bookmark-management/internal/repository"
	"github.com/lhducc/bookmark-management/internal/service"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type api struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
}

// New returns a new instance of the api, which implements the Engine interface.
// The returned api is used to start the HTTP server and listen for incoming requests on port 8080.
// The api is created with a gin.Engine instance, which is used to start the server.
// The registerEP method is called on the returned api to register the endpoints for the API.
// The returned api is ready to be used and does not require any additional setup before starting the server.
func New(cfg *Config, redisClient *redis.Client) Engine {
	a := &api{
		app:         gin.New(),
		cfg:         cfg,
		redisClient: redisClient,
	}
	a.registerEP()
	return a
}

// Start starts the HTTP server and listens for incoming requests on port 8080.
// It returns an error if there was an issue starting the server.
// The server is started using the gin.Engine instance stored in the api struct.
func (a *api) Start() error {
	return a.app.Run(fmt.Sprintf(":%s", a.cfg.AppPort))
}

// ServeHTTP serves HTTP requests to the gin.Engine instance.
// It implements the http.Handler interface and is used to serve HTTP requests.
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

func (a *api) registerEP() {
	//Repository
	urlRepo := repository.NewUrlStorage(a.redisClient)
	healthCheckRepo := repository.NewHealthCheck(a.redisClient)

	// Service
	passSvc := service.NewPassword()
	healthCheckSvc := service.NewHealthCheck(a.cfg.ServiceName, a.cfg.InstanceID, healthCheckRepo)
	urlShortenSvc := service.NewShortenUrl(urlRepo)

	// Handler
	passHandler := handler.NewPassword(passSvc)
	healthCheckHandler := handler.NewHealthCheckHandler(healthCheckSvc)
	urlShortenHandler := handler.NewUrlShortenHandler(urlShortenSvc)

	// Router
	a.app.GET("/gen-pass", passHandler.GenPass)
	a.app.GET("/health-check", healthCheckHandler.Check)
	a.app.POST("/v1/links/shorten", urlShortenHandler.ShortenUrl)

	// Swagger
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
