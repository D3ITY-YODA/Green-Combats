package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"green-compass-backend/internal/users"
	"green-compass-backend/pkg/clock"
	"green-compass-backend/pkg/security"
)

const refreshTokenBytes = 32

type Service struct {
	users      *users.Service
	repo       *Repository
	clk        clock.Clock
	secret     string
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
	dummyHash  string
}

type Options struct {
	Secret     string
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewService(userSvc *users.Service, repo *Repository, clk clock.Clock, opts Options) (*Service, error) {
	if len(opts.Secret) < 32 {
		return nil, fmt.Errorf("auth: secret must be at least 32 characters")
	}
	if opts.AccessTTL <= 0 || opts.RefreshTTL <= 0 {
		return nil, fmt.Errorf("auth: token TTLs must be positive")
	}
	if opts.RefreshTTL <= opts.AccessTTL {
		return nil, fmt.Errorf("auth: refresh TTL must exceed access TTL")
	}
	if opts.Issuer == "" {
		opts.Issuer = "green-compass"
	}
	dummy, err := security.HashPassword("timing-equalizer-dummy-password")
	if err != nil {
		return nil, fmt.Errorf("auth: precompute dummy hash: %w", err)
	}
	return &Service{
		users:      userSvc,
		repo:       repo,
		clk:        clk,
		secret:     opts.Secret,
		issuer:     opts.Issuer,
		accessTTL:  opts.AccessTTL,
		refreshTTL: opts.RefreshTTL,
		dummyHash:  dummy,
	}, nil
}

type RegisterResult struct {
	User   *users.User
	Tokens TokenPair
}

func (s *Service) Register(ctx context.Context, in users.RegisterInput) (*RegisterResult, error) {
	u, err := s.users.Register(ctx, in)
	if err != nil {
		return nil, err
	}
	pair, err := s.issuePair(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	return &RegisterResult{User: u, Tokens: pair}, nil
}

func (s *Service) Login(ctx context.Context, identifier, password string) (*TokenPair, error) {
	u, err := s.lookupUser(ctx, identifier)
	if err != nil {
		_, _ = security.VerifyPassword(password, s.dummyHash)
		return nil, ErrInvalidCredentials
	}
	ok, err := security.VerifyPassword(password, u.PasswordHash)
	if err != nil || !ok {
		return nil, ErrInvalidCredentials
	}
	pair, err := s.issuePair(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	return &pair, nil
}

func (s *Service) Refresh(ctx context.Context, rawRefreshToken string) (*TokenPair, error) {
	raw := strings.TrimSpace(rawRefreshToken)
	if raw == "" {
		return nil, ErrInvalidToken
	}
	hash := hashRefreshToken(raw)

	rt, err := s.repo.ByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, ErrRefreshTokenNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	now := s.clk.Now()

	if rt.RevokedAt != nil {
		if err := s.repo.RevokeAllForUser(ctx, rt.UserID, now); err != nil {
			return nil, fmt.Errorf("revoke all tokens after reuse: %w", err)
		}
		return nil, ErrTokenReuse
	}

	if !now.Before(rt.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	pair, err := s.issuePair(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	newHash := hashRefreshToken(pair.RefreshToken)
	newRow, err := s.repo.ByHash(ctx, newHash)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Revoke(ctx, rt.ID, &newRow.ID, now); err != nil {
		return nil, fmt.Errorf("mark rotated token used: %w", err)
	}
	return &pair, nil
}

func (s *Service) Logout(ctx context.Context, rawRefreshToken string) error {
	raw := strings.TrimSpace(rawRefreshToken)
	if raw == "" {
		return ErrInvalidToken
	}
	rt, err := s.repo.ByHash(ctx, hashRefreshToken(raw))
	if err != nil {
		if errors.Is(err, ErrRefreshTokenNotFound) {
			return ErrInvalidToken
		}
		return err
	}
	if rt.RevokedAt != nil {
		return nil
	}
	return s.repo.Revoke(ctx, rt.ID, nil, s.clk.Now())
}

func (s *Service) VerifyAccessToken(rawAccessToken string) (uuid.UUID, error) {
	userID, err := security.VerifyAccessToken(rawAccessToken, s.secret, s.issuer, s.clk.Now())
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}
	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}
	return id, nil
}

func (s *Service) SessionUser(ctx context.Context, id uuid.UUID) (*SessionUser, error) {
	u, err := s.users.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &SessionUser{
		ID:              u.ID,
		DisplayName:     u.DisplayName,
		Email:           u.Email,
		PhoneNumber:     u.PhoneNumber,
		Language:        u.Language,
		IsPlatformAdmin: u.IsPlatformAdmin,
	}, nil
}

func (s *Service) lookupUser(ctx context.Context, identifier string) (*users.User, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, ErrInvalidCredentials
	}
	if strings.Contains(identifier, "@") {
		return s.users.ByEmail(ctx, identifier)
	}
	return s.users.ByPhone(ctx, identifier)
}

func (s *Service) issuePair(ctx context.Context, userID uuid.UUID) (TokenPair, error) {
	now := s.clk.Now()

	access, accessExpiry, err := security.IssueAccessToken(userID.String(), s.issuer, s.secret, s.accessTTL, now)
	if err != nil {
		return TokenPair{}, err
	}

	raw := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(raw); err != nil {
		return TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}
	refresh := base64.RawURLEncoding.EncodeToString(raw)

	rt := &RefreshToken{
		UserID:    userID,
		TokenHash: hashRefreshToken(refresh),
		IssuedAt:  now,
		ExpiresAt: now.Add(s.refreshTTL),
	}
	if err := s.repo.Create(ctx, rt); err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:      access,
		RefreshToken:     refresh,
		AccessExpiresAt:  accessExpiry,
		RefreshExpiresAt: rt.ExpiresAt,
	}, nil
}

func hashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
