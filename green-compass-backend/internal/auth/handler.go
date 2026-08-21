package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/users"
	"green-compass-backend/pkg/httpx"
)

type API interface {
	Register(ctx context.Context, in users.RegisterInput) (*RegisterResult, error)
	Login(ctx context.Context, identifier, password string) (*TokenPair, error)
	Refresh(ctx context.Context, raw string) (*TokenPair, error)
	Logout(ctx context.Context, raw string) error
	VerifyAccessToken(raw string) (uuid.UUID, error)
	SessionUser(ctx context.Context, id uuid.UUID) (*SessionUser, error)
}

type Handler struct {
	svc API
}

func NewHandler(svc API) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg gin.IRouter) {
	group := rg.Group("/v1/auth")
	group.POST("/register", h.register)
	group.POST("/login", h.login)
	group.POST("/refresh", h.refresh)
	group.POST("/logout", h.logout)
	group.GET("/me", Middleware(h.svc), h.me)
}

type registerRequest struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	Language    string `json:"language"`
}

type tokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	AccessExpiresAt  string `json:"access_expires_at"`
	RefreshExpiresAt string `json:"refresh_expires_at"`
}

func (h *Handler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.HandleError(c, httpx.InvalidParam("body", "invalid JSON"))
		return
	}

	result, err := h.svc.Register(c.Request.Context(), users.RegisterInput{
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Password:    req.Password,
		Language:    req.Language,
	})
	if err != nil {
		writeServiceError(c, err)
		return
	}

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusCreated, httpx.Success(gin.H{
		"user":   sessionUserJSON(result.User.ID, result.User.DisplayName, result.User.Email, result.User.PhoneNumber, result.User.Language, result.User.IsPlatformAdmin),
		"tokens": toTokenResponse(&result.Tokens),
	}, requestID))
}

func (h *Handler) login(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.HandleError(c, httpx.InvalidParam("body", "invalid JSON"))
		return
	}

	pair, err := h.svc.Login(c.Request.Context(), req.Identifier, req.Password)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"tokens": toTokenResponse(pair)}, requestID))
}

func (h *Handler) refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.HandleError(c, httpx.InvalidParam("body", "invalid JSON"))
		return
	}

	pair, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"tokens": toTokenResponse(pair)}, requestID))
}

func (h *Handler) logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.HandleError(c, httpx.InvalidParam("body", "invalid JSON"))
		return
	}

	if err := h.svc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		writeServiceError(c, err)
		return
	}
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusNoContent, httpx.Success(nil, requestID))
}

func (h *Handler) me(c *gin.Context) {
	identity, ok := IdentityFrom(c.Request.Context())
	if !ok {
		httpx.HandleError(c, httpx.ErrUnauthorized)
		return
	}

	sessionUser, err := h.svc.SessionUser(c.Request.Context(), identity.UserID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"user": sessionUserJSON(
		sessionUser.ID, sessionUser.DisplayName, sessionUser.Email, sessionUser.PhoneNumber, sessionUser.Language, sessionUser.IsPlatformAdmin,
	)}, requestID))
}

func toTokenResponse(pair *TokenPair) tokenResponse {
	return tokenResponse{
		AccessToken:      pair.AccessToken,
		RefreshToken:     pair.RefreshToken,
		AccessExpiresAt:  pair.AccessExpiresAt.UTC().Format(timeFormat),
		RefreshExpiresAt: pair.RefreshExpiresAt.UTC().Format(timeFormat),
	}
}

const timeFormat = time.RFC3339

func sessionUserJSON(id uuid.UUID, displayName string, email, phone *string, language string, isPlatformAdmin bool) gin.H {
	return gin.H{
		"id":                id,
		"display_name":      displayName,
		"email":             email,
		"phone_number":      phone,
		"language":          language,
		"is_platform_admin": isPlatformAdmin,
	}
}

func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, users.ErrInvalidData):
		httpx.HandleError(c, httpx.ErrBadRequest)
	case errors.Is(err, users.ErrEmailTaken), errors.Is(err, users.ErrPhoneTaken):
		httpx.HandleError(c, httpx.ErrConflict)
	case errors.Is(err, ErrInvalidCredentials):
		httpx.HandleError(c, &httpx.AppError{Code: "INVALID_CREDENTIALS", Message: "invalid credentials", Status: http.StatusUnauthorized})
	case errors.Is(err, ErrTokenReuse):
		httpx.HandleError(c, &httpx.AppError{Code: "TOKEN_REUSE", Message: "refresh token reuse detected; all sessions revoked", Status: http.StatusUnauthorized})
	case errors.Is(err, ErrInvalidToken):
		httpx.HandleError(c, &httpx.AppError{Code: "INVALID_TOKEN", Message: "invalid or expired token", Status: http.StatusUnauthorized})
	default:
		httpx.HandleError(c, httpx.ErrInternal)
	}
}
