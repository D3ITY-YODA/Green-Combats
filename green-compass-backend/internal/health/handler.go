package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/health", h.status)
}

func (h *Handler) status(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.Status(c.Request.Context()))
}
