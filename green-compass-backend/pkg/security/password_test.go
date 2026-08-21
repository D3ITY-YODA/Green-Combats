package security_test

import (
	"strings"
	"testing"

	"green-compass-backend/pkg/security"
)

var fastParams = security.Argon2Params{MemoryKiB: 1024, TimeCost: 1, Threads: 1, SaltLen: 8, KeyLen: 16}

func TestHashAndVerifyPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		attempt  string
		wantOK   bool
	}{
		{name: "correct password", password: "correct horse battery staple", attempt: "correct horse battery staple", wantOK: true},
		{name: "wrong password", password: "correct horse battery staple", attempt: "incorrect horse battery staple"},
		{name: "case differs", password: "SecretPass", attempt: "secretpass"},
		{name: "empty attempt against nonempty password", password: "some-password", attempt: ""},
		{name: "unicode roundtrip", password: "påsswörd-日本語", attempt: "påsswörd-日本語", wantOK: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := hashForTest(tt.password)
			if err != nil {
				t.Fatalf("hash: %v", err)
			}
			ok, err := security.VerifyPassword(tt.attempt, encoded)
			if err != nil {
				t.Fatalf("verify returned unexpected error: %v", err)
			}
			if ok != tt.wantOK {
				t.Fatalf("verify = %v, want %v", ok, tt.wantOK)
			}
		})
	}
}

func TestHashPassword_UniqueSalts(t *testing.T) {
	first, err := hashForTest("same-password")
	if err != nil {
		t.Fatalf("first hash: %v", err)
	}
	second, err := hashForTest("same-password")
	if err != nil {
		t.Fatalf("second hash: %v", err)
	}
	if first == second {
		t.Fatal("two hashes of the same password are identical; salts are not random")
	}
	for _, encoded := range []string{first, second} {
		if ok, err := security.VerifyPassword("same-password", encoded); err != nil || !ok {
			t.Fatalf("verify of %q failed: ok=%v err=%v", encoded, ok, err)
		}
	}
}

func TestVerifyPassword_TamperedHashRejected(t *testing.T) {
	encoded, err := hashForTest("hunter2")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}

	tests := []struct {
		name    string
		mutate  func(string) string
		wantErr bool
	}{
		{
			name:   "flipped hash character",
			mutate: func(s string) string { return flipLastB64Char(s) },
		},
		{
			name:   "truncated hash",
			mutate: func(s string) string { return s[:len(s)-4] },
		},
		{
			name:    "wrong algorithm",
			mutate:  func(s string) string { return strings.Replace(s, "$argon2id$", "$argon2i$", 1) },
			wantErr: true,
		},
		{
			name:    "unsupported version",
			mutate:  func(s string) string { return strings.Replace(s, "v=19", "v=16", 1) },
			wantErr: true,
		},
		{
			name:    "garbage parameters",
			mutate:  func(s string) string { return replaceParamField(s, "m=notanumber,t=1,p=1") },
			wantErr: true,
		},
		{
			name:    "zero memory parameter",
			mutate:  func(s string) string { return replaceParamField(s, "m=0,t=1,p=1") },
			wantErr: true,
		},
		{
			name:    "invalid base64 salt",
			mutate:  corruptSalt,
			wantErr: true,
		},
		{
			name:    "missing fields",
			mutate:  func(s string) string { return "$argon2id$v=19$m=1024,t=1,p=1$" },
			wantErr: true,
		},
		{
			name:    "no leading delimiter",
			mutate:  func(s string) string { return strings.TrimPrefix(s, "$") },
			wantErr: true,
		},
		{
			name:    "completely empty",
			mutate:  func(string) string { return "" },
			wantErr: true,
		},
		{
			name:    "not a hash at all",
			mutate:  func(string) string { return "password123" },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corrupted := tt.mutate(encoded)
			ok, err := security.VerifyPassword("hunter2", corrupted)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error for %q, got ok=%v", corrupted, ok)
			}
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if ok {
					t.Fatal("tampered hash verified successfully")
				}
			}
		})
	}
}

func TestVerifyPassword_KnownVector(t *testing.T) {
	const knownVector = "$argon2id$v=19$m=19456,t=2,p=1$c4vYiLtzF800ehKVPCJ3Vw$3GPKV3G29SIAobaeeN3hLZnwpIpA4QNxggUOAu4Yaps"

	ok, err := security.VerifyPassword("green-compass-known-vector", knownVector)
	if err != nil {
		t.Fatalf("known vector verify error: %v", err)
	}
	if !ok {
		t.Fatal("known vector no longer verifies; encoded format or derivation changed")
	}

	ok, err = security.VerifyPassword("wrong", knownVector)
	if err != nil {
		t.Fatalf("known vector wrong-password error: %v", err)
	}
	if ok {
		t.Fatal("known vector verified with wrong password")
	}
}

func hashForTest(password string) (string, error) {
	return security.HashPasswordWithParams(password, fastParams)
}

func flipLastB64Char(s string) string {
	last := s[len(s)-1]
	replacement := byte('A')
	if last == 'A' {
		replacement = 'B'
	}
	return s[:len(s)-1] + string(replacement)
}

func replaceParamField(encoded, newField string) string {
	parts := strings.Split(encoded, "$")
	parts[3] = newField
	return strings.Join(parts, "$")
}

func corruptSalt(encoded string) string {
	parts := strings.Split(encoded, "$")
	parts[4] = "!!!not-base64!!!"
	return strings.Join(parts, "$")
}
