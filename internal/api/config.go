package api

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark-management" envconfig:"SERVICE_NAME"`
	InstanceID  string `default:"" envconfig:"INSTANCE_ID"`
}

// NewConfig returns a new instance of Config, which is used to configure the API.
// It populates the fields of the returned Config instance with values from environment variables.
// If an error occurs while populating the fields, it returns an error immediately.
// The returned Config instance is ready to be used and does not require any additional setup before starting the API.
// If the InstanceID field is empty, it generates a random UUID and assigns it to the field.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.NewString()
	}

	return cfg, nil
}
