package notifications

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
)

// --- Mock repository ---

type mockRepo struct {
	prefs       map[string]Preference
	notifications []Notification
	logs        []DeliveryLog
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		prefs: make(map[string]Preference),
	}
}

func (m *mockRepo) UpsertPreference(_ context.Context, p *Preference) error {
	key := string(p.Channel) + ":" + p.EventType
	m.prefs[key] = *p
	return nil
}

func (m *mockRepo) GetPreferences(_ context.Context, userID uuid.UUID) ([]Preference, error) {
	var result []Preference
	for _, p := range m.prefs {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockRepo) IsEnabled(_ context.Context, _ uuid.UUID, channel Channel, eventType string) (bool, error) {
	key := string(channel) + ":" + eventType
	if p, ok := m.prefs[key]; ok {
		return p.Enabled, nil
	}
	return true, nil // default opt-out model
}

func (m *mockRepo) Create(_ context.Context, n *Notification) error {
	n.ID = uuid.New()
	n.CreatedAt = time.Now()
	m.notifications = append(m.notifications, *n)
	return nil
}

func (m *mockRepo) GetByID(_ context.Context, id uuid.UUID) (*Notification, error) {
	for _, n := range m.notifications {
		if n.ID == id {
			return &n, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockRepo) List(_ context.Context, req ListRequest) ([]Notification, int, error) {
	var result []Notification
	for _, n := range m.notifications {
		if n.UserID == req.UserID {
			result = append(result, n)
		}
	}
	total := len(result)
	start := (req.Page - 1) * req.Limit
	if start >= total {
		return []Notification{}, total, nil
	}
	end := start + req.Limit
	if end > total {
		end = total
	}
	return result[start:end], total, nil
}

func (m *mockRepo) ListPending(_ context.Context, limit int) ([]Notification, error) {
	var result []Notification
	for _, n := range m.notifications {
		if n.Status == StatusPending {
			result = append(result, n)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *mockRepo) MarkSent(_ context.Context, id uuid.UUID) error {
	for i := range m.notifications {
		if m.notifications[i].ID == id {
			now := time.Now()
			m.notifications[i].Status = StatusSent
			m.notifications[i].SentAt = &now
			return nil
		}
	}
	return ErrNotFound
}

func (m *mockRepo) MarkFailed(_ context.Context, id uuid.UUID, errMsg string) error {
	for i := range m.notifications {
		if m.notifications[i].ID == id {
			m.notifications[i].Status = StatusFailed
			m.notifications[i].ErrorMessage = &errMsg
			return nil
		}
	}
	return ErrNotFound
}

func (m *mockRepo) LogDelivery(_ context.Context, log *DeliveryLog) error {
	log.ID = uuid.New()
	log.AttemptedAt = time.Now()
	m.logs = append(m.logs, *log)
	return nil
}

// --- Tests ---

func TestCreateNotification_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo, slog.Default())

	n, err := svc.CreateNotification(context.Background(), CreateRequest{
		UserID:    uuid.New(),
		Channel:   ChannelPush,
		EventType: "content.generated",
		Title:     "New update available",
		Body:      "Rain expected in your area",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected notification, got nil")
	}
	if n.Status != StatusPending {
		t.Errorf("status: got %s, want pending", n.Status)
	}
	if len(repo.notifications) != 1 {
		t.Errorf("stored count: got %d, want 1", len(repo.notifications))
	}
}

func TestCreateNotification_EmptyTitle(t *testing.T) {
	svc := NewService(newMockRepo(), slog.Default())

	_, err := svc.CreateNotification(context.Background(), CreateRequest{
		UserID:    uuid.New(),
		Channel:   ChannelPush,
		EventType: "test",
		Body:      "body",
	})
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestCreateNotification_EmptyBody(t *testing.T) {
	svc := NewService(newMockRepo(), slog.Default())

	_, err := svc.CreateNotification(context.Background(), CreateRequest{
		UserID:    uuid.New(),
		Channel:   ChannelPush,
		EventType: "test",
		Title:     "title",
	})
	if err == nil {
		t.Fatal("expected error for empty body")
	}
}

func TestCreateNotification_OptedOut(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo, slog.Default())

	userID := uuid.New()
	// Opt out of push notifications for content.generated
	repo.prefs["push:content.generated"] = Preference{
		UserID:    userID,
		Channel:   ChannelPush,
		EventType: "content.generated",
		Enabled:   false,
	}

	n, err := svc.CreateNotification(context.Background(), CreateRequest{
		UserID:    userID,
		Channel:   ChannelPush,
		EventType: "content.generated",
		Title:     "test",
		Body:      "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != nil {
		t.Fatal("expected nil notification (user opted out)")
	}
	if len(repo.notifications) != 0 {
		t.Error("should not have stored notification for opted-out user")
	}
}

func TestListNotifications(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo, slog.Default())
	userID := uuid.New()

	// Create 3 notifications
	for i := 0; i < 3; i++ {
		_, _ = svc.CreateNotification(context.Background(), CreateRequest{
			UserID:    userID,
			Channel:   ChannelPush,
			EventType: "test",
			Title:     "title",
			Body:      "body",
		})
	}

	resp, err := svc.ListNotifications(context.Background(), ListRequest{
		UserID: userID,
		Page:   1,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 3 {
		t.Errorf("total: got %d, want 3", resp.Total)
	}
	if len(resp.Notifications) != 3 {
		t.Errorf("notifications count: got %d, want 3", len(resp.Notifications))
	}
}

func TestGetNotification_Found(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo, slog.Default())
	userID := uuid.New()

	n, _ := svc.CreateNotification(context.Background(), CreateRequest{
		UserID:    userID,
		Channel:   ChannelPush,
		EventType: "test",
		Title:     "title",
		Body:      "body",
	})

	got, err := svc.GetNotification(context.Background(), userID, n.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != n.ID {
		t.Errorf("ID mismatch: got %v, want %v", got.ID, n.ID)
	}
}

func TestGetNotification_WrongUser(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo, slog.Default())
	userID := uuid.New()

	n, _ := svc.CreateNotification(context.Background(), CreateRequest{
		UserID:    userID,
		Channel:   ChannelPush,
		EventType: "test",
		Title:     "title",
		Body:      "body",
	})

	_, err := svc.GetNotification(context.Background(), uuid.New(), n.ID)
	if err == nil {
		t.Fatal("expected error for wrong user")
	}
}

func TestMarkSent(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo, slog.Default())

	n, _ := svc.CreateNotification(context.Background(), CreateRequest{
		UserID:    uuid.New(),
		Channel:   ChannelSMS,
		EventType: "alert",
		Title:     "Alert",
		Body:      "Emergency",
	})

	if err := svc.MarkSent(context.Background(), n.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := svc.GetNotification(context.Background(), n.UserID, n.ID)
	if got.Status != StatusSent {
		t.Errorf("status: got %s, want sent", got.Status)
	}
	if got.SentAt == nil {
		t.Error("sent_at should not be nil")
	}
}

func TestMarkFailed(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo, slog.Default())

	n, _ := svc.CreateNotification(context.Background(), CreateRequest{
		UserID:    uuid.New(),
		Channel:   ChannelSMS,
		EventType: "alert",
		Title:     "Alert",
		Body:      "Emergency",
	})

	if err := svc.MarkFailed(context.Background(), n.ID, "provider timeout"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := svc.GetNotification(context.Background(), n.UserID, n.ID)
	if got.Status != StatusFailed {
		t.Errorf("status: got %s, want failed", got.Status)
	}
	if got.ErrorMessage == nil || *got.ErrorMessage != "provider timeout" {
		t.Errorf("error_message: got %v, want 'provider timeout'", got.ErrorMessage)
	}
}

func TestUpdatePreference(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo, slog.Default())
	userID := uuid.New()

	err := svc.UpdatePreference(context.Background(), UpdatePreferenceRequest{
		UserID:    userID,
		Channel:   ChannelSMS,
		EventType: "alert",
		Enabled:   false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it's stored
	prefs, err := svc.GetPreferences(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prefs) != 1 {
		t.Fatalf("prefs count: got %d, want 1", len(prefs))
	}
	if prefs[0].Enabled {
		t.Error("expected enabled=false")
	}
}

func TestLogDelivery(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo, slog.Default())

	err := svc.LogDelivery(context.Background(), &DeliveryLog{
		NotificationID: uuid.New(),
		Attempt:        1,
		Status:         "sent",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.logs) != 1 {
		t.Errorf("logs count: got %d, want 1", len(repo.logs))
	}
}
