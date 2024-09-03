package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type IncidentServer struct {
	UnimplementedIncidentServiceServer
	dbClient *dynamodb.Client
}

func NewIncidentServer(client *dynamodb.Client) *IncidentServer {
	return &IncidentServer{dbClient: client}
}

func (s *IncidentServer) CreateIncident(ctx context.Context, req *CreateIncidentRequest) (*IncidentResponse, error) {
	incidentID := fmt.Sprintf("INC%d", time.Now().UnixNano())
	incident := Incident{
		IncidentID:   incidentID,
		Title:        req.Title,
		Description:  req.Description,
		Status:       req.Status,
		CreationDate: req.CreationDate,
	}

	err := createIncident(s.dbClient, incident)
	if err != nil {
		return nil, err
	}

	return &IncidentResponse{Incident: &IncidentProto{
		IncidentID:   incident.IncidentID,
		Title:        incident.Title,
		Description:  incident.Description,
		Status:       incident.Status,
		CreationDate: incident.CreationDate,
	}}, nil
}

func (s *IncidentServer) GetIncident(ctx context.Context, req *GetIncidentRequest) (*IncidentResponse, error) {
	incident, err := getIncident(s.dbClient, req.IncidentID)
	if err != nil {
		return nil, err
	}

	return &IncidentResponse{Incident: &IncidentProto{
		IncidentID:   incident.IncidentID,
		Title:        incident.Title,
		Description:  incident.Description,
		Status:       incident.Status,
		CreationDate: incident.CreationDate,
	}}, nil
}

func (s *IncidentServer) UpdateIncident(ctx context.Context, req *UpdateIncidentRequest) (*IncidentResponse, error) {
	err := updateIncident(s.dbClient, req.IncidentID, req.Status)
	if err != nil {
		return nil, err
	}

	updatedIncident, err := getIncident(s.dbClient, req.IncidentID)
	if err != nil {
		return nil, err
	}

	return &IncidentResponse{Incident: &IncidentProto{
		IncidentID:   updatedIncident.IncidentID,
		Title:        updatedIncident.Title,
		Description:  updatedIncident.Description,
		Status:       updatedIncident.Status,
		CreationDate: updatedIncident.CreationDate,
	}}, nil
}

func (s *IncidentServer) DeleteIncident(ctx context.Context, req *DeleteIncidentRequest) (*DeleteIncidentResponse, error) {
	err := deleteIncident(s.dbClient, req.IncidentID)
	if err != nil {
		return nil, err
	}

	return &DeleteIncidentResponse{Success: true}, nil
}
