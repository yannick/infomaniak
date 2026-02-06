package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cfg         ClientConfig
		wantBaseURL string
	}{
		{
			name:        "default base URL",
			cfg:         ClientConfig{Token: "tok"},
			wantBaseURL: defaultBaseURL,
		},
		{
			name:        "custom base URL",
			cfg:         ClientConfig{Token: "tok", BaseURL: "https://custom.api"},
			wantBaseURL: "https://custom.api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := NewClient(tt.cfg)
			if c.baseURL != tt.wantBaseURL {
				t.Errorf("baseURL = %q, want %q", c.baseURL, tt.wantBaseURL)
			}
			if c.token != tt.cfg.Token {
				t.Errorf("token = %q, want %q", c.token, tt.cfg.Token)
			}
		})
	}
}

func TestDecodeResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		status   int
		body     Response[string]
		wantData string
		wantErr  bool
	}{
		{
			name:     "success",
			status:   http.StatusOK,
			body:     Response[string]{Result: "success", Data: "hello"},
			wantData: "hello",
		},
		{
			name:   "api error",
			status: http.StatusUnprocessableEntity,
			body: Response[string]{
				Result: "error",
				Error:  &ErrorBody{Code: "validation_failed", Description: "bad input"},
			},
			wantErr: true,
		},
		{
			name:    "error without body",
			status:  http.StatusInternalServerError,
			body:    Response[string]{Result: "error"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			w.WriteHeader(tt.status)
			if err := json.NewEncoder(w).Encode(tt.body); err != nil {
				t.Fatal(err)
			}

			result, err := decodeResponse[string](w.Result())
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Data != tt.wantData {
				t.Errorf("data = %q, want %q", result.Data, tt.wantData)
			}
		})
	}
}

func TestDoRequestSetsHeaders(t *testing.T) {
	t.Parallel()

	var gotAuth, gotContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	c := NewClient(ClientConfig{Token: "test-token", BaseURL: srv.URL})
	resp, err := c.doRequest(context.Background(), "GET", "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	if gotAuth != "Bearer test-token" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer test-token")
	}
	if gotContentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", gotContentType, "application/json")
	}
}
