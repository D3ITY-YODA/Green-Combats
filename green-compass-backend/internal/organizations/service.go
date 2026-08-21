package organizations

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"green-compass-backend/pkg/database"
)

type Service struct {
	pool *database.Pool
	repo *Repository
}

func NewService(pool *database.Pool, repo *Repository) *Service {
	return &Service{pool: pool, repo: repo}
}

type CreateInput struct {
	Name    string
	OrgType string
	OwnerID uuid.UUID
}

func (s *Service) CreateWithOwner(ctx context.Context, in CreateInput) (*Organization, error) {
	name := strings.TrimSpace(in.Name)

	var problems []error
	if name == "" {
		problems = append(problems, fmt.Errorf("%w: name is required", ErrInvalidData))
	}
	if len(name) > 200 {
		problems = append(problems, fmt.Errorf("%w: name exceeds 200 characters", ErrInvalidData))
	}
	if !IsValidType(in.OrgType) {
		problems = append(problems, fmt.Errorf("%w: unknown organization type %q", ErrInvalidData, in.OrgType))
	}
	if in.OwnerID == uuid.Nil {
		problems = append(problems, fmt.Errorf("%w: owner is required", ErrInvalidData))
	}
	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}

	org := &Organization{Name: name, OrgType: in.OrgType}
	err := database.WithTx(ctx, s.pool, func(ctx context.Context, tx pgx.Tx) error {
		if err := s.repo.WithTx(tx).Create(ctx, org); err != nil {
			return err
		}
		return s.repo.WithTx(tx).AddMember(ctx, &Member{
			OrgID:  org.ID,
			UserID: in.OwnerID,
			Role:   RoleAdmin,
		})
	})
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (s *Service) ByID(ctx context.Context, id uuid.UUID) (*Organization, error) {
	return s.repo.ByID(ctx, id)
}

func (s *Service) SetStatus(ctx context.Context, orgID uuid.UUID, status string) (*Organization, error) {
	if !IsValidStatus(status) {
		return nil, fmt.Errorf("%w: unknown status %q", ErrInvalidData, status)
	}
	return s.repo.SetStatus(ctx, orgID, status)
}

func (s *Service) AddMember(ctx context.Context, orgID, userID uuid.UUID, role string) error {
	if !IsValidRole(role) {
		return fmt.Errorf("%w: unknown role %q", ErrInvalidData, role)
	}
	return s.repo.AddMember(ctx, &Member{OrgID: orgID, UserID: userID, Role: role})
}

func (s *Service) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	return s.repo.RemoveMember(ctx, orgID, userID)
}

func (s *Service) MemberRole(ctx context.Context, orgID, userID uuid.UUID) (string, error) {
	return s.repo.MemberRole(ctx, orgID, userID)
}
