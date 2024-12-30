package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type server struct {
	db *sql.DB
}

func NewLifeguardServer(db *sql.DB) *server {
	return &server{db: db}
}

func (s *server) mustEmbedUnimplementedLifeguardServiceServer() {
	panic("nie zaimplementowano metody")
}

func (s *server) CreateLifeguard(ctx context.Context, req *CreateLifeguardRequest) (*CreateLifeguardResponse, error) {
	id, err := CreateLifeguard(s.db, req.Name, req.Login, req.PasswordHash, int(req.YearsOfExperience), req.Specialization, req.OnMission)
	if err != nil {
		log.Printf("Nie udało się utworzyć wiersza w tabeli lifeguards: %v\n", err)
		return nil, fmt.Errorf("Nie udało się utworzyć wiersza w tabeli lifeguards: %w", err)
	}

	log.Printf("Utworzono wiersz w tabeli lifeguards, id wiersza: %d\n", id)

	return &CreateLifeguardResponse{Id: id}, nil
}

func (s *server) GetLifeguard(ctx context.Context, req *GetLifeguardRequest) (*GetLifeguardResponse, error) {
	lifeguard, err := GetLifeguardByID(s.db, int(req.Id))
	if err != nil {
		log.Printf("Nie udało się pobrać wiersza z tabeli lifeguards, id wiersza: %d, błąd: %v\n", req.Id, err)
		return nil, fmt.Errorf("Nie udało się pobrać wiersza z tabeli lifeguards: %w", err)
	}

	log.Printf("Pobrano wiersz z tabeli lifeguards: %+v\n", lifeguard)

	return &GetLifeguardResponse{
		Id:                int64(lifeguard.ID),
		Name:              lifeguard.Name,
		Login:             lifeguard.Login,
		PasswordHash:      lifeguard.PasswordHash,
		YearsOfExperience: int32(lifeguard.YearsOfExperience),
		Specialization:    lifeguard.Specialization,
		OnMission:         lifeguard.OnMission,
		CreatedAt:         lifeguard.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *server) UpdateLifeguard(ctx context.Context, req *UpdateLifeguardRequest) (*UpdateLifeguardResponse, error) {
	err := UpdateLifeguard(s.db, int(req.Id), req.Name, req.Login, req.PasswordHash, int(req.YearsOfExperience), req.Specialization, req.OnMission)
	if err != nil {
		log.Printf("Nie udało się zaktualizować wiersza w tabeli lifeguards, id wiersza: %d, błąd: %v\n", req.Id, err)
		return nil, fmt.Errorf("Nie udało się zaktualizować wiersza w tabeli lifeguards: %w", err)
	}

	log.Printf("Zaktualizowano wiersz w tabeli lifeguards, id wiersza: %d\n", req.Id)

	return &UpdateLifeguardResponse{Success: true}, nil
}

func (s *server) DeleteLifeguard(ctx context.Context, req *DeleteLifeguardRequest) (*DeleteLifeguardResponse, error) {
	err := DeleteLifeguard(s.db, int(req.Id))
	if err != nil {
		log.Printf("Nie udało się usunąć wiersza z tabeli lifeguards, id wiersza: %d, błąd: %v\n", req.Id, err)
		return nil, fmt.Errorf("Nie udało się usunąć wiersza z tabeli lifeguards: %w", err)
	}

	log.Printf("Usunięto wiersz w tabeli lifeguards, id wiersza: %d\n", req.Id)

	return &DeleteLifeguardResponse{Success: true}, nil
}
