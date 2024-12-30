package main

import (
	"log"
	"net"

	grpc "google.golang.org/grpc"
)

func main() {
	dataSourceName := "root:new_password@tcp(127.0.0.1:3306)/mydb"
	db, err := ConnectToDB(dataSourceName)
	if err != nil {
		log.Fatalf("Nie udało się nawiązać połączenia z bazą danych: %v", err)
	}
	defer db.Close()

	CreateLifeguardsTable(db)
	CreateVehiclesTable(db)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Nie udało się uruchomić serwera gRPC: %v", err)
	}

	s := grpc.NewServer()
	RegisterLifeguardServiceServer(s, NewLifeguardServer(db))
	RegisterVehicleServiceServer(s, NewVehicleServer(db))

	log.Printf("Serwer nasłuchuje na adresie %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Błąd podczas działania serwera: %v", err)
	}
}
