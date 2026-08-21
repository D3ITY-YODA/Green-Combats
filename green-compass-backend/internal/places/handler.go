package places

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
	Create(context.Context, CreateInput) (*Place, error)
	ByID(context.Context, uuid.UUID) (*Place, error)
	Nearby(context.Context, float64, float64, float64, int) ([]WithDistance, error)
	Update(context.Context, uuid.UUID, UpdateInput, Caller) (*Place, error)
}

type Handler struct {
	svc  API
	auth auth.API
}

func NewHandler(svc API, authService auth.API) *Handler {
	return &Handler{svc: svc, auth: authService}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/v1/places/:id", h.byID)
	r.GET("/v1/places/nearby", h.nearby)
	protected := r.Group("/v1/places", auth.Middleware(h.auth))
	protected.POST("", h.create)
	protected.PATCH("/:id", h.update)
}

type placeRequest struct {
	Name      *string  `json:"name"`
	PlaceType *string  `json:"place_type"`
	Lat       *float64 `json:"lat"`
	Lon       *float64 `json:"lon"`
}

func (h *Handler) create(c *gin.Context) {
	var req placeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok || req.Name == nil || req.PlaceType == nil || req.Lat == nil || req.Lon == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, place_type, lat, and lon are required"})
		return
	}
	p, err := h.svc.Create(c.Request.Context(), CreateInput{
		Name: *req.Name, PlaceType: *req.PlaceType, Lat: *req.Lat, Lon: *req.Lon,
		Caller: Caller{UserID: identity.UserID, IsPlatformAdmin: identity.IsPlatformAdmin},
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, placeJSON(p))
}

func (h *Handler) byID(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	p, err := h.svc.ByID(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, placeJSON(p))
}

func (h *Handler) nearby(c *gin.Context) {
	lat, ok := floatQuery(c, "lat", true)
	if !ok {
		return
	}
	lon, ok := floatQuery(c, "lon", true)
	if !ok {
		return
	}
	radius, ok := floatQuery(c, "radius_m", false)
	if !ok {
		return
	}
	limit, ok := intQuery(c, "limit")
	if !ok {
		return
	}
	result, err := h.svc.Nearby(c.Request.Context(), lat, lon, radius, limit)
	if err != nil {
		writeError(c, err)
		return
	}
	response := make([]gin.H, 0, len(result))
	for i := range result {
		item := placeJSON(&result[i].Place)
		item["distance_m"] = result[i].DistanceMeters
		response = append(response, item)
	}
	c.JSON(http.StatusOK, gin.H{"places": response})
}

func (h *Handler) update(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var req placeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	if req.PlaceType != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "place_type cannot be changed"})
		return
	}
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	p, err := h.svc.Update(c.Request.Context(), id, UpdateInput{Name: req.Name, Lat: req.Lat, Lon: req.Lon}, Caller{UserID: identity.UserID, IsPlatformAdmin: identity.IsPlatformAdmin})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, placeJSON(p))
}

func pathID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place id"})
		return uuid.Nil, false
	}
	return id, true
}

func floatQuery(c *gin.Context, key string, required bool) (float64, bool) {
	raw, exists := c.GetQuery(key)
	if !exists {
		if !required {
			return 0, true
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": key + " is required"})
		return 0, false
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + key})
		return 0, false
	}
	return value, true
}

func intQuery(c *gin.Context, key string) (int, bool) {
	raw, exists := c.GetQuery(key)
	if !exists {
		return 0, true
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + key})
		return 0, false
	}
	return value, true
}

func placeJSON(p *Place) gin.H {
	return gin.H{"id": p.ID, "name": p.Name, "place_type": p.PlaceType, "lat": p.Lat, "lon": p.Lon, "external_code": p.ExternalCode, "created_by": p.CreatedBy, "created_at": p.CreatedAt, "updated_at": p.UpdatedAt}
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidData):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "place not found"})
	case errors.Is(err, ErrNotAllowed):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
