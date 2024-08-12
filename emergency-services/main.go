package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Data source name (DSN) format: <username>:<password>@tcp(<hostname>:<port>)/<dbname>
	dsn := "root:test@tcp(127.0.0.1:3306)/mydb"
	db, err := ConnectToDB(dsn)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
}
