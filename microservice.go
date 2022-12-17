package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Passenger struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	PhoneNo   int    `json:"PhoneNo"`
	Email     string `json:"Email"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/passenger/{userid}", findPassenger).Methods("GET", "POST", "PATCH")
	router.HandleFunc("/api/v1/passenger", allPassenger)

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func findPassenger(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if val, ok := isPassengerExist(params["userid"]); ok {
		json.NewEncoder(w).Encode(val)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid User ID")
	}
}

func allPassenger(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	var allPassengers, _ = getPassenger()
	//found := false
	results := map[string]Passenger{}

	if value := query.Get("q"); len(value) > 0 {
		for k, v := range allPassengers {
			if strings.Contains(strings.ToLower(v.FirstName), strings.ToLower(value)) {
				results[k] = v
				//found = true
			}
		}
	} else {
		passengersWrapper := struct {
			Passengers map[string]Passenger `json:"Passengers"`
		}{allPassengers}
		json.NewEncoder(w).Encode(passengersWrapper)
		return
	}
}

func getPassenger() (map[string]Passenger, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}

	results, err := db.Query("SELECT * from Passenger")
	if err != nil {
		panic(err.Error())
	}

	var allPassengers map[string]Passenger = map[string]Passenger{}

	for results.Next() {
		var p Passenger
		var id string
		err = results.Scan(&id, &p.FirstName, &p.LastName, &p.PhoneNo, &p.Email)
		if err != nil {
			panic(err.Error())
		}
		allPassengers[id] = p
	}
	return allPassengers, true
}

func isPassengerExist(id string) (Passenger, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var p Passenger
	result := db.QueryRow("SELECT * from Passenger WHERE PassengerID=?", id)
	err = result.Scan(&id, &p.FirstName, &p.LastName, &p.PhoneNo, &p.Email)
	if err == sql.ErrNoRows {
		return p, false
	}
	return p, true
}
