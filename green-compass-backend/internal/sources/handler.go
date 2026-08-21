package sources

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"green-compass-backend/internal/auth"
	"green-compass-backend/pkg/httpx"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers data source routes.
func (h *Handler) RegisterRoutes(rg gin.IRouter) {
	group := rg.Group("/v1/sources")
	group.GET("", h.ListSources)
}

// ListSources returns all data sources.
// GET /v1/sources
func (h *Handler) ListSources(c *gin.Context) {
	_, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	sources, err := h.svc.Enabled(c.Request.Context())
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"sources": sources}, requestID))
}
