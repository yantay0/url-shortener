package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Urls UrlModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Urls: UrlModel{DB: db},
	}
}
