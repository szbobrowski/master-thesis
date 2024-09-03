package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const tableName = "Incidents"

func main() {
	// Load the default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolver(
		aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           "http://localhost:4566", // LocalStack URL
					SigningRegion: "us-west-2",             // Region for signing, can be any valid AWS region
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		}),
	))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	// Check if the Incidents table exists, and create it if not
	exists, err := tableExists(client)
	if err != nil {
		log.Fatalf("unable to check if table exists, %v", err)
	}

	if !exists {
		fmt.Println("Table does not exist, creating table...")
		err = createTable(client)
		if err != nil {
			log.Fatalf("unable to create table, %v", err)
		}
	} else {
		fmt.Println("Table already exists.")
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	incidentServer := NewIncidentServer(client)

	RegisterIncidentServiceServer(grpcServer, incidentServer)

	// Enable reflection for gRPC server
	reflection.Register(grpcServer)

	log.Printf("Server is running at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
