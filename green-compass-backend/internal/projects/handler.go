package projects

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
	"green-compass-backend/pkg/httpx"
)

type API interface {
	Create(ctx context.Context, req CreateRequest) (*Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Project, error)
	List(ctx context.Context, req ListRequest) (*ListResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	svc API
}

func NewHandler(svc API) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all project-related routes.
func (h *Handler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	group := router.Group("/v1/projects")
	group.Use(authMiddleware)

	group.POST("", h.Create)
	group.GET("", h.List)
	group.GET("/:id", h.GetByID)
	group.PUT("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create creates a new project.
// POST /v1/projects
func (h *Handler) Create(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}
	_ = identity

	var req struct {
		OrgID       string  `json:"org_id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		StartDate   *string `json:"start_date"`
		EndDate     *string `json:"end_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	orgID, err := uuid.Parse(req.OrgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org_id"})
		return
	}

	createReq := CreateRequest{
		OrgID:       orgID,
		Name:        req.Name,
		Description: req.Description,
	}

	project, err := h.svc.Create(c.Request.Context(), createReq)
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusCreated, httpx.Success(gin.H{"project": project}, requestID))
}

// List returns paginated projects for an organization.
// GET /v1/projects?org_id=<uuid>&page=1&limit=20
func (h *Handler) List(c *gin.Context) {
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

	page := 1
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}
	limit := 20
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 {
			limit = l
		}
	}

	resp, err := h.svc.List(c.Request.Context(), ListRequest{
		OrgID: orgID,
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	pagination := httpx.CalculatePaginationMeta(resp.Page, resp.Limit, resp.Total)
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.SuccessWithPagination(gin.H{
		"projects": resp.Projects,
	}, pagination, requestID))
}

// GetByID returns a single project.
// GET /v1/projects/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("id", "must be a valid UUID"))
		return
	}

	project, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		httpx.HandleError(c, httpx.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{"project": project}, httpx.GetRequestID(c)))
}

// Update modifies an existing project.
// PUT /v1/projects/:id
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("id", "must be a valid UUID"))
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.HandleError(c, httpx.InvalidParam("body", "invalid JSON"))
		return
	}

	project, err := h.svc.Update(c.Request.Context(), id, UpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{"project": project}, httpx.GetRequestID(c)))
}

// Delete removes a project.
// DELETE /v1/projects/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("id", "must be a valid UUID"))
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		httpx.HandleError(c, httpx.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{"message": "project deleted"}, httpx.GetRequestID(c)))
}
