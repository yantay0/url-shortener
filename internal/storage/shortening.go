package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/yantay0/url-shortener/internal/model"
)

var (
	ErrIdentifierExists = errors.New("identifier already exists")
	ErrInvalidURL       = errors.New("invalid url")
)

type ShorteningsStorage struct {
	DB *sql.DB
}

func (s *ShorteningsStorage) Insert(shortening *model.Shortening) error {
	query := `
		INSERT INTO url (original_url, identifier, user_id)
		VALUES ($1, $2, $3)
		RETURNING identifier, created_at, version`

	args := []interface{}{shortening.OriginalURL, shortening.Identifier, shortening.UserID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&shortening.Identifier, &shortening.CreatedAt, &shortening.Version)
}

func (s *ShorteningsStorage) Get(identifier string) (*model.Shortening, error) {
	if identifier != "" {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT indentifer, created_at, original_url, short_url, version
		FROM shortening 
		WHERE identifier = $1`

	var shortening model.Shortening
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, identifier).Scan(
		shortening.Identifier,
		shortening.CreatedAt,
		shortening.OriginalURL,
		shortening.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &shortening, err
}

func (s *ShorteningsStorage) Update(shortening *model.Shortening) error {
	query := `
		UPDATE url
		SET original_url = $1, version = version + 1
		WHERE identifier = $2 AND version = $3
		RETURNING version`

	args := []interface{}{
		shortening.OriginalURL,
		shortening.Identifier,
		shortening.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, args...).Scan(&shortening.Version)

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

func (s *ShorteningsStorage) Delete(Identifier string) error {
	if Identifier != "" {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM url
		WHERE identifier = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, query, Identifier)
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

func (s *ShorteningsStorage) GetAll(OriginalURL string, filters model.Filters) ([]*model.Shortening, model.Metadata, error) {
	// order by id for the consistent ordering
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), identifier, created_at, original_url, version, user_id
		FROM shortening
		WHERE (LOWER(original_url) = LOWER($1) OR $1 = '')
		ORDER BY %s %s, identifier ASC
		LIMIT $2 OFFSET $3`, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{OriginalURL, filters.Limit(), filters.Offset()}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, model.Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	shortenings := []*model.Shortening{}

	for rows.Next() {
		var shortening model.Shortening
		err := rows.Scan(
			&totalRecords,
			&shortening.Identifier,
			&shortening.CreatedAt,
			&shortening.OriginalURL,
			&shortening.Version,
			&shortening.UserID,
		)

		if err != nil {
			return nil, model.Metadata{}, err
		}
		shortenings = append(shortenings, &shortening)
	}

	if err = rows.Err(); err != nil {
		return nil, model.Metadata{}, err
	}

	metadata := model.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return shortenings, metadata, nil
}

func (s *ShorteningsStorage) GetUserAllShortenings(userID int64) ([]*model.Shortening, error) {
	query := `
	SELECT created_at, original_url, identifier, version, user_id
	FROM shortening
	WHERE user_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{userID}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	shortenings := []*model.Shortening{}

	for rows.Next() {
		var shortening model.Shortening
		err := rows.Scan(
			&shortening.CreatedAt,
			&shortening.OriginalURL,
			&shortening.Identifier,
			&shortening.Version,
			&shortening.UserID,
		)

		if err != nil {
			return nil, err
		}
		shortenings = append(shortenings, &shortening)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shortenings, nil
}
