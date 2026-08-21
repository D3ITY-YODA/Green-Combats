package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

// AccessTokenClaims deliberately carries only the user id — no roles, org
// memberships, or trust status. Permissions are resolved from the database at
// request time (see internal/permissions) so that revocation and institutional
// status changes (org verified -> rejected, admin demoted) take effect
// immediately instead of living on for the life of the token. Do not embed
// authorization data here; it would reintroduce stale-permission bugs for up
// to the access-token TTL.
type AccessTokenClaims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func IssueAccessToken(userID, issuer, secret string, ttl time.Duration, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(ttl)
	claims := AccessTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}
	return signed, expiresAt, nil
}

func VerifyAccessToken(token, secret, issuer string, now time.Time) (string, error) {
	parser := jwt.NewParser(
		jwt.WithTimeFunc(func() time.Time { return now }),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
	)

	var claims AccessTokenClaims
	parsed, err := parser.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return "", ErrInvalidToken
	}
	if claims.UserID == "" {
		return "", ErrInvalidToken
	}
	return claims.UserID, nil
}
