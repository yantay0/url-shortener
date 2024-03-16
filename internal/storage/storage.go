package storage

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Storage struct {
	Urls UrlStorage
}

func NewModels(db *sql.DB) Storage {
	return Storage{
		Urls: UrlStorage{DB: db},
	}
}
