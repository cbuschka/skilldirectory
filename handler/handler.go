package handler

import (
	"log"
	"net/http"
	"regexp"
	"skilldirectory/data"
)

var validPath = regexp.MustCompile("^/(skills)|(/([a-zA-Z0-9]+$))")
var skillsConnector = data.NewAccessor(data.NewFileWriter("skills/"))

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
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
