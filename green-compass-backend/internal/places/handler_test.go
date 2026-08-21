package places_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"green-compass-backend/internal/places"
)

type stubPlaceAPI struct {
	byID    *places.Place
	byIDErr error
	nearby  []places.WithDistance
	nearErr error
}

func (s *stubPlaceAPI) Create(context.Context, places.CreateInput) (*places.Place, error) {
	return nil, nil
}
func (s *stubPlaceAPI) ByID(context.Context, uuid.UUID) (*places.Place, error) {
	return s.byID, s.byIDErr
}
func (s *stubPlaceAPI) Nearby(context.Context, float64, float64, float64, int) ([]places.WithDistance, error) {
	return s.nearby, s.nearErr
}
func (s *stubPlaceAPI) Update(context.Context, uuid.UUID, places.UpdateInput, places.Caller) (*places.Place, error) {
	return nil, nil
}
func (s *stubPlaceAPI) Search(context.Context, string, int) ([]places.Place, error) {
	return nil, nil
}

func TestHandler_PublicPlaceRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	placeID := uuid.New()

	tests := []struct {
		name       string
		path       string
		stub       stubPlaceAPI
		wantStatus int
	}{
		{name: "get place", path: "/v1/places/" + placeID.String(), stub: stubPlaceAPI{byID: &places.Place{ID: placeID, Name: "Kisumu"}}, wantStatus: http.StatusOK},
		{name: "not found", path: "/v1/places/" + placeID.String(), stub: stubPlaceAPI{byIDErr: places.ErrNotFound}, wantStatus: http.StatusNotFound},
		{name: "invalid id", path: "/v1/places/nope", wantStatus: http.StatusBadRequest},
		{name: "nearby needs coordinates", path: "/v1/places/nearby", wantStatus: http.StatusBadRequest},
		{name: "nearby validation", path: "/v1/places/nearby?lat=bad&lon=36", wantStatus: http.StatusBadRequest},
		{name: "nearby success", path: "/v1/places/nearby?lat=-1.2&lon=36.8", stub: stubPlaceAPI{nearby: []places.WithDistance{{Place: places.Place{ID: placeID, Name: "Kisumu"}, DistanceMeters: 123}}}, wantStatus: http.StatusOK},
		{name: "nearby service validation", path: "/v1/places/nearby?lat=-1.2&lon=36.8&radius_m=999999", stub: stubPlaceAPI{nearErr: errors.Join(places.ErrInvalidData, errors.New("radius"))}, wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			places.NewHandler(&tt.stub, nil).RegisterRoutes(router)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}
