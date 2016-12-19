package controller

import (
	"net/http"
	"net/url"
	"path"
	"skilldirectory/data"
)

type BaseController struct {
	w       http.ResponseWriter
	r       *http.Request
	session data.DataAccess
}

func (bc *BaseController) Init(w http.ResponseWriter, r *http.Request, session data.DataAccess) {
	bc.w = w
	bc.r = r
	bc.session = session
}

// checkForID checks to see if an ID (e.g. 59317629-bcc3-11e6-9f43-6c4008bcfa84)
// has been appended to the end of the specified URL. If one has, then that ID
// will be returned. If not, then an empty string is returned ("").
func checkForID(url *url.URL) string {
	base := path.Base(url.RequestURI())
	if isValidEndpoint(url.EscapedPath()) {
		return ""
	}
	return base
}

// isValidEndpoint() returns true if endpoint is an endpoint being
// served by the SkillDirectory server AND doesn't contain an ID.
func isValidEndpoint(endpoint string) bool {
	endpoints := []string {
		"/skills", "/skills/",
		"/teammembers", "/teammembers/",
	}
	if stringSliceContains(endpoints, endpoint) {
		return true
	}
	return false
}

// stringSliceContains() returns true if slice contains target, false if not.
func stringSliceContains(slice []string, target string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}