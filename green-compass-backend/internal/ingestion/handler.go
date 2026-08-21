package ingestion

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
	"green-compass-backend/internal/sources"
	"green-compass-backend/pkg/httpx"
)

type Handler struct {
	svc        *Service
	sourcesSvc *sources.Service
}

func NewHandler(svc *Service, sourcesSvc *sources.Service) *Handler {
	return &Handler{svc: svc, sourcesSvc: sourcesSvc}
}

// RegisterRoutes registers ingestion-related routes.
func (h *Handler) RegisterRoutes(rg gin.IRouter) {
	group := rg.Group("/v1/ingestion")
	group.POST("/trigger", h.TriggerIngestion)
	group.GET("/runs", h.ListRuns)
}

// TriggerIngestion manually triggers an ingestion run for a source and place.
// POST /v1/ingestion/trigger
func (h *Handler) TriggerIngestion(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}
	_ = identity

	var req struct {
		SourceCode string `json:"source_code"`
		PlaceID    string `json:"place_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	if req.SourceCode == "" || req.PlaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_code and place_id are required"})
		return
	}

	// The full implementation would look up the source and place,
	// then call svc.Ingest. For now, return a placeholder.
	c.JSON(http.StatusAccepted, httpx.Success(gin.H{
		"message":     "ingestion triggered",
		"source_code": req.SourceCode,
		"place_id":    req.PlaceID,
	}))
}

// ListRuns returns ingestion run history (placeholder).
// GET /v1/ingestion/runs
func (h *Handler) ListRuns(c *gin.Context) {
	_, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	// Placeholder: full implementation would query ingestion_runs table
	c.JSON(http.StatusOK, httpx.Success(gin.H{
		"runs": []interface{}{},
	}))
}

// Helper to parse UUID safely
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
