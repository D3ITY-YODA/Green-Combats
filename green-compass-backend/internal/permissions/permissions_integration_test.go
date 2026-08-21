package permissions_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"

	"green-compass-backend/internal/organizations"
	"green-compass-backend/internal/permissions"
	"green-compass-backend/internal/testutil"
	"green-compass-backend/internal/users"
)

type orgFixture struct {
	svc      *permissions.Service
	orgs     *organizations.Service
	users    *users.Service
	org      *organizations.Organization
	admin    *users.User
	reviewer *users.User
	member   *users.User
	outsider *users.User
}

func newOrgFixture(t *testing.T, status string) *orgFixture {
	t.Helper()
	pool := testutil.TestPool(t)
	ctx := context.Background()

	userSvc := users.NewService(users.NewRepository(pool))
	register := func(name string) *users.User {
		u, err := userSvc.Register(ctx, users.RegisterInput{
			Email:       fmt.Sprintf("%s-%s@perm.test", name, uniqueSuffix(t)),
			DisplayName: name,
			Password:    "long-enough-password",
		})
		if err != nil {
			t.Fatalf("seed %s: %v", name, err)
		}
		return u
	}

	admin := register("Admin")
	reviewer := register("Reviewer")
	member := register("Member")
	outsider := register("Outsider")

	orgSvc := organizations.NewService(pool, organizations.NewRepository(pool))
	org, err := orgSvc.CreateWithOwner(ctx, organizations.CreateInput{
		Name:    "Perm Org " + uniqueSuffix(t),
		OrgType: organizations.TypeNGO,
		OwnerID: admin.ID,
	})
	if err != nil {
		t.Fatalf("create org: %v", err)
	}
	if status != organizations.StatusPending {
		if _, err := orgSvc.SetStatus(ctx, org.ID, status); err != nil {
			t.Fatalf("set org status: %v", err)
		}
	}
	if err := orgSvc.AddMember(ctx, org.ID, reviewer.ID, organizations.RoleReviewer); err != nil {
		t.Fatalf("add reviewer: %v", err)
	}
	if err := orgSvc.AddMember(ctx, org.ID, member.ID, organizations.RoleMember); err != nil {
		t.Fatalf("add member: %v", err)
	}

	return &orgFixture{
		svc:      permissions.New(orgSvc),
		orgs:     orgSvc,
		users:    userSvc,
		org:      org,
		admin:    admin,
		reviewer: reviewer,
		member:   member,
		outsider: outsider,
	}
}

func TestRequireOrgRole(t *testing.T) {
	f := newOrgFixture(t, organizations.StatusVerified)
	ctx := context.Background()

	if err := f.svc.RequireOrgRole(ctx, f.admin.ID, f.org.ID, organizations.RoleAdmin); err != nil {
		t.Errorf("admin should pass admin check: %v", err)
	}
	if err := f.svc.RequireOrgRole(ctx, f.reviewer.ID, f.org.ID, organizations.RoleAdmin, organizations.RoleReviewer); err != nil {
		t.Errorf("reviewer should pass reviewer check: %v", err)
	}
	if err := f.svc.RequireOrgRole(ctx, f.member.ID, f.org.ID, organizations.RoleAdmin); !errors.Is(err, permissions.ErrPermissionDenied) {
		t.Errorf("member err = %v, want ErrPermissionDenied", err)
	}
	if err := f.svc.RequireOrgRole(ctx, f.outsider.ID, f.org.ID, organizations.RoleAdmin); !errors.Is(err, organizations.ErrNotMember) {
		t.Errorf("outsider err = %v, want ErrNotMember", err)
	}
}

func TestCanVerifyObservations_RequiresVerifiedOrgAndTrustedRole(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		status  string
		user    func(f *orgFixture) *users.User
		allowed bool
	}{
		{name: "verified org admin", status: organizations.StatusVerified, user: func(f *orgFixture) *users.User { return f.admin }, allowed: true},
		{name: "verified org reviewer", status: organizations.StatusVerified, user: func(f *orgFixture) *users.User { return f.reviewer }, allowed: true},
		{name: "verified org member", status: organizations.StatusVerified, user: func(f *orgFixture) *users.User { return f.member }},
		{name: "pending org admin", status: organizations.StatusPending, user: func(f *orgFixture) *users.User { return f.admin }},
		{name: "rejected org reviewer", status: organizations.StatusRejected, user: func(f *orgFixture) *users.User { return f.reviewer }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newOrgFixture(t, tt.status)
			err := f.svc.CheckCanVerifyObservations(ctx, tt.user(f).ID, f.org.ID)
			if tt.allowed && err != nil {
				t.Fatalf("expected allowed, got error: %v", err)
			}
			if !tt.allowed && !errors.Is(err, permissions.ErrPermissionDenied) {
				t.Fatalf("err = %v, want ErrPermissionDenied", err)
			}
		})
	}
}

func TestCheckCanVerifyObservations_UnknownOrg(t *testing.T) {
	f := newOrgFixture(t, organizations.StatusVerified)
	ctx := context.Background()

	err := f.svc.CheckCanVerifyObservations(ctx, f.admin.ID, uuid.New())
	if !errors.Is(err, organizations.ErrNotFound) {
		t.Fatalf("err = %v, want organizations.ErrNotFound", err)
	}
}

func TestPureCanVerifyObservations(t *testing.T) {
	tests := []struct {
		name      string
		role      string
		orgStatus string
		want      bool
	}{
		{name: "reviewer in verified", role: organizations.RoleReviewer, orgStatus: organizations.StatusVerified, want: true},
		{name: "admin in verified", role: organizations.RoleAdmin, orgStatus: organizations.StatusVerified, want: true},
		{name: "reviewer in pending", role: organizations.RoleReviewer, orgStatus: organizations.StatusPending},
		{name: "member in verified", role: organizations.RoleMember, orgStatus: organizations.StatusVerified},
		{name: "no role in verified", role: "", orgStatus: organizations.StatusVerified},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := permissions.CanVerifyObservations(tt.role, tt.orgStatus)
			if got != tt.want {
				t.Errorf("CanVerifyObservations(%q, %q) = %v, want %v", tt.role, tt.orgStatus, got, tt.want)
			}
		})
	}
}
