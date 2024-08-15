package main

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateVehicle inserts a new vehicle record into the vehicles table.
func CreateVehicle(db *sql.DB, vehicleType, location string, fuelLevelInLiters int, onMission bool, lifeguardInChargeID int) (int64, error) {
	query := `
		INSERT INTO vehicles (Type, Location, FuelLevelInLiters, OnMission, LifeguardInChargeID)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, vehicleType, location, fuelLevelInLiters, onMission, lifeguardInChargeID)
	if err != nil {
		return 0, fmt.Errorf("error creating vehicle: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %w", err)
	}

	fmt.Printf("Vehicle with ID %d created successfully!\n", id)
	return id, nil
}

// GetVehicleByID retrieves a vehicle record by its ID.
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
			return nil, fmt.Errorf("vehicle with ID %d not found", id)
		}
		return nil, fmt.Errorf("error retrieving vehicle: %w", err)
	}

	// Convert the []uint8 to a time.Time
	vehicle.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing CreatedAt: %w", err)
	}

	return &vehicle, nil
}

// UpdateVehicle updates an existing vehicle record in the vehicles table.
func UpdateVehicle(db *sql.DB, id int, vehicleType, location string, fuelLevelInLiters int, onMission bool, lifeguardInChargeID int) error {
	query := `
		UPDATE vehicles
		SET Type = ?, Location = ?, FuelLevelInLiters = ?, OnMission = ?, LifeguardInChargeID = ?
		WHERE ID = ?
	`
	_, err := db.Exec(query, vehicleType, location, fuelLevelInLiters, onMission, lifeguardInChargeID, id)
	if err != nil {
		return fmt.Errorf("error updating vehicle: %w", err)
	}

	fmt.Printf("Vehicle with ID %d updated successfully!\n", id)
	return nil
}

// DeleteVehicle deletes a vehicle record from the vehicles table by its ID.
func DeleteVehicle(db *sql.DB, id int) error {
	query := `DELETE FROM vehicles WHERE ID = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting vehicle: %w", err)
	}

	fmt.Printf("Vehicle with ID %d deleted successfully!\n", id)
	return nil
}
