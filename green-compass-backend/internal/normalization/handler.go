package normalization

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"green-compass-backend/internal/ingestion"
)

type Handler struct {
	service *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{service: svc}
}

type NormalizeRequest struct {
	SourceCode string             `json:"source_code" binding:"required"`
	RawRecord  ingestion.RawRecord `json:"raw_record" binding:"required"`
}

func (h *Handler) Normalize(c *gin.Context) {
	var req NormalizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := h.service.Normalize(c.Request.Context(), req.SourceCode, req.RawRecord)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"observations": results, "count": len(results)})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/normalize", h.Normalize)
}
