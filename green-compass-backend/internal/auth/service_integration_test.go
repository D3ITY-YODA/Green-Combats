package auth_test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"green-compass-backend/internal/auth"
	"green-compass-backend/internal/testutil"
	"green-compass-backend/internal/users"
	"green-compass-backend/pkg/clock"
)

const (
	testSecret = "auth-test-secret-0123456789abcdef0123456789abcdef"
	testIssuer = "green-compass-test"
)

var (
	accessTTL  = 15 * time.Minute
	refreshTTL = 24 * time.Hour
	startTime  = time.Date(2026, 8, 21, 9, 0, 0, 0, time.UTC)
)

type fixture struct {
	svc   *auth.Service
	clk   *clock.Fixed
	user  *users.User
	email string
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	pool := testutil.TestPool(t)
	ctx := context.Background()

	userSvc := users.NewService(users.NewRepository(pool))
	email := "auth-" + uniqueSuffix(t) + "@greencompass.test"
	u, err := userSvc.Register(ctx, users.RegisterInput{
		Email:       email,
		DisplayName: "Auth Test User",
		Password:    "long-enough-password",
	})
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	clk := clock.NewFixed(startTime)
	svc, err := auth.NewService(userSvc, auth.NewRepository(pool), clk, auth.Options{
		Secret:     testSecret,
		Issuer:     testIssuer,
		AccessTTL:  accessTTL,
		RefreshTTL: refreshTTL,
	})
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	return &fixture{svc: svc, clk: clk, user: u, email: email}
}

func TestRegister_IssuesWorkingTokenPair(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	result, err := f.svc.Register(ctx, users.RegisterInput{
		Email:       "reg-" + uniqueSuffix(t) + "@greencompass.test",
		DisplayName: "Fresh User",
		Password:    "long-enough-password",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if result.Tokens.AccessToken == "" || result.Tokens.RefreshToken == "" {
		t.Fatal("Register returned empty tokens")
	}
	if !result.Tokens.RefreshExpiresAt.Equal(startTime.Add(refreshTTL)) {
		t.Errorf("refresh expiry = %v, want %v", result.Tokens.RefreshExpiresAt, startTime.Add(refreshTTL))
	}

	id, err := f.svc.VerifyAccessToken(result.Tokens.AccessToken)
	if err != nil {
		t.Fatalf("VerifyAccessToken: %v", err)
	}
	if id != result.User.ID {
		t.Errorf("token subject = %s, want %s", id, result.User.ID)
	}
}

func TestLogin_EmailAndPhone(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	phone := "+25471" + randomDigits(t, 8)
	phoneResult, err := f.svc.Register(ctx, users.RegisterInput{
		PhoneNumber: phone,
		DisplayName: "Phone Login User",
		Password:    "long-enough-password",
	})
	if err != nil {
		t.Fatalf("seed phone user: %v", err)
	}

	pairEmail, err := f.svc.Login(ctx, f.email, "long-enough-password")
	if err != nil {
		t.Fatalf("Login by email: %v", err)
	}
	id, err := f.svc.VerifyAccessToken(pairEmail.AccessToken)
	if err != nil {
		t.Fatalf("VerifyAccessToken after email login: %v", err)
	}
	if id != f.user.ID {
		t.Errorf("email login subject = %s, want %s", id, f.user.ID)
	}

	pairPhone, err := f.svc.Login(ctx, phone, "long-enough-password")
	if err != nil {
		t.Fatalf("Login by phone: %v", err)
	}
	id, err = f.svc.VerifyAccessToken(pairPhone.AccessToken)
	if err != nil {
		t.Fatalf("VerifyAccessToken after phone login: %v", err)
	}
	if id != phoneResult.User.ID {
		t.Errorf("phone login subject = %s, want %s", id, phoneResult.User.ID)
	}
}

func TestLogin_BadCredentials(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	tests := []struct {
		name       string
		identifier string
		password   string
	}{
		{name: "wrong password", identifier: f.email, password: "wrong-password-1"},
		{name: "unknown email", identifier: "nobody-" + uniqueSuffix(t) + "@greencompass.test", password: "long-enough-password"},
		{name: "unknown phone", identifier: "+254700000000", password: "long-enough-password"},
		{name: "empty identifier", identifier: "", password: "long-enough-password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := f.svc.Login(ctx, tt.identifier, tt.password); !errors.Is(err, auth.ErrInvalidCredentials) {
				t.Fatalf("Login err = %v, want ErrInvalidCredentials", err)
			}
		})
	}
}

func TestRefresh_RotatesAndInvalidatesOld(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	first, err := f.svc.Login(ctx, f.email, "long-enough-password")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	second, err := f.svc.Refresh(ctx, first.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if second.RefreshToken == first.RefreshToken {
		t.Error("rotation returned the same refresh token")
	}

	id, err := f.svc.VerifyAccessToken(second.AccessToken)
	if err != nil || id != f.user.ID {
		t.Fatalf("new access token invalid: id=%s err=%v", id, err)
	}

	if _, err := f.svc.Refresh(ctx, first.RefreshToken); !errors.Is(err, auth.ErrTokenReuse) {
		t.Fatalf("presenting rotated token err = %v, want ErrTokenReuse", err)
	}
}

func TestRefresh_ReuseDetectionRevokesEverything(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	first, err := f.svc.Login(ctx, f.email, "long-enough-password")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	second, err := f.svc.Refresh(ctx, first.RefreshToken)
	if err != nil {
		t.Fatalf("first rotation: %v", err)
	}
	third, err := f.svc.Refresh(ctx, second.RefreshToken)
	if err != nil {
		t.Fatalf("second rotation: %v", err)
	}

	if _, err := f.svc.Refresh(ctx, first.RefreshToken); !errors.Is(err, auth.ErrTokenReuse) {
		t.Fatalf("stale token replay err = %v, want ErrTokenReuse", err)
	}

	if _, err := f.svc.Refresh(ctx, third.RefreshToken); !errors.Is(err, auth.ErrTokenReuse) {
		t.Fatalf("newest token was revoked by reuse detection, err = %v, want ErrTokenReuse", err)
	}

	if _, err := f.svc.Login(ctx, f.email, "long-enough-password"); err != nil {
		t.Fatalf("fresh login after compromise response should work: %v", err)
	}
}

func TestRefresh_ExpiredTokenRejected(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	pair, err := f.svc.Login(ctx, f.email, "long-enough-password")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	f.clk.Advance(refreshTTL + time.Minute)

	if _, err := f.svc.Refresh(ctx, pair.RefreshToken); !errors.Is(err, auth.ErrInvalidToken) {
		t.Fatalf("expired refresh err = %v, want ErrInvalidToken", err)
	}
}

func TestLogout_RevokesAndIsIdempotent(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	pair, err := f.svc.Login(ctx, f.email, "long-enough-password")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if err := f.svc.Logout(ctx, pair.RefreshToken); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if _, err := f.svc.Refresh(ctx, pair.RefreshToken); !errors.Is(err, auth.ErrTokenReuse) {
		t.Fatalf("refresh after logout err = %v, want ErrTokenReuse", err)
	}

	if err := f.svc.Logout(ctx, pair.RefreshToken); err != nil {
		t.Fatalf("second Logout should be idempotent, got %v", err)
	}

	if err := f.svc.Logout(ctx, "totally-unknown-token"); !errors.Is(err, auth.ErrInvalidToken) {
		t.Errorf("unknown token logout err = %v, want ErrInvalidToken", err)
	}
}

func TestVerifyAccessToken_RejectsGarbage(t *testing.T) {
	f := newFixture(t)

	for _, raw := range []string{"", "garbage"} {
		if _, err := f.svc.VerifyAccessToken(raw); !errors.Is(err, auth.ErrInvalidToken) {
			t.Errorf("VerifyAccessToken(%q) err = %v, want ErrInvalidToken", raw, err)
		}
	}
}

func uniqueSuffix(t *testing.T) string {
	t.Helper()
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("rand: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func randomDigits(t *testing.T, n int) string {
	t.Helper()
	digits := make([]byte, n)
	for i := range digits {
		b := make([]byte, 1)
		if _, err := rand.Read(b); err != nil {
			t.Fatalf("rand: %v", err)
		}
		digits[i] = byte('0') + b[0]%10
	}
	return string(digits)
}
