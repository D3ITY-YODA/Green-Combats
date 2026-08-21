package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"green-compass-backend/pkg/security"
)

const (
	minPasswordLen  = 8
	maxPasswordLen  = 128
	maxDisplayName  = 100
	defaultLanguage = "en"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type RegisterInput struct {
	PhoneNumber string
	Email       string
	DisplayName string
	Password    string
	Language    string
}

func (s *Service) Register(ctx context.Context, in RegisterInput) (*User, error) {
	if err := validateRegisterInput(in); err != nil {
		return nil, err
	}

	hash, err := security.HashPassword(in.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	language := in.Language
	if strings.TrimSpace(language) == "" {
		language = defaultLanguage
	}

	u := &User{
		PasswordHash: hash,
		DisplayName:  strings.TrimSpace(in.DisplayName),
		Language:     language,
	}
	if phone := strings.TrimSpace(in.PhoneNumber); phone != "" {
		u.PhoneNumber = &phone
	}
	if email := strings.TrimSpace(in.Email); email != "" {
		u.Email = &email
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) ByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.ByID(ctx, id)
}

func (s *Service) ByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.ByEmail(ctx, email)
}

func (s *Service) ByPhone(ctx context.Context, phone string) (*User, error) {
	return s.repo.ByPhone(ctx, phone)
}

func (s *Service) UpdateProfile(ctx context.Context, id uuid.UUID, displayName, language string) (*User, error) {
	displayName = strings.TrimSpace(displayName)
	language = strings.TrimSpace(language)

	var problems []error
	if displayName == "" {
		problems = append(problems, fmt.Errorf("%w: display name is required", ErrInvalidData))
	}
	if len(displayName) > maxDisplayName {
		problems = append(problems, fmt.Errorf("%w: display name exceeds %d characters", ErrInvalidData, maxDisplayName))
	}
	if language == "" {
		language = defaultLanguage
	}
	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	return s.repo.UpdateProfile(ctx, id, displayName, language)
}

func validateRegisterInput(in RegisterInput) error {
	phone := strings.TrimSpace(in.PhoneNumber)
	email := strings.TrimSpace(in.Email)
	displayName := strings.TrimSpace(in.DisplayName)

	var problems []error
	if phone != "" && email != "" {
		problems = append(problems, fmt.Errorf("%w: provide exactly one of phone number or email, not both", ErrInvalidData))
	}
	if phone == "" && email == "" {
		problems = append(problems, fmt.Errorf("%w: one of phone number or email is required", ErrInvalidData))
	}
	if displayName == "" {
		problems = append(problems, fmt.Errorf("%w: display name is required", ErrInvalidData))
	}
	if len(displayName) > maxDisplayName {
		problems = append(problems, fmt.Errorf("%w: display name exceeds %d characters", ErrInvalidData, maxDisplayName))
	}
	if len(in.Password) < minPasswordLen {
		problems = append(problems, fmt.Errorf("%w: password must be at least %d characters", ErrInvalidData, minPasswordLen))
	}
	if len(in.Password) > maxPasswordLen {
		problems = append(problems, fmt.Errorf("%w: password must be at most %d characters", ErrInvalidData, maxPasswordLen))
	}
	if len(problems) > 0 {
		return errors.Join(problems...)
	}
	return nil
}
