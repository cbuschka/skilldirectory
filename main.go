package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"skilldir/data"
	"skilldir/model"
)

var templates = template.Must(template.ParseGlob("templates/*"))
var validPath = regexp.MustCompile("^/(edit|save|view|skills)/([a-zA-Z0-9]+$)")
var skillsConnector = data.NewAccessor(data.NewFileWriter("skills/"))

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
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

func makeFileHandler(fn func(http.ResponseWriter, *http.Request, string), fileName string) http.HandlerFunc {
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

func loadPage(title string) (*model.Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &model.Page{Title: title, Body: body}, nil
}

func loadSkill(title string) (*model.Skill, error) {
	skill := model.Skill{}
	err := skillsConnector.Read(title, &skill)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusNotFound)
		return
	}
	renderTemplate(w, "view", p)
}

func skillsHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Printf("Handling Skills Request")
	p, err := loadSkill(title)
	if err != nil {
		fmt.Println(err)
		return
	}
	renderTemplate(w, "skilltemplate", p)
}

func serveFile(w http.ResponseWriter, r *http.Request, title string) {
	err := handleURL(w, r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	title = "static/view/" + title + ".html"
	http.ServeFile(w, r, title)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("Test")
	p, err := loadPage(title)
	if err != nil {
		p = &model.Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &model.Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, object interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", object)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	addSkillsToView()
	http.HandleFunc("/", makeFileHandler(serveFile, "index"))
	http.HandleFunc("/index", makeFileHandler(serveFile, "index"))
	http.HandleFunc("/skills/", makeHandler(skillsHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}

func addSkillsToView() {
	skillsConnector.Save("Golang", model.NewSkill("Golang", "language"))
}
