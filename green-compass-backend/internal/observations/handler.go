package observations

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
	"green-compass-backend/pkg/httpx"
)

type API interface {
	Create(ctx context.Context, in CreateInput) (*Observation, error)
	ByID(ctx context.Context, id uuid.UUID, caller Caller) (*Observation, error)
	ListMine(ctx context.Context, userID uuid.UUID, limit int) ([]Observation, error)
	ListPending(ctx context.Context, caller Caller, limit int) ([]Observation, error)
	Verify(ctx context.Context, id uuid.UUID, caller Caller, status string) (*Observation, error)
}

type Handler struct {
	svc  API
	auth auth.API
}

func NewHandler(svc API, authService auth.API) *Handler {
	return &Handler{svc: svc, auth: authService}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	group := r.Group("/v1/observations", auth.Middleware(h.auth))
	group.POST("", h.create)
	group.GET("", h.listMine)
	group.GET("/:id", h.byID)
	inst := group.Group("/institutional")
	inst.GET("/pending", h.listPending)
	inst.PATCH("/:id/verify", h.verify)
}

type createRequest struct {
	PlaceID        *uuid.UUID `json:"place_id"`
	Lat            *float64   `json:"lat"`
	Lon            *float64   `json:"lon"`
	Category       string     `json:"category"`
	Description    string     `json:"description"`
	PhotoObjectKey *string    `json:"photo_object_key"`
}

func (h *Handler) create(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.HandleError(c, httpx.InvalidParam("body", "invalid JSON"))
		return
	}

	o, err := h.svc.Create(c.Request.Context(), CreateInput{
		PlaceID:        req.PlaceID,
		Lat:            req.Lat,
		Lon:            req.Lon,
		Category:       req.Category,
		Description:    req.Description,
		PhotoObjectKey: req.PhotoObjectKey,
		Caller:         Caller{UserID: identity.UserID, IsPlatformAdmin: identity.IsPlatformAdmin},
	})
	if err != nil {
		writeError(c, err)
		return
	}
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusCreated, httpx.Success(observationJSON(o), requestID))
}

func (h *Handler) listMine(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	limit, ok := intQuery(c, "limit")
	if !ok {
		return
	}

	items, err := h.svc.ListMine(c.Request.Context(), identity.UserID, limit)
	if err != nil {
		writeError(c, err)
		return
	}
	response := make([]gin.H, 0, len(items))
	for _, item := range items {
		response = append(response, observationJSON(&item))
	}
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"observations": response}, requestID))
}

func (h *Handler) byID(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	o, err := h.svc.ByID(c.Request.Context(), id, Caller{UserID: identity.UserID, IsPlatformAdmin: identity.IsPlatformAdmin})
	if err != nil {
		writeError(c, err)
		return
	}
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(observationJSON(o), requestID))
}

func (h *Handler) listPending(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	limit, ok := intQuery(c, "limit")
	if !ok {
		return
	}

	items, err := h.svc.ListPending(c.Request.Context(), Caller{UserID: identity.UserID, IsPlatformAdmin: identity.IsPlatformAdmin}, limit)
	if err != nil {
		writeError(c, err)
		return
	}
	response := make([]gin.H, 0, len(items))
	for _, item := range items {
		response = append(response, observationJSON(&item))
	}
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"observations": response}, requestID))
}

type verifyRequest struct {
	Status string  `json:"status"`
	Notes  *string `json:"notes"`
}

func (h *Handler) verify(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	var req verifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.HandleError(c, httpx.InvalidParam("body", "invalid JSON"))
		return
	}
	if req.Status == "" {
		httpx.HandleError(c, httpx.InvalidParam("status", "status is required"))
		return
	}

	o, err := h.svc.Verify(c.Request.Context(), id, Caller{UserID: identity.UserID, IsPlatformAdmin: identity.IsPlatformAdmin}, req.Status)
	if err != nil {
		writeError(c, err)
		return
	}
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(observationJSON(o), requestID))
}

func pathID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid observation id"})
		return uuid.Nil, false
	}
	return id, true
}

func intQuery(c *gin.Context, key string) (int, bool) {
	raw, exists := c.GetQuery(key)
	if !exists {
		return 25, true
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + key})
		return 0, false
	}
	if value <= 0 {
		value = 25
	}
	if value > 100 {
		value = 100
	}
	return value, true
}

func observationJSON(o *Observation) gin.H {
	item := gin.H{
		"id":          o.ID,
		"reporter_id": o.ReporterID,
		"category":    o.Category,
		"description": o.Description,
		"status":      o.Status,
		"created_at":  o.CreatedAt,
		"updated_at":  o.UpdatedAt,
	}
	if o.PlaceID != nil {
		item["place_id"] = o.PlaceID
	}
	if o.Lat != nil {
		item["lat"] = o.Lat
	}
	if o.Lon != nil {
		item["lon"] = o.Lon
	}
	if o.PhotoObjectKey != nil {
		item["photo_object_key"] = o.PhotoObjectKey
	}
	if o.VerifiedBy != nil {
		item["verified_by"] = o.VerifiedBy
	}
	if o.VerifiedAt != nil {
		item["verified_at"] = o.VerifiedAt
	}
	return item
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidData):
		httpx.HandleError(c, httpx.ErrBadRequest)
	case errors.Is(err, ErrNotFound):
		httpx.HandleError(c, httpx.ErrNotFound)
	case errors.Is(err, ErrNotAllowed):
		httpx.HandleError(c, httpx.ErrForbidden)
	default:
		httpx.HandleError(c, httpx.ErrInternal)
	}
}
