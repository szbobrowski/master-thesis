package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello from client-handler!")

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := NewLifeguardServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Example: Calling CreateLifeguard
	response, err := client.CreateLifeguard(ctx, &CreateLifeguardRequest{
		Name:              "John Doe",
		Login:             "johndoe123",
		PasswordHash:      "hashedpassword",
		YearsOfExperience: 5,
		Specialization:    "Swimming",
		OnMission:         false,
	})
	if err != nil {
		log.Fatalf("could not create lifeguard: %v", err)
	}
	log.Printf("Lifeguard created, ID: %v", response.Id)
}

// Ensure to generate client code as described or copy it from the existing emergency-services module.
