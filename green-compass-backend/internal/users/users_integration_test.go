package users_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"

	"green-compass-backend/internal/testutil"
	"green-compass-backend/internal/users"
	"green-compass-backend/pkg/security"
)

func TestRepository_CreateAndLookup_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	repo := users.NewRepository(pool)
	ctx := context.Background()

	email := uniqueEmail(t)

	created := &users.User{
		Email:        &email,
		PasswordHash: "test-hash",
		DisplayName:  "Asha Kimani",
		Language:     "sw",
	}
	if err := repo.Create(ctx, created); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == uuid.Nil {
		t.Fatal("Create did not populate ID")
	}
	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Fatal("Create did not populate timestamps")
	}

	byID, err := repo.ByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("ByID: %v", err)
	}
	if byID.Email == nil || *byID.Email != email {
		t.Errorf("ByEmail roundtrip: got %+v, want email %q", byID, email)
	}
	if byID.PhoneNumber != nil {
		t.Errorf("phone should be NULL, got %q", *byID.PhoneNumber)
	}
	if byID.DisplayName != "Asha Kimani" || byID.Language != "sw" || byID.PasswordHash != "test-hash" {
		t.Errorf("roundtrip mismatch: %+v", byID)
	}
	if byID.IsPlatformAdmin {
		t.Error("IsPlatformAdmin should default to false")
	}

	byEmail, err := repo.ByEmail(ctx, email)
	if err != nil {
		t.Fatalf("ByEmail: %v", err)
	}
	if byEmail.ID != created.ID {
		t.Errorf("ByEmail returned %s, want %s", byEmail.ID, created.ID)
	}
}

func TestRepository_PhoneUserRoundtrip_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	repo := users.NewRepository(pool)
	ctx := context.Background()

	phone := uniquePhone(t)

	created := &users.User{
		PhoneNumber:  &phone,
		PasswordHash: "test-hash",
		DisplayName:  "Phone User",
		Language:     "en",
	}
	if err := repo.Create(ctx, created); err != nil {
		t.Fatalf("Create: %v", err)
	}

	byPhone, err := repo.ByPhone(ctx, phone)
	if err != nil {
		t.Fatalf("ByPhone: %v", err)
	}
	if byPhone.ID != created.ID {
		t.Errorf("ByPhone returned %s, want %s", byPhone.ID, created.ID)
	}
	if byPhone.Email != nil {
		t.Errorf("email should be NULL, got %q", *byPhone.Email)
	}
	if byPhone.Language != "en" {
		t.Errorf("language = %q, want en", byPhone.Language)
	}
}

func TestRepository_NotFound_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	repo := users.NewRepository(pool)
	ctx := context.Background()

	if _, err := repo.ByID(ctx, uuid.New()); !errors.Is(err, users.ErrNotFound) {
		t.Errorf("ByID unknown id: err = %v, want ErrNotFound", err)
	}
	if _, err := repo.ByEmail(ctx, uniqueEmail(t)); !errors.Is(err, users.ErrNotFound) {
		t.Errorf("ByEmail unknown email: err = %v, want ErrNotFound", err)
	}
	if _, err := repo.ByPhone(ctx, uniquePhone(t)); !errors.Is(err, users.ErrNotFound) {
		t.Errorf("ByPhone unknown phone: err = %v, want ErrNotFound", err)
	}
}

func TestRepository_DuplicateIdentifiers_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	repo := users.NewRepository(pool)
	ctx := context.Background()

	email := uniqueEmail(t)
	phone := uniquePhone(t)

	phoneUser := &users.User{PhoneNumber: &phone, PasswordHash: "h", DisplayName: "Phone First"}
	if err := repo.Create(ctx, phoneUser); err != nil {
		t.Fatalf("seed phone user: %v", err)
	}

	dupPhone := &users.User{PhoneNumber: &phone, PasswordHash: "h", DisplayName: "Dup"}
	if err := repo.Create(ctx, dupPhone); !errors.Is(err, users.ErrPhoneTaken) {
		t.Errorf("duplicate phone: err = %v, want ErrPhoneTaken", err)
	}

	email = uniqueEmail(t)
	emailUser := &users.User{Email: &email, PasswordHash: "h", DisplayName: "Email First"}
	if err := repo.Create(ctx, emailUser); err != nil {
		t.Fatalf("seed email user: %v", err)
	}

	dupEmail := &users.User{Email: &email, PasswordHash: "h", DisplayName: "Dup"}
	if err := repo.Create(ctx, dupEmail); !errors.Is(err, users.ErrEmailTaken) {
		t.Errorf("duplicate email: err = %v, want ErrEmailTaken", err)
	}
}

func TestService_Register_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	svc := users.NewService(users.NewRepository(pool))
	ctx := context.Background()

	email := uniqueEmail(t)
	u, err := svc.Register(ctx, users.RegisterInput{
		Email:       email,
		DisplayName: "  Registered User  ",
		Password:    "long-enough-password",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if u.DisplayName != "Registered User" {
		t.Errorf("display name = %q, want trimmed value", u.DisplayName)
	}
	if u.PasswordHash == "long-enough-password" || len(u.PasswordHash) < 60 {
		t.Errorf("password not hashed: %q", u.PasswordHash)
	}
	if ok, err := verifyPassword(u.PasswordHash, "long-enough-password"); err != nil || !ok {
		t.Errorf("stored hash does not verify against original password: ok=%v err=%v", ok, err)
	}
}

func TestService_RegisterValidation_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	svc := users.NewService(users.NewRepository(pool))
	ctx := context.Background()

	tests := []struct {
		name    string
		input   users.RegisterInput
		wantErr error
	}{
		{
			name:    "no identifier",
			input:   users.RegisterInput{DisplayName: "X", Password: "long-enough"},
			wantErr: users.ErrInvalidData,
		},
		{
			name:    "both identifiers",
			input:   users.RegisterInput{Email: uniqueEmail(t), PhoneNumber: uniquePhone(t), DisplayName: "X", Password: "long-enough"},
			wantErr: users.ErrInvalidData,
		},
		{
			name:    "short password",
			input:   users.RegisterInput{Email: uniqueEmail(t), DisplayName: "X", Password: "short"},
			wantErr: users.ErrInvalidData,
		},
		{
			name:    "empty display name",
			input:   users.RegisterInput{Email: uniqueEmail(t), DisplayName: "   ", Password: "long-enough"},
			wantErr: users.ErrInvalidData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := svc.Register(ctx, tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Register err = %v, want %v", err, tt.wantErr)
			}
			if u != nil {
				t.Fatalf("Register returned user %+v on validation failure", u)
			}
		})
	}
}

func TestService_UpdateProfile_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	svc := users.NewService(users.NewRepository(pool))
	ctx := context.Background()

	email := uniqueEmail(t)
	u, err := svc.Register(ctx, users.RegisterInput{
		Email:       email,
		DisplayName: "Before",
		Password:    "long-enough-password",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	updated, err := svc.UpdateProfile(ctx, u.ID, "After", "sw")
	if err != nil {
		t.Fatalf("UpdateProfile: %v", err)
	}
	if updated.DisplayName != "After" || updated.Language != "sw" {
		t.Errorf("UpdateProfile result = %+v, want name After language sw", updated)
	}
	if !updated.UpdatedAt.After(updated.CreatedAt) {
		t.Errorf("updated_at %v should be after created_at %v", updated.UpdatedAt, updated.CreatedAt)
	}

	fresh, err := svc.ByID(ctx, u.ID)
	if err != nil {
		t.Fatalf("ByID after update: %v", err)
	}
	if fresh.DisplayName != "After" {
		t.Errorf("update not persisted: display name = %q", fresh.DisplayName)
	}
}

func verifyPassword(hash, password string) (bool, error) {
	return security.VerifyPassword(password, hash)
}

func uniqueEmail(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("user-%s@greencompass.test", uuid.NewString())
}

func uniquePhone(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("+2547%09d", uuid.New().ID()%1_000_000_000)
}
