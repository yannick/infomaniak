package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// ListDomains returns all domains for the given account.
func (c *Client) ListDomains(ctx context.Context, accountID string) ([]Domain, error) {
	path := fmt.Sprintf("/2/domains/accounts/%s/domains", accountID)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("list domains for account %s: %w", accountID, err)
	}

	result, err := decodeResponse[[]Domain](resp)
	if err != nil {
		return nil, fmt.Errorf("list domains for account %s: %w", accountID, err)
	}

	return result.Data, nil
}

// ShowDomain returns details for a single domain.
func (c *Client) ShowDomain(ctx context.Context, domain string) (*Domain, error) {
	path := fmt.Sprintf("/2/domains/%s", domain)

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("show domain %s: %w", domain, err)
	}

	result, err := decodeResponse[Domain](resp)
	if err != nil {
		return nil, fmt.Errorf("show domain %s: %w", domain, err)
	}

	return &result.Data, nil
}

// UpdateNameservers sets the nameservers for a domain.
func (c *Client) UpdateNameservers(ctx context.Context, domain string, input UpdateNameserversInput) error {
	path := fmt.Sprintf("/2/domains/%s/nameservers", domain)

	body, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("marshal nameserver update for %s: %w", domain, err)
	}

	resp, err := c.doRequest(ctx, "PUT", path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("update nameservers for %s: %w", domain, err)
	}

	if _, err := decodeResponse[any](resp); err != nil {
		return fmt.Errorf("update nameservers for %s: %w", domain, err)
	}

	return nil
}
