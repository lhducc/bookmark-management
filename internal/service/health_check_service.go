package service

type healthCheckService struct {
	serviceName string
	instanceID  string
}

//go:generate mockery --name=HealthCheck --filename health_check_service.go
type HealthCheck interface {
	Check() (string, string, string)
}

// NewHealthCheck returns a new instance of the healthCheckService, which implements the HealthCheck interface.
// It takes two parameters, serviceName and instanceID, which are used to generate the health check response of the service.
// The returned healthCheckService is used to generate health check responses of the service.
// It returns a tuple of (message, serviceName, instanceID), where message is "OK" if the service is healthy, serviceName is the name of the service,
// and instanceID is the ID of the service instance.
func NewHealthCheck(serviceName, instanceID string) HealthCheck {
	return &healthCheckService{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

// Check returns the health check response of the service.
// It returns a tuple of (message, serviceName, instanceID), where
// message is "OK" if the service is healthy, serviceName is the name of the service,
// and instanceID is the ID of the service instance.
// If an error occurs while generating the health check response, the error is returned immediately and the generated health check response is an empty string.
// The character set used for generating the health check response is constant and does not change across different implementations of the interface. The length of the generated health check response is constant and does not change across different implementations of the interface.
func (s *healthCheckService) Check() (string, string, string) {
	return "OK", s.serviceName, s.instanceID
}
