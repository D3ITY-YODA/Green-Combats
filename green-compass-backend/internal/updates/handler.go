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
	GetByID(ctx context.Context, updateID uuid.UUID) (*Update, error)
	Acknowledge(ctx context.Context, updateID, userID uuid.UUID) error
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
	group.GET("/updates/:update_id", h.GetUpdate)
	group.POST("/updates/:update_id/acknowledge", h.AcknowledgeUpdate)

	group.GET("/explore", h.ExploreIndicators)
	group.GET("/places/:id/explore", h.ExploreByPlace)
	group.GET("/places/:id/local-outlook", h.LocalOutlook)
	group.GET("/places/:id/seasonal", h.Seasonal)
	group.GET("/places/:id/water", h.Water)
	group.GET("/places/:id/land-ecosystems", h.LandEcosystems)
	group.GET("/places/:id/food-agriculture", h.FoodAgriculture)
	group.GET("/places/:id/community", h.Community)
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
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.SuccessWithPagination(gin.H{
		"updates": resp.Updates,
	}, pagination, requestID))
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
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.SuccessWithPagination(gin.H{
		"indicators": resp.Indicators,
	}, pagination, requestID))
}

// GetUpdate returns a single update by ID
// GET /v1/updates/{update_id}
func (h *Handler) GetUpdate(c *gin.Context) {
	if _, ok := auth.IdentityFrom(c.Request.Context()); !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	updateIDStr := c.Param("update_id")
	updateID, err := uuid.Parse(updateIDStr)
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("update_id", "must be a valid UUID"))
		return
	}

	update, err := h.svc.GetByID(c.Request.Context(), updateID)
	if err != nil {
		httpx.HandleError(c, httpx.ErrNotFound)
		return
	}

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(update, requestID))
}

// AcknowledgeUpdate marks an update as acknowledged by the user
// POST /v1/updates/{update_id}/acknowledge
func (h *Handler) AcknowledgeUpdate(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	updateIDStr := c.Param("update_id")
	updateID, err := uuid.Parse(updateIDStr)
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("update_id", "must be a valid UUID"))
		return
	}

	if err := h.svc.Acknowledge(c.Request.Context(), updateID, identity.UserID); err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"message": "update acknowledged"}, requestID))
}

// ExploreByPlace returns indicators for a specific place
// GET /v1/places/{place_id}/explore
func (h *Handler) ExploreByPlace(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	placeIDStr := c.Param("id")
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("place_id", "must be a valid UUID"))
		return
	}

	userID := identity.UserID
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

	var category *string
	if cat := c.Query("category"); cat != "" {
		category = &cat
	}

	resp, err := h.svc.Explore(c.Request.Context(), ExploreRequest{
		UserID:   userID,
		Category: category,
		PlaceID:  &placeID,
		Page:     page,
		Limit:    limit,
	})
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	pagination := httpx.CalculatePaginationMeta(resp.Page, resp.Limit, resp.Total)
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.SuccessWithPagination(gin.H{
		"indicators": resp.Indicators,
	}, pagination, requestID))
}

// LocalOutlook returns local outlook indicators for a place
// GET /v1/places/{place_id}/local-outlook
func (h *Handler) LocalOutlook(c *gin.Context) {
	h.exploreByTopic(c, "local_outlook")
}

// Seasonal returns seasonal information for a place
// GET /v1/places/{place_id}/seasonal
func (h *Handler) Seasonal(c *gin.Context) {
	h.exploreByTopic(c, "seasonal_information")
}

// Water returns water outlook for a place
// GET /v1/places/{place_id}/water
func (h *Handler) Water(c *gin.Context) {
	h.exploreByTopic(c, "water_outlook")
}

// LandEcosystems returns land and ecosystems for a place
// GET /v1/places/{place_id}/land-ecosystems
func (h *Handler) LandEcosystems(c *gin.Context) {
	h.exploreByTopic(c, "land_and_ecosystems")
}

// FoodAgriculture returns food and agriculture for a place
// GET /v1/places/{place_id}/food-agriculture
func (h *Handler) FoodAgriculture(c *gin.Context) {
	h.exploreByTopic(c, "food_and_agriculture")
}

// Community returns community updates for a place
// GET /v1/places/{place_id}/community
func (h *Handler) Community(c *gin.Context) {
	h.exploreByTopic(c, "community_updates")
}

func (h *Handler) exploreByTopic(c *gin.Context, topicKey string) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	placeIDStr := c.Param("id")
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("place_id", "must be a valid UUID"))
		return
	}

	userID := identity.UserID
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

	resp, err := h.svc.Explore(c.Request.Context(), ExploreRequest{
		UserID:   userID,
		Category: &topicKey,
		PlaceID:  &placeID,
		Page:     page,
		Limit:    limit,
	})
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	pagination := httpx.CalculatePaginationMeta(resp.Page, resp.Limit, resp.Total)
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.SuccessWithPagination(gin.H{
		"indicators": resp.Indicators,
	}, pagination, requestID))
}
