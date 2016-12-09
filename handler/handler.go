package handler

import (
	"net/http"
	"skilldirectory/controller"
	"skilldirectory/data"
)

var skillsConnector = data.NewAccessor(data.NewFileWriter("skills/"))

func MakeHandler(fn func(http.ResponseWriter, *http.Request, controller.RESTController), cont controller.RESTController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, cont)
	}
}

func Handler(w http.ResponseWriter, r *http.Request, cont controller.RESTController) {
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
	}
}
