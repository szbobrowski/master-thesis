package main

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc"
)

func main() {
	restConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Nie udało się połączyć z serwerem gRPC emergency-services: %v", err)
	}
	log.Println("Połączono z serwerem gRPC emergency-services na localhost:50051")
	defer restConn.Close()

	lifeguardClient = NewLifeguardServiceClient(restConn)
	vehicleClient = NewVehicleServiceClient(restConn)

	graphqlConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Nie udało się połączyć z serwerem gRPC incident-notifier: %v", err)
	}
	log.Println("Połączono z serwerem gRPC incident-notifier na localhost:50052")
	defer graphqlConn.Close()

	incidentClient = NewIncidentServiceClient(graphqlConn)
	initIncidentGrpcClient()

	mux := http.NewServeMux()

	mux.HandleFunc("/graphql", graphqlHandler)

	mux.HandleFunc("/lifeguard", CreateLifeguardHandler)
	mux.HandleFunc("/lifeguard/get", GetLifeguardHandler)
	mux.HandleFunc("/lifeguard/update", UpdateLifeguardHandler)
	mux.HandleFunc("/lifeguard/delete", DeleteLifeguardHandler)

	mux.HandleFunc("/vehicle", CreateVehicleHandler)
	mux.HandleFunc("/vehicle/get", GetVehicleHandler)
	mux.HandleFunc("/vehicle/update", UpdateVehicleHandler)
	mux.HandleFunc("/vehicle/delete", DeleteVehicleHandler)

	fmt.Println("Serwer obsługujący zapytania klienta nasłuchuje na adresie http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Nie udało się uruchomić serwera http: %v", err)
	}
}
