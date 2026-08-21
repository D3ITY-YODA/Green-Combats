package reports

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
)

type API interface {
	ListPending(ctx context.Context, caller Caller, limit int) ([]Report, error)
	Get(ctx context.Context, id uuid.UUID, caller Caller) (*Report, error)
	Verify(ctx context.Context, id uuid.UUID, caller Caller, status string) (*Report, error)
}

type Handler struct {
	svc  API
	auth auth.API
}

func NewHandler(svc API, authService auth.API) *Handler {
	return &Handler{svc: svc, auth: authService}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	inst := r.Group("/v1/reports", auth.Middleware(h.auth))
	inst.GET("/pending", h.listPending)
	inst.GET("/:id", h.get)
	inst.PATCH("/:id/verify", h.verify)
}

func (h *Handler) listPending(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
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
		response = append(response, reportJSON(&item))
	}
	c.JSON(http.StatusOK, gin.H{"reports": response})
}

func (h *Handler) get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	report, err := h.svc.Get(c.Request.Context(), id, Caller{UserID: identity.UserID, IsPlatformAdmin: identity.IsPlatformAdmin})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, reportJSON(report))
}

type verifyRequest struct {
	Status string `json:"status"`
}

func (h *Handler) verify(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var req verifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	if req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}

	report, err := h.svc.Verify(c.Request.Context(), id, Caller{UserID: identity.UserID, IsPlatformAdmin: identity.IsPlatformAdmin}, req.Status)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, reportJSON(report))
}

func pathID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report id"})
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

func reportJSON(r *Report) gin.H {
	item := gin.H{
		"id":          r.ID,
		"reporter_id": r.ReporterID,
		"category":    r.Category,
		"description": r.Description,
		"status":      r.Status,
		"created_at":  r.CreatedAt,
		"updated_at":  r.UpdatedAt,
	}
	if r.PlaceID != nil {
		item["place_id"] = r.PlaceID
	}
	if r.Lat != nil {
		item["lat"] = r.Lat
	}
	if r.Lon != nil {
		item["lon"] = r.Lon
	}
	if r.PhotoObjectKey != nil {
		item["photo_object_key"] = r.PhotoObjectKey
	}
	if r.VerifiedBy != nil {
		item["verified_by"] = r.VerifiedBy
	}
	if r.VerifiedAt != nil {
		item["verified_at"] = r.VerifiedAt
	}
	return item
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidData):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
	case errors.Is(err, ErrNotAllowed):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
