package route

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"skilldir/model"
)

type Router struct {
	validPath string
	templates *template.Template
}

func (rr Router) MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := rr.validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func (rr Router) makeFileHandler(fn func(http.ResponseWriter, *http.Request, string), fileName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, fileName)
	}
}

func (rr Router) handleURL(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path
	switch path {
	case "/", "/index", "/index.html":
		return nil
	}
	return fmt.Errorf("Bad path: %s", r.URL.Path)
}

func (rr Router) getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := rr.validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}

func (r Router) loadPage(title string) (*model.Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &model.Page{Title: title, Body: body}, nil
}

func (r Router) viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusNotFound)
		return
	}
	renderTemplate(w, "view", p)
}

func (r Router) serveFile(w http.ResponseWriter, r *http.Request, title string) {
	err := handleURL(w, r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	title = "static/view/" + title + ".html"
	http.ServeFile(w, r, title)
}

func (r Router) editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &model.Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func (r Router) saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &model.Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func (r Router) renderTemplate(w http.ResponseWriter, tmpl string, p *model.Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
