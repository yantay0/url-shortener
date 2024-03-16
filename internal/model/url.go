package model

import (
	"time"
)

type Url struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	OriginalUrl string    `json:"original_url"`
	ShortUrl    string    `json:"short_url"`
	Version     int32     `json:"version"` // The version number starts at 1 and is incremented each
	// time the url information is updated.
	// User         *User `json:"user"` // after adding seralization

}
