package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSManager struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSManager(client *sqs.Client) (*SQSManager, error) {
	manager := &SQSManager{
		client: client,
	}

	err := manager.ensureQueueExists()
	if err != nil {
		return nil, err
	}

	return manager, nil
}

func (m *SQSManager) ensureQueueExists() error {
	result, err := m.client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		var notFoundErr *types.QueueDoesNotExist
		if errors.As(err, &notFoundErr) {
			createResult, err := m.client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
				QueueName: aws.String(queueName),
			})
			if err != nil {
				return fmt.Errorf("Nie udało się utworzyć kolejki SQS: %v", err)
			}
			m.queueURL = *createResult.QueueUrl
			log.Printf("Utworzono kolejkę SQS %s z URL:\n %s\n", queueName, m.queueURL)
		} else {
			return fmt.Errorf("Nie udało się pobrać kolejki SQS: %v", err)
		}
	} else {
		m.queueURL = *result.QueueUrl
		log.Printf("Kolejka SQS %s już istnieje z URL:\n %s\n", queueName, m.queueURL)
	}

	return nil
}

func (m *SQSManager) SendMessage(incident Incident, operation string) error {
	messageBody := fmt.Sprintf(
		`{"operation":"%s","incidentId":"%s","title":"%s","description":"%s","status":"%s","creationDate":"%s"}`,
		operation, incident.IncidentID, incident.Title, incident.Description, incident.Status, incident.CreationDate,
	)

	_, err := m.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    &m.queueURL,
		MessageBody: &messageBody,
	})

	if err != nil {
		return fmt.Errorf("Nie udało się wysłać wiadomości na kolejkę SQS: %v", err)
	}

	log.Printf("Incydent o ID %s z operacją %s wysłany do kolejki SQS %s\n", incident.IncidentID, operation, queueName)
	return nil
}
