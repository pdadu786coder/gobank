package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrExiredToken = errors.New("token has expired")

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"username"`
	IssuedAt  time.Time `json:"issues_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type JWTPayload struct {
	Payload
	jwt.RegisteredClaims
}

// NewPayload creates a new token payload
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		UserName:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func NewJWTPayload(username string, duration time.Duration) (*JWTPayload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return nil, err
	}
	jwtPayload := &JWTPayload{Payload: *payload}
	return jwtPayload, nil
}

func (payload *JWTPayload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExiredToken
	}
	return nil
}
