package helper

import (
	"log"
	"os"
)

func GetProperty(key string) string {
	log.Printf("Getting Env: %s", key)
	return os.Getenv(key)
}
