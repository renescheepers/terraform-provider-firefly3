// Copyright (c) HashiCorp, Inc.

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Rule struct {
	ID             string        `json:"id,omitempty"`
	CreatedAt      string        `json:"created_at,omitempty"`
	UpdatedAt      string        `json:"updated_at,omitempty"`
	Title          string        `json:"title"`
	Description    string        `json:"description,omitempty"`
	RuleGroupID    string        `json:"rule_group_id"`
	RuleGroupTitle string        `json:"rule_group_title,omitempty"`
	Order          int32         `json:"order,omitempty"`
	Trigger        string        `json:"trigger"`
	Active         bool          `json:"active"`
	Strict         bool          `json:"strict"`
	StopProcessing bool          `json:"stop_processing"`
	Triggers       []RuleTrigger `json:"triggers"`
	Actions        []RuleAction  `json:"actions"`
}

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

type RuleSingle struct {
	Data RuleData `json:"data"`
}

type RuleData struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes Rule   `json:"attributes"`
}

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

func (c *Client) DeleteRule(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/api/v1/rules/"+id, nil)
	return err
}
