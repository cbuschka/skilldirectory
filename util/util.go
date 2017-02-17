package util

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/url"
	"os"
	"path"
	"skilldirectory/errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	// Uncomment below line(s) to recognize these image types as valid
	// _ "image/gif"
)

// GetProperty returns the value from the environment or key/value store
func GetProperty(key string) string {
	log.Printf("Getting Env: %s", key)
	return os.Getenv(key)
}

// CheckForID checks to see if an ID (e.g. 59317629-bcc3-11e6-9f43-6c4008bcfa84)
// has been appended to the end of the specified URL. If one has, then that ID
// will be returned. If not, then an empty string is returned ("").
func CheckForID(url *url.URL) string {
	base := path.Base(url.RequestURI())
	if IsValidEndpoint(url.EscapedPath()) {
		return ""
	}
	return base
}

// IsValidEndpoint returns true if endpoint is an endpoint being
// served by the SkillDirectory server AND doesn't contain an ID.
func IsValidEndpoint(endpoint string) bool {
	endpoints := []string{
		"/api/skills", "/api/skills/",
		"/api/teammembers", "/api/teammembers/",
		"/api/tmskills", "/api/tmskills/",
		"/api/links", "/api/links/",
		"/api/skillreviews", "/api/skillreviews/",
		"/api/skillicons", "/api/skillicons/",
	}
	if StringSliceContains(endpoints, endpoint) {
		return true
	}
	return false
}

// StringSliceContains returns true if slice contains target, false if not.
func StringSliceContains(slice []string, target string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

// Returns the directory at the root of the specified path. Ignores starting slashes (regards
// "/skills/files" as "skills/files". Calling with "skills/files/whatever/1234-5678-9101" would return "skills/".
func getRootDir(path string) string {
	if path[0] == '/' {
		path = path[1:]
	}
	if path[0:4] == "api/" {
		path = path[4:]
	}
	var rootDir string
	if strings.Index(path, "/") != -1 {
		rootDir = path[:strings.Index(path, "/")+1]
	} else {
		rootDir = path + "/"
	}
	return rootDir
}

// ValidateIcon returns a non-nil error if the icon cannot be decoded into a
// recognized image format. If image is successfully decoded, then return string
// contains the image encoding format (e.g. PNG, JPG, GIF)
func ValidateIcon(icon io.Reader) (string, error) {
	_, format, err := image.Decode(icon)
	if err != nil {
		return "", errors.InvalidDataModelState(fmt.Errorf(
			"failed to decode image data in %q field: %s", "Icon", err.Error()))
	}
	return format, nil
}

// SanitizeInput escapes the single quotes in the Cassandra query
func SanitizeInput(input string) string {
	return strings.Replace(input, "'", "''", -1)
}
