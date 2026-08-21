package organizations

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
	CreateWithOwner(ctx context.Context, in CreateInput) (*Organization, error)
	ByID(ctx context.Context, id uuid.UUID) (*Organization, error)
	SetStatus(ctx context.Context, orgID uuid.UUID, status string) (*Organization, error)
	AddMember(ctx context.Context, orgID, userID uuid.UUID, role string) error
	RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error
}

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all organization-related routes.
func (h *Handler) RegisterRoutes(rg gin.IRouter) {
	group := rg.Group("/v1/organizations")
	group.POST("", h.Create)
	group.GET("", h.List)
	group.GET("/:id", h.GetByID)
	group.PUT("/:id/status", h.SetStatus)
	group.POST("/:id/members", h.AddMember)
	group.DELETE("/:id/members/:user_id", h.RemoveMember)
}

// Create creates a new organization.
// POST /v1/organizations
func (h *Handler) Create(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	var req struct {
		Name    string `json:"name"`
		OrgType string `json:"org_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	org, err := h.svc.CreateWithOwner(c.Request.Context(), CreateInput{
		Name:    req.Name,
		OrgType: req.OrgType,
		OwnerID: identity.UserID,
	})
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusCreated, httpx.Success(gin.H{"organization": org}, requestID))
}

// List returns organizations (placeholder for full implementation).
// GET /v1/organizations?page=1&limit=20
func (h *Handler) List(c *gin.Context) {
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

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{
		"organizations": []interface{}{},
		"page":          page,
		"limit":         limit,
	}, requestID))
}

// GetByID returns a single organization.
// GET /v1/organizations/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("id", "must be a valid UUID"))
		return
	}

	org, err := h.svc.ByID(c.Request.Context(), id)
	if err != nil {
		httpx.HandleError(c, httpx.ErrNotFound)
		return
	}

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"organization": org}, requestID))
}

// SetStatus updates an organization's status (admin only).
// PUT /v1/organizations/:id/status
func (h *Handler) SetStatus(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}
	if !identity.IsPlatformAdmin {
		httpx.HandleError(c, httpx.ErrForbidden)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("id", "must be a valid UUID"))
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.HandleError(c, httpx.InvalidParam("body", "invalid JSON"))
		return
	}

	org, err := h.svc.SetStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"organization": org}, requestID))
}

// AddMember adds a user to an organization.
var req struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.HandleError(c, httpx.InvalidParam("body", "invalid JSON"))
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("user_id", "invalid user_id"))
		return
	}

	if err := h.svc.AddMember(c.Request.Context(), id, userID, req.Role); err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"message": "member added"}, requestID))
}

// RemoveMember removes a user from an organization.
// DELETE /v1/organizations/:id/members/:user_id
func (h *Handler) RemoveMember(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("id", "must be a valid UUID"))
		return
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("user_id", "must be a valid UUID"))
		return
	}

	if err := h.svc.RemoveMember(c.Request.Context(), orgID, userID); err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"message": "member removed"}, requestID))
}
