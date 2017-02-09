package util

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
)

// LogInit intializes a logger with debug flag
func LogInit() *log.Logger {
	logger := log.New()
	// Log as JSON instead of the default ASCII formatter.
	logger.Formatter = &log.TextFormatter{
		FullTimestamp: true,
	}

	// Output to stderr instead of stdout, could also be a file.
	logger.Out = os.Stderr

	debugFlag := flag.Lookup("debug")
	if debugFlag == nil {
		logger.Level = log.InfoLevel
		return logger
	}
	if debugFlag.Value.String() == "true" {
		logger.Level = log.DebugLevel
	} else {
		logger.Level = log.InfoLevel
	}

	return logger
}
