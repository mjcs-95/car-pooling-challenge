package server

import (
	"net/http"
)

func initRoutes() {
	// Done

	// Performance test and improves required
	http.HandleFunc("/status", statusHandler)

	http.HandleFunc("/cars", carsHandler)

	http.HandleFunc("/journey", journeyHandler)

	http.HandleFunc("/locate", locateHandler)

	http.HandleFunc("/dropoff", dropoffHandler)
}
