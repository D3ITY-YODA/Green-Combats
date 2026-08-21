package observations_test

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
	"green-compass-backend/internal/users"
)

type stubObsAPI struct {
	create      *observations.Observation
	createErr   error
	byID        *observations.Observation
	byIDErr     error
	listMine    []observations.Observation
	listMineErr error
	listPending []observations.Observation
	listPendErr error
	verify      *observations.Observation
	verifyErr   error
}

func (s *stubObsAPI) Create(context.Context, observations.CreateInput) (*observations.Observation, error) {
	return s.create, s.createErr
}
func (s *stubObsAPI) ByID(context.Context, uuid.UUID, observations.Caller) (*observations.Observation, error) {
	return s.byID, s.byIDErr
}
func (s *stubObsAPI) ListMine(context.Context, uuid.UUID, int) ([]observations.Observation, error) {
	return s.listMine, s.listMineErr
}
func (s *stubObsAPI) ListPending(context.Context, observations.Caller, int) ([]observations.Observation, error) {
	return s.listPending, s.listPendErr
}
func (s *stubObsAPI) Verify(context.Context, uuid.UUID, observations.Caller, string) (*observations.Observation, error) {
	return s.verify, s.verifyErr
}

func TestHandler_PublicRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	obsID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		headers    map[string]string
		stub       stubObsAPI
		wantStatus int
	}{
		{
			name:   "create observation",
			method: http.MethodPost,
			path:   "/v1/observations",
			body:   `{"category":"flood","description":"River rising near market","lat":-1.2,"lon":36.8}`,
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			stub: stubObsAPI{create: &observations.Observation{ID: obsID, Category: "flood", Status: observations.StatusPending}},
			wantStatus: http.StatusCreated,
		},
		{
			name:   "create invalid category",
			method: http.MethodPost,
			path:   "/v1/observations",
			body:   `{"category":"invalid","description":"test"}`,
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			stub:   stubObsAPI{createErr: observations.ErrInvalidData},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing auth",
			method:     http.MethodGet,
			path:       "/v1/observations",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:   "list mine",
			method: http.MethodGet,
			path:   "/v1/observations",
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			stub: stubObsAPI{listMine: []observations.Observation{{ID: obsID, Category: "flood"}}},
			wantStatus: http.StatusOK,
		},
		{
			name:   "get observation",
			method: http.MethodGet,
			path:   "/v1/observations/" + obsID.String(),
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			stub:   stubObsAPI{byID: &observations.Observation{ID: obsID, Category: "flood", ReporterID: userID}},
			wantStatus: http.StatusOK,
		},
		{
			name:   "observation not found",
			method: http.MethodGet,
			path:   "/v1/observations/" + obsID.String(),
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			stub:   stubObsAPI{byIDErr: observations.ErrNotFound},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "invalid id",
			method: http.MethodGet,
			path:   "/v1/observations/nope",
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			observations.NewHandler(&tt.stub, &stubAuthAPI{userID: userID}).RegisterRoutes(router)
			rec := doRequest(router, tt.method, tt.path, tt.body, tt.headers)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
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
		stub       stubObsAPI
		wantStatus int
	}{
		{
			name:   "list pending as admin",
			method: http.MethodGet,
			path:   "/v1/observations/institutional/pending",
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			stub: stubObsAPI{listPending: []observations.Observation{{ID: obsID, Status: observations.StatusPending}}},
			wantStatus: http.StatusOK,
		},
		{
			name:   "verify observation",
			method: http.MethodPatch,
			path:   "/v1/observations/institutional/" + obsID.String() + "/verify",
			body:   `{"status":"verified"}`,
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			stub: stubObsAPI{verify: &observations.Observation{ID: obsID, Status: observations.StatusVerified}},
			wantStatus: http.StatusOK,
		},
		{
			name:   "verify invalid status",
			method: http.MethodPatch,
			path:   "/v1/observations/institutional/" + obsID.String() + "/verify",
			body:   `{"status":"pending"}`,
			headers: map[string]string{"Authorization": "Bearer " + uuid.New().String()},
			stub:   stubObsAPI{verifyErr: observations.ErrInvalidData},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			observations.NewHandler(&tt.stub, &stubAuthAPI{userID: adminID, isAdmin: true}).RegisterRoutes(router)
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
