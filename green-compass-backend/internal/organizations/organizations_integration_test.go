package organizations_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"

	"green-compass-backend/internal/organizations"
	"green-compass-backend/internal/testutil"
	"green-compass-backend/internal/users"
)

type harness struct {
	svc    *organizations.Service
	users  *users.Service
	owner  *users.User
	member *users.User
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	pool := testutil.TestPool(t)
	ctx := context.Background()

	userSvc := users.NewService(users.NewRepository(pool))
	owner, err := userSvc.Register(ctx, users.RegisterInput{
		Email:       uniqueEmail(t),
		DisplayName: "Org Owner",
		Password:    "long-enough-password",
	})
	if err != nil {
		t.Fatalf("seed owner: %v", err)
	}
	plain, err := userSvc.Register(ctx, users.RegisterInput{
		Email:       uniqueEmail(t),
		DisplayName: "Plain Member",
		Password:    "long-enough-password",
	})
	if err != nil {
		t.Fatalf("seed member: %v", err)
	}

	return &harness{
		svc:    organizations.NewService(pool, organizations.NewRepository(pool)),
		users:  userSvc,
		owner:  owner,
		member: plain,
	}
}

func TestService_CreateWithOwner_Integration(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	org, err := h.svc.CreateWithOwner(ctx, organizations.CreateInput{
		Name:    uniqueOrgName(t),
		OrgType: organizations.TypeWaterAuthority,
		OwnerID: h.owner.ID,
	})
	if err != nil {
		t.Fatalf("CreateWithOwner: %v", err)
	}
	if org.Status != organizations.StatusPending {
		t.Errorf("new org status = %q, want pending", org.Status)
	}

	fresh, err := h.svc.ByID(ctx, org.ID)
	if err != nil {
		t.Fatalf("ByID: %v", err)
	}
	if fresh.Name != org.Name || fresh.OrgType != organizations.TypeWaterAuthority {
		t.Errorf("roundtrip mismatch: %+v", fresh)
	}

	role, err := h.svc.MemberRole(ctx, org.ID, h.owner.ID)
	if err != nil {
		t.Fatalf("MemberRole for owner: %v", err)
	}
	if role != organizations.RoleAdmin {
		t.Errorf("owner role = %q, want admin", role)
	}
}

func TestService_CreateValidation_Integration(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	tests := []struct {
		name    string
		input   organizations.CreateInput
		wantErr error
	}{
		{
			name:    "empty name",
			input:   organizations.CreateInput{Name: "   ", OrgType: organizations.TypeNGO, OwnerID: h.owner.ID},
			wantErr: organizations.ErrInvalidData,
		},
		{
			name:    "unknown type",
			input:   organizations.CreateInput{Name: uniqueOrgName(t), OrgType: "cartel", OwnerID: h.owner.ID},
			wantErr: organizations.ErrInvalidData,
		},
		{
			name:    "missing owner",
			input:   organizations.CreateInput{Name: uniqueOrgName(t), OrgType: organizations.TypeNGO},
			wantErr: organizations.ErrInvalidData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := h.svc.CreateWithOwner(ctx, tt.input); !errors.Is(err, tt.wantErr) {
				t.Fatalf("CreateWithOwner err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_DuplicateName_Integration(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()
	name := uniqueOrgName(t)

	if _, err := h.svc.CreateWithOwner(ctx, organizations.CreateInput{
		Name: name, OrgType: organizations.TypeNGO, OwnerID: h.owner.ID,
	}); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := h.svc.CreateWithOwner(ctx, organizations.CreateInput{
		Name: name, OrgType: organizations.TypeNGO, OwnerID: h.member.ID,
	}); !errors.Is(err, organizations.ErrNameTaken) {
		t.Fatalf("duplicate name err = %v, want ErrNameTaken", err)
	}
}

func TestService_SetStatus_Integration(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	org, err := h.svc.CreateWithOwner(ctx, organizations.CreateInput{
		Name: uniqueOrgName(t), OrgType: organizations.TypeGovernment, OwnerID: h.owner.ID,
	})
	if err != nil {
		t.Fatalf("create org: %v", err)
	}

	for _, status := range []string{organizations.StatusVerified, organizations.StatusRejected, organizations.StatusPending} {
		updated, err := h.svc.SetStatus(ctx, org.ID, status)
		if err != nil {
			t.Fatalf("SetStatus(%q): %v", status, err)
		}
		if updated.Status != status {
			t.Errorf("SetStatus(%q) result = %q", status, updated.Status)
		}
	}

	if _, err := h.svc.SetStatus(ctx, org.ID, "elite"); !errors.Is(err, organizations.ErrInvalidData) {
		t.Errorf("invalid status err = %v, want ErrInvalidData", err)
	}
	if _, err := h.svc.SetStatus(ctx, uuid.New(), organizations.StatusVerified); !errors.Is(err, organizations.ErrNotFound) {
		t.Errorf("unknown org err = %v, want ErrNotFound", err)
	}
}

func TestService_MembershipLifecycle_Integration(t *testing.T) {
	h := newHarness(t)
	ctx := context.Background()

	org, err := h.svc.CreateWithOwner(ctx, organizations.CreateInput{
		Name: uniqueOrgName(t), OrgType: organizations.TypeNGO, OwnerID: h.owner.ID,
	})
	if err != nil {
		t.Fatalf("create org: %v", err)
	}

	if _, err := h.svc.MemberRole(ctx, org.ID, h.member.ID); !errors.Is(err, organizations.ErrNotMember) {
		t.Errorf("non-member lookup err = %v, want ErrNotMember", err)
	}

	if err := h.svc.AddMember(ctx, org.ID, h.member.ID, organizations.RoleReviewer); err != nil {
		t.Fatalf("AddMember: %v", err)
	}
	role, err := h.svc.MemberRole(ctx, org.ID, h.member.ID)
	if err != nil {
		t.Fatalf("MemberRole after add: %v", err)
	}
	if role != organizations.RoleReviewer {
		t.Errorf("role = %q, want reviewer", role)
	}

	if err := h.svc.AddMember(ctx, org.ID, h.member.ID, organizations.RoleMember); !errors.Is(err, organizations.ErrAlreadyMember) {
		t.Errorf("re-add err = %v, want ErrAlreadyMember", err)
	}

	if err := h.svc.AddMember(ctx, org.ID, uuid.New(), "supreme-leader"); !errors.Is(err, organizations.ErrInvalidData) {
		t.Errorf("invalid role err = %v, want ErrInvalidData", err)
	}

	if err := h.svc.RemoveMember(ctx, org.ID, h.member.ID); err != nil {
		t.Fatalf("RemoveMember: %v", err)
	}
	if _, err := h.svc.MemberRole(ctx, org.ID, h.member.ID); !errors.Is(err, organizations.ErrNotMember) {
		t.Errorf("after removal err = %v, want ErrNotMember", err)
	}
	if err := h.svc.RemoveMember(ctx, org.ID, h.member.ID); !errors.Is(err, organizations.ErrNotMember) {
		t.Errorf("double remove err = %v, want ErrNotMember", err)
	}
}

func uniqueOrgName(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("Test Org %s", uuid.NewString())
}

func uniqueEmail(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("org-user-%s@greencompass.test", uuid.NewString())
}
