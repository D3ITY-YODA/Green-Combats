package reports_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
	"green-compass-backend/internal/observations"
	"green-compass-backend/internal/reports"
	"green-compass-backend/internal/users"
)

type stubReportsAPI struct {
	listPending []reports.Report
	listPendErr error
	byID        *reports.Report
	byIDErr     error
	verify      *reports.Report
	verifyErr   error
}

func (s *stubReportsAPI) ListPending(context.Context, reports.Caller, int) ([]reports.Report, error) {
	return s.listPending, s.listPendErr
}
func (s *stubReportsAPI) Get(context.Context, uuid.UUID, reports.Caller) (*reports.Report, error) {
	return s.byID, s.byIDErr
}
func (s *stubReportsAPI) Verify(context.Context, uuid.UUID, reports.Caller, string) (*reports.Report, error) {
	return s.verify, s.verifyErr
}

func TestHandler_InstitutionalRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	obsID := uuid.New()
	adminID := uuid.New()

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		headers    map[string]string
		stub       stubReportsAPI
		wantStatus int
	}{
		{
			name:   "list pending as admin",
			method: http.MethodGet,
			path:   "/v1/reports/pending",
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			stub: stubReportsAPI{listPending: []reports.Report{{ID: obsID, Status: observations.StatusPending}}},
			wantStatus: http.StatusOK,
		},
		{
			name:   "get report",
			method: http.MethodGet,
			path:   "/v1/reports/" + obsID.String(),
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			stub:   stubReportsAPI{byID: &reports.Report{ID: obsID, Status: observations.StatusPending}},
			wantStatus: http.StatusOK,
		},
		{
			name:   "verify report",
			method: http.MethodPatch,
			path:   "/v1/reports/" + obsID.String() + "/verify",
			body:   `{"status":"verified"}`,
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			stub:   stubReportsAPI{verify: &reports.Report{ID: obsID, Status: observations.StatusVerified}},
			wantStatus: http.StatusOK,
		},
		{
			name:   "missing auth",
			method: http.MethodGet,
			path:   "/v1/reports/pending",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "invalid id",
			method: http.MethodGet,
			path:   "/v1/reports/nope",
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			reports.NewHandler(&tt.stub, &stubAuthAPI{userID: adminID, isAdmin: true}).RegisterRoutes(router)
			rec := doRequest(router, tt.method, tt.path, tt.body, tt.headers)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}

type stubAuthAPI struct {
	userID  uuid.UUID
	isAdmin bool
}

func (s *stubAuthAPI) Register(context.Context, users.RegisterInput) (*auth.RegisterResult, error) {
	return nil, nil
}
func (s *stubAuthAPI) Login(context.Context, string, string) (*auth.TokenPair, error) {
	return nil, nil
}
func (s *stubAuthAPI) Refresh(context.Context, string) (*auth.TokenPair, error) {
	return nil, nil
}
func (s *stubAuthAPI) Logout(context.Context, string) error { return nil }
func (s *stubAuthAPI) VerifyAccessToken(raw string) (uuid.UUID, error) {
	if raw == "bad" {
		return uuid.Nil, errors.New("invalid token")
	}
	return s.userID, nil
}
func (s *stubAuthAPI) SessionUser(context.Context, uuid.UUID) (*auth.SessionUser, error) {
	return &auth.SessionUser{ID: s.userID, IsPlatformAdmin: s.isAdmin}, nil
}

func doRequest(router *gin.Engine, method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}
