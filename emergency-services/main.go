package main

import (
	"log"
	"net"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcAPI = "grpc"
	restAPI = "rest"
)

// Set this to either "grpc" or "rest"
const apiType = restAPI

func main() {
	// Data source name (DSN) format: <username>:<password>@tcp(<hostname>:<port>)/<dbname>
	dsn := "root:new_password@tcp(127.0.0.1:3306)/mydb"
	db, err := ConnectToDB(dsn)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	CreateLifeguardsTable(db)
	CreateVehiclesTable(db)

	// Start the appropriate server based on the apiType constant
	switch apiType {
	case grpcAPI:
		// Create the gRPC server
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		s := grpc.NewServer()

		// Register your services
		RegisterLifeguardServiceServer(s, NewLifeguardServer(db))
		RegisterVehicleServiceServer(s, NewVehicleServer(db))

		// Enable gRPC reflection for tools like grpcurl
		reflection.Register(s)

		log.Printf("Server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}

	case restAPI:
		// Start REST API
		if err := startRESTServer(db); err != nil {
			log.Fatalf("Failed to start REST server: %v", err)
		}

	default:
		log.Fatalf("Unknown API type: %v", apiType)
	}

}
