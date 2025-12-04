package entities

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/port-experimental/port-go-sdk/pkg/porter"
)

// Service handles entity endpoints.
type Service struct {
	doer Doer
}

// New creates an entity service.
func New(doer Doer) *Service {
	return &Service{doer: doer}
}

// Doer matches client.Client for dependency injection.
type Doer interface {
	Do(ctx context.Context, method, path string, body any, out any) error
}

// Entity represents a Port entity.
type Entity struct {
	Identifier string                 `json:"identifier"`
	Blueprint  string                 `json:"blueprint"`
	Title      string                 `json:"title,omitempty"`
	Properties map[string]any         `json:"properties,omitempty"`
	Relations  map[string][]string    `json:"relations,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ListOptions control pagination/filtering.
type ListOptions struct {
	Include []string
	Exclude []string
	Query   map[string]any
	From    string
	Limit   int
}

// ListResponse wraps entity lists returned from the search endpoint.
type ListResponse struct {
	Entities []Entity `json:"entities"`
	Next     string   `json:"next,omitempty"`
	OK       bool     `json:"ok"`
}

// Create creates a new entity.
func (s *Service) Create(ctx context.Context, blueprint string, ent Entity) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities", url.PathEscape(blueprint))
	payload := map[string]any{
		"identifier": ent.Identifier,
		"properties": ent.Properties,
	}
	if len(ent.Relations) > 0 {
		rel := make(map[string]any, len(ent.Relations))
		for k, v := range ent.Relations {
			rel[k] = v
		}
		payload["relations"] = rel
	}
	return s.doer.Do(ctx, "POST", path, payload, nil)
}

// Upsert creates or updates an entity.
func (s *Service) Upsert(ctx context.Context, blueprint string, ent Entity) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities?upsert=true&merge=true", url.PathEscape(blueprint))
	payload := map[string]any{
		"identifier": ent.Identifier,
		"properties": ent.Properties,
	}
	if len(ent.Relations) > 0 {
		rel := make(map[string]any, len(ent.Relations))
		for k, v := range ent.Relations {
			rel[k] = v
		}
		payload["relations"] = rel
	}
	return s.doer.Do(ctx, "POST", path, payload, nil)
}

// Get fetches an entity by identifier.
func (s *Service) Get(ctx context.Context, blueprint, identifier string) (Entity, error) {
	var out Entity
	path := fmt.Sprintf("/v1/blueprints/%s/entities/%s", url.PathEscape(blueprint), url.PathEscape(identifier))
	err := s.doer.Do(ctx, "GET", path, nil, &out)
	return out, err
}

// Delete removes an entity.
func (s *Service) Delete(ctx context.Context, blueprint, identifier string) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities/%s", url.PathEscape(blueprint), url.PathEscape(identifier))
	return s.doer.Do(ctx, "DELETE", path, nil, nil)
}

// List returns entities for a blueprint with optional filters.
func (s *Service) List(ctx context.Context, blueprint string, opts *ListOptions) (ListResponse, error) {
	var out ListResponse
	var body map[string]any
	if opts != nil {
		body = map[string]any{}
		if len(opts.Include) > 0 {
			body["include"] = opts.Include
		}
		if len(opts.Exclude) > 0 {
			body["exclude"] = opts.Exclude
		}
		if opts.Query != nil {
			body["query"] = opts.Query
		}
		if opts.From != "" {
			body["from"] = opts.From
		}
		if opts.Limit > 0 {
			body["limit"] = opts.Limit
		}
		if len(body) == 0 {
			body = nil
		}
	}
	path := fmt.Sprintf("/v1/blueprints/%s/entities/search", url.PathEscape(blueprint))
	err := s.doer.Do(ctx, "POST", path, body, &out)
	return out, err
}

// Update applies a partial update to entity properties (merge=true).
func (s *Service) Update(ctx context.Context, blueprint, identifier string, properties map[string]any) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities/%s", url.PathEscape(blueprint), url.PathEscape(identifier))
	payload := map[string]any{
		"properties": properties,
	}
	if err := s.doer.Do(ctx, "PATCH", path, payload, nil); err != nil {
		var perr *porter.Error
		if errors.As(err, &perr) && perr.StatusCode == 422 {
			return s.doer.Do(ctx, "PUT", path, payload, nil)
		}
		return err
	}
	return nil
}

// LinkRelation links targets to a relation.
func (s *Service) LinkRelation(ctx context.Context, blueprint, identifier, relation string, targets []string) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities/%s/relations/%s", url.PathEscape(blueprint), url.PathEscape(identifier), url.PathEscape(relation))
	payload := map[string]any{"identifiers": targets}
	return s.doer.Do(ctx, "POST", path, payload, nil)
}

// UnlinkRelation removes relation links.
func (s *Service) UnlinkRelation(ctx context.Context, blueprint, identifier, relation string, targets []string) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities/%s/relations/%s", url.PathEscape(blueprint), url.PathEscape(identifier), url.PathEscape(relation))
	payload := map[string]any{"identifiers": targets}
	return s.doer.Do(ctx, "DELETE", path, payload, nil)
}
