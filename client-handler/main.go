package main

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc"
)

// main function sets up the HTTP server and connects to the gRPC servers.
func main() {
	// Set up a connection to the gRPC server for the REST APIs (lifeguard/vehicle services)
	restConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect to gRPC server for REST: %v", err)
	}
	log.Println("Connected to gRPC server at localhost:50051 for REST")
	defer restConn.Close()

	// Initialize gRPC clients for REST API usage
	lifeguardClient = NewLifeguardServiceClient(restConn)
	vehicleClient = NewVehicleServiceClient(restConn)

	// Set up a connection to the gRPC server for the GraphQL API (incident service)
	graphqlConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect to gRPC server for GraphQL: %v", err)
	}
	defer graphqlConn.Close()

	// Initialize the gRPC client for the GraphQL API
	incidentClient = NewIncidentServiceClient(graphqlConn)
	initIncidentGrpcClient() // Set the global client for GraphQL resolvers
	// Create a new ServeMux to handlNewIncidentServiceClient(graphqlConn)e different routes
	mux := http.NewServeMux()

	// GraphQL handler at /graphql
	mux.HandleFunc("/graphql", graphqlHandler)

	// REST API endpoints for Lifeguards
	mux.HandleFunc("/lifeguard", CreateLifeguardHandler)
	mux.HandleFunc("/lifeguard/get", GetLifeguardHandler)
	mux.HandleFunc("/lifeguard/update", UpdateLifeguardHandler)
	mux.HandleFunc("/lifeguard/delete", DeleteLifeguardHandler)

	// REST API endpoints for Vehicles
	mux.HandleFunc("/vehicle", CreateVehicleHandler)
	mux.HandleFunc("/vehicle/get", GetVehicleHandler)
	mux.HandleFunc("/vehicle/update", UpdateVehicleHandler)
	mux.HandleFunc("/vehicle/delete", DeleteVehicleHandler)

	// Start the server on port 8080 and serve both GraphQL and REST
	fmt.Println("GraphQL and REST API server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
