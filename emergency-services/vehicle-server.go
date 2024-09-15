package main

import (
	"context"
	"database/sql"
	"fmt"
	"log" // Added log package
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
	// Log the start of vehicle creation
	log.Printf("Creating vehicle with type: %s, location: %s\n", req.Type, req.Location)

	id, err := CreateVehicle(s.db, req.Type, req.Location, int(req.FuelLevelInLiters), req.OnMission, int(req.LifeguardInChargeId))
	if err != nil {
		log.Printf("Failed to create vehicle: %v\n", err) // Log error
		return nil, fmt.Errorf("failed to create vehicle: %w", err)
	}

	// Log success
	log.Printf("Vehicle created successfully with ID: %d\n", id)

	return &CreateVehicleResponse{Id: id}, nil
}

// GetVehicle implements VehicleServiceServer.GetVehicle.
func (s *server) GetVehicle(ctx context.Context, req *GetVehicleRequest) (*GetVehicleResponse, error) {
	// Log the start of vehicle retrieval
	log.Printf("Retrieving vehicle with ID: %d\n", req.Id)

	vehicle, err := GetVehicleByID(s.db, int(req.Id))
	if err != nil {
		log.Printf("Failed to get vehicle with ID: %d, error: %v\n", req.Id, err) // Log error
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	// Log success
	log.Printf("Vehicle retrieved successfully: %+v\n", vehicle)

	// Convert VehicleDTO to GetVehicleResponse
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

// UpdateVehicle implements VehicleServiceServer.UpdateVehicle.
func (s *server) UpdateVehicle(ctx context.Context, req *UpdateVehicleRequest) (*UpdateVehicleResponse, error) {
	// Log the start of vehicle update
	log.Printf("Updating vehicle with ID: %d\n", req.Id)

	err := UpdateVehicle(s.db, int(req.Id), req.Type, req.Location, int(req.FuelLevelInLiters), req.OnMission, int(req.LifeguardInChargeId))
	if err != nil {
		log.Printf("Failed to update vehicle with ID: %d, error: %v\n", req.Id, err) // Log error
		return nil, fmt.Errorf("failed to update vehicle: %w", err)
	}

	// Log success
	log.Printf("Vehicle updated successfully with ID: %d\n", req.Id)

	return &UpdateVehicleResponse{Success: true}, nil
}

// DeleteVehicle implements VehicleServiceServer.DeleteVehicle.
func (s *server) DeleteVehicle(ctx context.Context, req *DeleteVehicleRequest) (*DeleteVehicleResponse, error) {
	// Log the start of vehicle deletion
	log.Printf("Deleting vehicle with ID: %d\n", req.Id)

	err := DeleteVehicle(s.db, int(req.Id))
	if err != nil {
		log.Printf("Failed to delete vehicle with ID: %d, error: %v\n", req.Id, err) // Log error
		return nil, fmt.Errorf("failed to delete vehicle: %w", err)
	}

	// Log success
	log.Printf("Vehicle deleted successfully with ID: %d\n", req.Id)

	return &DeleteVehicleResponse{Success: true}, nil
}
