package main

import (
	handle "first-webapp-go/handlers"
	"log"
	"net/http"
)

// Regular expression to validate URL ensuring vague paths cannot be accessed on server

func main() {
	http.HandleFunc("/save/", handle.SaveHandler)
	http.HandleFunc("/edit/", handle.EditHandler)
	http.HandleFunc("/view/", handle.ViewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}