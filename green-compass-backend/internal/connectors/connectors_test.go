package connectors

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOpenMeteo_Fetch_PinsRequestFormatAndEmitsHourlyRecords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timezone") != "UTC" || q.Get("temperature_unit") != "celsius" || q.Get("precipitation_unit") != "mm" || q.Get("wind_speed_unit") != "ms" {
			t.Errorf("units/timezone query = %v", q)
		}
		if q.Get("start_date") != "2026-01-01" || q.Get("end_date") != "2026-01-01" || q.Get("latitude") != "-1.200000" || q.Get("longitude") != "36.800000" {
			t.Errorf("location/window query = %v", q)
		}
		_, _ = w.Write([]byte(`{"hourly":{"time":["2026-01-01T00:00","2026-01-01T01:00"],"temperature_2m":[21.5,22.0],"relative_humidity_2m":[50,51],"precipitation":[0,1.2],"rain":[0,1.2],"showers":[0,0],"weather_code":[0,61],"wind_speed_10m":[2.1,2.2],"soil_moisture_0_to_1cm":[0.2,0.21]}}`))
	}))
	defer server.Close()

	connector := NewOpenMeteo(HTTPOptions{BaseURL: server.URL, Client: server.Client(), Attempts: 1})
	records, err := connector.Fetch(context.Background(), FetchRequest{Lat: -1.2, Lon: 36.8, From: date(2026, 1, 1), To: date(2026, 1, 1)})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(records) != 2 || !records[0].ObservedAt.Equal(date(2026, 1, 1)) || records[0].SourceURL == "" || records[0].ContentType != "application/json" {
		t.Errorf("records = %+v", records)
	}
}

func TestNASAPower_Fetch_PinsRequestFormatAndEmitsDailyRecords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("community") != "AG" || q.Get("format") != "JSON" || q.Get("time-standard") != "UTC" || q.Get("start") != "20260101" || q.Get("end") != "20260102" {
			t.Errorf("NASA POWER query = %v", q)
		}
		_, _ = w.Write([]byte(`{"properties":{"parameter":{"T2M":{"20260101":20,"20260102":21},"PRECTOTCORR":{"20260101":1,"20260102":0},"RH2M":{"20260101":55,"20260102":56},"WS2M":{"20260101":2,"20260102":3},"ALLSKY_SFC_SW_DWN":{"20260101":5,"20260102":6}}}}`))
	}))
	defer server.Close()

	connector := NewNASAPower(HTTPOptions{BaseURL: server.URL, Client: server.Client(), Attempts: 1})
	records, err := connector.Fetch(context.Background(), FetchRequest{Lat: -1.2, Lon: 36.8, From: date(2026, 1, 1), To: date(2026, 1, 2)})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(records) != 2 || !records[0].ObservedAt.Equal(date(2026, 1, 1)) || !records[1].ObservedAt.Equal(date(2026, 1, 2)) {
		t.Errorf("records = %+v", records)
	}
}

func TestConnectorParsers_RejectIncompleteProviderData(t *testing.T) {
	tests := []struct {
		name  string
		parse func([]byte, string) ([]RawObservation, error)
		body  string
	}{
		{name: "Open-Meteo missing aligned variable", parse: parseOpenMeteo, body: `{"hourly":{"time":["2026-01-01T00:00"]}}`},
		{name: "NASA POWER missing required variable", parse: parseNASAPower, body: `{"properties":{"parameter":{"T2M":{"20260101":20}}}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.parse([]byte(tt.body), "https://provider.test"); err == nil {
				t.Fatal("parse error = nil, want error")
			}
		})
	}
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
