package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func populateCarsList(w http.ResponseWriter, r *http.Request) error {
	var carsArr []Car = []Car{}
	err := json.NewDecoder(r.Body).Decode(&carsArr)
	if err != nil {
		return err
	}
	for _, car := range carsArr {
		// To Do Add checking for unique ID
		if _, ok := carsMap[car.Id]; ok {
			//do something here
			cleanJourneysAndCars()
			err = fmt.Errorf("cars Ids must be unique")
			return err
		}
		carsMap[car.Id] = car.Seats
		carsSize[car.Id] = car.Seats
		if capacitiesMap[car.Seats] == nil {
			capacitiesMap[car.Seats] = make(map[uint]struct{})
		}
		capacitiesMap[car.Seats][car.Id] = struct{}{}
	}
	return nil
}

func cleanJourneysAndCars() {
	for k := range carsMap {
		delete(carsMap, k)
		delete(carsSize, k)
	}
	for k := range capacitiesMap {
		delete(capacitiesMap, k)
	}
	for k := range groupsMap {
		delete(groupsMap, k)
	}
	for k := range journeysMap {
		delete(journeysMap, k)
	}
	waitingGroups = []uint{}
}

func searchValidCapacity(people uint) uint {
	var validCars = capacitiesMap[people]
	if len(validCars) != 0 {
		return people
	} else {
		for freeCap := range capacitiesMap {
			if freeCap > people && len(capacitiesMap[freeCap]) > 0 {
				return freeCap
			}
		}
	}
	return 0
}

func addNewGroup(group Group) int {
	// To Do Add checking for unique ID
	groupsMap[group.Id] = group.People
	availableCarSize := searchValidCapacity(group.People)
	if availableCarSize == 0 {
		journeysMap[group.Id] = 0
		waitingGroups = append(waitingGroups, group.Id)
		return http.StatusAccepted
	}
	for firstCarId := range capacitiesMap[availableCarSize] {
		assignCar(firstCarId, group)
		delete(capacitiesMap[availableCarSize], firstCarId)
		break
	}
	return http.StatusOK
}

func assignCar(chosenCarID uint, group Group) {
	newFreeCap := carsMap[chosenCarID] - group.People
	carsMap[chosenCarID] = carsMap[chosenCarID] - group.People
	if capacitiesMap[newFreeCap] == nil {
		capacitiesMap[newFreeCap] = make(map[uint]struct{})
	}
	capacitiesMap[newFreeCap][chosenCarID] = struct{}{}
	journeysMap[group.Id] = chosenCarID
}

func checkOrDeleteGroupWithoutCar(groupId uint, w http.ResponseWriter) bool {
	if _, exists := groupsMap[groupId]; !exists {
		w.WriteHeader(http.StatusNotFound)
		return true
	}
	if journeysMap[groupId] == 0 {
		delete(journeysMap, groupId)
		delete(groupsMap, groupId)
		found := false
		if len(waitingGroups) != 0 {
			if len(waitingGroups) == 1 {
				if waitingGroups[0] == groupId {
					waitingGroups = waitingGroups[:0]
				}
			} else {
				for idx := 0; idx < len(waitingGroups) && !found; idx++ {
					if waitingGroups[idx] == groupId {
						waitingGroups = append(waitingGroups[:idx], waitingGroups[idx+1:]...)
						found = true
					}
				}
			}
		}
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func removeGroup(groupId uint) (uint, uint) {
	// We won't update the waiting list, since it is a expensive operation untill
	// we assign a new group
	carId := journeysMap[groupId]
	delete(journeysMap, groupId)

	currCarSeats := carsMap[carId]
	delete(capacitiesMap[currCarSeats], carId)

	carsMap[carId] = carsMap[carId] + groupsMap[groupId]

	newFreeSeats := carsMap[carId]
	capacitiesMap[newFreeSeats][carId] = struct{}{}

	delete(groupsMap, groupId)
	return carId, newFreeSeats
}

func tryAssignWaitingGroupsToCar(carId uint, newFreeSeats uint) {
	for idx := 0; idx < len(waitingGroups) && newFreeSeats > 0; idx++ {
		groupId := waitingGroups[idx]
		//check for dropped groups
		if groupsMap[groupId] <= newFreeSeats {
			assignCar(carId, Group{groupId, groupsMap[groupId]})
			waitingGroups = append(waitingGroups[:idx], waitingGroups[idx+1:]...)
			newFreeSeats -= groupsMap[groupId]
		}
	}
}

/*
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}*/
