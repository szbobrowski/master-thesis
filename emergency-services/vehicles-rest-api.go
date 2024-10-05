package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// VehicleRestDTO represents the vehicle model used in the REST API.
type VehicleRestDTO struct {
	ID                  int    `json:"id"`
	Type                string `json:"type"`
	Location            string `json:"location"`
	FuelLevelInLiters   int    `json:"fuel_level_in_liters"`
	OnMission           bool   `json:"on_mission"`
	LifeguardInChargeID int    `json:"lifeguard_in_charge_id"`
}

// Handler to create a new vehicle.
func createVehicleHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vehicle VehicleRestDTO
		if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		id, err := CreateVehicle(db, vehicle.Type, vehicle.Location, vehicle.FuelLevelInLiters, vehicle.OnMission, vehicle.LifeguardInChargeID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create vehicle: %v", err), http.StatusInternalServerError)
			return
		}

		vehicle.ID = int(id)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(vehicle)
	}
}

// Handler to retrieve a vehicle by ID.
func getVehicleByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get vehicle ID from query parameters.
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
			return
		}

		vehicle, err := GetVehicleByID(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Vehicle not found: %v", err), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(vehicle)
	}
}

// Handler to update an existing vehicle.
func updateVehicleHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var vehicle VehicleRestDTO
		if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
			return
		}

		// Update vehicle in the database
		err = UpdateVehicle(db, id, vehicle.Type, vehicle.Location, vehicle.FuelLevelInLiters, vehicle.OnMission, vehicle.LifeguardInChargeID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update vehicle: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Vehicle with ID %d updated successfully!\n", id)
	}
}

// Handler to delete a vehicle by ID.
func deleteVehicleHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get vehicle ID from query parameters
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
			return
		}

		err = DeleteVehicle(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete vehicle: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Vehicle with ID %d deleted successfully!\n", id)
	}
}
