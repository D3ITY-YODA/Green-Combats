package notifications

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"green-compass-backend/pkg/httpx"
)

type repo interface {
	UpsertPreference(ctx context.Context, p *Preference) error
	GetPreferences(ctx context.Context, userID uuid.UUID) ([]Preference, error)
	IsEnabled(ctx context.Context, userID uuid.UUID, channel Channel, eventType string) (bool, error)
	Create(ctx context.Context, n *Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*Notification, error)
	List(ctx context.Context, req ListRequest) ([]Notification, int, error)
	ListPending(ctx context.Context, limit int) ([]Notification, error)
	MarkSent(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error
	LogDelivery(ctx context.Context, log *DeliveryLog) error
}

type Service struct {
	repo   repo
	logger *slog.Logger
}

func NewService(repo repo, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// CreateNotification creates a new pending notification record.
func (s *Service) CreateNotification(ctx context.Context, req CreateRequest) (*Notification, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("%w: title is required", httpx.ErrBadRequest)
	}
	if req.Body == "" {
		return nil, fmt.Errorf("%w: body is required", httpx.ErrBadRequest)
	}

	// Check if user has this channel+event enabled (opt-out model)
	enabled, err := s.repo.IsEnabled(ctx, req.UserID, req.Channel, req.EventType)
	if err != nil {
		return nil, fmt.Errorf("check preference: %w", err)
	}
	if !enabled {
		s.logger.Info("notification skipped: user opted out",
			"user_id", req.UserID,
			"channel", req.Channel,
			"event_type", req.EventType,
		)
		return nil, nil
	}

	n := &Notification{
		UserID:    req.UserID,
		Channel:   req.Channel,
		EventType: req.EventType,
		Title:     req.Title,
		Body:      req.Body,
		Metadata:  req.Metadata,
		Status:    StatusPending,
	}

	if err := s.repo.Create(ctx, n); err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}

	s.logger.Info("notification created",
		"notification_id", n.ID,
		"user_id", n.UserID,
		"channel", n.Channel,
		"event_type", n.EventType,
	)

	return n, nil
}

// ListNotifications returns paginated notifications for a user.
func (s *Service) ListNotifications(ctx context.Context, req ListRequest) (*ListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	notifications, total, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}

	return &ListResponse{
		Notifications: notifications,
		Page:          req.Page,
		Limit:         req.Limit,
		Total:         total,
	}, nil
}

// GetNotification returns a single notification by ID.
func (s *Service) GetNotification(ctx context.Context, userID, notificationID uuid.UUID) (*Notification, error) {
	n, err := s.repo.GetByID(ctx, notificationID)
	if err != nil {
		return nil, err
	}
	if n.UserID != userID {
		return nil, httpx.ErrNotFound
	}
	return n, nil
}

// MarkSent marks a notification as successfully sent.
func (s *Service) MarkSent(ctx context.Context, id uuid.UUID) error {
	return s.repo.MarkSent(ctx, id)
}

// MarkFailed marks a notification as failed with an error message.
func (s *Service) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return s.repo.MarkFailed(ctx, id, errMsg)
}

// LogDelivery records a delivery attempt.
func (s *Service) LogDelivery(ctx context.Context, log *DeliveryLog) error {
	return s.repo.LogDelivery(ctx, log)
}

// --- Preferences ---

// UpdatePreference sets a user's notification preference for a channel and event type.
func (s *Service) UpdatePreference(ctx context.Context, req UpdatePreferenceRequest) error {
	p := &Preference{
		UserID:    req.UserID,
		Channel:   req.Channel,
		EventType: req.EventType,
		Enabled:   req.Enabled,
	}
	return s.repo.UpsertPreference(ctx, p)
}

// GetPreferences returns all notification preferences for a user.
func (s *Service) GetPreferences(ctx context.Context, userID uuid.UUID) ([]Preference, error) {
	return s.repo.GetPreferences(ctx, userID)
}
