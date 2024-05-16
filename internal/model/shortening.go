package model

import (
	"net/url"
	"strings"
	"time"

	"github.com/yantay0/url-shortener/internal/util"
)

const alphabet = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"

var alphabetLen = uint32(len(alphabet))

type Shortening struct {
	Identifier  string    `json:"identifier"`
	OriginalURL string    `json:"original_url"`
	Version     int32     `json:"version"` // The version number starts at 1 and is incremented each time the url information is updated.
	UserID      int64     `json:"user_id"` // after adding seralization
	Visits      int64     `json:"visits"`
	CreatedAt   time.Time `json:"created_at"`
}

func GenerateShortURL(id uint32) string {
	var (
		digits  []uint32
		num     = id
		builder strings.Builder
	)

	for num > 0 {
		digits = append(digits, num%alphabetLen)
		num /= alphabetLen
	}

	util.Reverse(digits)

	for _, digit := range digits {
		builder.WriteString(string(alphabet[digit]))
	}

	return builder.String()
}

func PrependBaseURL(baseURL, identifier string) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	parsed.Path = identifier

	return parsed.String(), nil
}
