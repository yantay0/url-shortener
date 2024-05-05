package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/yantay0/url-shortener/internal/model"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type TokenStorage struct {
	DB *sql.DB
}

func (s TokenStorage) New(userID int64, ttl time.Duration, scope string) (*model.Token, error) {
	token, err := model.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = s.Insert(token)
	return token, err
}

func (s TokenStorage) Insert(token *model.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`
	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := s.DB.ExecContext(ctx, query, args...)
	return err
}

func (s TokenStorage) DeleteAllForUser(scope string, userID int64) error {
	query := `
	DELETE FROM tokens
	WHERE scope = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := s.DB.ExecContext(ctx, query, scope, userID)
	return err
}
