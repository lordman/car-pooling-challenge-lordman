package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"

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
	Seats     int `json:"seats" validate:"required,min=4,max=6"`
	Available int // Seats available
}

// Journey is the model for the journeys
type Journey struct {
	ID int `json:"id" validate:"required"`
	// People requests cars in groups of 1 to 6
	// People in the same group want to ride on the same car
	People int `json:"people" validate:"required,min=1,max=6"`
	CarID  int // ID of the car assigned (default 0)
}

// Since we won’t be using a database for now, we are initiating our vars as slices
var cars []Car
var journeys []Journey

// This function validate if all the cars match the requirements
func validateCars(c []Car) error {
	for _, c := range c {
		validate := validator.New()
		err := validate.Struct(c)
		if err != nil {
			log.Printf("%+v", err)
			return err
		}
	}
	return nil
}

// This function validate if the journey matches the requirements
func validateJourney(j Journey) error {
	validate := validator.New()
	err := validate.Struct(j)
	if err != nil {
		log.Printf("%+v", err)
	}
	return err
}

/*
This function assign unassigned groups to available cars.
It would be triggered every time a new group request a journey
of if a group request a drop off (a group could be waiting
and a car could have seats available).
*/
func assignJourneysToCars() {
	// We are sorting the slice based on the total number of seats
	// in order to assign cars with less seats before
	// (and leaving free those with more seats for the bigger groups)
	sort.SliceStable(cars, func(i, j int) bool {
		return cars[i].Seats < cars[j].Seats
	})
	// We are iterating over the entire slice of journeys to find all of them that are unassigned
	for i, j := range journeys {
		if journeys[i].CarID == 0 {
			for k, c := range cars {
				if c.Available >= j.People {
					journeys[i].CarID = cars[k].ID
					cars[k].Available -= journeys[i].People
					log.Printf("Group %d has been assigned to car %d",
						journeys[i].ID, cars[k].ID)
					break
				} else {
					log.Printf("Group %d can not be assigned to car %d",
						journeys[i].ID, cars[k].ID)
				}
			}
		}
	}
	log.Printf("%+v", cars)
	log.Printf("%+v", journeys)
}

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
	for i := 0; i < len(cars); i++ {
		cars[i].Available = cars[i].Seats
	}
	// We need to check if the data provided match the requirements
	valid := validateCars(cars)

	if decoded != nil || valid != nil {
		w.WriteHeader(http.StatusBadRequest)
		cars = nil
		log.Println("There is a failure in the request format, expected headers, or the payload can't be unmarshaled")
	} else {
		w.WriteHeader(http.StatusOK)
		log.Printf("%+v", cars)
		log.Println("The list has been registered correctly")
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
	w.Header().Set("Content Type", "application/json")
	var newJourney Journey
	// Decode JSON
	decoded := json.NewDecoder(router.Body).Decode(&newJourney)
	// We need to check if the data provided match the requirements
	valid := validateJourney(newJourney)

	if decoded != nil || valid != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("There is a failure in the request format or the payload can't be unmarshaled")
	} else {
		w.WriteHeader(http.StatusAccepted)
		journeys = append(journeys, newJourney)
		assignJourneysToCars()
		log.Println("The group has been registered correctly")
	}
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
	w.Header().Set("Content Type", "application/x-www-form-urlencoded")
	parsed := router.ParseForm()
	if parsed != nil { // The form can not be parsed
		w.WriteHeader(http.StatusBadRequest)
		log.Println("There is a failure in the request format or the payload can't be unmarshaled")
	} else { // The form can be parsed
		idstring := router.FormValue("ID")
		if idstring != "" { // The key "ID" is present in the form
			idint, converted := strconv.Atoi(idstring)
			if converted != nil { // The id can not be converted from string to int
				w.WriteHeader(http.StatusBadRequest)
				log.Println("There is a failure in the request format or the payload can't be unmarshaled")
			} else { // The id has been converted from string to int
				// Find the journey
				k := sort.Search(len(journeys), func(k int) bool { return journeys[k].ID >= idint })
				log.Printf("%+v", journeys[k])
				if k < len(journeys) && journeys[k].ID == idint {
					// Find the car assiociated to the journey
					sort.SliceStable(cars, func(l, m int) bool {
						return cars[l].ID < cars[m].ID
					})
					n := sort.Search(len(cars), func(n int) bool { return cars[n].ID >= journeys[k].CarID })
					if n < len(cars) && cars[n].ID == journeys[k].CarID {
						log.Printf("%+v", cars[n])
						// Free seats in the car
						cars[n].Available += journeys[k].People
						// Remove the journey for the journeys slice
						journeys = append(journeys[:k], journeys[k+1:]...)
						// Once the journey has been removed and the car seats released
						// we should try to assign unassigned groups
						assignJourneysToCars()
						w.WriteHeader(http.StatusNoContent)
						log.Println("The group has been unregistered correctly.")
					} else {
						// This would mean the group has no car assigned
						// so we would not net to release the car or
						// try to assign unassigned groups
						journeys = append(journeys[:k], journeys[k+1:]...)
						w.WriteHeader(http.StatusNoContent)
						log.Println("The group has been unregistered correctly.")
						log.Printf("%+v", journeys)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
					log.Println("The group has not been found.")
				}
			}
		} else { // The key "ID" is not present in the form
			w.WriteHeader(http.StatusBadRequest)
			log.Println("There is a failure in the request format or the payload can't be unmarshaled")
		}
	}
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
	router.HandleFunc("/journey", requestJourney).Methods("POST")
	router.HandleFunc("/dropoff", dropOff).Methods("POST")
	router.HandleFunc("/locate", locateCar).Methods("POST")

	log.Fatal(http.ListenAndServe(":9091", router))

}
