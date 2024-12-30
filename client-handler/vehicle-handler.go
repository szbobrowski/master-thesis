package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Vehicle struct {
	Type                string `json:"type"`
	Location            string `json:"location"`
	FuelLevelInLiters   int32  `json:"fuel_level_in_liters"`
	OnMission           bool   `json:"on_mission"`
	LifeguardInChargeId int64  `json:"lifeguard_in_charge_id"`
}

var vehicleClient VehicleServiceClient

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

	log.Printf("Utworzono nowy wiersz w tabeli vehicles: %+v\n", vehicle)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vehicleResponse)
}

func GetVehicleHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Niepoprawny format id podany przez użytkownika", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	vehicleResponse, err := vehicleClient.GetVehicle(ctx, &GetVehicleRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Pobrano wiersz z tabeli vehicles, id wiersza: %d\n", id)

	json.NewEncoder(w).Encode(vehicleResponse)
}

func UpdateVehicleHandler(w http.ResponseWriter, r *http.Request) {
	var vehicle Vehicle
	err := json.NewDecoder(r.Body).Decode(&vehicle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Niepoprawny format id podany przez użytkownika", http.StatusBadRequest)
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

	log.Printf("Zaktualizowano wiersz w tabeli vehicles, id wiersza: %d, nowy wiersz: %+v\n", id, vehicle)

	json.NewEncoder(w).Encode(vehicleResponse)
}

func DeleteVehicleHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Niepoprawny format id podany przez użytkownika", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = vehicleClient.DeleteVehicle(ctx, &DeleteVehicleRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Usunięto wiersz z tabeli vehicles, id wiersza: %d\n", id)

	w.WriteHeader(http.StatusNoContent)
}
