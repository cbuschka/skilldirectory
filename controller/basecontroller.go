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
func checkForId(url *url.URL) string {
	base := path.Base(url.RequestURI())
	if base == "skills" {
		return ""
	}
	return base
}
