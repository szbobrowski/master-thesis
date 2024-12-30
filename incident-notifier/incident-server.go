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

	err := createIncident(s.dbClient, incident)
	if err != nil {
		log.Printf("Nie udało się utworzyć incydentu: %v\n", err)
		return nil, err
	}
	log.Printf("Incydent został utworzony: %s\n", incidentID)

	s.sqsManager.SendMessage(incident, "CREATE")

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
		log.Printf("Nie udało się pobrać incydentu o ID: %s, błąd: %v\n", req.IncidentID, err)
		return nil, err
	}

	log.Printf("Pobrano incydent: %+v\n", incident)

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
		log.Printf("Nie udało się zaktualizować incydentu o ID: %s, błąd: %v\n", req.IncidentID, err)
		return nil, err
	}

	updatedIncident, err := getIncident(s.dbClient, req.IncidentID)
	if err != nil {
		log.Printf("Nie udało się pobrać incydentu o ID: %s, błąd: %v\n", req.IncidentID, err)
		return nil, err
	}
	log.Printf("Zaktualizowano incydent: %+v\n", updatedIncident)

	s.sqsManager.SendMessage(*updatedIncident, "UPDATE")

	return &IncidentResponse{Incident: &IncidentProto{
		IncidentID:   updatedIncident.IncidentID,
		Title:        updatedIncident.Title,
		Description:  updatedIncident.Description,
		Status:       updatedIncident.Status,
		CreationDate: updatedIncident.CreationDate,
	}}, nil
}

func (s *IncidentServer) DeleteIncident(ctx context.Context, req *DeleteIncidentRequest) (*DeleteIncidentResponse, error) {
	incident, err := getIncident(s.dbClient, req.IncidentID)
	if err != nil {
		log.Printf("Nie udało się pobrać incydentu do usunięcia, ID incydentu: %s, błąd: %v\n", req.IncidentID, err)
		return nil, err
	}

	err = deleteIncident(s.dbClient, req.IncidentID)
	if err != nil {
		log.Printf("Nie udało się usunąć incydentu o ID: %s, błąd: %v\n", req.IncidentID, err)
		return nil, err
	}
	log.Printf("Usunięgo incydent o ID %s\n", req.IncidentID)

	s.sqsManager.SendMessage(*incident, "DELETE")

	return &DeleteIncidentResponse{Success: true}, nil
}
