package main

import (
	"log"
	"net"

	grpc "google.golang.org/grpc"
)

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

	// Create the gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterLifeguardServiceServer(s, NewLifeguardServer(db))
	RegisterVehicleServiceServer(s, NewVehicleServer(db))

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
