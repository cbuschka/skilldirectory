package handler

import (
	"log"
	"net/http"
	"skilldirectory/controller"
	"skilldirectory/data"
	"skilldirectory/errors"
)

var skillsConnector = data.NewAccessor(data.NewFileWriter("skills/"))

func MakeHandler(fn func(http.ResponseWriter, *http.Request, controller.RESTController), cont controller.RESTController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, cont)
	}
}

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
		case *errors.MarshalingError, *errors.InvalidSkillTypeError, *errors.MissingSkillIDError:
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
