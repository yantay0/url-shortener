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
	Urls  UrlStorage
	Users UserStorage
}

func New(db *sql.DB) Storage {
	return Storage{
		Urls:  UrlStorage{DB: db},
		Users: UserStorage{DB: db},
	}
}
