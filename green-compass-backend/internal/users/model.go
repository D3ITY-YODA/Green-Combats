package users

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound    = errors.New("user not found")
	ErrEmailTaken  = errors.New("email already registered")
	ErrPhoneTaken  = errors.New("phone number already registered")
	ErrInvalidData = errors.New("invalid user data")
)

type User struct {
	ID              uuid.UUID
	PhoneNumber     *string
	Email           *string
	PasswordHash    string
	DisplayName     string
	Language        string
	IsPlatformAdmin bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (u *User) HasEmail() bool { return u.Email != nil }

func (u *User) HasPhone() bool { return u.PhoneNumber != nil }
