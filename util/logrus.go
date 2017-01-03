package util

import (
	"flag"
	"fmt"
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
	fmt.Println(debugFlag)
	// Only log the warning severity or above.
	logger.Level = log.DebugLevel

	return logger
}
