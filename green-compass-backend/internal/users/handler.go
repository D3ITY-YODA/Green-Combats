package users

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/pkg/httpx"
)

type API interface {
	Register(ctx context.Context, in RegisterInput) (*User, error)
	ByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, displayName, language string) (*User, error)
}

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all user-related routes.
// Auth protection is applied via middleware in cmd/api/main.go.
func (h *Handler) RegisterRoutes(rg gin.IRouter) {
	group := rg.Group("/v1/users")
	group.GET("", h.ListUsers)
	group.GET("/:id", h.GetUser)
}

// ListUsers returns paginated users.
// GET /v1/users?page=1&limit=20
// Note: Authorization (admin-only) should be enforced via middleware.
func (h *Handler) ListUsers(c *gin.Context) {
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

	c.JSON(http.StatusOK, httpx.Success(gin.H{
		"users": []interface{}{},
		"page":  page,
		"limit": limit,
	}))
}

// GetUser returns a single user by ID.
// GET /v1/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.HandleError(c, httpx.InvalidParam("id", "must be a valid UUID"))
		return
	}

	user, err := h.svc.ByID(c.Request.Context(), id)
	if err != nil {
		httpx.HandleError(c, httpx.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{
		"user": gin.H{
			"id":                user.ID,
			"display_name":      user.DisplayName,
			"email":             user.Email,
			"phone_number":      user.PhoneNumber,
			"language":          user.Language,
			"is_platform_admin": user.IsPlatformAdmin,
		},
	}))
}
