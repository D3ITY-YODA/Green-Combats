package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

const nasaPowerURL = "https://power.larc.nasa.gov/api/temporal/daily/point"

var nasaPowerVariables = []string{"T2M", "PRECTOTCORR", "RH2M", "WS2M", "ALLSKY_SFC_SW_DWN"}

type NASAPower struct {
	baseURL  string
	client   httpDoer
	attempts int
}

func NewNASAPower(options HTTPOptions) *NASAPower {
	baseURL := options.BaseURL
	if baseURL == "" {
		baseURL = nasaPowerURL
	}
	attempts := options.Attempts
	if attempts <= 0 {
		attempts = defaultAttempts
	}
	return &NASAPower{baseURL: baseURL, client: configuredClient(options), attempts: attempts}
}

func (*NASAPower) Code() string { return "nasa_power" }

func (c *NASAPower) Fetch(ctx context.Context, request FetchRequest) ([]RawObservation, error) {
	if err := validateWindow(request); err != nil {
		return nil, err
	}
	endpoint, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse NASA POWER URL: %w", err)
	}
	query := endpoint.Query()
	query.Set("parameters", strings.Join(nasaPowerVariables, ","))
	query.Set("community", "AG")
	query.Set("longitude", fmt.Sprintf("%.6f", request.Lon))
	query.Set("latitude", fmt.Sprintf("%.6f", request.Lat))
	query.Set("start", request.From.UTC().Format("20060102"))
	query.Set("end", request.To.UTC().Format("20060102"))
	query.Set("format", "JSON")
	query.Set("time-standard", "UTC")
	endpoint.RawQuery = query.Encode()
	body, err := requestJSON(ctx, c.client, endpoint.String(), c.attempts)
	if err != nil {
		return nil, fmt.Errorf("fetch NASA POWER: %w", err)
	}
	return parseNASAPower(body, endpoint.String())
}

func parseNASAPower(body []byte, sourceURL string) ([]RawObservation, error) {
	var response struct {
		Properties struct {
			Parameter map[string]map[string]json.RawMessage `json:"parameter"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("decode NASA POWER response: %w", err)
	}
	if response.Properties.Parameter == nil {
		return nil, fmt.Errorf("decode NASA POWER response: parameter data is required")
	}
	for _, variable := range nasaPowerVariables {
		if len(response.Properties.Parameter[variable]) == 0 {
			return nil, fmt.Errorf("decode NASA POWER response: %s is required", variable)
		}
	}
	dates := make([]string, 0, len(response.Properties.Parameter[nasaPowerVariables[0]]))
	for date := range response.Properties.Parameter[nasaPowerVariables[0]] {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	result := make([]RawObservation, 0, len(dates))
	for _, date := range dates {
		observedAt, err := time.Parse("20060102", date)
		if err != nil {
			return nil, fmt.Errorf("decode NASA POWER response: invalid date %q: %w", date, err)
		}
		payload := map[string]json.RawMessage{"date": mustJSON(date)}
		for _, variable := range nasaPowerVariables {
			value, ok := response.Properties.Parameter[variable][date]
			if !ok {
				return nil, fmt.Errorf("decode NASA POWER response: %s missing %s", variable, date)
			}
			payload[variable] = value
		}
		encoded, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("encode NASA POWER observation: %w", err)
		}
		result = append(result, RawObservation{ObservedAt: observedAt.UTC(), SourceURL: sourceURL, ContentType: "application/json", Payload: encoded})
	}
	return result, nil
}
