package main

import (
	context "context"
	"encoding/json"
	"log" // Added log package
	"net/http"
	"strconv"
	"time"
)

// Lifeguard represents the structure used in the REST API.
type Lifeguard struct {
	Name              string `json:"name"`
	Login             string `json:"login"`
	PasswordHash      string `json:"password_hash"`
	YearsOfExperience int32  `json:"years_of_experience"`
	Specialization    string `json:"specialization"`
	OnMission         bool   `json:"on_mission"`
}

// Global gRPC client variable
var lifeguardClient LifeguardServiceClient

// CreateLifeguardHandler handles the creation of a new lifeguard.
func CreateLifeguardHandler(w http.ResponseWriter, r *http.Request) {
	var lifeguard Lifeguard
	err := json.NewDecoder(r.Body).Decode(&lifeguard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	lifeguardResponse, err := lifeguardClient.CreateLifeguard(ctx, &CreateLifeguardRequest{
		Name:              lifeguard.Name,
		Login:             lifeguard.Login,
		PasswordHash:      lifeguard.PasswordHash,
		YearsOfExperience: lifeguard.YearsOfExperience,
		Specialization:    lifeguard.Specialization,
		OnMission:         lifeguard.OnMission,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the event
	log.Printf("Lifeguard created: %+v\n", lifeguard)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lifeguardResponse)
}

// GetLifeguardHandler handles retrieving a lifeguard by ID.
func GetLifeguardHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	// Convert the ID from string to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	lifeguardResponse, err := lifeguardClient.GetLifeguard(ctx, &GetLifeguardRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the event
	log.Printf("Lifeguard retrieved with ID: %d\n", id)

	json.NewEncoder(w).Encode(lifeguardResponse)
}

// UpdateLifeguardHandler handles updating a lifeguard.
func UpdateLifeguardHandler(w http.ResponseWriter, r *http.Request) {
	var lifeguard Lifeguard
	err := json.NewDecoder(r.Body).Decode(&lifeguard)
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

	lifeguardResponse, err := lifeguardClient.UpdateLifeguard(ctx, &UpdateLifeguardRequest{
		Id:                id,
		Name:              lifeguard.Name,
		Login:             lifeguard.Login,
		PasswordHash:      lifeguard.PasswordHash,
		YearsOfExperience: lifeguard.YearsOfExperience,
		Specialization:    lifeguard.Specialization,
		OnMission:         lifeguard.OnMission,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the event
	log.Printf("Lifeguard updated with ID: %d, New Data: %+v\n", id, lifeguard)

	json.NewEncoder(w).Encode(lifeguardResponse)
}

// DeleteLifeguardHandler handles deleting a lifeguard.
func DeleteLifeguardHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	// Convert the ID from string to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = lifeguardClient.DeleteLifeguard(ctx, &DeleteLifeguardRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the event
	log.Printf("Lifeguard deleted with ID: %d\n", id)

	w.WriteHeader(http.StatusNoContent)
}
