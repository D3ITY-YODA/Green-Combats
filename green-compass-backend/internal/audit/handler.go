package audit

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
	Record(ctx context.Context, req CreateRequest) (*Entry, error)
	List(ctx context.Context, req ListRequest) (*ListResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Entry, error)
}

type Handler struct {
	svc API
}

func NewHandler(svc API) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all audit-related routes.
func (h *Handler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	group := router.Group("/v1/audit")
	group.Use(authMiddleware)

	group.GET("", h.ListEntries)
	group.GET("/:id", h.GetEntry)
}

// ListEntries returns paginated audit entries.
// GET /v1/audit?page=1&limit=20&action=<optional>&resource_type=<optional>
func (h *Handler) ListEntries(c *gin.Context) {
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

	var action *Action
	if v := c.Query("action"); v != "" {
		a := Action(v)
		action = &a
	}

	var resourceType *string
	if v := c.Query("resource_type"); v != "" {
		resourceType = &v
	}

	// Non-admin users can only see their own audit entries
	req := ListRequest{
		Action:       action,
		ResourceType: resourceType,
		Page:         page,
		Limit:        limit,
	}

	if !identity.IsPlatformAdmin {
		req.ActorID = &identity.UserID
	}

	resp, err := h.svc.List(c.Request.Context(), req)
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	pagination := httpx.CalculatePaginationMeta(resp.Page, resp.Limit, resp.Total)
	c.JSON(http.StatusOK, httpx.SuccessWithPagination(gin.H{
		"entries": resp.Entries,
	}, pagination))
}

// GetEntry returns a single audit entry by ID.
// GET /v1/audit/:id
func (h *Handler) GetEntry(c *gin.Context) {
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

	entry, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		httpx.HandleError(c, httpx.ErrNotFound)
		return
	}

	// Non-admin users can only view their own audit entries
	if !identity.IsPlatformAdmin && entry.ActorID != nil && *entry.ActorID != identity.UserID {
		httpx.HandleError(c, httpx.ErrForbidden)
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{"entry": entry}))
}
