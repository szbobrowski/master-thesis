package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

type VehicleDTO struct {
	ID                  int
	Type                string
	Location            string
	FuelLevelInLiters   int
	OnMission           bool
	LifeguardInChargeID int
	CreatedAt           time.Time
}

// CreateVehiclesTable checks if the vehicles table exists, and creates it if it doesn't.
func CreateVehiclesTable(db *sql.DB) error {
	// SQL query to check if the table exists
	checkTableQuery := `
		SELECT COUNT(*)
		FROM information_schema.tables 
		WHERE table_schema = DATABASE() AND table_name = 'vehicles';
	`

	var count int
	err := db.QueryRow(checkTableQuery).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking if table exists: %w", err)
	}

	// If the table does not exist, create it
	if count == 0 {
		createTableQuery := `
			CREATE TABLE vehicles (
				ID INT AUTO_INCREMENT PRIMARY KEY,
				Type TEXT,
				Location TEXT,
				FuelLevelInLiters INT,
				OnMission BOOLEAN,
				LifeguardInChargeID INT,
				CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (LifeguardInChargeID) REFERENCES lifeguards(ID)
			);
		`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			return fmt.Errorf("error creating vehicles table: %w", err)
		}

		fmt.Println("Vehicles table created successfully!")
	} else {
		fmt.Println("Vehicles table already exists.")
	}

	return nil
}
