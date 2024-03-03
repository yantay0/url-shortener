package model

import (
	"time"
)

type Url struct {
	ID           int64     `json:"id"`
	Created_at   time.Time `json:"-"`
	Original_url string    `json:"original_url"`
	Short_url    string    `json:"short_url"`
	Version      int32     `json:"version"` // The version number starts at 1 and is incremented each
	// time the url information is updated.
	// User         *User `json:"user"` // after adding seralization

}
