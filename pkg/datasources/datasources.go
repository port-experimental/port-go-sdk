package datasources

import (
	"context"
	"fmt"
	"net/url"
)

type Doer interface {
	Do(ctx context.Context, method, path string, body any, out any) error
}

// Service manages Port data sources and webhooks.
type Service struct {
	doer Doer
}

func New(doer Doer) *Service {
	return &Service{doer: doer}
}

type DataSource struct {
	Identifier  string                 `json:"identifier"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

func (s *Service) List(ctx context.Context) ([]DataSource, error) {
	var out []DataSource
	err := s.doer.Do(ctx, "GET", "/v1/data_sources", nil, &out)
	return out, err
}

func (s *Service) Get(ctx context.Context, identifier string) (DataSource, error) {
	var out DataSource
	path := fmt.Sprintf("/v1/data_sources/%s", url.PathEscape(identifier))
	err := s.doer.Do(ctx, "GET", path, nil, &out)
	return out, err
}

func (s *Service) Create(ctx context.Context, ds DataSource) error {
	return s.doer.Do(ctx, "POST", "/v1/data_sources", ds, nil)
}

func (s *Service) Delete(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("/v1/data_sources/%s", url.PathEscape(identifier))
	return s.doer.Do(ctx, "DELETE", path, nil, nil)
}

// RotateSecret rotates webhook secret if supported.
func (s *Service) RotateSecret(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("/v1/data_sources/%s/rotate_secret", url.PathEscape(identifier))
	return s.doer.Do(ctx, "POST", path, nil, nil)
}

// SetMapping uploads a mapping definition.
func (s *Service) SetMapping(ctx context.Context, identifier string, mapping any) error {
	path := fmt.Sprintf("/v1/data_sources/%s/mapping", url.PathEscape(identifier))
	return s.doer.Do(ctx, "PUT", path, mapping, nil)
}
