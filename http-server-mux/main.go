package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func testHandler(w http.ResponseWriter, _ *http.Request) {
	tpl, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to parse templates: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	tpl.Execute(w, struct {
		Title string
	}{
		Title: "Home page",
	})
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", testHandler)

	// Start lambda event handler
	http.ListenAndServe(":8080", r)
}
