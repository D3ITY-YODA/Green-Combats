package projects

import (
	"time"

	"github.com/google/uuid"
)

// Project represents an institutional project or program.
type Project struct {
	ID          uuid.UUID
	OrgID       uuid.UUID
	Name        string
	Description string
	Status      string // active, archived
	StartDate   *time.Time
	EndDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateRequest is the input for creating a project.
type CreateRequest struct {
	OrgID       uuid.UUID
	Name        string
	Description string
	StartDate   *time.Time
	EndDate     *time.Time
}

// UpdateRequest is the input for updating a project.
type UpdateRequest struct {
	Name        *string
	Description *string
	Status      *string
	StartDate   *time.Time
	EndDate     *time.Time
}

// ListRequest is the input for listing projects.
type ListRequest struct {
	OrgID uuid.UUID
	Page  int
	Limit int
}

// ListResponse is the output for listing projects.
type ListResponse struct {
	Projects []Project
	Page     int
	Limit    int
	Total    int
}
