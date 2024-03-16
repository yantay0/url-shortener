package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

	args := []interface{}{url.OriginalUrl, url.ShortUrl}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&url.ID, &url.CreatedAt, &url.Version)
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, id).Scan(
		&url.ID,
		&url.CreatedAt,
		&url.OriginalUrl,
		&url.ShortUrl,
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
		WHERE id = $3 AND version = $4
		RETURNING version`

	args := []interface{}{
		url.OriginalUrl,
		url.ShortUrl,
		url.ID,
		url.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&url.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (u *UrlStorage) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM url
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := u.DB.ExecContext(ctx, query, id)
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

func (u *UrlStorage) GetAll(originalUrl, shortUrl string) ([]*model.Url, error) {
	query := `
		SELECT id, created_at, original_url, short_url, version
		FROM url
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := []*model.Url{}

	for rows.Next() {
		var url model.Url
		err := rows.Scan(
			&url.ID,
			&url.CreatedAt,
			&url.OriginalUrl,
			&url.ShortUrl,
			&url.Version,
		)
		if err != nil {
			return nil, err
		}
		urls = append(urls, &url)
	}

	return urls, nil
}
