package logger

import (
	"github.com/rs/zerolog"
	"os"
)

// SetLogLevel sets the global log level based on the LOG_LEVEL environment variable.
// If the variable is not set, or if the value cannot be parsed into a log level,
// the global log level is set to InfoLevel.
func SetLogLevel() {
	level, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}
