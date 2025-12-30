package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lhducc/bookmark-management/internal/handler"
	"github.com/lhducc/bookmark-management/internal/service"
	"net/http"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type api struct {
	app *gin.Engine
	cfg *Config
}

// New returns a new instance of the api, which implements the Engine interface.
// The returned api is used to start the HTTP server and listen for incoming requests on port 8080.
// The api is created with a gin.Engine instance, which is used to start the server.
// The registerEP method is called on the returned api to register the endpoints for the API.
// The returned api is ready to be used and does not require any additional setup before starting the server.
func New(cfg *Config) Engine {
	a := &api{
		app: gin.New(),
		cfg: cfg,
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
	// Service
	passSvc := service.NewPassword()
	healthCheckSvc := service.NewHealthCheck(a.cfg.ServiceName, a.cfg.InstanceID)

	// Handler
	passHandler := handler.NewPassword(passSvc)
	healthCheckHandler := handler.NewHealthCheckHandler(healthCheckSvc)

	// Router
	a.app.GET("/gen-pass", passHandler.GenPass)
	a.app.GET("/health-check", healthCheckHandler.Check)
}
