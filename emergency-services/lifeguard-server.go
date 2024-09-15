package main

import (
	"context"
	"database/sql"
	"fmt"
	"log" // Added log package
	"time"
)

// server is used to implement VehicleServiceServer.
type server struct {
	db *sql.DB
}

// NewLifeguardServer returns a new vehicle server.
func NewLifeguardServer(db *sql.DB) *server {
	return &server{db: db}
}

// mustEmbedUnimplementedLifeguardServiceServer implements LifeguardServiceServer.
func (s *server) mustEmbedUnimplementedLifeguardServiceServer() {
	panic("unimplemented")
}

// CreateLifeguard implements LifeguardServiceServer.CreateLifeguard.
func (s *server) CreateLifeguard(ctx context.Context, req *CreateLifeguardRequest) (*CreateLifeguardResponse, error) {
	// Log the start of lifeguard creation
	log.Printf("Creating lifeguard with name: %s, login: %s\n", req.Name, req.Login)

	id, err := CreateLifeguard(s.db, req.Name, req.Login, req.PasswordHash, int(req.YearsOfExperience), req.Specialization, req.OnMission)
	if err != nil {
		log.Printf("Failed to create lifeguard: %v\n", err) // Log error
		return nil, fmt.Errorf("failed to create lifeguard: %w", err)
	}

	// Log success
	log.Printf("Lifeguard created successfully with ID: %d\n", id)

	return &CreateLifeguardResponse{Id: id}, nil
}

// GetLifeguard implements LifeguardServiceServer.GetLifeguard.
func (s *server) GetLifeguard(ctx context.Context, req *GetLifeguardRequest) (*GetLifeguardResponse, error) {
	// Log the start of lifeguard retrieval
	log.Printf("Retrieving lifeguard with ID: %d\n", req.Id)

	lifeguard, err := GetLifeguardByID(s.db, int(req.Id))
	if err != nil {
		log.Printf("Failed to get lifeguard with ID: %d, error: %v\n", req.Id, err) // Log error
		return nil, fmt.Errorf("failed to get lifeguard: %w", err)
	}

	// Log success
	log.Printf("Lifeguard retrieved successfully: %+v\n", lifeguard)

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

// UpdateLifeguard implements LifeguardServiceServer.UpdateLifeguard.
func (s *server) UpdateLifeguard(ctx context.Context, req *UpdateLifeguardRequest) (*UpdateLifeguardResponse, error) {
	// Log the start of lifeguard update
	log.Printf("Updating lifeguard with ID: %d\n", req.Id)

	err := UpdateLifeguard(s.db, int(req.Id), req.Name, req.Login, req.PasswordHash, int(req.YearsOfExperience), req.Specialization, req.OnMission)
	if err != nil {
		log.Printf("Failed to update lifeguard with ID: %d, error: %v\n", req.Id, err) // Log error
		return nil, fmt.Errorf("failed to update lifeguard: %w", err)
	}

	// Log success
	log.Printf("Lifeguard updated successfully with ID: %d\n", req.Id)

	return &UpdateLifeguardResponse{Success: true}, nil
}

// DeleteLifeguard implements LifeguardServiceServer.DeleteLifeguard.
func (s *server) DeleteLifeguard(ctx context.Context, req *DeleteLifeguardRequest) (*DeleteLifeguardResponse, error) {
	// Log the start of lifeguard deletion
	log.Printf("Deleting lifeguard with ID: %d\n", req.Id)

	err := DeleteLifeguard(s.db, int(req.Id))
	if err != nil {
		log.Printf("Failed to delete lifeguard with ID: %d, error: %v\n", req.Id, err) // Log error
		return nil, fmt.Errorf("failed to delete lifeguard: %w", err)
	}

	// Log success
	log.Printf("Lifeguard with ID: %d deleted successfully\n", req.Id)

	return &DeleteLifeguardResponse{Success: true}, nil
}
