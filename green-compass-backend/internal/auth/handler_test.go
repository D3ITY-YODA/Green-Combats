package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/auth"
	"green-compass-backend/internal/users"
)

type stubAPI struct {
	registerResult *auth.RegisterResult
	registerErr    error
	loginPair      *auth.TokenPair
	loginErr       error
	refreshPair    *auth.TokenPair
	refreshErr     error
	logoutErr      error
	sessionUser    *auth.SessionUser
	verifyUserID   uuid.UUID
	verifyErr      error

	lastRegisterInput users.RegisterInput
}

func (s *stubAPI) Register(_ context.Context, in users.RegisterInput) (*auth.RegisterResult, error) {
	s.lastRegisterInput = in
	if s.registerErr != nil {
		return nil, s.registerErr
	}
	return s.registerResult, nil
}

func (s *stubAPI) Login(context.Context, string, string) (*auth.TokenPair, error) {
	return s.loginPair, s.loginErr
}

func (s *stubAPI) Refresh(context.Context, string) (*auth.TokenPair, error) {
	return s.refreshPair, s.refreshErr
}

func (s *stubAPI) Logout(context.Context, string) error { return s.logoutErr }

func (s *stubAPI) VerifyAccessToken(string) (uuid.UUID, error) {
	return s.verifyUserID, s.verifyErr
}

func (s *stubAPI) SessionUser(context.Context, uuid.UUID) (*auth.SessionUser, error) {
	return s.sessionUser, nil
}

func newTestRouter(svc auth.API) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	auth.NewHandler(svc).RegisterRoutes(router)
	return router
}

func doJSON(router *gin.Engine, method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func okPair() *auth.TokenPair {
	return &auth.TokenPair{
		AccessToken:      "access.jwt.value",
		RefreshToken:     "refresh-token-value",
		AccessExpiresAt:  timeAt(1000),
		RefreshExpiresAt: timeAt(2000),
	}
}

func TokenPairValue() auth.TokenPair {
	return *okPair()
}

func timeAt(unix int64) time.Time {
	return time.Unix(unix, 0).UTC()
}

func TestHandler_Register(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name       string
		body       string
		stubErr    error
		wantStatus int
	}{
		{
			name:       "success",
			body:       `{"email":"a@b.test","display_name":"A","password":"long-enough"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid data",
			body:       `{"display_name":"A","password":"long-enough"}`,
			stubErr:    users.ErrInvalidData,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "email taken",
			body:       `{"email":"taken@b.test","display_name":"A","password":"long-enough"}`,
			stubErr:    users.ErrEmailTaken,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "phone taken",
			body:       `{"phone_number":"+254700000001","display_name":"A","password":"long-enough"}`,
			stubErr:    users.ErrPhoneTaken,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "malformed json",
			body:       `{not-json`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &stubAPI{
				registerErr: tt.stubErr,
				registerResult: &auth.RegisterResult{
					User:   &users.User{ID: userID, DisplayName: "A", Language: "en"},
					Tokens: TokenPairValue(),
				},
			}
			router := newTestRouter(stub)

			rec := doJSON(router, http.MethodPost, "/v1/auth/register", tt.body, nil)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantStatus == http.StatusCreated {
				var resp struct {
					Data struct {
						User struct {
							ID string `json:"id"`
						} `json:"user"`
						Tokens struct {
							AccessToken string `json:"access_token"`
						} `json:"tokens"`
					} `json:"data"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.Data.User.ID != userID.String() || resp.Data.Tokens.AccessToken != "access.jwt.value" {
					t.Errorf("response body mismatch: %s", rec.Body.String())
				}
			}
		})
	}
}

func TestHandler_Register_PassesFieldsToService(t *testing.T) {
	stub := &stubAPI{registerResult: &auth.RegisterResult{
		User:   &users.User{ID: uuid.New()},
		Tokens: TokenPairValue(),
	}}
	router := newTestRouter(stub)

	doJSON(router, http.MethodPost, "/v1/auth/register",
		`{"phone_number":"+254711222333","display_name":"Zawadi","password":"long-enough","language":"sw"}`, nil)

	in := stub.lastRegisterInput
	if in.PhoneNumber != "+254711222333" || in.DisplayName != "Zawadi" || in.Password != "long-enough" || in.Language != "sw" {
		t.Errorf("service received %+v", in)
	}
}

func TestHandler_LoginAndRefreshAndLogout(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		stubErr    error
		wantStatus int
	}{
		{name: "login success", method: http.MethodPost, path: "/v1/auth/login", body: `{"identifier":"a@b.test","password":"pw-12345678"}`, wantStatus: http.StatusOK},
		{name: "login bad credentials", method: http.MethodPost, path: "/v1/auth/login", body: `{"identifier":"a@b.test","password":"nope12345"}`, stubErr: auth.ErrInvalidCredentials, wantStatus: http.StatusUnauthorized},
		{name: "refresh success", method: http.MethodPost, path: "/v1/auth/refresh", body: `{"refresh_token":"r"}`, wantStatus: http.StatusOK},
		{name: "refresh reuse", method: http.MethodPost, path: "/v1/auth/refresh", body: `{"refresh_token":"r"}`, stubErr: auth.ErrTokenReuse, wantStatus: http.StatusUnauthorized},
		{name: "refresh invalid", method: http.MethodPost, path: "/v1/auth/refresh", body: `{"refresh_token":"r"}`, stubErr: auth.ErrInvalidToken, wantStatus: http.StatusUnauthorized},
		{name: "logout success", method: http.MethodPost, path: "/v1/auth/logout", body: `{"refresh_token":"r"}`, wantStatus: http.StatusNoContent},
		{name: "logout invalid", method: http.MethodPost, path: "/v1/auth/logout", body: `{"refresh_token":"r"}`, stubErr: auth.ErrInvalidToken, wantStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &stubAPI{
				loginErr: tt.stubErr, loginPair: okPair(),
				refreshErr: tt.stubErr, refreshPair: okPair(),
				logoutErr: tt.stubErr,
			}
			router := newTestRouter(stub)

			rec := doJSON(router, tt.method, tt.path, tt.body, nil)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}

func TestHandler_Me_RequiresBearer(t *testing.T) {
	stub := &stubAPI{
		sessionUser: &auth.SessionUser{ID: uuid.New(), DisplayName: "Me"},
	}
	router := newTestRouter(stub)

	tests := []struct {
		name       string
		headers    map[string]string
		badToken   bool
		wantStatus int
	}{
		{name: "no authorization header", wantStatus: http.StatusUnauthorized},
		{name: "wrong scheme", headers: map[string]string{"Authorization": "Basic abc"}, wantStatus: http.StatusUnauthorized},
		{name: "invalid token", headers: map[string]string{"Authorization": "Bearer bad"}, badToken: true, wantStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.badToken {
				stub.verifyErr = errors.New("bad token")
			} else {
				stub.verifyErr = nil
				stub.verifyUserID = uuid.New()
			}
			rec := doJSON(router, http.MethodGet, "/v1/auth/me", "", tt.headers)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_Me_ReturnsSessionUser(t *testing.T) {
	id := uuid.New()
	stub := &stubAPI{
		sessionUser:  &auth.SessionUser{ID: id, DisplayName: "Me", Language: "en", IsPlatformAdmin: true},
		verifyUserID: id,
	}
	router := newTestRouter(stub)

	rec := doJSON(router, http.MethodGet, "/v1/auth/me", "", map[string]string{"Authorization": "Bearer good"})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var resp struct {
		Data struct {
			User struct {
				ID              string `json:"id"`
				IsPlatformAdmin bool   `json:"is_platform_admin"`
			} `json:"user"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Data.User.ID != id.String() || !resp.Data.User.IsPlatformAdmin {
		t.Errorf("me response mismatch: %s", rec.Body.String())
	}
}
