package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"skilldir/data"
	"skilldir/model"
)

var templates = template.Must(template.ParseGlob("templates/*"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+$)")
var dataConnector = data.FileWriter{}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusNotFound)
		return
	}
	renderTemplate(w, "view", p)
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *model.Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	testAddSkills()
	http.HandleFunc("/", makeFileHandler(serveFile, "index"))
	http.HandleFunc("/index", makeFileHandler(serveFile, "index"))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}

func testAddSkills() {
	fmt.Println("Adding Skill")
	testSkillName := "Test"
	testSkillType := "language"
	newSkill := model.NewSkill(testSkillName, testSkillType)
	err := dataConnector.Save("skills", newSkill.Name, newSkill)
	if err != nil {
		fmt.Println(err)
	}
	readSkill := model.Skill{}
	err = dataConnector.Read("skills", newSkill.Name, &readSkill)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Name: %s Type: %s", readSkill.Name, readSkill.SkillType)
	}
}
