package context

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
	ResolveContext(ctx context.Context, req ResolveContextRequest) (*ContextResponse, error)
}

type Handler struct {
	svc API
}

func NewHandler(svc API) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers context-related routes
func (h *Handler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	group := router.Group("/v1")
	group.Use(authMiddleware)
	group.GET("/context/today", h.GetToday)
}

// GetToday resolves the user's current context and returns today's content
// GET /v1/context/today?place_id=<optional-uuid>
func (h *Handler) GetToday(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	var placeID *uuid.UUID
	if placeIDStr := c.Query("place_id"); placeIDStr != "" {
		if pID, err := uuid.Parse(placeIDStr); err == nil {
			placeID = &pID
		} else {
			httpx.HandleError(c, httpx.InvalidParam("place_id", "must be a valid UUID"))
			return
		}
	}

	resp, err := h.svc.ResolveContext(c.Request.Context(), ResolveContextRequest{
		UserID:  identity.UserID,
		PlaceID: placeID,
	})
	if err != nil {
		var appErr *httpx.AppError
		if errors.As(err, &appErr) {
			httpx.HandleError(c, err)
			return
		}
		switch err {
		case ErrNoPlace:
			httpx.HandleError(c, httpx.NotFound("no place assigned to user"))
		case ErrPlaceNotFound:
			httpx.HandleError(c, httpx.ErrNotFound)
		default:
			httpx.HandleError(c, httpx.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, httpx.Success(resp))
}
