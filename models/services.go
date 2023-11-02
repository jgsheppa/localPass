package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
	PassService
	db *sql.DB
}

func NewServices() (*Service, error) {
	db, err := sql.Open("sqlite3", "locapass.db")
	if err != nil {
		return &Service{}, err
	}

	_, err = db.Exec(`PRAGMA journal_mode = wal;`)
	if err != nil {
		return &Service{}, err
	}
	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		return &Service{}, err
	}

	return &Service{
		PassService: NewPassService(db),
		db:          db,
	}, nil
}
