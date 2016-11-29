package handler

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"skilldir/data"
	"skilldir/model"
)

var validPath = regexp.MustCompile("^/(edit|save|view|skills)/([a-zA-Z0-9]+$)")
var templates = template.Must(template.ParseGlob("templates/*"))
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

func loadSkill(title string) (*model.Skill, error) {
	skill := model.Skill{}
	err := skillsConnector.Read(title, &skill)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func SkillsHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Printf("Handling Skills Request")
	p, err := loadSkill(title)
	if err != nil {
		fmt.Println(err)
		return
	}
	renderTemplate(w, "skilltemplate", p)
}

func ServeFile(w http.ResponseWriter, r *http.Request, title string) {
	err := handleURL(w, r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	title = "static/view/" + title + ".html"
	http.ServeFile(w, r, title)
}

func renderTemplate(w http.ResponseWriter, tmpl string, object interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", object)
	if err != nil {
		log.Panicf("Template Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
