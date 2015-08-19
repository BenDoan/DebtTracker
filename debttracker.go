package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))

type DebtData struct {
	Ower       string
	OwedAmount float64
}

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	data := DebtData{"Ben", 14.5}
	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", BaseHandler)
	//r.HandleFunc("/user/get", HandleUserGet)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/").HandlerFunc(BaseHandler)

	http.Handle("/", r)

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
