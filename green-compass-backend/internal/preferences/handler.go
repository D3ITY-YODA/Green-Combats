package preferences

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
)

type API interface {
	List(context.Context, uuid.UUID) ([]SavedPlace, error)
	Save(context.Context, uuid.UUID, uuid.UUID, *string) error
	Unsave(context.Context, uuid.UUID, uuid.UUID) error
	SetPrimary(context.Context, uuid.UUID, uuid.UUID) error
}

type Handler struct{ svc API }

func NewHandler(svc API) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterRoutes(r gin.IRouter, authService auth.API) {
	group := r.Group("/v1/me/places", auth.Middleware(authService))
	group.GET("", h.list)
	group.PUT("/:place_id", h.save)
	group.DELETE("/:place_id", h.unsave)
	group.PUT("/:place_id/primary", h.setPrimary)
}

func (h *Handler) list(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	items, err := h.svc.List(c.Request.Context(), identity.UserID)
	if err != nil {
		writeError(c, err)
		return
	}
	response := make([]gin.H, 0, len(items))
	for _, item := range items {
		response = append(response, gin.H{"id": item.ID, "name": item.Name, "place_type": item.PlaceType, "lat": item.Lat, "lon": item.Lon, "label": item.Label, "is_primary": item.IsPrimary, "saved_at": item.SavedAt})
	}
	c.JSON(http.StatusOK, gin.H{"places": response})
}

func (h *Handler) save(c *gin.Context) {
	placeID, identity, ok := requestIdentity(c)
	if !ok {
		return
	}
	var req struct {
		Label *string `json:"label"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	if err := h.svc.Save(c.Request.Context(), identity.UserID, placeID, req.Label); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) unsave(c *gin.Context) {
	placeID, identity, ok := requestIdentity(c)
	if !ok {
		return
	}
	if err := h.svc.Unsave(c.Request.Context(), identity.UserID, placeID); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) setPrimary(c *gin.Context) {
	placeID, identity, ok := requestIdentity(c)
	if !ok {
		return
	}
	if err := h.svc.SetPrimary(c.Request.Context(), identity.UserID, placeID); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func requestIdentity(c *gin.Context) (uuid.UUID, auth.Identity, bool) {
	placeID, err := uuid.Parse(c.Param("place_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place id"})
		return uuid.Nil, auth.Identity{}, false
	}
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return uuid.Nil, auth.Identity{}, false
	}
	return placeID, identity, true
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidData):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrPlaceNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "place not found"})
	case errors.Is(err, ErrNotSaved):
		c.JSON(http.StatusNotFound, gin.H{"error": "saved place not found"})
	case errors.Is(err, ErrSavedPlaceLimit):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
