package reporting

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
	"green-compass-backend/pkg/httpx"
)

type API interface {
	GetDeliveryStats(ctx context.Context, placeID uuid.UUID, from, to string) (*DeliveryStats, error)
	GetReportStats(ctx context.Context, orgID uuid.UUID, from, to string) (*ReportStats, error)
	GetDashboard(ctx context.Context, req QueryRequest) (*DashboardResponse, error)
}

type Handler struct {
	svc API
}

func NewHandler(svc API) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all reporting-related routes.
func (h *Handler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	group := router.Group("/v1/reporting")
	group.Use(authMiddleware)

	group.GET("/delivery", h.DeliveryStats)
	group.GET("/reports", h.ReportStats)
	group.GET("/dashboard", h.Dashboard)
}

// DeliveryStats returns notification delivery metrics for a place.
// GET /v1/reporting/delivery?place_id=<uuid>&period_start=<rfc3339>&period_end=<rfc3339>
func (h *Handler) DeliveryStats(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}
	_ = identity

	placeIDStr := c.Query("place_id")
	if placeIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "place_id is required"})
		return
	}
	placeID, err := uuid.Parse(placeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place_id"})
		return
	}

	periodStart := c.Query("period_start")
	periodEnd := c.Query("period_end")
	if periodStart == "" || periodEnd == "" {
		// Default to last 30 days
		periodEnd = time.Now().UTC().Format(time.RFC3339)
		periodStart = time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)
	}

	stats, err := h.svc.GetDeliveryStats(c.Request.Context(), placeID, periodStart, periodEnd)
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{"delivery_stats": stats}))
}

// ReportStats returns community report metrics for an organization.
// GET /v1/reporting/reports?org_id=<uuid>&period_start=<rfc3339>&period_end=<rfc3339>
func (h *Handler) ReportStats(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}
	_ = identity

	orgIDStr := c.Query("org_id")
	if orgIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org_id is required"})
		return
	}
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org_id"})
		return
	}

	periodStart := c.Query("period_start")
	periodEnd := c.Query("period_end")
	if periodStart == "" || periodEnd == "" {
		periodEnd = time.Now().UTC().Format(time.RFC3339)
		periodStart = time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)
	}

	stats, err := h.svc.GetReportStats(c.Request.Context(), orgID, periodStart, periodEnd)
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{"report_stats": stats}))
}

// Dashboard returns a combined dashboard view.
// GET /v1/reporting/dashboard?org_id=<uuid>&place_id=<uuid>&period_start=<rfc3339>&period_end=<rfc3339>
func (h *Handler) Dashboard(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}
	_ = identity

	var orgID *uuid.UUID
	if v := c.Query("org_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org_id"})
			return
		}
		orgID = &id
	}

	var placeID *uuid.UUID
	if v := c.Query("place_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place_id"})
			return
		}
		placeID = &id
	}

	periodStart := c.Query("period_start")
	periodEnd := c.Query("period_end")
	if periodStart == "" || periodEnd == "" {
		periodEnd = time.Now().UTC().Format(time.RFC3339)
		periodStart = time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)
	}

	start, _ := time.Parse(time.RFC3339, periodStart)
	end, _ := time.Parse(time.RFC3339, periodEnd)

	dashboard, err := h.svc.GetDashboard(c.Request.Context(), QueryRequest{
		OrgID:       orgID,
		PlaceID:     placeID,
		PeriodStart: start,
		PeriodEnd:   end,
	})
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{"dashboard": dashboard}))
}
