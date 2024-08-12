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
func GetLifeguardByID(db *sql.DB, id int) (*Lifeguard, error) {
	query := `SELECT ID, Name, Login, PasswordHash, YearsOfExperience, Specialization, OnMission, CreatedAt FROM lifeguards WHERE ID = ?`

	var lifeguard Lifeguard
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

// GetAllLifeguards retrieves all lifeguard records from the lifeguards table.
func GetAllLifeguards(db *sql.DB) ([]Lifeguard, error) {
	query := `SELECT ID, Name, Login, PasswordHash, YearsOfExperience, Specialization, OnMission, CreatedAt FROM lifeguards`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving lifeguards: %w", err)
	}
	defer rows.Close()

	var lifeguards []Lifeguard
	for rows.Next() {
		var lifeguard Lifeguard
		var createdAt []byte

		err := rows.Scan(
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
			return nil, fmt.Errorf("error scanning lifeguard: %w", err)
		}

		// Convert the []uint8 to a time.Time
		lifeguard.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, fmt.Errorf("error parsing CreatedAt: %w", err)
		}

		lifeguards = append(lifeguards, lifeguard)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing rows: %w", err)
	}

	return lifeguards, nil
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
func GetVehicleByID(db *sql.DB, id int) (*Vehicle, error) {
	query := `SELECT ID, Type, Location, FuelLevelInLiters, OnMission, LifeguardInChargeID, CreatedAt FROM vehicles WHERE ID = ?`

	var vehicle Vehicle
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

// GetAllVehicles retrieves all vehicle records from the vehicles table.
func GetAllVehicles(db *sql.DB) ([]Vehicle, error) {
	query := `SELECT ID, Type, Location, FuelLevelInLiters, OnMission, LifeguardInChargeID, CreatedAt FROM vehicles`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var vehicle Vehicle
		var createdAt []byte

		err := rows.Scan(
			&vehicle.ID,
			&vehicle.Type,
			&vehicle.Location,
			&vehicle.FuelLevelInLiters,
			&vehicle.OnMission,
			&vehicle.LifeguardInChargeID,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning vehicle: %w", err)
		}

		// Convert the []uint8 to a time.Time
		vehicle.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, fmt.Errorf("error parsing CreatedAt: %w", err)
		}

		vehicles = append(vehicles, vehicle)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing rows: %w", err)
	}

	return vehicles, nil
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
