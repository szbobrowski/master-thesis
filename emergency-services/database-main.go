package main

import (
	"database/sql"
	"fmt"
)

func ConnectToDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Błąd podczas otwierania bazy danych: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("Błąd podczas próby połączenia z bazą danych: %w", err)
	}

	fmt.Println("Udało się połączyć z bazą danych!")
	return db, nil
}
