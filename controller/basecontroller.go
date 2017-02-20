package controller

import (
	"net/http"
	"skilldirectory/data"

	"github.com/Sirupsen/logrus"
)

type BaseController struct {
	w http.ResponseWriter
	r *http.Request
	*logrus.Logger
	session    data.DataAccess
	fileSystem data.FileSystem
}

func (bc *BaseController) Init(w http.ResponseWriter, r *http.Request,
	session data.DataAccess, fs data.FileSystem, logger *logrus.Logger) {
	bc.w = w
	bc.r = r
	bc.Logger = logger
	bc.session = session
	bc.fileSystem = fs
}

// SetAllowDefaultMethods adds the following methods to the
// "Access-Control-Allow-Methods" header of w:
// GET, POST, DELETE, OPTIONS
func SetAllowDefaultMethods(w http.ResponseWriter) {
	headers := w.Header().Get("Access-Control-Allow-Methods")
	w.Header().Set("Access-Control-Allow-Methods",
		headers+", GET, POST, DELETE, OPTIONS")
}

// SetAllowDefaultHeaders adds the following headers to the
// "Access-Control-Allow-Headers" header of w:
// Origin, Accept, X-Requested-With, Content-Type, Access-Control-Request-Method,
// Access-Control-Request-Headers
func SetAllowDefaultHeaders(w http.ResponseWriter) {
	headers := w.Header().Get("Access-Control-Allow-Headers")
	w.Header().Set("Access-Control-Allow-Headers",
		headers+", Origin, Accept, X-Requested-With, Content-Type, "+
			"Access-Control-Request-Method, Access-Control-Request-Headers")
}
