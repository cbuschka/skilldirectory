package handler

import (
	"errors"
	"fmt"
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

func MakeFileHandler(fn func(http.ResponseWriter, *http.Request, string), fileName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, fileName)
	}
}

func handleURL(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path
	switch path {
	case "/", "/index", "/index.html":
		return nil
	}
	return fmt.Errorf("Bad path: %s", r.URL.Path)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}
