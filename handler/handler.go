package handler

import (
	"log"
	"net/http"
	"skilldirectory/controller"
	"skilldirectory/data"
	"skilldirectory/errors"
	"strings"
)

/*
MakeHandler() returns a new function of the adapter type http.HandlerFunc using the passed-in function, fn.
*/
func MakeHandler(fn func(http.ResponseWriter, *http.Request, controller.RESTController, data.DataAccess), cont controller.RESTController, session data.DataAccess) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, cont, session)
	}
}

/*
Handler() should be invoked to handle responding to the passed-in HTTP request. Responses are sent via
the passed-in http.ResponseWriter.

The passed-in RESTController is first initialized using the specified http.ResponseWriter and http.Request,
and is connected to the Skills filesystem/database. Once initialized, the RESTController is used to handle
responses to the passed-in HTTP request.

If the RESTController generates any errors, then Handler() will
log them, and respond to the request with the appropriate error.
*/
func Handler(w http.ResponseWriter, r *http.Request, cont controller.RESTController, session data.DataAccess) {
	log.Printf("Handling Request: %s", r.Method)
	cont.Base().Init(w, r, session)

	var err error
	switch r.Method {
	case http.MethodGet:
		err = cont.Get()
	case http.MethodPost:
		err = cont.Post()
	case http.MethodDelete:
		err = cont.Delete()
	}

	var statusCode int
	if err != nil {
		switch err.(type) {
		case *errors.MarshalingError, *errors.InvalidSkillTypeError, *errors.MissingIDError,
			*errors.IncompletePOSTBodyError, *errors.InvalidPOSTBodyError:
			statusCode = http.StatusBadRequest
		case *errors.SavingError:
			statusCode = http.StatusInternalServerError
		case *errors.NoSuchIDError:
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		log.Printf("Handler Method: %s, Err: %v", r.Method, err)
		http.Error(w, err.Error(), statusCode)
	}
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
