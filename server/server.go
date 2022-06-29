package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const ContentTypeJSON = "application/json"
const ContentTypeURLENCODED = "application/x-www-form-urlencoded"

const MinSeats uint = 4
const MaxSeats uint = 6

const MinPeople uint = 1
const MaxPeople uint = 6

var carsMap map[uint]uint
var carsSize map[uint]uint
var groupsMap map[uint]uint
var capacitiesMap map[uint]map[uint]struct{}
var journeysMap map[uint]uint
var waitingGroups []uint

// Cars
type Car struct {
	Id    uint `json:"id"`
	Seats uint `json:"seats"`
}

func (car *Car) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Id    *uint `json:"id"`
		Seats *uint `json:"seats"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if required.Id == nil || required.Seats == nil {
		err = fmt.Errorf("id and seats are required for all cars")
		return err
	} else if *required.Id == 0 {
		err = fmt.Errorf("id must be different from 0")
		return err
	} else if *required.Seats < MinSeats {
		err = fmt.Errorf("seats must be > %d", MinSeats-1)
		return err
	} else if *required.Seats > MaxSeats {
		err = fmt.Errorf("seats must be < %d", MaxSeats+1)
		return err
	}

	car.Id = *required.Id
	car.Seats = *required.Seats
	return
}

// groups
type Group struct {
	Id     uint `json:"id"`
	People uint `json:"people"`
}

func (group Group) toJSON() string {
	return fmt.Sprintf("{ \"id\": %d, \"people\": %d }", group.Id, group.People)
}

func (group *Group) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Id     *uint `json:"id"`
		People *uint `json:"People"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return err
	} else if required.Id == nil || required.People == nil {
		err = fmt.Errorf("id and people are required for all groups")
		return err
	} else if *required.People < MinPeople {
		err = fmt.Errorf("number of people should be between %d or %d", MinPeople, MaxPeople)
		return err
		// To Do Add checking for unique ID
	} else if *required.People > MaxPeople {
		err = fmt.Errorf("number of people should be %d at most", MaxPeople)
		return err
		// To Do Add checking for unique ID
	} else {
		group.Id = *required.Id
		group.People = *required.People
	}
	return nil
}

// Service variables

func startStorage() {
	carsMap = make(map[uint]uint)
	carsSize = make(map[uint]uint)
	groupsMap = make(map[uint]uint)
	capacitiesMap = make(map[uint]map[uint]struct{})
	journeysMap = make(map[uint]uint)
	waitingGroups = []uint{}
}

func New(addr string) *http.Server {
	startStorage()
	initRoutes()
	return &http.Server{
		Addr: addr,
	}
}

//check if adding a variable to group to store car can save the journey Map
