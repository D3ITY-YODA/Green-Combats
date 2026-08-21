package organizations

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("organization not found")
	ErrNotMember     = errors.New("not a member of this organization")
	ErrAlreadyMember = errors.New("already a member of this organization")
	ErrNameTaken     = errors.New("organization name already taken")
	ErrInvalidData   = errors.New("invalid organization data")
)

const (
	StatusPending  = "pending"
	StatusVerified = "verified"
	StatusRejected = "rejected"

	RoleAdmin    = "admin"
	RoleReviewer = "reviewer"
	RoleMember   = "member"

	TypeGovernment     = "government"
	TypeWaterAuthority = "water_authority"
	TypeNGO            = "ngo"
	TypeOther          = "other"
)

type Organization struct {
	ID        uuid.UUID
	Name      string
	OrgType   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Member struct {
	OrgID    uuid.UUID
	UserID   uuid.UUID
	Role     string
	JoinedAt time.Time
}

func IsValidStatus(s string) bool {
	switch s {
	case StatusPending, StatusVerified, StatusRejected:
		return true
	}
	return false
}

func IsValidRole(r string) bool {
	switch r {
	case RoleAdmin, RoleReviewer, RoleMember:
		return true
	}
	return false
}

func IsValidType(t string) bool {
	switch t {
	case TypeGovernment, TypeWaterAuthority, TypeNGO, TypeOther:
		return true
	}
	return false
}
