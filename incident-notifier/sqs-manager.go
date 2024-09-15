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

// SQSManager manages operations related to SQS queues
type SQSManager struct {
	client   *sqs.Client
	queueURL string
}

// NewSQSManager creates a new SQSManager and ensures the queue exists
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

// ensureQueueExists checks if the queue exists and creates it if it doesn't
func (m *SQSManager) ensureQueueExists() error {
	// Check if the queue exists
	result, err := m.client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		var notFoundErr *types.QueueDoesNotExist
		if errors.As(err, &notFoundErr) {
			// Queue does not exist, create it
			createResult, err := m.client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
				QueueName: aws.String(queueName),
			})
			if err != nil {
				return fmt.Errorf("unable to create queue: %v", err)
			}
			m.queueURL = *createResult.QueueUrl
			log.Printf("Created SQS queue %s with URL: %s\n", queueName, m.queueURL)
		} else {
			return fmt.Errorf("unable to get queue URL: %v", err)
		}
	} else {
		m.queueURL = *result.QueueUrl
		log.Printf("SQS queue %s already exists with URL: %s\n", queueName, m.queueURL)
	}

	return nil
}

// SendMessage sends a message to the SQS queue with the given incident and operation type
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
		return fmt.Errorf("failed to send message to SQS queue: %v", err)
	}

	log.Printf("Incident %s with operation %s sent to SQS queue %s\n", incident.IncidentID, operation, queueName)
	return nil
}
