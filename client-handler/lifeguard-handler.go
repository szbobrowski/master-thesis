package main

import (
	context "context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Lifeguard struct {
	Name              string `json:"name"`
	Login             string `json:"login"`
	PasswordHash      string `json:"password_hash"`
	YearsOfExperience int32  `json:"years_of_experience"`
	Specialization    string `json:"specialization"`
	OnMission         bool   `json:"on_mission"`
}

var lifeguardClient LifeguardServiceClient

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

	log.Printf("Utworzono nowy wiersz w tabeli lifeguards: %+v\n", lifeguard)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lifeguardResponse)
}

func GetLifeguardHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Niepoprawny format id podany przez użytkownika", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	lifeguardResponse, err := lifeguardClient.GetLifeguard(ctx, &GetLifeguardRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Pobrano wiersz z tabeli lifeguards, id wiersza: %d\n", id)
	json.NewEncoder(w).Encode(lifeguardResponse)
}

func UpdateLifeguardHandler(w http.ResponseWriter, r *http.Request) {
	var lifeguard Lifeguard
	err := json.NewDecoder(r.Body).Decode(&lifeguard)
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

	log.Printf("Zaktualizowano wiersz w tabeli lifeguards, id wiersza: %d, nowy wiersz: %+v\n", id, lifeguard)

	json.NewEncoder(w).Encode(lifeguardResponse)
}

func DeleteLifeguardHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Niepoprawny format id podany przez użytkownika", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = lifeguardClient.DeleteLifeguard(ctx, &DeleteLifeguardRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Usunięto wiersz z tabeli lifeguards, id wiersza: %d\n", id)

	w.WriteHeader(http.StatusNoContent)
}
