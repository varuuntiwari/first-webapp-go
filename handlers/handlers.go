package handlers

import (
	"errors"
	"html/template"
	io "io/ioutil"
	"net/http"
	reg "regexp"
)


var validateURL = reg.MustCompile("^/(save|edit|view)/([a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))

type Page struct {
	Title string
	Body  []byte
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	u := validateURL.FindStringSubmatch(r.URL.Path)
	if u == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid url")
	}
	return u[2], nil
}

func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := io.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmp string, p *Page) {
	err := templates.ExecuteTemplate(w, tmp + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Page) save() error {
	filename := p.Title + ".html"
	return io.WriteFile(filename, p.Body, 0644)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	val := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(val)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}