package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// LifeguardRestDTO represents the lifeguard model used in the REST API.
type LifeguardRestDTO struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Login             string    `json:"login"`
	PasswordHash      string    `json:"password_hash"`
	YearsOfExperience int       `json:"years_of_experience"`
	Specialization    string    `json:"specialization"`
	OnMission         bool      `json:"on_mission"`
	CreatedAt         time.Time `json:"created_at"`
}

// Handler to create a new lifeguard.
func createLifeguardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var lifeguard LifeguardRestDTO
		if err := json.NewDecoder(r.Body).Decode(&lifeguard); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		id, err := CreateLifeguard(db, lifeguard.Name, lifeguard.Login, lifeguard.PasswordHash, lifeguard.YearsOfExperience, lifeguard.Specialization, lifeguard.OnMission)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create lifeguard: %v", err), http.StatusInternalServerError)
			return
		}

		lifeguard.ID = int(id)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(lifeguard)
	}
}

// Handler to retrieve a lifeguard by ID.
func getLifeguardByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get lifeguard ID from query parameters
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid lifeguard ID", http.StatusBadRequest)
			return
		}

		lifeguard, err := GetLifeguardByID(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Lifeguard not found: %v", err), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(lifeguard)
	}
}

// Handler to update an existing lifeguard.
func updateLifeguardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var lifeguard LifeguardRestDTO
		if err := json.NewDecoder(r.Body).Decode(&lifeguard); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid lifeguard ID", http.StatusBadRequest)
			return
		}

		// Update lifeguard in the database
		err = UpdateLifeguard(db, id, lifeguard.Name, lifeguard.Login, lifeguard.PasswordHash, lifeguard.YearsOfExperience, lifeguard.Specialization, lifeguard.OnMission)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update lifeguard: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Lifeguard with ID %d updated successfully!\n", id)
	}
}

// Handler to delete a lifeguard by ID.
func deleteLifeguardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get lifeguard ID from query parameters
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid lifeguard ID", http.StatusBadRequest)
			return
		}

		err = DeleteLifeguard(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete lifeguard: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Lifeguard with ID %d deleted successfully!\n", id)
	}
}
