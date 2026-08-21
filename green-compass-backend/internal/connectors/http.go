package connectors

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout  = 15 * time.Second
	defaultAttempts = 3
)

type httpDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPOptions struct {
	BaseURL  string
	Client   httpDoer
	Attempts int
	Timeout  time.Duration
}

func requestJSON(ctx context.Context, client httpDoer, rawURL string, attempts int) ([]byte, error) {
	if attempts <= 0 {
		attempts = defaultAttempts
	}
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, fmt.Errorf("build provider request: %w", err)
		}
		resp, err := client.Do(req)
		if err == nil {
			body, readErr := io.ReadAll(resp.Body)
			closeErr := resp.Body.Close()
			if readErr != nil {
				return nil, fmt.Errorf("read provider response: %w", readErr)
			}
			if closeErr != nil {
				return nil, fmt.Errorf("close provider response: %w", closeErr)
			}
			if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
				return body, nil
			}
			lastErr = fmt.Errorf("provider returned status %d", resp.StatusCode)
			if resp.StatusCode != http.StatusTooManyRequests && resp.StatusCode < http.StatusInternalServerError {
				return nil, lastErr
			}
		} else {
			lastErr = fmt.Errorf("request provider: %w", err)
		}
		if attempt+1 < attempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt+1) * 100 * time.Millisecond):
			}
		}
	}
	return nil, lastErr
}

func configuredClient(options HTTPOptions) httpDoer {
	if options.Client != nil {
		return options.Client
	}
	timeout := options.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &http.Client{Timeout: timeout}
}
