package controller

import (
	"fmt"
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
func checkForId(url *url.URL) string {
	base := path.Base(url.RequestURI())
	if url.EscapedPath() != "/skills" && url.EscapedPath() != "/skills/" {
		return base
	}
	return ""
}

func extractSkillFilter(url *url.URL) (string, error) {
	// Extract the URL path's base and make sure it is "/skills".
	// Return error if it's not, because it doesn't make sense to filter
	// by skill type if the base isn't "/skills".
	base := path.Base(url.Path)
	if base != "skills" {
		return "", fmt.Errorf("URL path base must be \"skills\" to filter by skill type")
	}

	// Extract the query string from the URL as a key, value map.
	// Then search the map for a "skilltype" filter. If the query
	// string contains this filter, return the value. If not, return
	// an empty string ("").
	query := url.Query()
	for key, val := range query {
		if key == "skilltype" {
			return val[0], nil
		}
	}
	return "", nil
}
