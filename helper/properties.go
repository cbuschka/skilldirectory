package helper

import (
	"log"
	"os"
)

// GetProperty returns the value from the environment or key/value store
func GetProperty(key string) string {
	log.Printf("Getting Env: %s", key)
	return os.Getenv(key)
}
