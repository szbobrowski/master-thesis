package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQSConsumer manages consuming messages from an SQS queue
type SQSConsumer struct {
	client   *sqs.Client
	queueURL string
}

// NewSQSConsumer creates a new SQSConsumer with the given client and queue URL
func NewSQSConsumer(client *sqs.Client) (*SQSConsumer, error) {
	consumer := &SQSConsumer{
		client: client,
	}

	err := consumer.ensureQueueExists()
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// ensureQueueExists checks if the queue exists and retrieves its URL
func (c *SQSConsumer) ensureQueueExists() error {
	result, err := c.client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String("IncidentsQueue"), // Replace with your actual queue name
	})
	if err != nil {
		return fmt.Errorf("unable to get queue URL: %v", err)
	}

	c.queueURL = *result.QueueUrl
	log.Printf("SQS queue %s exists with URL: %s\n", "IncidentsQueue", c.queueURL)
	return nil
}

// StartPolling starts polling the SQS queue for messages and processing them
func (c *SQSConsumer) StartPolling() {
	for {
		output, err := c.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            &c.queueURL,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     5,
			VisibilityTimeout:   20,
		})

		if err != nil {
			log.Printf("Error receiving messages: %v", err)
			continue
		}

		if len(output.Messages) == 0 {
			// log.Println("No messages received, waiting...")
			continue
		}

		for _, message := range output.Messages {
			err := c.processMessage(message)
			if err != nil {
				log.Printf("Failed to process message: %v", err)
				continue
			}

			// Delete the message from the queue after successful processing
			_, err = c.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      &c.queueURL,
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Printf("Failed to delete message: %v", err)
			}
		}
	}
}

// processMessage processes a single SQS message and displays it on the console
func (c *SQSConsumer) processMessage(message types.Message) error {
	fmt.Printf("Received message: %s\n", *message.Body)
	return nil
}

func main() {
	// Load the default AWS configuration with LocalStack endpoint resolver
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolver(
		aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			switch service {
			case sqs.ServiceID:
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           "http://localhost:4566", // LocalStack URL for SQS
					SigningRegion: "us-west-2",
				}, nil
			default:
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested for service: %s", service)
			}
		}),
	))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	consumer, err := NewSQSConsumer(client)
	if err != nil {
		log.Fatalf("Failed to create SQS consumer: %v", err)
	}

	log.Println("Starting SQS consumer...")
	consumer.StartPolling()
}
