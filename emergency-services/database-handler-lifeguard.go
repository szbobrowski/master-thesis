package main

import (
	"database/sql"
	"fmt"
	"time"
)

func CreateLifeguard(db *sql.DB, name, login, passwordHash string, yearsOfExperience int, specialization string, onMission bool) (int64, error) {
	query := `
		INSERT INTO lifeguards (Name, Login, PasswordHash, YearsOfExperience, Specialization, OnMission)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query, name, login, passwordHash, yearsOfExperience, specialization, onMission)
	if err != nil {
		return 0, fmt.Errorf("Nie udało się utworzyć ratownika: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Błąd podczas pobierania ID ostatniego wiersza: %w", err)
	}

	return id, nil
}

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
			return nil, fmt.Errorf("Ratownik o ID %d nie znaleziony", id)
		}
		return nil, fmt.Errorf("Błąd podczas pobierania ratownika: %w", err)
	}

	lifeguard.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, fmt.Errorf("Błąd podczas parsowania pola CreatedAt: %w", err)
	}

	return &lifeguard, nil
}

func UpdateLifeguard(db *sql.DB, id int, name, login, passwordHash string, yearsOfExperience int, specialization string, onMission bool) error {
	query := `
		UPDATE lifeguards
		SET Name = ?, Login = ?, PasswordHash = ?, YearsOfExperience = ?, Specialization = ?, OnMission = ?
		WHERE ID = ?
	`

	_, err := db.Exec(query, name, login, passwordHash, yearsOfExperience, specialization, onMission, id)
	if err != nil {
		return fmt.Errorf("Błąd podczas aktualizowania ratownika: %w", err)
	}

	fmt.Printf("Zaktualizowano ratownika o ID %d!\n", id)
	return nil
}

func DeleteLifeguard(db *sql.DB, id int) error {
	query := `DELETE FROM lifeguards WHERE ID = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Błąd podczas usuwania ratownika: %w", err)
	}

	fmt.Printf("Ratownik o ID %d został usunięty!\n", id)
	return nil
}
