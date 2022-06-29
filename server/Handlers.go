package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if !isSameMethod(w, r, "GET") {
		return
	}
	w.WriteHeader(http.StatusOK)
}

// /cars
func carsHandler(w http.ResponseWriter, r *http.Request) {
	if !isSameMethod(w, r, "PUT") || isBodyEmpty(w, r) || !isContentJson(w, r) {
		return
	}
	cleanJourneysAndCars()
	err := populateCarsList(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Input(JSON) format, %s", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	//PrintMemUsage()
}

// /journey

func journeyHandler(w http.ResponseWriter, r *http.Request) {
	if !isSameMethod(w, r, "POST") || isBodyEmpty(w, r) || !isContentJson(w, r) {
		return
	}
	group := Group{}
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Input(JSON) format, %s", err.Error())
		return
	}
	if _, ok := groupsMap[group.Id]; ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error, group Id already exists")
		return
	}
	status := addNewGroup(group)
	w.WriteHeader(status)
}

// /dropoff
func dropoffHandler(w http.ResponseWriter, r *http.Request) {
	if !urlEncReqHasValidSettings(w, r) {
		return
	}
	id, _ := strconv.Atoi(r.PostForm["ID"][0])
	groupId := uint(id)
	noChangeInJourneys := checkOrDeleteGroupWithoutCar(groupId, w)
	if noChangeInJourneys {
		return
	}
	carId, newFreeSeats := removeGroup(groupId)
	tryAssignWaitingGroupsToCar(carId, newFreeSeats)
	w.WriteHeader(http.StatusOK)
}

// /locate
func locateHandler(w http.ResponseWriter, r *http.Request) {
	if !urlEncReqHasValidSettings(w, r) {
		w.Header().Set("Content-Type", ContentTypeJSON)
		return
	}

	id, _ := strconv.Atoi(r.PostForm["ID"][0])
	groupId := uint(id)
	if _, exists := groupsMap[groupId]; !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", ContentTypeJSON)
		return
	}
	if journeysMap[groupId] == 0 {
		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", ContentTypeJSON)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", ContentTypeJSON)
	carId := journeysMap[groupId]
	fmt.Fprintf(w, "{ \"id\": %d, \"seats\": %d }", carId, carsSize[carId])
}
