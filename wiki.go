package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	reg "regexp"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

type Page struct {
	Title string
	Body  []byte
}

// Regular expression to validate URL ensuring vague paths cannot be accessed on server
var validateURL = reg.MustCompile("^/(save|edit|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	u := validateURL.FindStringSubmatch(r.URL.Path)
	if u == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid URL")
	}
	return u[2], nil
}

func (p *Page) save() error {
	filename := p.Title + ".html"
	return os.WriteFile(filename, p.Body, 0644)
}

func renderTemplate(w http.ResponseWriter, tmp string, p *Page) {
	err := templates.ExecuteTemplate(w, tmp + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
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

func editHandler(w http.ResponseWriter, r *http.Request) {
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

func viewHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}