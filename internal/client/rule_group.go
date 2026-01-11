// Copyright (c) HashiCorp, Inc.

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type RuleGroup struct {
	ID          string `json:"id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Order       int32  `json:"order,omitempty"`
	Active      bool   `json:"active"`
}

type RuleGroupSingle struct {
	Data RuleGroupData `json:"data"`
}

type RuleGroupData struct {
	Type       string    `json:"type"`
	ID         string    `json:"id"`
	Attributes RuleGroup `json:"attributes"`
}

func (c *Client) CreateRuleGroup(ctx context.Context, ruleGroup *RuleGroup) (*RuleGroup, error) {
	respBody, err := c.doRequest(ctx, http.MethodPost, "/api/v1/rule-groups", ruleGroup)
	if err != nil {
		return nil, err
	}

	var result RuleGroupSingle
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	createdRuleGroup := result.Data.Attributes
	createdRuleGroup.ID = result.Data.ID
	return &createdRuleGroup, nil
}

// GetRuleGroup retrieves a rule group by ID
func (c *Client) GetRuleGroup(ctx context.Context, id string) (*RuleGroup, error) {
	respBody, err := c.doRequest(ctx, http.MethodGet, "/api/v1/rule-groups/"+id, nil)
	if err != nil {
		return nil, err
	}

	var result RuleGroupSingle
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	ruleGroup := result.Data.Attributes
	ruleGroup.ID = result.Data.ID
	return &ruleGroup, nil
}

func (c *Client) UpdateRuleGroup(ctx context.Context, id string, ruleGroup *RuleGroup) (*RuleGroup, error) {
	respBody, err := c.doRequest(ctx, http.MethodPut, "/api/v1/rule-groups/"+id, ruleGroup)
	if err != nil {
		return nil, err
	}

	var result RuleGroupSingle
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	updatedRuleGroup := result.Data.Attributes
	updatedRuleGroup.ID = result.Data.ID
	return &updatedRuleGroup, nil
}

func (c *Client) DeleteRuleGroup(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/api/v1/rule-groups/"+id, nil)
	return err
}
