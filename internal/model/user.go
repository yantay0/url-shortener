package model

import (
	"errors"
	"time"

	"github.com/yantay0/url-shortener/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	Plaintext *string
	Hash      []byte
}

// calculate the bcrypt hash of Plaintext password, stores both
func (p *password) Set(plaintexPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintexPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintexPassword
	p.Hash = hash

	return nil
}

// check if provided Plaintext password matches to hashed one
func (p *password) Matches(plaintexPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintexPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.Plaintext != nil {
		ValidatePassword(v, *user.Password.Plaintext)
	}

	// If the password hash is ever nil
	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
