package main

import (
	"database/sql"
	"log"
	"net/http"
)

// startRESTServer starts the REST server to manage both vehicles and lifeguards.
func startRESTServer(db *sql.DB) error {
	// Vehicles routes
	http.HandleFunc("/vehicle", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createVehicleHandler(db)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/vehicle/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getVehicleByIDHandler(db)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/vehicle/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			updateVehicleHandler(db)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/vehicle/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deleteVehicleHandler(db)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Lifeguards routes
	http.HandleFunc("/lifeguard", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createLifeguardHandler(db)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/lifeguard/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getLifeguardByIDHandler(db)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/lifeguard/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			updateLifeguardHandler(db)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/lifeguard/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deleteLifeguardHandler(db)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("REST server listening on port 8080")
	return http.ListenAndServe(":8080", nil)
}
