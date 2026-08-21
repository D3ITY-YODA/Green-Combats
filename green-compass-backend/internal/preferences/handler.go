package preferences

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
	"green-compass-backend/pkg/httpx"
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
		httpx.HandleError(c, httpx.ErrUnauthorized)
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
	c.JSON(http.StatusOK, httpx.Success(gin.H{"places": response}, httpx.GetRequestID(c)))
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
		httpx.HandleError(c, httpx.InvalidParam("body", "invalid JSON"))
		return
	}
	if err := h.svc.Save(c.Request.Context(), identity.UserID, placeID, req.Label); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, httpx.Success(gin.H{"message": "place saved"}, httpx.GetRequestID(c)))
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
	c.JSON(http.StatusOK, httpx.Success(gin.H{"message": "place unsaved"}, httpx.GetRequestID(c)))
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
	c.JSON(http.StatusOK, httpx.Success(gin.H{"message": "primary place updated"}, httpx.GetRequestID(c)))
}

func requestIdentity(c *gin.Context) (uuid.UUID, auth.Identity, bool) {
	placeID, err := uuid.Parse(c.Param("place_id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("place_id", "invalid place id"))
		return uuid.Nil, auth.Identity{}, false
	}
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return uuid.Nil, auth.Identity{}, false
	}
	return placeID, identity, true
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidData):
		httpx.HandleError(c, httpx.ErrBadRequest)
	case errors.Is(err, ErrPlaceNotFound):
		httpx.HandleError(c, httpx.ErrNotFound)
	case errors.Is(err, ErrNotSaved):
		httpx.HandleError(c, httpx.ErrNotFound)
	case errors.Is(err, ErrSavedPlaceLimit):
		httpx.HandleError(c, httpx.ErrConflict)
	default:
		httpx.HandleError(c, httpx.ErrInternal)
	}
}
