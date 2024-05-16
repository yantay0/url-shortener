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
	Shortenings ShorteningsStorage
	Permissions PermissionsStorage
	Tokens      TokenStorage
	Users       UserStorage
}

func New(db *sql.DB) Storage {
	return Storage{
		Shortenings: ShorteningsStorage{DB: db},
		Permissions: PermissionsStorage{DB: db},
		Tokens:      TokenStorage{DB: db},
		Users:       UserStorage{DB: db},
	}
}
