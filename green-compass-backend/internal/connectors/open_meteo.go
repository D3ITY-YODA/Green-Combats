package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const openMeteoURL = "https://api.open-meteo.com/v1/forecast"

var openMeteoVariables = []string{
	"temperature_2m", "relative_humidity_2m", "precipitation", "rain", "showers",
	"weather_code", "wind_speed_10m", "soil_moisture_0_to_1cm",
}

type OpenMeteo struct {
	baseURL  string
	client   httpDoer
	attempts int
}

func NewOpenMeteo(options HTTPOptions) *OpenMeteo {
	baseURL := options.BaseURL
	if baseURL == "" {
		baseURL = openMeteoURL
	}
	attempts := options.Attempts
	if attempts <= 0 {
		attempts = defaultAttempts
	}
	return &OpenMeteo{baseURL: baseURL, client: configuredClient(options), attempts: attempts}
}

func (*OpenMeteo) Code() string { return "open_meteo" }

func (c *OpenMeteo) Fetch(ctx context.Context, request FetchRequest) ([]RawObservation, error) {
	if err := validateWindow(request); err != nil {
		return nil, err
	}
	endpoint, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse Open-Meteo URL: %w", err)
	}
	query := endpoint.Query()
	query.Set("latitude", fmt.Sprintf("%.6f", request.Lat))
	query.Set("longitude", fmt.Sprintf("%.6f", request.Lon))
	query.Set("hourly", strings.Join(openMeteoVariables, ","))
	query.Set("start_date", request.From.UTC().Format("2006-01-02"))
	query.Set("end_date", request.To.UTC().Format("2006-01-02"))
	query.Set("timezone", "UTC")
	query.Set("temperature_unit", "celsius")
	query.Set("precipitation_unit", "mm")
	query.Set("wind_speed_unit", "ms")
	endpoint.RawQuery = query.Encode()

	body, err := requestJSON(ctx, c.client, endpoint.String(), c.attempts)
	if err != nil {
		return nil, fmt.Errorf("fetch Open-Meteo: %w", err)
	}
	return parseOpenMeteo(body, endpoint.String())
}

func parseOpenMeteo(body []byte, sourceURL string) ([]RawObservation, error) {
	var response struct {
		Hourly map[string]json.RawMessage `json:"hourly"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("decode Open-Meteo response: %w", err)
	}
	if response.Hourly == nil {
		return nil, fmt.Errorf("decode Open-Meteo response: hourly data is required")
	}
	var timestamps []string
	if err := json.Unmarshal(response.Hourly["time"], &timestamps); err != nil || len(timestamps) == 0 {
		return nil, fmt.Errorf("decode Open-Meteo response: hourly time array is required")
	}
	values := make(map[string][]json.RawMessage, len(openMeteoVariables))
	for _, variable := range openMeteoVariables {
		var series []json.RawMessage
		if err := json.Unmarshal(response.Hourly[variable], &series); err != nil || len(series) != len(timestamps) {
			return nil, fmt.Errorf("decode Open-Meteo response: %s must align with time", variable)
		}
		values[variable] = series
	}
	result := make([]RawObservation, 0, len(timestamps))
	for i, rawTime := range timestamps {
		observedAt, err := time.Parse("2006-01-02T15:04", rawTime)
		if err != nil {
			return nil, fmt.Errorf("decode Open-Meteo response: invalid time %q: %w", rawTime, err)
		}
		payload := map[string]json.RawMessage{"time": mustJSON(rawTime)}
		for _, variable := range openMeteoVariables {
			payload[variable] = values[variable][i]
		}
		encoded, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("encode Open-Meteo observation: %w", err)
		}
		result = append(result, RawObservation{ObservedAt: observedAt.UTC(), SourceURL: sourceURL, ContentType: "application/json", Payload: encoded})
	}
	return result, nil
}

func validateWindow(request FetchRequest) error {
	if request.From.IsZero() || request.To.IsZero() || request.To.Before(request.From) {
		return fmt.Errorf("invalid fetch time window")
	}
	return nil
}

func mustJSON(value string) json.RawMessage {
	encoded, _ := json.Marshal(value)
	return encoded
}
