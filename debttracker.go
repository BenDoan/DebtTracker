package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
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
	Cents int
}

func NewMoney(v string) (Money, error) {
	dollars := 0
	cents := 0
	_, err := fmt.Sscanf(v, "%d.%02d", &dollars, &cents)
	if err != nil {
		return Money{0}, err
	}
	return Money{Cents: dollars*100 + cents}, nil
}

func (v Money) Add(lhs Money) Money {
	return Money{Cents: v.Cents + lhs.Cents}
}

func (v Money) Subtract(lhs Money) Money {
	return Money{Cents: v.Cents - lhs.Cents}
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
	ower, amount := CalculateOwed(debtStore)
	data := DebtData{ower, amount, debtStore}

	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CalculateOwed(debtStore []DebtItem) (ower string, amount Money) {
	aAmount := Money{0}
	bAmount := Money{0}

	for _, item := range debtStore {
		if item.Person == "ben" {
			aAmount = aAmount.Add(item.Amount)
		} else if item.Person == "mitchell" {
			bAmount = bAmount.Add(item.Amount)
		}
	}

	if aAmount.Cents > bAmount.Cents {
		return "Ben", aAmount.Subtract(bAmount)
	} else {
		return "Mitchell", bAmount.Subtract(aAmount)
	}
}

func HandleAddDebt(w http.ResponseWriter, r *http.Request) {
	person := r.FormValue("person")
	moneyAmount, _ := NewMoney(r.FormValue("amount"))

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
