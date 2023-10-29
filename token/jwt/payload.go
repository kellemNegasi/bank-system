package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Different types of error used in the token package.
var (
	ErrInvalidToken = errors.New("token is invalid.")
	ErrExpiredToken = errors.New("token is expired.")
)

// Payload contains the data of a token.
type JWTPayload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`

	jwt.RegisteredClaims
}

// New returns a new payload object with username and duration.
func NewPayLoad(userName string, duration time.Duration) (*JWTPayload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	rc := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
	}

	return &JWTPayload{
		ID:               id,
		Username:         userName,
		RegisteredClaims: rc,
	}, nil
}
