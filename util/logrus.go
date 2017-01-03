package util

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func LogInit() *log.Logger {
	logger := log.New()
	// Log as JSON instead of the default ASCII formatter.
	logger.Formatter = &log.TextFormatter{
		FullTimestamp: true,
	}

	// Output to stderr instead of stdout, could also be a file.
	logger.Out = os.Stderr

	// Only log the warning severity or above.
	logger.Level = log.InfoLevel

	return logger
}
