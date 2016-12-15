package handler

import (
	"log"
	"net/http"
	"skilldirectory/controller"
	"skilldirectory/data"
	"skilldirectory/errors"
)

/*
MakeHandler() returns a new function of the adapter type http.HandlerFunc using the passed-in function, fn.
*/
func MakeHandler(fn func(http.ResponseWriter, *http.Request, controller.RESTController), cont controller.RESTController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, cont)
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
func Handler(w http.ResponseWriter, r *http.Request, cont controller.RESTController) {
	log.Printf("Handling Skills Request: %s", r.Method)
	cont.Base().Init(w, r, data.NewAccessor(data.NewFileWriter("skills/")))

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
		case *errors.MarshalingError, *errors.InvalidSkillTypeError, *errors.MissingSkillIDError,
			*errors.IncompletePOSTBodyError:
			statusCode = http.StatusBadRequest
		case *errors.SavingError:
			statusCode = http.StatusInternalServerError
		case *errors.NoSuchIDError:
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		log.Printf("SkillsHandler Method: %s, Err: %v", r.Method, err)
		http.Error(w, err.Error(), statusCode)
	}
}
