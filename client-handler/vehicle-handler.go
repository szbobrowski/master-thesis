package main

import (
	"context"
	"encoding/json"
	"log" // Added log package
	"net/http"
	"strconv"
	"time"
)

// Vehicle represents the structure used in the REST API.
type Vehicle struct {
	Type                string `json:"type"`
	Location            string `json:"location"`
	FuelLevelInLiters   int32  `json:"fuel_level_in_liters"`
	OnMission           bool   `json:"on_mission"`
	LifeguardInChargeId int64  `json:"lifeguard_in_charge_id"`
}

// Global gRPC client variable
var vehicleClient VehicleServiceClient

// CreateVehicleHandler handles the creation of a new vehicle.
func CreateVehicleHandler(w http.ResponseWriter, r *http.Request) {
	var vehicle Vehicle
	err := json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	vehicleResponse, err := vehicleClient.CreateVehicle(ctx, &CreateVehicleRequest{
		Type:                vehicle.Type,
		Location:            vehicle.Location,
		FuelLevelInLiters:   vehicle.FuelLevelInLiters,
		OnMission:           vehicle.OnMission,
		LifeguardInChargeId: vehicle.LifeguardInChargeId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the event
	log.Printf("Vehicle created: %+v\n", vehicle)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vehicleResponse)
}

// GetVehicleHandler handles retrieving a vehicle by ID.
func GetVehicleHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	// Convert the ID from string to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	vehicleResponse, err := vehicleClient.GetVehicle(ctx, &GetVehicleRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the event
	log.Printf("Vehicle retrieved with ID: %d\n", id)

	json.NewEncoder(w).Encode(vehicleResponse)
}

// UpdateVehicleHandler handles updating a vehicle.
func UpdateVehicleHandler(w http.ResponseWriter, r *http.Request) {
	var vehicle Vehicle
	err := json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idStr := r.URL.Query().Get("id")

	// Convert the ID from string to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	vehicleResponse, err := vehicleClient.UpdateVehicle(ctx, &UpdateVehicleRequest{
		Id:                  id,
		Type:                vehicle.Type,
		Location:            vehicle.Location,
		FuelLevelInLiters:   vehicle.FuelLevelInLiters,
		OnMission:           vehicle.OnMission,
		LifeguardInChargeId: vehicle.LifeguardInChargeId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the event
	log.Printf("Vehicle updated with ID: %d, New Data: %+v\n", id, vehicle)

	json.NewEncoder(w).Encode(vehicleResponse)
}

// DeleteVehicleHandler handles deleting a vehicle.
func DeleteVehicleHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	// Convert the ID from string to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = vehicleClient.DeleteVehicle(ctx, &DeleteVehicleRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the event
	log.Printf("Vehicle deleted with ID: %d\n", id)

	w.WriteHeader(http.StatusNoContent)
}
