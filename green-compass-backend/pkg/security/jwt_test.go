package security_test

import (
	"strings"
	"testing"
	"time"

	"green-compass-backend/pkg/security"
)

const (
	testSecret = "test-secret-0123456789abcdef0123456789abcdef"
	testIssuer = "green-compass-test"
)

var testNow = time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)

func TestIssueAndVerifyAccessToken(t *testing.T) {
	token, expiresAt, err := security.IssueAccessToken("user-123", testIssuer, testSecret, 15*time.Minute, testNow)
	if err != nil {
		t.Fatalf("IssueAccessToken: %v", err)
	}
	if !expiresAt.Equal(testNow.Add(15 * time.Minute)) {
		t.Errorf("expiresAt = %v, want %v", expiresAt, testNow.Add(15*time.Minute))
	}

	userID, err := security.VerifyAccessToken(token, testSecret, testIssuer, testNow.Add(time.Minute))
	if err != nil {
		t.Fatalf("VerifyAccessToken: %v", err)
	}
	if userID != "user-123" {
		t.Errorf("userID = %q, want user-123", userID)
	}
}

func TestVerifyAccessToken_Rejections(t *testing.T) {
	token, _, err := security.IssueAccessToken("user-123", testIssuer, testSecret, 15*time.Minute, testNow)
	if err != nil {
		t.Fatalf("IssueAccessToken: %v", err)
	}

	expired, _, err := security.IssueAccessToken("user-123", testIssuer, testSecret, time.Minute, testNow)
	if err != nil {
		t.Fatalf("IssueAccessToken(expired): %v", err)
	}

	otherKey := strings.Repeat("x", len(testSecret))

	tests := []struct {
		name    string
		token   string
		secret  string
		issuer  string
		now     time.Time
		wantUID string
	}{
		{name: "expired token", token: expired, secret: testSecret, issuer: testIssuer, now: testNow.Add(2 * time.Minute)},
		{name: "wrong secret", token: token, secret: otherKey, issuer: testIssuer, now: testNow},
		{name: "wrong issuer", token: token, secret: testSecret, issuer: "someone-else", now: testNow},
		{name: "tampered payload", token: tamperSegment(token, 1), secret: testSecret, issuer: testIssuer, now: testNow},
		{name: "tampered signature", token: tamperSegment(token, 2), secret: testSecret, issuer: testIssuer, now: testNow},
		{name: "garbage string", token: "not.a.jwt", secret: testSecret, issuer: testIssuer, now: testNow},
		{name: "empty string", token: "", secret: testSecret, issuer: testIssuer, now: testNow},
		{name: "none algorithm header", token: noneAlgToken(t), secret: testSecret, issuer: testIssuer, now: testNow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := security.VerifyAccessToken(tt.token, tt.secret, tt.issuer, tt.now)
			if err == nil {
				t.Fatalf("expected rejection, got userID %q", userID)
			}
			if userID != "" {
				t.Errorf("userID = %q on failure, want empty", userID)
			}
			if err != security.ErrInvalidToken {
				t.Errorf("err = %v, want ErrInvalidToken sentinel", err)
			}
		})
	}
}

func tamperSegment(token string, segment int) string {
	parts := strings.Split(token, ".")
	candidate := "A"
	if parts[segment][0] == 'A' {
		candidate = "B"
	}
	parts[segment] = candidate + parts[segment][1:]
	return strings.Join(parts, ".")
}

func noneAlgToken(t *testing.T) string {
	t.Helper()
	return "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1aWQiOiJoYWNrZXIifQ."
}
