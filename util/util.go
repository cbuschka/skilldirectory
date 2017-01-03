package util

import (
	"net/url"
	"os"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
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

// IsValidEndpoint() returns true if endpoint is an endpoint being
// served by the SkillDirectory server AND doesn't contain an ID.
func IsValidEndpoint(endpoint string) bool {
	endpoints := []string{
		"/skills", "/skills/",
		"/teammembers", "/teammembers/",
		"/tmskills", "/tmskills/",
		"/links", "/links/",
		"/skillreviews", "/skillreviews/",
	}
	if StringSliceContains(endpoints, endpoint) {
		return true
	}
	return false
}

// StringSliceContains() returns true if slice contains target, false if not.
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
	var rootDir string
	if strings.Index(path, "/") != -1 {
		rootDir = path[:strings.Index(path, "/")+1]
	} else {
		rootDir = path + "/"
	}
	return rootDir
}
