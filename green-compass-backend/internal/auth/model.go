package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenReuse         = errors.New("refresh token reuse detected")
)

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

type SessionUser struct {
	ID              uuid.UUID
	DisplayName     string
	Email           *string
	PhoneNumber     *string
	Language        string
	IsPlatformAdmin bool
}
