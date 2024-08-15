package main

import (
	"log"
	"net/http"

	"google.golang.org/grpc"
)

// main function sets up the HTTP server and connects to the gRPC server.
func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Set up HTTP server
	lifeguardClient = NewLifeguardServiceClient(conn)
	vehicleClient = NewVehicleServiceClient(conn)

	http.HandleFunc("/lifeguard", CreateLifeguardHandler)
	http.HandleFunc("/lifeguard/get", GetLifeguardHandler)
	http.HandleFunc("/lifeguard/update", UpdateLifeguardHandler)
	http.HandleFunc("/lifeguard/delete", DeleteLifeguardHandler)

	http.HandleFunc("/vehicle", CreateVehicleHandler)
	http.HandleFunc("/vehicle/get", GetVehicleHandler)
	http.HandleFunc("/vehicle/update", UpdateVehicleHandler)
	http.HandleFunc("/vehicle/delete", DeleteVehicleHandler)

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
