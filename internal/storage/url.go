package storage

import (
	"database/sql"
	"errors"

	"github.com/yantay0/url-shortener/internal/model"
)

type UrlStorage struct {
	DB *sql.DB
}

func (u *UrlStorage) Insert(url *model.Url) error {
	query := `
		INSERT INTO url (original_url, short_url)
		VALUES ($1, $2)
		RETURNING id, created_at, version`

	args := []interface{}{url.Original_url, url.Short_url}
	return u.DB.QueryRow(query, args...).Scan(&url.ID, &url.Created_at, &url.Version)
}

func (u *UrlStorage) Get(id int64) (*model.Url, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, original_url, short_url, version
		FROM url 
		WHERE id = $1`

	var url model.Url
	err := u.DB.QueryRow(query, id).Scan(
		&url.ID,
		&url.Created_at,
		&url.Original_url,
		&url.Short_url,
		&url.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &url, err
}

func (u *UrlStorage) Update(url *model.Url) error {
	query := `
		UPDATE url
		SET original_url = $1, short_url = $2, version = version + 1
		WHERE id = $3
		RETURNING version`

	args := []interface{}{
		url.Original_url,
		url.Short_url,
		url.ID,
	}

	return u.DB.QueryRow(query, args...).Scan(&url.Version)
}

func (u *UrlStorage) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM url
		WHERE id = $1`

	result, err := u.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
