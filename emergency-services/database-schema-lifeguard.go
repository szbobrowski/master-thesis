package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type LifeguardDTO struct {
	ID                int
	Name              string
	Login             string
	PasswordHash      string
	YearsOfExperience int
	Specialization    string
	OnMission         bool
	CreatedAt         time.Time
}

func CreateLifeguardsTable(db *sql.DB) error {
	checkTableQuery := `
		SELECT COUNT(*)
		FROM information_schema.tables 
		WHERE table_schema = DATABASE() AND table_name = 'lifeguards';
	`

	var count int
	err := db.QueryRow(checkTableQuery).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking if table exists: %w", err)
	}

	if count == 0 {
		createTableQuery := `
			CREATE TABLE lifeguards (
				ID INT AUTO_INCREMENT PRIMARY KEY,
				Name TEXT,
				Login TEXT,
				PasswordHash TEXT,
				YearsOfExperience INT,
				Specialization TEXT,
				OnMission BOOLEAN,
				CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
		`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			return fmt.Errorf("Błąd podczas tworzenia tabeli lifeguards: %w", err)
		}

		fmt.Println("Pomyślnie utworzono tabelę lifeguards!")
	} else {
		fmt.Println("Tabela lifeguards już istnieje.")
	}

	return nil
}
