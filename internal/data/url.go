package data

import (
	"database/sql"
	"time"
)

type Url struct {
	ID           int64     `json:"id"`
	Created_at   time.Time `json:"-"`
	Short_url    string    `json:"short_url"`
	Original_url string    `json:"original_url"`
	Version      int32     `json:"version"` // The version number starts at 1 and is incremented each
	// time the url information is updated.
	// User         *User `json:"user"` // after adding seralization

}

type UrlModel struct {
	DB *sql.DB
}

func (m UrlModel) Insert(url *Url) error {
	return nil
}

func (m UrlModel) Get(id int64) error {
	return nil
}

func (m UrlModel) Update(url *Url) error {
	return nil
}

func (m UrlModel) Delete(id int64) error {
	return nil
}
