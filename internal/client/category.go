// Copyright (c) HashiCorp, Inc.

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

type Category struct {
	ID        string `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Name      string `json:"name"`
	Notes     string `json:"notes"`
}

type CategorySingle struct {
	Data CategoryData `json:"data"`
}

type CategoryData struct {
	Type       string   `json:"type"`
	ID         string   `json:"id"`
	Attributes Category `json:"attributes"`
}

// unescapeHTML decodes HTML entities in all string fields
func (c *Category) unescapeHTML() {
	c.Name = html.UnescapeString(c.Name)
	c.Notes = html.UnescapeString(c.Notes)
}

func (c *Client) CreateCategory(ctx context.Context, category *Category) (*Category, error) {
	respBody, err := c.doRequest(ctx, http.MethodPost, "/api/v1/categories", category)
	if err != nil {
		return nil, err
	}

	var result CategorySingle
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	createdCategory := result.Data.Attributes
	createdCategory.ID = result.Data.ID
	createdCategory.unescapeHTML()
	return &createdCategory, nil
}

func (c *Client) GetCategory(ctx context.Context, id string) (*Category, error) {
	respBody, err := c.doRequest(ctx, http.MethodGet, "/api/v1/categories/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w, url: %s", err, "/api/v1/categories/"+id)
	}

	var result CategorySingle
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	category := result.Data.Attributes
	category.ID = result.Data.ID
	category.unescapeHTML()
	return &category, nil
}

func (c *Client) UpdateCategory(ctx context.Context, id string, category *Category) (*Category, error) {
	respBody, err := c.doRequest(ctx, http.MethodPut, "/api/v1/categories/"+id, category)
	if err != nil {
		return nil, err
	}

	var result CategorySingle
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	updatedCategory := result.Data.Attributes
	updatedCategory.ID = result.Data.ID
	updatedCategory.unescapeHTML()
	return &updatedCategory, nil
}

func (c *Client) DeleteCategory(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/api/v1/categories/"+id, nil)
	return err
}
