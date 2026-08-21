package notifications

import (
	"time"

	"github.com/google/uuid"
)

// Channel represents a notification delivery channel.
type Channel string

const (
	ChannelPush  Channel = "push"
	ChannelSMS   Channel = "sms"
	ChannelUSSD  Channel = "ussd"
	ChannelEmail Channel = "email"
)

// Status represents the delivery status of a notification.
type Status string

const (
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// Preference stores a user's opt-in/opt-out per channel and event type.
type Preference struct {
	UserID    uuid.UUID
	Channel   Channel
	EventType string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Notification is a single notification record to be delivered.
type Notification struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Channel      Channel
	EventType    string
	Title        string
	Body         string
	Metadata     map[string]interface{}
	Status       Status
	CreatedAt    time.Time
	SentAt       *time.Time
	ErrorMessage *string
}

// DeliveryLog records a single delivery attempt.
type DeliveryLog struct {
	ID               uuid.UUID
	NotificationID   uuid.UUID
	Attempt          int
	Status           string
	ProviderResponse *string
	AttemptedAt      time.Time
}

// CreateRequest is the input for creating a notification.
type CreateRequest struct {
	UserID    uuid.UUID
	Channel   Channel
	EventType string
	Title     string
	Body      string
	Metadata  map[string]interface{}
}

// ListRequest is the input for listing notifications.
type ListRequest struct {
	UserID uuid.UUID
	Status *Status
	Page   int
	Limit  int
}

// ListResponse is the output for listing notifications.
type ListResponse struct {
	Notifications []Notification
	Page          int
	Limit         int
	Total         int
}

// UpdatePreferenceRequest is the input for updating a preference.
type UpdatePreferenceRequest struct {
	UserID    uuid.UUID
	Channel   Channel
	EventType string
	Enabled   bool
}
