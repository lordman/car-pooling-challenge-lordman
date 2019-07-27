package main

import (
	"encoding/json"
	"log"
	"net/http"

	// Package gorilla/mux implements a request router and dispatcher for matching incoming requests to their respective handler
	"github.com/gorilla/mux"
	// Package validator implements value validations for structs and individual fields based on tags
	"gopkg.in/go-playground/validator.v9"
)

// Car is the model for the cars
type Car struct {
	ID int `json:"id" validate:"required"`
	// Cars have a different amount of seats available,
	// they can accommodate groups of up to 4, 5 or 6 people
	Seats int8 `json:"seats" validate:"required,min=4,max=6"`
}

// Journey is the model for the journeys
type Journey struct {
	ID     string `json:"id" validate:"required"`
	People int8   `json:"people" validate:"required,min=1,max=6"`
}

// Since we won’t be using a database for now, we are initiating our vars as slices
var cars []Car
var journeys []Journey

/*
GET /status
Indicate the service has started up correctly and is ready to accept requests.
Responses:
200 OK When the service is ready to receive requests.
*/
func getStatus(w http.ResponseWriter, router *http.Request) {
	w.WriteHeader(http.StatusOK)
}

/*
PUT /cars
Load the list of available cars in the service and remove all previous data
(existing journeys and cars). This method may be called more than once during
the life cycle of the service.
Body required The list of cars to load.
Content Type application/json
Responses:
200 OK When the list is registered correctly.
400 Bad Request When there is a failure in the request format, expected
headers, or the payload can't be unmarshaled.
*/
func loadAvailableCars(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content Type", "application/json")
	// Remove all previous data (existing journeys and cars)
	cars = nil
	journeys = nil
	// Decode JSON
	decoded := json.NewDecoder(router.Body).Decode(&cars)
	// We need to check if the data provided match the requirements
	validate := validator.New()
	var valid error
	for _, cars := range cars {
		valid = validate.Struct(cars)
		if valid != nil {
			break // At least one car has invalid data
		}
	}

	if decoded != nil || valid != nil {
		w.WriteHeader(http.StatusBadRequest)
		cars = nil
		log.Println("There is a failure in the request format, expected headers, or the payload can't be unmarshaled")

	} else {
		w.WriteHeader(http.StatusOK)
		log.Printf("%+v", cars)
		log.Println("The list is registered correctly")
	}
}

/*
POST /journey
A group of people requests to perform a journey.
Body required The group of people that wants to perform the journey
Content Type application/json
Responses:
200 OK or 202 Accepted When the group is registered correctly
400 Bad Request When there is a failure in the request format or the
payload can't be unmarshaled.
*/
func requestJourney(w http.ResponseWriter, router *http.Request) {
	//TODO
}

/*
POST /dropoff
A group of people requests to be dropped off. Wether they traveled or not.
Body required A form with the group ID, such that ID=X
Content Type application/x-www-form-urlencoded
Responses:
200 OK or 204 No Content When the group is unregistered correctly.
404 Not Found When the group is not to be found.
400 Bad Request When there is a failure in the request format or the
payload can't be unmarshaled.
*/
func dropOff(w http.ResponseWriter, router *http.Request) {
	//TODO
}

/*
POST /locate
Given a group ID such that ID=X, return the car the group is traveling
with, or no car if they are still waiting to be served.
Body required A url encoded form with the group ID such that ID=X
Content Type application/x-www-form-urlencoded
Accept application/json
Responses:
200 OK With the car as the payload when the group is assigned to a car.
204 No Content When the group is waiting to be assigned to a car.
404 Not Found When the group is not to be found.
400 Bad Request When there is a failure in the request format or the
payload can't be unmarshaled.
*/
func locateCar(w http.ResponseWriter, router *http.Request) {
	//TODO
}

func main() {
	// Initialize router
	router := mux.NewRouter()

	// Endpoints
	router.HandleFunc("/status", getStatus).Methods("GET")
	router.HandleFunc("/cars", loadAvailableCars).Methods("PUT")
	router.HandleFunc("/journey/", requestJourney).Methods("POST")
	router.HandleFunc("/dropoff/", dropOff).Methods("POST")
	router.HandleFunc("/locate/", locateCar).Methods("POST")

	log.Fatal(http.ListenAndServe(":9091", router))

}
