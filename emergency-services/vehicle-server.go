package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func NewVehicleServer(db *sql.DB) *server {
	return &server{db: db}
}

func (s *server) mustEmbedUnimplementedVehicleServiceServer() {
	panic("nie zaimplementowano metody")
}

func (s *server) CreateVehicle(ctx context.Context, req *CreateVehicleRequest) (*CreateVehicleResponse, error) {
	id, err := CreateVehicle(s.db, req.Type, req.Location, int(req.FuelLevelInLiters), req.OnMission, int(req.LifeguardInChargeId))
	if err != nil {
		log.Printf("Nie udało się utworzyć wiersza w tabeli vehicles: %v\n", err)
		return nil, fmt.Errorf("Nie udało się utworzyć wiersza w tabeli vehicles: %w", err)
	}

	log.Printf("Utworzono wiersz w tabeli vehicles, id wiersza: %d\n", id)

	return &CreateVehicleResponse{Id: id}, nil
}

func (s *server) GetVehicle(ctx context.Context, req *GetVehicleRequest) (*GetVehicleResponse, error) {
	vehicle, err := GetVehicleByID(s.db, int(req.Id))
	if err != nil {
		log.Printf("Nie udało się pobrać wiersza z tabeli vehicles, id wiersza: %d, error: %v\n", req.Id, err)
		return nil, fmt.Errorf("Nie udało się pobrać wiersza z tabeli vehicles: %w", err)
	}

	log.Printf("Pobrano wiersz w tabeli vehicles: %+v\n", vehicle)

	return &GetVehicleResponse{
		Id:                  int64(vehicle.ID),
		Type:                vehicle.Type,
		Location:            vehicle.Location,
		FuelLevelInLiters:   int32(vehicle.FuelLevelInLiters),
		OnMission:           vehicle.OnMission,
		LifeguardInChargeId: int64(vehicle.LifeguardInChargeID),
		CreatedAt:           vehicle.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *server) UpdateVehicle(ctx context.Context, req *UpdateVehicleRequest) (*UpdateVehicleResponse, error) {
	err := UpdateVehicle(s.db, int(req.Id), req.Type, req.Location, int(req.FuelLevelInLiters), req.OnMission, int(req.LifeguardInChargeId))
	if err != nil {
		log.Printf("Nie udało się zaktualizować wiersza w tabeli vehicles, id wiersza: %d, błąd: %v\n", req.Id, err)
		return nil, fmt.Errorf("Nie udało się zaktualizować wiersza w tabeli vehicles: %w", err)
	}

	log.Printf("Zaktualizowano wiersz w tabeli vehicles, id wiersza: %d\n", req.Id)

	return &UpdateVehicleResponse{Success: true}, nil
}

func (s *server) DeleteVehicle(ctx context.Context, req *DeleteVehicleRequest) (*DeleteVehicleResponse, error) {
	err := DeleteVehicle(s.db, int(req.Id))
	if err != nil {
		log.Printf("Nie udało się usunąć wiersza z tabeli vehicles, id wiersza: %d, error: %v\n", req.Id, err)
		return nil, fmt.Errorf("Nie udało się usunąć wiersza z tabeli vehicles: %w", err)
	}

	log.Printf("Usunięto wiersz z tabeli vehicles, id wiersza: %d\n", req.Id)

	return &DeleteVehicleResponse{Success: true}, nil
}
