package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

var (
	templates = template.Must(template.ParseFiles("templates/index.html"))
	debtStore []DebtItem
)

type DebtItem struct {
	Person string
	Amount Money
	Note   string
}

type Money struct {
	Cents int64
}

func (v Money) Add(lhs Money) Money {
	return Money{Cents: v.Cents + lhs.Cents}
}
func (v Money) String() string {
	return fmt.Sprintf("$%d.%02d", v.Cents/100, v.Cents%100)
}

type DebtData struct {
	Ower       string
	OwedAmount Money
	DebtStore  []DebtItem
}

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	data := DebtData{"Ben", Money{14}, debtStore}
	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleAddDebt(w http.ResponseWriter, r *http.Request) {
	person := r.FormValue("person")
	amount, _ := strconv.ParseInt(r.FormValue("amount"), 10, 64)
	moneyAmount := Money{amount}

	debtStore = append(debtStore, DebtItem{Person: person, Amount: moneyAmount, Note: ""})

	http.Redirect(w, r, "/", 301)
}

func init() {
	debtStore = make([]DebtItem, 0)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", BaseHandler).Methods("GET")

	r.HandleFunc("/", HandleAddDebt).Methods("POST")
	//r.HandleFunc("/user/get", HandleUserGet)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/").HandlerFunc(BaseHandler)

	http.Handle("/", r)

	fmt.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v", err)
	}
}
