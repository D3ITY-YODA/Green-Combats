package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Algorithm = "argon2id"
	argon2Version   = 19
)

type Argon2Params struct {
	MemoryKiB uint32
	TimeCost  uint32
	Threads   uint8
	SaltLen   uint32
	KeyLen    uint32
}

var DefaultArgon2Params = Argon2Params{
	MemoryKiB: 19456,
	TimeCost:  2,
	Threads:   1,
	SaltLen:   16,
	KeyLen:    32,
}

func HashPassword(password string) (string, error) {
	return hashPasswordWith(password, DefaultArgon2Params)
}

func HashPasswordWithParams(password string, params Argon2Params) (string, error) {
	return hashPasswordWith(password, params)
}

func hashPasswordWith(password string, params Argon2Params) (string, error) {
	salt := make([]byte, params.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}
	key := argon2.IDKey([]byte(password), salt, params.TimeCost, params.MemoryKiB, params.Threads, params.KeyLen)
	enc := base64.RawStdEncoding
	return fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2Algorithm, argon2Version,
		params.MemoryKiB, params.TimeCost, params.Threads,
		enc.EncodeToString(salt), enc.EncodeToString(key)), nil
}

func VerifyPassword(password, encoded string) (bool, error) {
	params, salt, want, err := decodeEncodedHash(encoded)
	if err != nil {
		return false, err
	}
	got := argon2.IDKey([]byte(password), salt, params.TimeCost, params.MemoryKiB, params.Threads, uint32(len(want)))
	return subtle.ConstantTimeCompare(got, want) == 1, nil
}

func decodeEncodedHash(encoded string) (Argon2Params, []byte, []byte, error) {
	var zero Argon2Params
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return zero, nil, nil, fmt.Errorf("malformed encoded hash: expected 6 fields, got %d", len(parts))
	}
	if parts[0] != "" {
		return zero, nil, nil, errors.New("malformed encoded hash: missing leading delimiter")
	}
	if parts[1] != argon2Algorithm {
		return zero, nil, nil, fmt.Errorf("unsupported algorithm %q", parts[1])
	}
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return zero, nil, nil, fmt.Errorf("malformed version %q", parts[2])
	}
	if version != argon2Version {
		return zero, nil, nil, fmt.Errorf("unsupported version %d", version)
	}
	params, err := parseArgon2ParamField(parts[3])
	if err != nil {
		return zero, nil, nil, err
	}
	enc := base64.RawStdEncoding
	salt, err := enc.DecodeString(parts[4])
	if err != nil {
		return zero, nil, nil, fmt.Errorf("decode salt: %w", err)
	}
	key, err := enc.DecodeString(parts[5])
	if err != nil {
		return zero, nil, nil, fmt.Errorf("decode hash: %w", err)
	}
	if len(salt) == 0 || len(key) == 0 {
		return zero, nil, nil, errors.New("encoded hash has empty salt or key")
	}
	return params, salt, key, nil
}

func parseArgon2ParamField(field string) (Argon2Params, error) {
	var p Argon2Params
	n, err := fmt.Sscanf(field, "m=%d,t=%d,p=%d", &p.MemoryKiB, &p.TimeCost, &p.Threads)
	if err != nil || n != 3 {
		return p, fmt.Errorf("malformed parameters %q", field)
	}
	if p.MemoryKiB == 0 || p.TimeCost == 0 || p.Threads == 0 {
		return p, fmt.Errorf("invalid parameters %q", field)
	}
	return p, nil
}
