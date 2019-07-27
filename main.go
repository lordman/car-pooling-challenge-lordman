package main

import (
	"net/http"
	// Package validator implements value validations for structs and individual fields based on tags

	// Package gorilla/mux implements a request router and dispatcher for matching incoming requests to their respective handler
	"github.com/gorilla/mux"
)

// Car is the model for the cars
type Car struct {
	ID    string `json:"id" validate:"required"`
	Seats int8   `json:"seats" validate:"required,min=4,max=6"`
}

// Journey is the model for the journeys
type Journey struct {
	ID     string `json:"id" validate:"required"`
	People int8   `json:"people" validate:"required,min=1,max=6"`
}

/*
type Group struct {
	ID		string	`json:"id"`
	Seats	int8	`json:"people"`
}
*/

// Since we won’t be using a database for now, we are initiating our vars as slices
var cars []Car
var journeys []Journey

// var Group []Group

/*
GET /status
Indicate the service has started up correctly and is ready to accept requests.
Responses:
200 OK When the service is ready to receive requests.
*/
func getStatus(w http.ResponseWriter, r *http.Request) {
	//TODO
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
func loadAvailableCars(w http.ResponseWriter, r *http.Request) {
	//TODO
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
func requestJourney(w http.ResponseWriter, r *http.Request) {
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
func dropOff(w http.ResponseWriter, r *http.Request) {
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
func locateCar(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func main() {
	// Initialize router
	r := mux.NewRouter()

	// Endpoints
	r.HandleFunc("/status/", getStatus).Methods("GET")
	r.HandleFunc("/cars/", loadAvailableCars).Methods("PUT")
	r.HandleFunc("/journey/", requestJourney).Methods("POST")
	r.HandleFunc("/dropoff/", dropOff).Methods("POST")
	r.HandleFunc("/locate/", locateCar).Methods("POST")

}
