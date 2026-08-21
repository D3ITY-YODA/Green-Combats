package notifications

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
	"green-compass-backend/pkg/httpx"
)

type API interface {
	CreateNotification(ctx context.Context, req CreateRequest) (*Notification, error)
	ListNotifications(ctx context.Context, req ListRequest) (*ListResponse, error)
	GetNotification(ctx context.Context, userID, notificationID uuid.UUID) (*Notification, error)
	UpdatePreference(ctx context.Context, req UpdatePreferenceRequest) error
	GetPreferences(ctx context.Context, userID uuid.UUID) ([]Preference, error)
}

type Handler struct {
	svc API
}

func NewHandler(svc API) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all notification-related routes.
func (h *Handler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	group := router.Group("/v1/notifications")
	group.Use(authMiddleware)

	group.GET("", h.ListNotifications)
	group.GET("/:id", h.GetNotification)
	group.GET("/preferences", h.ListPreferences)
	group.PUT("/preferences", h.UpdatePreference)
}

// ListNotifications returns paginated notifications for the authenticated user.
// GET /v1/notifications?page=1&limit=20&status=<optional>
func (h *Handler) ListNotifications(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	page := 1
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 {
			limit = l
		}
	}

	var status *Status
	if v := c.Query("status"); v != "" {
		s := Status(v)
		status = &s
	}

	resp, err := h.svc.ListNotifications(c.Request.Context(), ListRequest{
		UserID: identity.UserID,
		Status: status,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	pagination := httpx.CalculatePaginationMeta(resp.Page, resp.Limit, resp.Total)
	c.JSON(http.StatusOK, httpx.SuccessWithPagination(gin.H{
		"notifications": resp.Notifications,
	}, pagination))
}

// GetNotification returns a single notification by ID.
// GET /v1/notifications/:id
func (h *Handler) GetNotification(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("id", "must be a valid UUID"))
		return
	}

	n, err := h.svc.GetNotification(c.Request.Context(), identity.UserID, id)
	if err != nil {
		httpx.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{"notification": n}))
}

// ListPreferences returns all notification preferences for the authenticated user.
// GET /v1/notifications/preferences
func (h *Handler) ListPreferences(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	prefs, err := h.svc.GetPreferences(c.Request.Context(), identity.UserID)
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{"preferences": prefs}))
}

// UpdatePreference sets a user's notification preference.
// PUT /v1/notifications/preferences
func (h *Handler) UpdatePreference(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	var req struct {
		Channel   string `json:"channel"`
		EventType string `json:"event_type"`
		Enabled   bool   `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	if req.Channel == "" || req.EventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channel and event_type are required"})
		return
	}

	err := h.svc.UpdatePreference(c.Request.Context(), UpdatePreferenceRequest{
		UserID:    identity.UserID,
		Channel:   Channel(req.Channel),
		EventType: req.EventType,
		Enabled:   req.Enabled,
	})
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{"message": "preference updated"}))
}
