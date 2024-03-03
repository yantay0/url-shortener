package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // init driver
	"github.com/yantay0/url-shortener/internal/config"
)

type Repository struct {
	Db *sql.DB
}

func OpenDB(cfg *config.Config) (*Repository, error) {

	const op = "repository.postgres.OpenDB"
	db, err := sql.Open("postgres", cfg.DB.Dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Repository{Db: db}, nil
}
