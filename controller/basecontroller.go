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

// GetDefaultMethods returns a string containing a ", " seperated list of the
// default HTTP methods for an endpoint.
func GetDefaultMethods() string {
	return "GET, POST, DELETE, OPTIONS"
}

// GetDefaultHeaders returns s string containing a ", " seperated list of the
// default HTTP methods for an endpoint.
func GetDefaultHeaders() string {
	return "Origin, Accept, X-Requested-With, Content-Type, " +
		"Access-Control-Request-Methods, Access-Control-Request-Headers, " +
		"Access-Control-Allow-Methods"
}
