package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

func CreateVehiclesTable(db *sql.DB) error {
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
			return fmt.Errorf("Błąd podczas tworzenia tabeli vehicles: %w", err)
		}

		fmt.Println("Pomyślnie utworzono tabelę vehicles!")
	} else {
		fmt.Println("Tabela vehicles już istnieje.")
	}

	return nil
}
