package util

import "github.com/gocql/gocql"

// NewID returns a uuid from the current library
func NewID() string {
	return gocql.TimeUUID().String()
}
