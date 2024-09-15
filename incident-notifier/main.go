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
	// Load the default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolver(
		aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			switch service {
			case dynamodb.ServiceID:
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           "http://localhost:4566", // LocalStack URL for DynamoDB
					SigningRegion: "us-west-2",             // Region for signing, can be any valid AWS region
				}, nil
			case sqs.ServiceID:
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           "http://localhost:4566", // LocalStack URL for SQS
					SigningRegion: "us-west-2",             // Region for signing, can be any valid AWS region
				}, nil
			default:
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested for service: %s", service)
			}
		}),
	))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	sqsClient := sqs.NewFromConfig(cfg)

	// Initialize SQSManager
	sqsManager, err := NewSQSManager(sqsClient)
	if err != nil {
		log.Fatalf("unable to initialize SQS manager, %v", err)
	}

	// Check if the Incidents table exists, and create it if not
	exists, err := tableExists(dynamoClient)
	if err != nil {
		log.Fatalf("unable to check if table exists, %v", err)
	}

	if !exists {
		fmt.Println("Table does not exist, creating table...")
		err = createTable(dynamoClient)
		if err != nil {
			log.Fatalf("unable to create table, %v", err)
		}
	} else {
		fmt.Println("Table already exists.")
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	incidentServer := NewIncidentServer(dynamoClient, sqsManager)

	RegisterIncidentServiceServer(grpcServer, incidentServer)

	// Enable reflection for gRPC server
	reflection.Register(grpcServer)

	log.Printf("Server is running at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
