package audit

import (
	"time"

	"github.com/google/uuid"
)

// Action represents the type of auditable action.
type Action string

const (
	ActionObservationVerified Action = "observation_verified"
	ActionObservationRejected Action = "observation_rejected"
	ActionContentPublished   Action = "content_published"
	ActionUserCreated        Action = "user_created"
	ActionUserDeactivated    Action = "user_deactivated"
	ActionOrgCreated         Action = "org_created"
	ActionOrgMemberAdded     Action = "org_member_added"
	ActionOrgMemberRemoved   Action = "org_member_removed"
	ActionSourceEnabled      Action = "source_enabled"
	ActionSourceDisabled     Action = "source_disabled"
	ActionConfigChanged      Action = "config_changed"
)

// Entry represents a single audit log record.
type Entry struct {
	ID           uuid.UUID
	ActorID      *uuid.UUID
	ActorOrgID   *uuid.UUID
	Action       Action
	ResourceType string
	ResourceID   *uuid.UUID
	Details      map[string]interface{}
	IPAddress    *string
	CreatedAt    time.Time
}

// CreateRequest is the input for creating an audit entry.
type CreateRequest struct {
	ActorID      *uuid.UUID
	ActorOrgID   *uuid.UUID
	Action       Action
	ResourceType string
	ResourceID   *uuid.UUID
	Details      map[string]interface{}
	IPAddress    *string
}

// ListRequest is the input for listing audit entries.
type ListRequest struct {
	ActorID      *uuid.UUID
	ActorOrgID   *uuid.UUID
	Action       *Action
	ResourceType *string
	ResourceID   *uuid.UUID
	Page         int
	Limit        int
}

// ListResponse is the output for listing audit entries.
type ListResponse struct {
	Entries []Entry
	Page    int
	Limit   int
	Total   int
}
