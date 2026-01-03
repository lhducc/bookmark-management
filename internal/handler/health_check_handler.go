package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lhducc/bookmark-management/internal/service"
	"net/http"
)

type healthCheckResponse struct {
	Message     string `json:"message"`
	ServiceName string `json:"serviceName"`
	InstanceID  string `json:"instanceID"`
}

type healthCheckErrorResponse struct {
	Error       string `json:"error"`
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}
type healthCheckHandler struct {
	svc service.HealthCheck
}

type HealthCheckHandler interface {
	Check(c *gin.Context)
}

func NewHealthCheckHandler(svc service.HealthCheck) HealthCheckHandler {
	return &healthCheckHandler{
		svc: svc,
	}
}

// Check health of the service
// @Summary Check health of the service
// @Description Check health of the service
// @Tags Health Check
// @Success 200 {object} healthCheckResponse
// @Router /health-check [get]
func (h *healthCheckHandler) Check(c *gin.Context) {
	message, serviceName, instanceID, err := h.svc.Check(c)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, healthCheckErrorResponse{
			Error:       "Internal Server Error",
			Message:     message,
			ServiceName: serviceName,
			InstanceID:  instanceID,
		})
		return
	}

	c.JSON(http.StatusOK, healthCheckResponse{
		Message:     message,
		ServiceName: serviceName,
		InstanceID:  instanceID,
	})
}
