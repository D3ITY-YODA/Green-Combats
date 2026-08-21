package permissions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
	"green-compass-backend/pkg/httpx"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers permission-related routes.
func (h *Handler) RegisterRoutes(rg gin.IRouter) {
	group := rg.Group("/v1/permissions")
	group.GET("/check", h.CheckPermission)
	group.GET("/roles", h.ListRoles)
}

// CheckPermission checks if the current user has a specific role in an organization.
// GET /v1/permissions/check?org_id=<uuid>&role=<string>
func (h *Handler) CheckPermission(c *gin.Context) {
	identity, ok := auth.IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

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

	role := c.Query("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role is required"})
		return
	}

	err = h.svc.RequireOrgRole(c.Request.Context(), identity.UserID, orgID, role)
	if err != nil {
		c.JSON(http.StatusOK, httpx.Success(gin.H{
			"has_permission": false,
			"user_id":        identity.UserID,
			"org_id":         orgID,
			"required_role":  role,
		}, httpx.GetRequestID(c)))
		return
	}

	c.JSON(http.StatusOK, httpx.Success(gin.H{
		"has_permission": true,
		"user_id":        identity.UserID,
		"org_id":         orgID,
		"required_role":  role,
	}, httpx.GetRequestID(c)))
}

// ListRoles returns the available roles.
// GET /v1/permissions/roles
func (h *Handler) ListRoles(c *gin.Context) {
	c.JSON(http.StatusOK, httpx.Success(gin.H{
		"roles": []string{"admin", "reviewer", "viewer"},
	}, httpx.GetRequestID(c)))
}
