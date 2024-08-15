package main

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateLifeguard inserts a new lifeguard record into the lifeguards table.
func CreateLifeguard(db *sql.DB, name, login, passwordHash string, yearsOfExperience int, specialization string, onMission bool) (int64, error) {
	query := `
		INSERT INTO lifeguards (Name, Login, PasswordHash, YearsOfExperience, Specialization, OnMission)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query, name, login, passwordHash, yearsOfExperience, specialization, onMission)
	if err != nil {
		return 0, fmt.Errorf("error inserting lifeguard: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %w", err)
	}

	fmt.Printf("Lifeguard with ID %d created successfully!\n", id)
	return id, nil
}

// GetLifeguardByID retrieves a lifeguard record by its ID.
func GetLifeguardByID(db *sql.DB, id int) (*LifeguardDTO, error) {
	query := `SELECT ID, Name, Login, PasswordHash, YearsOfExperience, Specialization, OnMission, CreatedAt FROM lifeguards WHERE ID = ?`

	var lifeguard LifeguardDTO
	var createdAt []byte

	err := db.QueryRow(query, id).Scan(
		&lifeguard.ID,
		&lifeguard.Name,
		&lifeguard.Login,
		&lifeguard.PasswordHash,
		&lifeguard.YearsOfExperience,
		&lifeguard.Specialization,
		&lifeguard.OnMission,
		&createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lifeguard with ID %d not found", id)
		}
		return nil, fmt.Errorf("error retrieving lifeguard: %w", err)
	}

	// Convert the []uint8 to a time.Time
	lifeguard.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing CreatedAt: %w", err)
	}

	return &lifeguard, nil
}

// UpdateLifeguard updates the details of a lifeguard in the lifeguards table.
func UpdateLifeguard(db *sql.DB, id int, name, login, passwordHash string, yearsOfExperience int, specialization string, onMission bool) error {
	query := `
		UPDATE lifeguards
		SET Name = ?, Login = ?, PasswordHash = ?, YearsOfExperience = ?, Specialization = ?, OnMission = ?
		WHERE ID = ?
	`

	result, err := db.Exec(query, name, login, passwordHash, yearsOfExperience, specialization, onMission, id)
	if err != nil {
		return fmt.Errorf("error updating lifeguard: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no lifeguard with ID %d found", id)
	}

	fmt.Printf("Lifeguard with ID %d updated successfully!\n", id)
	return nil
}

// DeleteLifeguard deletes a lifeguard record from the lifeguards table by ID.
func DeleteLifeguard(db *sql.DB, id int) error {
	query := `DELETE FROM lifeguards WHERE ID = ?`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting lifeguard: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no lifeguard with ID %d found", id)
	}

	fmt.Printf("Lifeguard with ID %d deleted successfully!\n", id)
	return nil
}
