package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListDomains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		accountID string
		response  Response[[]Domain]
		status    int
		wantLen   int
		wantErr   bool
	}{
		{
			name:      "success with domains",
			accountID: "123",
			status:    http.StatusOK,
			response: Response[[]Domain]{
				Result: "success",
				Data: []Domain{
					{Name: "example.ch", TLD: "ch"},
					{Name: "test.com", TLD: "com"},
				},
			},
			wantLen: 2,
		},
		{
			name:      "empty list",
			accountID: "456",
			status:    http.StatusOK,
			response:  Response[[]Domain]{Result: "success", Data: []Domain{}},
			wantLen:   0,
		},
		{
			name:      "api error",
			accountID: "999",
			status:    http.StatusForbidden,
			response: Response[[]Domain]{
				Result: "error",
				Error:  &ErrorBody{Code: "forbidden", Description: "access denied"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantPath := "/2/domains/accounts/" + tt.accountID + "/domains"
				if r.URL.Path != wantPath {
					t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
				}
				if r.Method != http.MethodGet {
					t.Errorf("method = %q, want GET", r.Method)
				}
				w.WriteHeader(tt.status)
				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			t.Cleanup(srv.Close)

			c := NewClient(ClientConfig{Token: "tok", BaseURL: srv.URL})
			domains, err := c.ListDomains(context.Background(), tt.accountID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(domains) != tt.wantLen {
				t.Errorf("got %d domains, want %d", len(domains), tt.wantLen)
			}
		})
	}
}

func TestShowDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		domain   string
		response Response[Domain]
		status   int
		wantName string
		wantErr  bool
	}{
		{
			name:   "success",
			domain: "example.ch",
			status: http.StatusOK,
			response: Response[Domain]{
				Result: "success",
				Data:   Domain{Name: "example.ch", TLD: "ch", ExpiresAt: 1734444000},
			},
			wantName: "example.ch",
		},
		{
			name:   "not found",
			domain: "nope.ch",
			status: http.StatusNotFound,
			response: Response[Domain]{
				Result: "error",
				Error:  &ErrorBody{Code: "object_not_found", Description: "Object not found"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantPath := "/2/domains/" + tt.domain
				if r.URL.Path != wantPath {
					t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
				}
				w.WriteHeader(tt.status)
				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			t.Cleanup(srv.Close)

			c := NewClient(ClientConfig{Token: "tok", BaseURL: srv.URL})
			domain, err := c.ShowDomain(context.Background(), tt.domain)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if domain.Name != tt.wantName {
				t.Errorf("name = %q, want %q", domain.Name, tt.wantName)
			}
		})
	}
}

func TestUpdateNameservers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		domain   string
		input    UpdateNameserversInput
		response Response[any]
		status   int
		wantErr  bool
	}{
		{
			name:   "success",
			domain: "example.ch",
			input: UpdateNameserversInput{
				Nameservers:          []string{"ns1.example.ch", "ns2.example.ch"},
				VerifyNSAvailability: false,
			},
			status:   http.StatusOK,
			response: Response[any]{Result: "success"},
		},
		{
			name:   "validation error",
			domain: "example.ch",
			input: UpdateNameserversInput{
				Nameservers: []string{"invalid"},
			},
			status: http.StatusUnprocessableEntity,
			response: Response[any]{
				Result: "error",
				Error:  &ErrorBody{Code: "validation_failed", Description: "invalid nameserver"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantPath := "/2/domains/" + tt.domain + "/nameservers"
				if r.URL.Path != wantPath {
					t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
				}
				if r.Method != http.MethodPut {
					t.Errorf("method = %q, want PUT", r.Method)
				}

				var body UpdateNameserversInput
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("decode request body: %v", err)
				}
				if len(body.Nameservers) != len(tt.input.Nameservers) {
					t.Errorf("nameservers count = %d, want %d", len(body.Nameservers), len(tt.input.Nameservers))
				}

				w.WriteHeader(tt.status)
				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			t.Cleanup(srv.Close)

			c := NewClient(ClientConfig{Token: "tok", BaseURL: srv.URL})
			err := c.UpdateNameservers(context.Background(), tt.domain, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
