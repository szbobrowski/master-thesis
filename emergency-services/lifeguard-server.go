package main

import (
	"context"
	"database/sql"
	"fmt"
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
	id, err := CreateLifeguard(s.db, req.Name, req.Login, req.PasswordHash, int(req.YearsOfExperience), req.Specialization, req.OnMission)
	if err != nil {
		return nil, fmt.Errorf("failed to create lifeguard: %w", err)
	}
	return &CreateLifeguardResponse{Id: id}, nil
}

// GetLifeguard implements LifeguardServiceServer.GetLifeguard.
func (s *server) GetLifeguard(ctx context.Context, req *GetLifeguardRequest) (*GetLifeguardResponse, error) {
	lifeguard, err := GetLifeguardByID(s.db, int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to get lifeguard: %w", err)
	}

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
	err := UpdateLifeguard(s.db, int(req.Id), req.Name, req.Login, req.PasswordHash, int(req.YearsOfExperience), req.Specialization, req.OnMission)
	if err != nil {
		return nil, fmt.Errorf("failed to update lifeguard: %w", err)
	}

	return &UpdateLifeguardResponse{Success: true}, nil
}

// DeleteLifeguard implements LifeguardServiceServer.DeleteLifeguard.
func (s *server) DeleteLifeguard(ctx context.Context, req *DeleteLifeguardRequest) (*DeleteLifeguardResponse, error) {
	err := DeleteLifeguard(s.db, int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to delete lifeguard: %w", err)
	}

	return &DeleteLifeguardResponse{Success: true}, nil
}
