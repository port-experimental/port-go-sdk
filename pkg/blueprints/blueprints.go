package blueprints

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/port-experimental/port-go-sdk/pkg/porter"
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
	Relations   map[string]Relation    `json:"relations,omitempty"`
}

// List returns all blueprints.
func (s *Service) List(ctx context.Context) ([]Blueprint, error) {
	var resp struct {
		Blueprints []Blueprint `json:"blueprints"`
	}
	if err := s.doer.Do(ctx, "GET", "/v1/blueprints", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Blueprints, nil
}

// Get fetches a blueprint by identifier.
func (s *Service) Get(ctx context.Context, identifier string) (Blueprint, error) {
	var resp struct {
		Blueprint Blueprint `json:"blueprint"`
	}
	path := fmt.Sprintf("/v1/blueprints/%s", url.PathEscape(identifier))
	if err := s.doer.Do(ctx, "GET", path, nil, &resp); err != nil {
		return Blueprint{}, err
	}
	return resp.Blueprint, nil
}

// Upsert creates or updates a blueprint using PUT, falling back to POST when the blueprint doesn't exist yet.
func (s *Service) Upsert(ctx context.Context, bp Blueprint) error {
	if bp.Identifier == "" {
		return fmt.Errorf("blueprint identifier required")
	}
	path := fmt.Sprintf("/v1/blueprints/%s", url.PathEscape(bp.Identifier))
	if err := s.doer.Do(ctx, "PUT", path, bp, nil); err != nil {
		var perr *porter.Error
		if errors.As(err, &perr) && perr.StatusCode == 404 {
			return s.doer.Do(ctx, "POST", "/v1/blueprints", bp, nil)
		}
		return err
	}
	return nil
}

// Delete removes a blueprint.
func (s *Service) Delete(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("/v1/blueprints/%s", url.PathEscape(identifier))
	return s.doer.Do(ctx, "DELETE", path, nil, nil)
}

// Relation defines a blueprint relation.
type Relation struct {
	Title    string `json:"title"`
	Target   string `json:"target"`
	Many     bool   `json:"many"`
	Required bool   `json:"required,omitempty"`
}
