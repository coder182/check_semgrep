package controller

import (
	"DronesData/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DronesRouter(router *mux.Router) {
	router.HandleFunc("/api/sensors", insertDronesData).Methods("POST")
	router.HandleFunc("/api/AllFlights", insertAllFlightsData).Methods("POST")
	router.HandleFunc("/api/FlightLocation/{serial_id}", insertFlightData).Methods("POST")
	router.HandleFunc("/api/trackingId", GetTrackingIDs).Methods("POST")
}

// this function will call the sensors API and will insert the data into DB
func insertDronesData(w http.ResponseWriter, r *http.Request) {
	// Call the service to fetch and insert drone data
	if err := services.FetchAndStoreSensors(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Drone data inserted successfully",
	})
}

// this function will call the paginated_flights api( having all types of drones data) and will insert into DB
func insertAllFlightsData(w http.ResponseWriter, r *http.Request) {

	if err := services.StoreAllFlightsData(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Flights data inserted successfully",
	})

}

// this function will call the paginated_flights api( having all types of drones data) and will insert into DB
func insertFlightData(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	serialID := vars["serial_id"]
	if serialID == "" {
		http.Error(w, "serial_id is required in the URL path", http.StatusBadRequest)
		return
	}

	if err := services.StoreFlightLocation(serialID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Flights data inserted successfully",
	})
}

// GET all Tracking Ids
func GetTrackingIDs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	limitStr := vars["limit"]
	var isFristRun bool

	if limitStr == "" {
		http.Error(w, "limit is required in the URL path", http.StatusBadRequest)
		return
	}

	// Convert string to int
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
		return
	}

	// Get IDs from database
	ids, err := services.GetLastTrackingIDs(isFristRun)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return as simple JSON array
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ids)

}
