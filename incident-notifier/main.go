package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	tableName = "Incidents"
	queueName = "IncidentsQueue"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolver(
		aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			switch service {
			case dynamodb.ServiceID:
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           "http://localhost:4566",
					SigningRegion: "us-west-2",
				}, nil
			case sqs.ServiceID:
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           "http://localhost:4566",
					SigningRegion: "us-west-2",
				}, nil
			default:
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested for service: %s", service)
			}
		}),
	))
	if err != nil {
		log.Fatalf("Nie udało się pobrać konfiguracji SDK, %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	sqsClient := sqs.NewFromConfig(cfg)

	sqsManager, err := NewSQSManager(sqsClient)
	if err != nil {
		log.Fatalf("Nie udało się utworzyć menedżera SQS, %v", err)
	}

	exists, err := tableExists(dynamoClient)
	if err != nil {
		log.Fatalf("Nie udało się sprawdzić czy tabela istnieje, %v", err)
	}

	if !exists {
		err = createTable(dynamoClient)
		if err != nil {
			log.Fatalf("Utwrzono tabelę Incidents, %v", err)
		}
	} else {
		fmt.Println("Tabela Incidents już istnieje.")
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Nie udało się rozpocząć nasłuchiwania na porcie 50052: %v", err)
	}

	grpcServer := grpc.NewServer()
	incidentServer := NewIncidentServer(dynamoClient, sqsManager)

	RegisterIncidentServiceServer(grpcServer, incidentServer)

	reflection.Register(grpcServer)

	log.Printf("Serwer nasłuchuje na adresie %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Błąd podczas działania serwera: %v", err)
	}
}
