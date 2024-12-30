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

type SQSConsumer struct {
	client   *sqs.Client
	queueURL string
}

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

func (c *SQSConsumer) ensureQueueExists() error {
	result, err := c.client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String("IncidentsQueue"),
	})
	if err != nil {
		return fmt.Errorf("Nie udało się pobrać URL kolejki: %v", err)
	}

	c.queueURL = *result.QueueUrl
	log.Printf("Kolejka SQS %s już istnieje, URL kolejki:\n %s\n", "IncidentsQueue", c.queueURL)
	return nil
}

func (c *SQSConsumer) StartPolling() {
	for {
		output, err := c.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            &c.queueURL,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     5,
			VisibilityTimeout:   20,
		})

		if err != nil {
			log.Printf("Błąd podczas odbierania wiadomości: %v", err)
			continue
		}

		if len(output.Messages) == 0 {
			continue
		}

		for _, message := range output.Messages {
			err := c.processMessage(message)
			if err != nil {
				log.Printf("Nie udało się przetworzyć wiadomości: %v", err)
				continue
			}

			_, err = c.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      &c.queueURL,
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Printf("Nie udało się usunąć wiadomości: %v", err)
			}
		}
	}
}

func (c *SQSConsumer) processMessage(message types.Message) error {
	fmt.Printf("Otrzymano wiadomość: %s\n", *message.Body)
	return nil
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolver(
		aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			switch service {
			case sqs.ServiceID:
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           "http://localhost:4566",
					SigningRegion: "us-west-2",
				}, nil
			default:
				return aws.Endpoint{}, fmt.Errorf("Nieznany endpoint dla serwisu: %s", service)
			}
		}),
	))
	if err != nil {
		log.Fatalf("Nie udało się pobrać konfiguracji SDK, %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	consumer, err := NewSQSConsumer(client)
	if err != nil {
		log.Fatalf("Nie udało się utworzyć odbiorcy SQS: %v", err)
	}

	log.Println("Serwis rozpoczyna nasłuchiwanie na wiadomości SQS...")
	consumer.StartPolling()
}
