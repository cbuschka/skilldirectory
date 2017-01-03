package util

import (
	"flag"
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

	debugFlag := flag.Lookup("debug")
	if debugFlag.Value.String() == "true" {
		logger.Level = log.DebugLevel
	} else {
		logger.Level = log.InfoLevel
	}

	return logger
}
