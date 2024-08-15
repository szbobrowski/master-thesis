package main

import (
	"database/sql"
	"fmt"
)

// ConnectToDB establishes a connection to the MySQL database and returns the database connection.
func ConnectToDB(dsn string) (*sql.DB, error) {
	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Verify the connection to the database
	if err := db.Ping(); err != nil {
		db.Close() // Close the database connection if ping fails
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	fmt.Println("Successfully connected to the database!")
	return db, nil
}
