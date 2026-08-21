package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ctxKey struct{}

type Identity struct {
	UserID          uuid.UUID
	IsPlatformAdmin bool
}

func Middleware(svc API) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			abortUnauthorized(c, "missing bearer token")
			return
		}

		userID, err := svc.VerifyAccessToken(strings.TrimPrefix(header, prefix))
		if err != nil {
			abortUnauthorized(c, "invalid or expired access token")
			return
		}

		sessionUser, err := svc.SessionUser(c.Request.Context(), userID)
		if err != nil {
			abortUnauthorized(c, "account no longer exists")
			return
		}

		identity := Identity{UserID: sessionUser.ID, IsPlatformAdmin: sessionUser.IsPlatformAdmin}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), ctxKey{}, identity))
		c.Next()
	}
}

func IdentityFrom(ctx context.Context) (Identity, bool) {
	identity, ok := ctx.Value(ctxKey{}).(Identity)
	return identity, ok
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
}
