package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	templates = template.Must(template.ParseFiles("templates/index.html"))
	debtStore []DebtItem
	debtFile  = "debt.csv"
)

type DebtItem struct {
	Person   string
	Amount   Money
	Note     string
	Creation time.Time
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

func LoadDebtData(filename string) ([]DebtItem, error) {
	csvfile, err := os.Open(filename)
	output := []DebtItem{}
	if err != nil {
		return output, err
	}
	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1
	for {
		entry, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return output, err
		}
		cents, err := strconv.Atoi(entry[1])
		if err != nil {
			return output, err
		}

		t, err := strconv.ParseInt(entry[3], 10, 64)
		if err != nil {
			return output, err
		}

		output = append(output, DebtItem{entry[0], Money{cents}, entry[2], time.Unix(int64(t), 0)})
	}
	return output, nil
}

func SaveDebtData(l []DebtItem, filename string) error {
	csvfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		if _, err := os.Stat(filename); err != nil {
			csvfile, _ = os.Create(filename)
		} else {
			fmt.Printf("Error opening debt file: %v", err)
			panic(err)
		}
	}
	defer csvfile.Close()
	writer := csv.NewWriter(csvfile)
	for _, item := range l {
		err = writer.Write([]string{item.Person,
			fmt.Sprintf("%d", item.Amount.Cents),
			item.Note,
			strconv.Itoa(int(item.Creation.Unix()))})
		if err != nil {
			fmt.Println(err)
		}
	}
	err = writer.Error()
	if err != nil {
		fmt.Println(err)
	}
	writer.Flush()
	return nil
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
	notes := r.FormValue("notes")
	moneyAmount, _ := NewMoney(r.FormValue("amount"))

	debtStore = append(debtStore, DebtItem{Person: person, Amount: moneyAmount, Note: notes, Creation: time.Now()})

	SaveDebtData(debtStore, debtFile)
	http.Redirect(w, r, "/", 301)
}

func init() {
	data, err := LoadDebtData(debtFile)
	if err != nil {
		fmt.Println(err)
		debtStore = make([]DebtItem, 0)
	} else {
		debtStore = data
	}
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
