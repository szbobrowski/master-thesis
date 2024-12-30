package main

import (
	"database/sql"
	"fmt"
	"time"
)

func CreateVehicle(db *sql.DB, vehicleType, location string, fuelLevelInLiters int, onMission bool, lifeguardInChargeID int) (int64, error) {
	query := `
		INSERT INTO vehicles (Type, Location, FuelLevelInLiters, OnMission, LifeguardInChargeID)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, vehicleType, location, fuelLevelInLiters, onMission, lifeguardInChargeID)
	if err != nil {
		return 0, fmt.Errorf("Nie udało się utworzyć pojazdu: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Błąd podczas pobierania ID ostatniego wiersza: %w", err)
	}

	return id, nil
}

func GetVehicleByID(db *sql.DB, id int) (*VehicleDTO, error) {
	query := `SELECT ID, Type, Location, FuelLevelInLiters, OnMission, LifeguardInChargeID, CreatedAt FROM vehicles WHERE ID = ?`

	var vehicle VehicleDTO
	var createdAt []byte

	err := db.QueryRow(query, id).Scan(
		&vehicle.ID,
		&vehicle.Type,
		&vehicle.Location,
		&vehicle.FuelLevelInLiters,
		&vehicle.OnMission,
		&vehicle.LifeguardInChargeID,
		&createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Pojazd o ID %d nie znaleziony", id)
		}
		return nil, fmt.Errorf("Błąd podczas pobierania pojazdu: %w", err)
	}

	vehicle.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, fmt.Errorf("Błąd podczas parsowania pola CreatedAt: %w", err)
	}

	return &vehicle, nil
}

func UpdateVehicle(db *sql.DB, id int, vehicleType, location string, fuelLevelInLiters int, onMission bool, lifeguardInChargeID int) error {
	query := `
		UPDATE vehicles
		SET Type = ?, Location = ?, FuelLevelInLiters = ?, OnMission = ?, LifeguardInChargeID = ?
		WHERE ID = ?
	`
	_, err := db.Exec(query, vehicleType, location, fuelLevelInLiters, onMission, lifeguardInChargeID, id)
	if err != nil {
		return fmt.Errorf("Błąd podczas aktualizowania pojazdu: %w", err)
	}

	fmt.Printf("Zaktualizowano pojazd o ID %d!\n", id)
	return nil
}

func DeleteVehicle(db *sql.DB, id int) error {
	query := `DELETE FROM vehicles WHERE ID = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Błąd podczas usuwania pojazdu: %w", err)
	}

	fmt.Printf("Pojazd o ID %d został usunięty!\n", id)
	return nil
}
