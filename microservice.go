package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Passenger struct {
	Username  string `json:"Username"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	PhoneNo   int    `json:"PhoneNo"`
	Email     string `json:"Email"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/passenger/{username}", findPassenger).Methods("GET", "POST", "PATCH")
	router.HandleFunc("/api/v1/passenger", allPassenger) //display all passengers

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func findPassenger(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if val, ok := isPassengerExist(params["username"]); ok {
		if r.Method == "GET" { //retrieve record of 1 passenger
			json.NewEncoder(w).Encode(val)
		}
	} else if r.Method == "POST" { //create new passenger account
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Passenger
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isPassengerExist(params["username"]); !ok {
					fmt.Println(data)
					insertPassenger(params["username"], data)
					w.WriteHeader(http.StatusAccepted)
				}
			} else {
				fmt.Println(err)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid Username")
	}
}

func allPassenger(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	var allPassengers, _ = getPassenger()
	//found := false
	results := map[string]Passenger{}

	if value := query.Get("q"); len(value) > 0 {
		for k, v := range allPassengers {
			if strings.Contains(strings.ToLower(v.Username), strings.ToLower(value)) {
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

func insertPassenger(username string, passenger Passenger) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO passenger VALUES (?, ?, ?, ?, ?)", passenger.Username, passenger.FirstName, passenger.LastName, passenger.PhoneNo, passenger.Email)
	if err != nil {
		panic(err.Error())
	}
}

func getPassenger() (map[string]Passenger, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT * from Passenger")
	if err != nil {
		panic(err.Error())
	}

	var allPassengers map[string]Passenger = map[string]Passenger{}

	for results.Next() {
		var p Passenger
		var username string
		err = results.Scan(&p.Username, &p.FirstName, &p.LastName, &p.PhoneNo, &p.Email)
		if err != nil {
			panic(err.Error())
		}
		allPassengers[username] = p
	}
	return allPassengers, true
}

func isPassengerExist(username string) (Passenger, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var p Passenger
	result := db.QueryRow("SELECT * from Passenger WHERE Username=?", username)
	err = result.Scan(&p.Username, &p.FirstName, &p.LastName, &p.PhoneNo, &p.Email)
	if err == sql.ErrNoRows {
		return p, false
	}
	return p, true
}
