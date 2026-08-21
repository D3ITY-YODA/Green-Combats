package updates

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
	List(ctx context.Context, req ListRequest) (*ListResponse, error)
	Explore(ctx context.Context, req ExploreRequest) (*ExploreResponse, error)
}

type Handler struct {
	svc API
}

func NewHandler(svc API) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all update-related routes
func (h *Handler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	group := router.Group("/v1")
	group.Use(authMiddleware)

	group.GET("/updates", h.ListUpdates)
	group.GET("/explore", h.ExploreIndicators)
}

// ListUpdates retrieves paginated updates for user's saved places
// GET /v1/updates?page=1&limit=20&place_id=<optional-uuid>
func (h *Handler) ListUpdates(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}
	userID := identity.UserID

	// Parse pagination
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Optional place filter
	var placeID *uuid.UUID
	if placeIDStr := c.Query("place_id"); placeIDStr != "" {
		if pID, err := uuid.Parse(placeIDStr); err == nil {
			placeID = &pID
		}
	}

	resp, err := h.svc.List(c.Request.Context(), ListRequest{
		UserID:  userID,
		PlaceID: placeID,
		Page:    page,
		Limit:   limit,
	})
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	pagination := httpx.CalculatePaginationMeta(resp.Page, resp.Limit, resp.Total)
	c.JSON(http.StatusOK, httpx.SuccessWithPagination(gin.H{
		"updates": resp.Updates,
	}, pagination))
}

// ExploreIndicators retrieves indicators for exploration/discovery
// GET /v1/explore?category=<optional>&place_id=<optional>&page=1&limit=50
func (h *Handler) ExploreIndicators(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}
	userID := identity.UserID

	// Parse pagination
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Optional filters
	var category *string
	if cat := c.Query("category"); cat != "" {
		category = &cat
	}

	var placeID *uuid.UUID
	if placeIDStr := c.Query("place_id"); placeIDStr != "" {
		if pID, err := uuid.Parse(placeIDStr); err == nil {
			placeID = &pID
		}
	}

	resp, err := h.svc.Explore(c.Request.Context(), ExploreRequest{
		UserID:   userID,
		Category: category,
		PlaceID:  placeID,
		Page:     page,
		Limit:    limit,
	})
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	pagination := httpx.CalculatePaginationMeta(resp.Page, resp.Limit, resp.Total)
	c.JSON(http.StatusOK, httpx.SuccessWithPagination(gin.H{
		"indicators": resp.Indicators,
	}, pagination))
}
