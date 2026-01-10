// Copyright 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is the Firefly3 API client
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new Firefly3 API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Rule represents a Firefly3 rule
type Rule struct {
	ID             string          `json:"id,omitempty"`
	CreatedAt      string          `json:"created_at,omitempty"`
	UpdatedAt      string          `json:"updated_at,omitempty"`
	Title          string          `json:"title"`
	Description    string          `json:"description,omitempty"`
	RuleGroupID    string          `json:"rule_group_id"`
	RuleGroupTitle string          `json:"rule_group_title,omitempty"`
	Order          int32           `json:"order,omitempty"`
	Trigger        string          `json:"trigger"`
	Active         bool            `json:"active"`
	Strict         bool            `json:"strict"`
	StopProcessing bool            `json:"stop_processing"`
	Triggers       []RuleTrigger   `json:"triggers"`
	Actions        []RuleAction    `json:"actions"`
}

// RuleTrigger represents a trigger within a rule
type RuleTrigger struct {
	ID             string `json:"id,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	Type           string `json:"type"`
	Value          string `json:"value"`
	Order          int32  `json:"order,omitempty"`
	Active         bool   `json:"active"`
	Prohibited     bool   `json:"prohibited"`
	StopProcessing bool   `json:"stop_processing"`
}

// RuleAction represents an action within a rule
type RuleAction struct {
	ID             string `json:"id,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	Type           string `json:"type"`
	Value          string `json:"value"`
	Order          int32  `json:"order,omitempty"`
	Active         bool   `json:"active"`
	StopProcessing bool   `json:"stop_processing"`
}

// RuleSingle represents the API response for a single rule
type RuleSingle struct {
	Data RuleData `json:"data"`
}

// RuleData wraps the rule attributes
type RuleData struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes Rule   `json:"attributes"`
}

// CreateRule creates a new rule
func (c *Client) CreateRule(ctx context.Context, rule *Rule) (*Rule, error) {
	respBody, err := c.doRequest(ctx, http.MethodPost, "/api/v1/rules", rule)
	if err != nil {
		return nil, err
	}

	var result RuleSingle
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	createdRule := result.Data.Attributes
	createdRule.ID = result.Data.ID
	return &createdRule, nil
}

// GetRule retrieves a rule by ID
func (c *Client) GetRule(ctx context.Context, id string) (*Rule, error) {
	respBody, err := c.doRequest(ctx, http.MethodGet, "/api/v1/rules/"+id, nil)
	if err != nil {
		return nil, err
	}

	var result RuleSingle
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	rule := result.Data.Attributes
	rule.ID = result.Data.ID
	return &rule, nil
}

// UpdateRule updates an existing rule
func (c *Client) UpdateRule(ctx context.Context, id string, rule *Rule) (*Rule, error) {
	respBody, err := c.doRequest(ctx, http.MethodPut, "/api/v1/rules/"+id, rule)
	if err != nil {
		return nil, err
	}

	var result RuleSingle
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	updatedRule := result.Data.Attributes
	updatedRule.ID = result.Data.ID
	return &updatedRule, nil
}

// DeleteRule deletes a rule by ID
func (c *Client) DeleteRule(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/api/v1/rules/"+id, nil)
	return err
}
