package repository

import (
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

// type UrlModel struct {
// 	DB *sql.DB
// }

// type Models struct {
// 	Urls UrlModel
// }

// func NewModels(db *sql.DB) Models {
// 	return Models{
// 		Urls: UrlModel{DB: db},
// 	}
// }
