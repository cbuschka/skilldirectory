package handler

import (
	"log"
	"net/http"
	"regexp"
	"skilldirectory/controller"
	"skilldirectory/data"
)

var validPath = regexp.MustCompile("^/(skills)|(/([a-zA-Z0-9]+$))")
var skillsConnector = data.NewAccessor(data.NewFileWriter("skills/"))

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string), cont controller.RESTController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			log.Panicf("Doesn't pass valid path test: %s\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func Handler(w http.ResponseWriter, r *http.Request, cont controller.RESTController) {
	controller.NewController(w, *r, data.NewAccessor(data.NewFileWriter("skills/")))
	switch r.Method {
	case http.MethodGet:
		cont.Get()

	case http.MethodPost:
		cont.Post()
	}
}
