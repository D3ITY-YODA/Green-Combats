package geo_test

import (
	"math"
	"strings"
	"testing"

	"green-compass-backend/pkg/geo"
)

func TestNewPoint(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lon     float64
		wantErr bool
	}{
		{name: "nairobi", lat: -1.2921, lon: 36.8219},
		{name: "equator prime meridian", lat: 0, lon: 0},
		{name: "north pole boundary", lat: 90, lon: 0},
		{name: "south pole boundary", lat: -90, lon: 0},
		{name: "antimeridian east", lat: 0, lon: 180},
		{name: "antimeridian west", lat: 0, lon: -180},
		{name: "latitude too high", lat: 90.000001, lon: 0, wantErr: true},
		{name: "latitude too low", lat: -90.5, lon: 0, wantErr: true},
		{name: "longitude too high", lat: 0, lon: 180.000001, wantErr: true},
		{name: "longitude too low", lat: 0, lon: -181, wantErr: true},
		{name: "nan latitude", lat: math.NaN(), lon: 0, wantErr: true},
		{name: "nan longitude", lat: 0, lon: math.NaN(), wantErr: true},
		{name: "infinite latitude", lat: math.Inf(1), lon: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := geo.NewPoint(tt.lat, tt.lon)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("NewPoint(%v, %v) expected error, got %+v", tt.lat, tt.lon, p)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewPoint(%v, %v) unexpected error: %v", tt.lat, tt.lon, err)
			}
			if p.Lat != tt.lat || p.Lon != tt.lon {
				t.Errorf("point = %+v, want lat=%v lon=%v", p, tt.lat, tt.lon)
			}
		})
	}
}

func TestMetersBetween(t *testing.T) {
	nairobi := mustPoint(t, -1.2921, 36.8219)
	mombasa := mustPoint(t, -4.0435, 39.6682)

	tests := []struct {
		name      string
		a         geo.Point
		b         geo.Point
		want      float64
		tolerance float64
	}{
		{
			name:      "same point is zero",
			a:         nairobi,
			b:         nairobi,
			want:      0,
			tolerance: 0.001,
		},
		{
			name:      "nairobi to mombasa",
			a:         nairobi,
			b:         mombasa,
			want:      440_200,
			tolerance: 2_500,
		},
		{
			name:      "one degree of longitude at equator",
			a:         mustPoint(t, 0, 36),
			b:         mustPoint(t, 0, 37),
			want:      111_195,
			tolerance: 100,
		},
		{
			name:      "one degree of latitude",
			a:         mustPoint(t, 0, 0),
			b:         mustPoint(t, 1, 0),
			want:      111_195,
			tolerance: 100,
		},
		{
			name:      "across the antimeridian",
			a:         mustPoint(t, 0, 179.5),
			b:         mustPoint(t, 0, -179.5),
			want:      111_195,
			tolerance: 100,
		},
		{
			name:      "pole to equator quarter circumference",
			a:         mustPoint(t, 90, 0),
			b:         mustPoint(t, 0, 0),
			want:      10_007_557,
			tolerance: 50,
		},
		{
			name:      "symmetric in either direction",
			a:         nairobi,
			b:         mombasa,
			want:      geo.MetersBetween(mombasa, nairobi),
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := geo.MetersBetween(tt.a, tt.b)
			diff := math.Abs(got - tt.want)
			if diff > tt.tolerance {
				t.Fatalf("MetersBetween = %.1f, want %.1f (tolerance %.1f)", got, tt.want, tt.tolerance)
			}
		})
	}
}

func TestWKT_PointOrderingIsLonLat(t *testing.T) {
	p := mustPoint(t, -1.2921, 36.8219)

	wkt := p.WKT()

	if !strings.HasPrefix(wkt, "POINT(") || !strings.HasSuffix(wkt, ")") {
		t.Fatalf("WKT() = %q, want POINT(x y) form", wkt)
	}
	inner := strings.TrimSuffix(strings.TrimPrefix(wkt, "POINT("), ")")
	parts := strings.Fields(inner)
	if len(parts) != 2 {
		t.Fatalf("WKT() = %q, want two coordinate fields", wkt)
	}
	if parts[0] != "36.8219" || parts[1] != "-1.2921" {
		t.Errorf("WKT() = %q, want longitude first then latitude", wkt)
	}
}

func TestSRID4326(t *testing.T) {
	if geo.SRID4326 != 4326 {
		t.Fatalf("SRID4326 = %d, want 4326", geo.SRID4326)
	}
}

func mustPoint(t *testing.T, lat, lon float64) geo.Point {
	t.Helper()
	p, err := geo.NewPoint(lat, lon)
	if err != nil {
		t.Fatalf("NewPoint(%v, %v) unexpected error: %v", lat, lon, err)
	}
	return p
}
