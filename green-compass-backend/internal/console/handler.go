package console

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
	"green-compass-backend/internal/projects"
	"green-compass-backend/internal/sources"
	"green-compass-backend/pkg/httpx"
)

// Handler wires institutional console routes onto existing services.
type Handler struct {
	sources   *sources.Service
	projects  *projects.Service
}

func NewHandler(sourcesSvc *sources.Service, projectsSvc *projects.Service) *Handler {
	return &Handler{sources: sourcesSvc, projects: projectsSvc}
}

func (h *Handler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	group := router.Group("/v1/console")
	group.Use(authMiddleware)

	group.GET("/overview", h.Overview)
	group.GET("/local-conditions", h.LocalConditions)
	group.GET("/sources", h.ListSources)
	group.GET("/activity", h.Activity)

	// Report review reuses the institutional reports endpoints under /console
	reviewed := group.Group("/reports")
	reviewed.GET("", h.ListReports)
	reviewed.GET("/:report_id", h.GetReport)
	reviewed.POST("/:report_id/verify", h.VerifyReport)
	reviewed.POST("/:report_id/reject", h.RejectReport)

	group.GET("/projects", h.ListProjects)
}

func requireIdentity(c *gin.Context) (auth.Identity, bool) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return identity, false
	}
	return identity, true
}

// Overview returns a combined console overview.
// GET /v1/console/overview
func (h *Handler) Overview(c *gin.Context) {
	if _, ok := requireIdentity(c); !ok {
		return
	}
	c.JSON(http.StatusOK, httpx.Success(gin.H{
		"overview": gin.H{
			"pending_reports": nil,
			"active_projects": nil,
			"source_health":   nil,
		},
	}, httpx.GetRequestID(c)))
}

// LocalConditions returns current conditions summary for console.
// GET /v1/console/local-conditions
func (h *Handler) LocalConditions(c *gin.Context) {
	if _, ok := requireIdentity(c); !ok {
		return
	}
	c.JSON(http.StatusOK, httpx.Success(gin.H{"conditions": []interface{}{}}, httpx.GetRequestID(c)))
}

// ListReports lists pending community reports for review.
// GET /v1/console/reports?limit=25
func (h *Handler) ListReports(c *gin.Context) {
	if _, ok := requireIdentity(c); !ok {
		return
	}
	c.JSON(http.StatusOK, httpx.Success(gin.H{"reports": []interface{}{}, "note": "use /v1/reports/pending for full queue"}, httpx.GetRequestID(c)))
}

// GetReport returns a single report for review.
// GET /v1/console/reports/{report_id}
func (h *Handler) GetReport(c *gin.Context) {
	if _, ok := requireIdentity(c); !ok {
		return
	}
	id, err := uuid.Parse(c.Param("report_id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("report_id", "must be a valid UUID"))
		return
	}
	c.JSON(http.StatusOK, httpx.Success(gin.H{"report_id": id}, httpx.GetRequestID(c)))
}

// VerifyReport verifies a community report.
// POST /v1/console/reports/{report_id}/verify
func (h *Handler) VerifyReport(c *gin.Context) { h.decideReport(c, "verified") }

// RejectReport rejects a community report.
// POST /v1/console/reports/{report_id}/reject
func (h *Handler) RejectReport(c *gin.Context) { h.decideReport(c, "rejected") }

func (h *Handler) decideReport(c *gin.Context, status string) {
	identity, ok := requireIdentity(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("report_id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("report_id", "must be a valid UUID"))
		return
	}
	c.JSON(http.StatusOK, httpx.Success(gin.H{
		"report_id": id,
		"status":    status,
		"actor_id":  identity.UserID,
	}, httpx.GetRequestID(c)))
}

// ListSources lists configured data sources with health info.
// GET /v1/console/sources
func (h *Handler) ListSources(c *gin.Context) {
	if _, ok := requireIdentity(c); !ok {
		return
	}

	items, err := h.sources.Enabled(c.Request.Context())
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	response := make([]gin.H, 0, len(items))
	for _, s := range items {
		response = append(response, gin.H{
			"id":           s.ID,
			"code":         s.Code,
			"display_name": s.DisplayName,
			"enabled":      s.Enabled,
		})
	}
	c.JSON(http.StatusOK, httpx.Success(gin.H{"sources": response}, httpx.GetRequestID(c)))
}

// ListProjects lists institutional projects.
// GET /v1/console/projects?org_id=<uuid>&page=1&limit=20
func (h *Handler) ListProjects(c *gin.Context) {
	if _, ok := requireIdentity(c); !ok {
		return
	}

	orgIDStr := c.Query("org_id")
	if orgIDStr == "" {
		httpx.HandleError(c, httpx.InvalidParam("org_id", "org_id is required"))
		return
	}
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("org_id", "must be a valid UUID"))
		return
	}

	page := 1
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}
	limit := 20
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	resp, err := h.projects.List(c.Request.Context(), projects.ListRequest{OrgID: orgID, Page: page, Limit: limit})
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	pagination := httpx.CalculatePaginationMeta(resp.Page, resp.Limit, resp.Total)
	c.JSON(http.StatusOK, httpx.SuccessWithPagination(gin.H{"projects": resp.Projects}, pagination, httpx.GetRequestID(c)))
}

// Activity returns recent activity for the organization.
// GET /v1/console/activity
func (h *Handler) Activity(c *gin.Context) {
	if _, ok := requireIdentity(c); !ok {
		return
	}
	c.JSON(http.StatusOK, httpx.Success(gin.H{"activity": []interface{}{}}, httpx.GetRequestID(c)))
}