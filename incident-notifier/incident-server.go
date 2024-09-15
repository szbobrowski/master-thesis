package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type IncidentServer struct {
	UnimplementedIncidentServiceServer
	dbClient   *dynamodb.Client
	sqsManager *SQSManager
}

func NewIncidentServer(client *dynamodb.Client, queue *SQSManager) *IncidentServer {
	return &IncidentServer{dbClient: client, sqsManager: queue}
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

	// Log the start of incident creation
	log.Printf("Creating incident: %+v\n", incident)

	err := createIncident(s.dbClient, incident)
	if err != nil {
		log.Printf("Failed to create incident: %v\n", err)
		return nil, err
	}

	// Log success of the incident creation in DynamoDB
	log.Printf("Incident created successfully: %s\n", incidentID)

	// Send message to SQS for incident creation
	s.sqsManager.SendMessage(incident, "CREATE")
	log.Printf("Incident message sent to SQS for incident creation: %s\n", incidentID)

	return &IncidentResponse{Incident: &IncidentProto{
		IncidentID:   incident.IncidentID,
		Title:        incident.Title,
		Description:  incident.Description,
		Status:       incident.Status,
		CreationDate: incident.CreationDate,
	}}, nil
}

func (s *IncidentServer) GetIncident(ctx context.Context, req *GetIncidentRequest) (*IncidentResponse, error) {
	// Log the start of incident retrieval
	log.Printf("Retrieving incident with ID: %s\n", req.IncidentID)

	incident, err := getIncident(s.dbClient, req.IncidentID)
	if err != nil {
		log.Printf("Failed to retrieve incident with ID: %s, error: %v\n", req.IncidentID, err)
		return nil, err
	}

	// Log the success of incident retrieval
	log.Printf("Incident retrieved successfully: %+v\n", incident)

	return &IncidentResponse{Incident: &IncidentProto{
		IncidentID:   incident.IncidentID,
		Title:        incident.Title,
		Description:  incident.Description,
		Status:       incident.Status,
		CreationDate: incident.CreationDate,
	}}, nil
}

func (s *IncidentServer) UpdateIncident(ctx context.Context, req *UpdateIncidentRequest) (*IncidentResponse, error) {
	// Log the start of incident update
	log.Printf("Updating incident with ID: %s to status: %s\n", req.IncidentID, req.Status)

	err := updateIncident(s.dbClient, req.IncidentID, req.Status)
	if err != nil {
		log.Printf("Failed to update incident with ID: %s, error: %v\n", req.IncidentID, err)
		return nil, err
	}

	updatedIncident, err := getIncident(s.dbClient, req.IncidentID)
	if err != nil {
		log.Printf("Failed to retrieve updated incident with ID: %s, error: %v\n", req.IncidentID, err)
		return nil, err
	}

	// Log the success of incident update
	log.Printf("Incident updated successfully: %+v\n", updatedIncident)

	// Send message to SQS for incident update
	s.sqsManager.SendMessage(*updatedIncident, "UPDATE")
	log.Printf("Incident message sent to SQS for incident update: %s\n", updatedIncident.IncidentID)

	return &IncidentResponse{Incident: &IncidentProto{
		IncidentID:   updatedIncident.IncidentID,
		Title:        updatedIncident.Title,
		Description:  updatedIncident.Description,
		Status:       updatedIncident.Status,
		CreationDate: updatedIncident.CreationDate,
	}}, nil
}

func (s *IncidentServer) DeleteIncident(ctx context.Context, req *DeleteIncidentRequest) (*DeleteIncidentResponse, error) {
	// Log the start of incident deletion
	log.Printf("Deleting incident with ID: %s\n", req.IncidentID)

	// Retrieve the incident before deletion to send its details to SQS
	incident, err := getIncident(s.dbClient, req.IncidentID)
	if err != nil {
		log.Printf("Failed to retrieve incident for deletion with ID: %s, error: %v\n", req.IncidentID, err)
		return nil, err
	}

	err = deleteIncident(s.dbClient, req.IncidentID)
	if err != nil {
		log.Printf("Failed to delete incident with ID: %s, error: %v\n", req.IncidentID, err)
		return nil, err
	}

	// Log the success of incident deletion
	log.Printf("Incident with ID: %s deleted successfully\n", req.IncidentID)

	// Send message to SQS for incident deletion
	s.sqsManager.SendMessage(*incident, "DELETE")
	log.Printf("Incident message sent to SQS for incident deletion: %s\n", incident.IncidentID)

	return &DeleteIncidentResponse{Success: true}, nil
}
