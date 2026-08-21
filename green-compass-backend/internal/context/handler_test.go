package context_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	gcauth "green-compass-backend/internal/auth"
	gcctx "green-compass-backend/internal/context"
	"green-compass-backend/internal/updates"
	"green-compass-backend/internal/users"
)

type stubAuthAPI struct {
	verifyUserID uuid.UUID
	verifyErr    error
	sessionUser  *gcauth.SessionUser
}

func (s *stubAuthAPI) Register(_ context.Context, _ users.RegisterInput) (*gcauth.RegisterResult, error) {
	return nil, nil
}

func (s *stubAuthAPI) Login(_ context.Context, _, _ string) (*gcauth.TokenPair, error) {
	return nil, nil
}

func (s *stubAuthAPI) Refresh(_ context.Context, _ string) (*gcauth.TokenPair, error) {
	return nil, nil
}

func (s *stubAuthAPI) Logout(_ context.Context, _ string) error { return nil }

func (s *stubAuthAPI) VerifyAccessToken(_ string) (uuid.UUID, error) {
	return s.verifyUserID, s.verifyErr
}

func (s *stubAuthAPI) SessionUser(_ context.Context, _ uuid.UUID) (*gcauth.SessionUser, error) {
	return s.sessionUser, nil
}

type stubContextAPI struct {
	resp *gcctx.ContextResponse
	err  error
}

func (s *stubContextAPI) ResolveContext(_ context.Context, _ gcctx.ResolveContextRequest) (*gcctx.ContextResponse, error) {
	return s.resp, s.err
}

func newTestRouter(ctxAPI gcctx.API, authAPI gcauth.API) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	gcctx.NewHandler(ctxAPI).RegisterRoutes(router, gcauth.Middleware(authAPI))
	return router
}

func TestGetToday_MissingAuth(t *testing.T) {
	authStub := &stubAuthAPI{}
	router := newTestRouter(&stubContextAPI{}, authStub)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/context/today", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestGetToday_InvalidToken(t *testing.T) {
	authStub := &stubAuthAPI{verifyErr: errInvalidToken}
	router := newTestRouter(&stubContextAPI{}, authStub)

	req := httptest.NewRequest(http.MethodGet, "/v1/context/today", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestGetToday_Success(t *testing.T) {
	userID := uuid.New()
	placeID := uuid.New()

	authStub := &stubAuthAPI{
		verifyUserID: userID,
		sessionUser:  &gcauth.SessionUser{ID: userID, DisplayName: "Test"},
	}

	ctxAPI := &stubContextAPI{
		resp: &gcctx.ContextResponse{
			Place: gcctx.PlaceInfo{
				ID:        placeID,
				Name:      "Nairobi",
				PlaceType: "community",
				Lat:       -1.2,
				Lon:       36.8,
				IsPrimary: true,
			},
			Update: &updates.Update{
				ID:          uuid.New(),
				PlaceID:     placeID,
				PlaceName:   "Nairobi",
				ContentType: "today",
				Headline:    "Today: Rain expected",
				BodyText:    "Conditions require attention.",
			},
		},
	}

	router := newTestRouter(ctxAPI, authStub)
	req := httptest.NewRequest(http.MethodGet, "/v1/context/today", nil)
	req.Header.Set("Authorization", "Bearer good")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d (body: %s)", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp struct {
		Status string `json:"status"`
		Data   struct {
			Place struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				PlaceType string `json:"place_type"`
				IsPrimary bool   `json:"is_primary"`
			} `json:"place"`
			Update struct {
				Headline string `json:"headline"`
			} `json:"update"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("status: got %q, want %q", resp.Status, "success")
	}
	if resp.Data.Place.Name != "Nairobi" {
		t.Errorf("place name: got %q, want %q", resp.Data.Place.Name, "Nairobi")
	}
	if resp.Data.Update.Headline != "Today: Rain expected" {
		t.Errorf("headline: got %q, want %q", resp.Data.Update.Headline, "Today: Rain expected")
	}
}

func TestGetToday_NoPlace(t *testing.T) {
	userID := uuid.New()

	authStub := &stubAuthAPI{
		verifyUserID: userID,
		sessionUser:  &gcauth.SessionUser{ID: userID, DisplayName: "Test"},
	}
	ctxAPI := &stubContextAPI{err: gcctx.ErrNoPlace}

	router := newTestRouter(ctxAPI, authStub)
	req := httptest.NewRequest(http.MethodGet, "/v1/context/today", nil)
	req.Header.Set("Authorization", "Bearer good")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d (body: %s)", rec.Code, http.StatusNotFound, rec.Body.String())
	}
}

func TestGetToday_InvalidPlaceID(t *testing.T) {
	userID := uuid.New()

	authStub := &stubAuthAPI{
		verifyUserID: userID,
		sessionUser:  &gcauth.SessionUser{ID: userID, DisplayName: "Test"},
	}
	ctxAPI := &stubContextAPI{resp: &gcctx.ContextResponse{}}

	router := newTestRouter(ctxAPI, authStub)
	req := httptest.NewRequest(http.MethodGet, "/v1/context/today?place_id=not-a-uuid", nil)
	req.Header.Set("Authorization", "Bearer good")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

type stubError struct{ msg string }

func (e *stubError) Error() string { return e.msg }

var errInvalidToken = &stubError{msg: "invalid token"}
