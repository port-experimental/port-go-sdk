package blueprints

import (
	"context"
	"fmt"
	"net/url"
)

// Doer matches client.Client.
type Doer interface {
	Do(ctx context.Context, method, path string, body any, out any) error
}

// Service manages blueprints.
type Service struct {
	doer Doer
}

// New returns a blueprint service.
func New(doer Doer) *Service {
	return &Service{doer: doer}
}

// Blueprint represents the Port blueprint object (subset).
type Blueprint struct {
	Identifier  string                 `json:"identifier"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Schema      map[string]interface{} `json:"schema"`
	Icon        string                 `json:"icon,omitempty"`
}

// List returns all blueprints.
func (s *Service) List(ctx context.Context) ([]Blueprint, error) {
	var out []Blueprint
	err := s.doer.Do(ctx, "GET", "/v1/blueprints", nil, &out)
	return out, err
}

// Get fetches a blueprint by identifier.
func (s *Service) Get(ctx context.Context, identifier string) (Blueprint, error) {
	var out Blueprint
	path := fmt.Sprintf("/v1/blueprints/%s", url.PathEscape(identifier))
	err := s.doer.Do(ctx, "GET", path, nil, &out)
	return out, err
}

// Upsert creates or updates a blueprint.
func (s *Service) Upsert(ctx context.Context, bp Blueprint) error {
	path := "/v1/blueprints"
	return s.doer.Do(ctx, "POST", path, bp, nil)
}

// Delete removes a blueprint.
func (s *Service) Delete(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("/v1/blueprints/%s", url.PathEscape(identifier))
	return s.doer.Do(ctx, "DELETE", path, nil, nil)
}
