package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// NewLifeguardServer returns a new vehicle server.
func NewVehicleServer(db *sql.DB) *server {
	return &server{db: db}
}

// mustEmbedUnimplementedVehicleServiceServer implements VehicleServiceServer.
func (s *server) mustEmbedUnimplementedVehicleServiceServer() {
	panic("unimplemented")
}

// CreateVehicle implements VehicleServiceServer.CreateVehicle.
func (s *server) CreateVehicle(ctx context.Context, req *CreateVehicleRequest) (*CreateVehicleResponse, error) {
	id, err := CreateVehicle(s.db, req.Type, req.Location, int(req.FuelLevelInLiters), req.OnMission, int(req.LifeguardInChargeId))
	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle: %w", err)
	}
	return &CreateVehicleResponse{Id: id}, nil
}

// GetVehicle implements VehicleServiceServer.GetVehicle.
func (s *server) GetVehicle(ctx context.Context, req *GetVehicleRequest) (*GetVehicleResponse, error) {
	vehicle, err := GetVehicleByID(s.db, int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	// Convert VehicleDTO to GetVehicleResponse
	return &GetVehicleResponse{
		Id:                  int64(vehicle.ID),
		Type:                vehicle.Type,
		Location:            vehicle.Location,
		FuelLevelInLiters:   int32(vehicle.FuelLevelInLiters),
		OnMission:           vehicle.OnMission,
		LifeguardInChargeId: int64(vehicle.LifeguardInChargeID),
		CreatedAt:           vehicle.CreatedAt.Format(time.RFC3339), // Convert time.Time to timestamp.Timestamp
	}, nil
}

// UpdateVehicle implements VehicleServiceServer.UpdateVehicle.
func (s *server) UpdateVehicle(ctx context.Context, req *UpdateVehicleRequest) (*UpdateVehicleResponse, error) {
	err := UpdateVehicle(s.db, int(req.Id), req.Type, req.Location, int(req.FuelLevelInLiters), req.OnMission, int(req.LifeguardInChargeId))
	if err != nil {
		return nil, fmt.Errorf("failed to update vehicle: %w", err)
	}

	return &UpdateVehicleResponse{Success: true}, nil
}

// DeleteVehicle implements VehicleServiceServer.DeleteVehicle.
func (s *server) DeleteVehicle(ctx context.Context, req *DeleteVehicleRequest) (*DeleteVehicleResponse, error) {
	err := DeleteVehicle(s.db, int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to delete vehicle: %w", err)
	}

	return &DeleteVehicleResponse{Success: true}, nil
}
