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
	Urls        UrlStorage
	Permissions PermissionsStorage
	Tokens      TokenStorage
	Users       UserStorage
}

func New(db *sql.DB) Storage {
	return Storage{
		Urls:        UrlStorage{DB: db},
		Permissions: PermissionsStorage{DB: db},
		Tokens:      TokenStorage{DB: db},
		Users:       UserStorage{DB: db},
	}
}
