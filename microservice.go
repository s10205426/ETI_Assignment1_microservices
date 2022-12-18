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
	Password  string `json:"Password"`
}

type Driver struct {
	Username  string `json:"Username"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	PhoneNo   int    `json:"PhoneNo"`
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	IDNo      string `json:"IDNo"`
	LicenseNo string `json:"LicenseNo"`
	IsBusy    string `json:"IsBusy"`
}

type CarTrip struct {
	ID                int    `json:"ID"`
	PassengerUsername string `json:"PassengerUsername"`
	DriverUsername    string `json:"DriverUsername"`
	Pickup            string `json:"Pickup"`
	Dropoff           string `json:"Dropoff"`
	PickupTime        string `json:"PickupTime"`
	IsCompleted       string `json:"IsCompleted"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/passenger/{username}", findPassenger).Methods("GET", "POST", "PUT")
	router.HandleFunc("/api/v1/passenger", allPassenger) //display all passengers
	router.HandleFunc("/api/v1/driver/{username}", findDriver).Methods("GET", "POST", "PUT")
	router.HandleFunc("/api/v1/driver", allDriver) //display all drivers

	router.HandleFunc("/api/v1/cartrip/{passengerUsername}", findCarTrip).Methods("GET", "POST", "PUT")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func findCarTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "POST" { //create new car trip record
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data CarTrip
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isPassengerExist(params["passengerUsername"]); ok {
					fmt.Println(data)
					insertCarTrip(data)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "Username does not exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PUT" { //update car trip record
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data CarTrip
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isCarTripExist(params["passengerUsername"]); ok {
					fmt.Println(data)
					updateCarTrip(params["passengerUsername"], data)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "Username does not exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if _, ok := isCarTripExist(params["passengerUsername"]); ok { //retrieve all car trip records for a passenger
		var allCarTrips, _ = getCarTrip(params["passengerUsername"])
		cartripsWrapper := struct {
			CarTrips map[string]CarTrip `json:"CarTrips"`
		}{allCarTrips}
		json.NewEncoder(w).Encode(cartripsWrapper)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid Username")
	}
}

func findPassenger(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "POST" { //create new passenger account
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Passenger
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isPassengerExist(params["username"]); !ok {
					fmt.Println(data)
					insertPassenger(params["username"], data)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "Username already exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PUT" { //update passenger information
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Passenger
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isPassengerExist(params["username"]); ok {
					fmt.Println(data)
					updatePassenger(params["username"], data)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "Username does not exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if val, ok := isPassengerExist(params["username"]); ok {
		json.NewEncoder(w).Encode(val) //retrieve record of 1 passenger
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid Username")
	}
}

func findDriver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Method == "POST" { //create new driver account
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Driver
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isDriverExist(params["username"]); !ok {
					fmt.Println(data)
					insertDriver(params["username"], data)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "Username already exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PUT" { //update driver information
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var data Driver
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isDriverExist(params["username"]); ok {
					fmt.Println(data)
					updateDriver(params["username"], data)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "Username does not exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if val, ok := isDriverExist(params["username"]); ok {
		json.NewEncoder(w).Encode(val) //retrieve record of 1 driver
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

func allDriver(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	var allDrivers, _ = getDriver()
	//found := false
	results := map[string]Driver{}

	if value := query.Get("q"); len(value) > 0 {
		for k, v := range allDrivers {
			if strings.Contains(strings.ToLower(v.Username), strings.ToLower(value)) {
				results[k] = v
				//found = true
			}
		}
	} else {
		driversWrapper := struct {
			Drivers map[string]Driver `json:"Drivers"`
		}{allDrivers}
		json.NewEncoder(w).Encode(driversWrapper)
		return
	}
}

func insertCarTrip(cartrip CarTrip) { //add new car trip record into database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	result := db.QueryRow("SELECT Username FROM driver WHERE IsBusy = '0' ORDER BY RAND() LIMIT 1") //randomly selects a driver that is available
	var newDriver string
	err = result.Scan(&newDriver) //retrieve username of driver that was chosen
	if err != nil {
		panic(err.Error())
	}

	//add Car Trip record with empty driver username
	_, err = db.Exec("INSERT INTO cartrip (PassengerUsername, DriverUsername, Pickup, Dropoff, PickupTime, IsCompleted) VALUES (?, ?, ?, ?, ?, ?)", cartrip.PassengerUsername, cartrip.DriverUsername, cartrip.Pickup, cartrip.Dropoff, cartrip.PickupTime, cartrip.IsCompleted)
	if err != nil {
		panic(err.Error())
	}

	//update Car Trip record with the newly added driver
	_, err = db.Exec("UPDATE cartrip SET DriverUsername = ? WHERE DriverUsername = ''", newDriver)
	if err != nil {
		panic(err.Error())
	}
}

func insertPassenger(username string, passenger Passenger) { //add new passenger record into database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO passenger VALUES (?, ?, ?, ?, ?, ?)", passenger.Username, passenger.FirstName, passenger.LastName, passenger.PhoneNo, passenger.Email, passenger.Password)
	if err != nil {
		panic(err.Error())
	}
}

func insertDriver(username string, driver Driver) { //add new driver record into database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO driver VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", driver.Username, driver.FirstName, driver.LastName, driver.PhoneNo, driver.Email, driver.Password, driver.IDNo, driver.LicenseNo, driver.IsBusy)
	if err != nil {
		panic(err.Error())
	}
}

func updateCarTrip(username string, cartrip CarTrip) { //update an existing car trip record into database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("UPDATE cartrip SET IsCompleted = ? WHERE PassengerUsername = ?", "1", username)
	if err != nil {
		panic(err.Error())
	}
}

func updatePassenger(username string, passenger Passenger) { //update an existing passenger record into database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("UPDATE passenger SET FirstName = ?, LastName = ?, PhoneNo = ?, Email = ?, Password = ? WHERE Username = ?", passenger.FirstName, passenger.LastName, passenger.PhoneNo, passenger.Email, passenger.Password, passenger.Username)
	if err != nil {
		panic(err.Error())
	}
}

func updateDriver(username string, driver Driver) { //update an existing driver record into database
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("UPDATE driver SET FirstName = ?, LastName = ?, PhoneNo = ?, Email = ?, Password = ?, LicenseNo = ?, IsBusy = ? WHERE Username = ?", driver.FirstName, driver.LastName, driver.PhoneNo, driver.Email, driver.Password, driver.LicenseNo, driver.IsBusy, driver.Username)
	if err != nil {
		panic(err.Error())
	}
}

func getCarTrip(username string) (map[string]CarTrip, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//retrieves all COMPLETED car trips in reverse chronological order
	results, err := db.Query("SELECT * from CarTrip WHERE PassengerUsername=? AND IsCompleted=? ORDER BY PickupTime DESC", username, "1")
	if err != nil {
		panic(err.Error())
	}

	var allCarTrips map[string]CarTrip = map[string]CarTrip{}

	for results.Next() { //group all car trips into a map as there might be more than 1 record
		var ct CarTrip
		var ID string
		err = results.Scan(&ID, &ct.PassengerUsername, &ct.DriverUsername, &ct.Pickup, &ct.Dropoff, &ct.PickupTime, &ct.IsCompleted)
		if err != nil {
			panic(err.Error())
		}
		allCarTrips[ID] = ct
	}
	return allCarTrips, true
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
		err = results.Scan(&username, &p.FirstName, &p.LastName, &p.PhoneNo, &p.Email, &p.Password)
		if err != nil {
			panic(err.Error())
		}
		allPassengers[username] = p
	}
	return allPassengers, true
}

func getDriver() (map[string]Driver, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT * from Driver")
	if err != nil {
		panic(err.Error())
	}

	var allDrivers map[string]Driver = map[string]Driver{}

	for results.Next() {
		var d Driver
		var username string
		err = results.Scan(&username, &d.FirstName, &d.LastName, &d.PhoneNo, &d.Email, &d.Password, &d.IDNo, &d.LicenseNo, &d.IsBusy)
		if err != nil {
			panic(err.Error())
		}
		allDrivers[username] = d
	}
	return allDrivers, true
}

func isCarTripExist(username string) (CarTrip, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var p Passenger
	var ct CarTrip

	//extra statement to make sure username entered is in the database
	result := db.QueryRow("SELECT * from Passenger WHERE Username=?", username)
	err = result.Scan(&p.Username, &p.FirstName, &p.LastName, &p.PhoneNo, &p.Email, &p.Password)

	//retrieves all car trips for username entered
	result2 := db.QueryRow("SELECT * from CarTrip WHERE PassengerUsername=?", username)
	err2 := result2.Scan(&ct.ID, &ct.PassengerUsername, &ct.DriverUsername, &ct.Pickup, &ct.Dropoff, &ct.PickupTime, &ct.IsCompleted)
	if err == sql.ErrNoRows || err2 == sql.ErrNoRows {
		return ct, false
	}
	return ct, true
}

func isPassengerExist(username string) (Passenger, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var p Passenger
	result := db.QueryRow("SELECT * from Passenger WHERE Username=?", username)
	err = result.Scan(&p.Username, &p.FirstName, &p.LastName, &p.PhoneNo, &p.Email, &p.Password)
	if err == sql.ErrNoRows {
		return p, false
	}
	return p, true
}

func isDriverExist(username string) (Driver, bool) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3307)/my_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var d Driver
	result := db.QueryRow("SELECT * from Driver WHERE Username=?", username)
	err = result.Scan(&d.Username, &d.FirstName, &d.LastName, &d.PhoneNo, &d.Email, &d.Password, &d.IDNo, &d.LicenseNo, &d.IsBusy)
	if err == sql.ErrNoRows {
		return d, false
	}
	return d, true
}
