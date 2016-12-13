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
	var statusCode int

	switch r.Method {
	case http.MethodGet:
		err = cont.Get()
		statusCode = http.StatusNotFound
		switch err.(type) {
		case *errors.NoSuchIDError:
			statusCode = http.StatusNotFound
		}
	case http.MethodPost:
		err = cont.Post()
		statusCode = http.StatusBadRequest
		switch err.(type) {
		case *errors.MarshalingError:
			statusCode = http.StatusBadRequest
		case *errors.InvalidSkillTypeError:
			statusCode = http.StatusBadRequest
		case *errors.SavingError:
			statusCode = http.StatusInternalServerError
		}
	case http.MethodDelete:
		err = cont.Delete()
		switch err.(type) {
		case *errors.NoSuchIDError:
			statusCode = http.StatusNotFound
		case *errors.MissingSkillIDError:
			statusCode = http.StatusBadRequest
		}
	}

	if err != nil {
		log.Printf("SkillsHandler Method: %s, Err: %v", r.Method, err)
		http.Error(w, err.Error(), statusCode)
	}
}
