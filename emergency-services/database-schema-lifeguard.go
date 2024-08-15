package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

// LifeguardDTO represents a lifeguard record in the lifeguards table.
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

// CreateLifeguardsTable checks if the lifeguards table exists, and creates it if it doesn't.
func CreateLifeguardsTable(db *sql.DB) error {
	// SQL query to check if the table exists
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

	// If the table does not exist, create it
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
			return fmt.Errorf("error creating lifeguards table: %w", err)
		}

		fmt.Println("Lifeguards table created successfully!")
	} else {
		fmt.Println("Lifeguards table already exists.")
	}

	return nil
}
