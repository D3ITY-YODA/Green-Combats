package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"green-compass-backend/pkg/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/health/live", h.live)
	r.GET("/health/ready", h.ready)
	r.GET("/health/startup", h.startup)
	r.GET("/health", h.status) // legacy
}

func (h *Handler) live(c *gin.Context) {
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(h.service.Status(c.Request.Context()), requestID))
}

func (h *Handler) ready(c *gin.Context) {
	requestID := httpx.GetRequestID(c)
	// TODO: Add database/redis connectivity checks
	c.JSON(http.StatusOK, httpx.Success(h.service.Status(c.Request.Context()), requestID))
}

func (h *Handler) startup(c *gin.Context) {
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(h.service.Status(c.Request.Context()), requestID))
}

func (h *Handler) status(c *gin.Context) {
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(h.service.Status(c.Request.Context()), requestID))
}
