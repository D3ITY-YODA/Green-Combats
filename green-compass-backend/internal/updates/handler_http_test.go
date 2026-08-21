package updates_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
	"green-compass-backend/internal/updates"
	"green-compass-backend/internal/users"
)

type stubAuthAPI struct {
	verifyUserID uuid.UUID
	verifyErr    error
	sessionUser  *auth.SessionUser
}

func (s *stubAuthAPI) Register(_ context.Context, _ users.RegisterInput) (*auth.RegisterResult, error) {
	return nil, nil
}

func (s *stubAuthAPI) Login(_ context.Context, _, _ string) (*auth.TokenPair, error) {
	return nil, nil
}

func (s *stubAuthAPI) Refresh(_ context.Context, _ string) (*auth.TokenPair, error) {
	return nil, nil
}

func (s *stubAuthAPI) Logout(_ context.Context, _ string) error { return nil }

func (s *stubAuthAPI) VerifyAccessToken(_ string) (uuid.UUID, error) {
	return s.verifyUserID, s.verifyErr
}

func (s *stubAuthAPI) SessionUser(_ context.Context, _ uuid.UUID) (*auth.SessionUser, error) {
	return s.sessionUser, nil
}

type stubUpdatesAPI struct {
	listResp    *updates.ListResponse
	listErr     error
	exploreResp *updates.ExploreResponse
	exploreErr  error

	lastListReq    *updates.ListRequest
	lastExploreReq *updates.ExploreRequest
}

func (s *stubUpdatesAPI) List(_ context.Context, req updates.ListRequest) (*updates.ListResponse, error) {
	s.lastListReq = &req
	if s.listErr != nil {
		return nil, s.listErr
	}
	return s.listResp, nil
}

func (s *stubUpdatesAPI) Explore(_ context.Context, req updates.ExploreRequest) (*updates.ExploreResponse, error) {
	s.lastExploreReq = &req
	if s.exploreErr != nil {
		return nil, s.exploreErr
	}
	return s.exploreResp, nil
}

func (s *stubUpdatesAPI) GetByID(context.Context, uuid.UUID) (*updates.Update, error) {
	return nil, updates.ErrNotFound
}

func (s *stubUpdatesAPI) Acknowledge(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}

func newTestRouter(api updates.API, authAPI auth.API) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	updates.NewHandler(api).RegisterRoutes(router, auth.Middleware(authAPI))
	return router
}

func TestListUpdates_MissingAuth(t *testing.T) {
	router := newTestRouter(&stubUpdatesAPI{}, &stubAuthAPI{})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/updates", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestListUpdates_Success(t *testing.T) {
	userID := uuid.New()
	authStub := &stubAuthAPI{
		verifyUserID: userID,
		sessionUser:  &auth.SessionUser{ID: userID, DisplayName: "Test"},
	}

	stub := &stubUpdatesAPI{
		listResp: &updates.ListResponse{
			Updates: []updates.Update{
				{ID: uuid.New(), PlaceName: "Nairobi", Headline: "Update 1"},
			},
			Total: 1,
			Page:  1,
			Limit: 20,
		},
	}

	router := newTestRouter(stub, authStub)
	req := httptest.NewRequest(http.MethodGet, "/v1/updates", nil)
	req.Header.Set("Authorization", "Bearer test")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}

	if stub.lastListReq == nil {
		t.Fatal("expected List to be called")
	}
	if stub.lastListReq.UserID != userID {
		t.Errorf("UserID: got %v, want %v", stub.lastListReq.UserID, userID)
	}
	if stub.lastListReq.Page != 1 {
		t.Errorf("Page: got %d, want 1", stub.lastListReq.Page)
	}
	if stub.lastListReq.Limit != 20 {
		t.Errorf("Limit: got %d, want 20", stub.lastListReq.Limit)
	}
}

func TestListUpdates_PaginationParams(t *testing.T) {
	userID := uuid.New()
	authStub := &stubAuthAPI{
		verifyUserID: userID,
		sessionUser:  &auth.SessionUser{ID: userID, DisplayName: "Test"},
	}

	stub := &stubUpdatesAPI{
		listResp: &updates.ListResponse{
			Updates: []updates.Update{},
			Total:   0,
			Page:    2,
			Limit:   25,
		},
	}

	router := newTestRouter(stub, authStub)
	req := httptest.NewRequest(http.MethodGet, "/v1/updates?page=2&limit=25", nil)
	req.Header.Set("Authorization", "Bearer test")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if stub.lastListReq.Page != 2 {
		t.Errorf("Page: got %d, want 2", stub.lastListReq.Page)
	}
	if stub.lastListReq.Limit != 25 {
		t.Errorf("Limit: got %d, want 25", stub.lastListReq.Limit)
	}
}

func TestListUpdates_LimitCappedByService(t *testing.T) {
	userID := uuid.New()
	authStub := &stubAuthAPI{
		verifyUserID: userID,
		sessionUser:  &auth.SessionUser{ID: userID, DisplayName: "Test"},
	}

	stub := &stubUpdatesAPI{
		listResp: &updates.ListResponse{
			Updates: []updates.Update{},
			Total:   0,
			Page:    1,
			Limit:   100,
		},
	}

	router := newTestRouter(stub, authStub)
	req := httptest.NewRequest(http.MethodGet, "/v1/updates?limit=500", nil)
	req.Header.Set("Authorization", "Bearer test")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	// The handler passes the raw limit; the service is responsible for capping.
	// The stub response shows the capped limit of 100.
	if stub.lastListReq.Limit != 500 {
		t.Errorf("Request Limit: got %d, want 500", stub.lastListReq.Limit)
	}
}

func TestListUpdates_ServiceError(t *testing.T) {
	userID := uuid.New()
	authStub := &stubAuthAPI{
		verifyUserID: userID,
		sessionUser:  &auth.SessionUser{ID: userID, DisplayName: "Test"},
	}

	stub := &stubUpdatesAPI{
		listErr: &updatesListError{},
	}

	router := newTestRouter(stub, authStub)
	req := httptest.NewRequest(http.MethodGet, "/v1/updates", nil)
	req.Header.Set("Authorization", "Bearer test")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
}

func TestExplore_Success(t *testing.T) {
	userID := uuid.New()
	authStub := &stubAuthAPI{
		verifyUserID: userID,
		sessionUser:  &auth.SessionUser{ID: userID, DisplayName: "Test"},
	}

	stub := &stubUpdatesAPI{
		exploreResp: &updates.ExploreResponse{
			Indicators: []updates.ExploreIndicator{
				{ID: uuid.New(), PlaceName: "Nairobi", Code: "rain_intensity_trend", DisplayName: "Rain Intensity"},
			},
			Total: 1,
			Page:  1,
			Limit: 50,
		},
	}

	router := newTestRouter(stub, authStub)
	req := httptest.NewRequest(http.MethodGet, "/v1/explore?category=weather", nil)
	req.Header.Set("Authorization", "Bearer test")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	if stub.lastExploreReq == nil {
		t.Fatal("expected Explore to be called")
	}
	if stub.lastExploreReq.Category == nil || *stub.lastExploreReq.Category != "weather" {
		t.Errorf("Category: got %v, want weather", stub.lastExploreReq.Category)
	}
}

func TestExplore_MissingAuth(t *testing.T) {
	router := newTestRouter(&stubUpdatesAPI{}, &stubAuthAPI{})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/explore", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

type updatesListError struct{}

func (e *updatesListError) Error() string { return "list error" }
