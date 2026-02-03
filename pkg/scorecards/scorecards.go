package scorecards

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Doer matches the client.Client contract.
type Doer interface {
	Do(ctx context.Context, method, path string, body any, out any) error
}

// Service exposes Port scorecard endpoints.
type Service struct {
	doer Doer
}

// New creates a scorecard service bound to the given HTTP client.
func New(doer Doer) *Service {
	return &Service{doer: doer}
}

// Scorecard represents a Port scorecard including metadata returned by the API.
type Scorecard struct {
	Identifier  string           `json:"identifier"`
	Title       string           `json:"title"`
	Description string           `json:"description,omitempty"`
	Filter      map[string]any   `json:"filter,omitempty"`
	Properties  map[string]any   `json:"properties,omitempty"`
	Levels      []ScorecardLevel `json:"levels,omitempty"`
	Rules       []ScorecardRule  `json:"rules"`
	Blueprint   string           `json:"blueprint,omitempty"`
	CreatedAt   string           `json:"createdAt,omitempty"`
	CreatedBy   string           `json:"createdBy,omitempty"`
	UpdatedAt   string           `json:"updatedAt,omitempty"`
	UpdatedBy   string           `json:"updatedBy,omitempty"`
}

// ScorecardDefinition is used for create/update operations.
type ScorecardDefinition struct {
	Identifier  string           `json:"identifier"`
	Title       string           `json:"title"`
	Description string           `json:"description,omitempty"`
	Filter      map[string]any   `json:"filter,omitempty"`
	Properties  map[string]any   `json:"properties,omitempty"`
	Levels      []ScorecardLevel `json:"levels,omitempty"`
	Rules       []ScorecardRule  `json:"rules"`
}

// ScorecardRule defines an individual condition block inside a scorecard.
type ScorecardRule struct {
	Identifier  string         `json:"identifier"`
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	Level       string         `json:"level"`
	Query       map[string]any `json:"query"`
}

// ScorecardLevel customizes the color/name that show up in Port's UI.
type ScorecardLevel struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

// List returns every scorecard in the organization.
func (s *Service) List(ctx context.Context) ([]Scorecard, error) {
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", "/v1/scorecards", nil, &raw); err != nil {
		return nil, err
	}
	return decodeScorecardList(raw)
}

// ListBlueprint returns all scorecards that belong to a specific blueprint.
func (s *Service) ListBlueprint(ctx context.Context, blueprint string) ([]Scorecard, error) {
	if blueprint == "" {
		return nil, fmt.Errorf("scorecards: blueprint identifier required")
	}
	path := fmt.Sprintf("/v1/blueprints/%s/scorecards", url.PathEscape(blueprint))
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", path, nil, &raw); err != nil {
		return nil, err
	}
	return decodeScorecardList(raw)
}

// Get fetches a single scorecard by identifier for the given blueprint.
func (s *Service) Get(ctx context.Context, blueprint, identifier string) (Scorecard, error) {
	if blueprint == "" {
		return Scorecard{}, fmt.Errorf("scorecards: blueprint identifier required")
	}
	if identifier == "" {
		return Scorecard{}, fmt.Errorf("scorecards: scorecard identifier required")
	}
	path := fmt.Sprintf(
		"/v1/blueprints/%s/scorecards/%s",
		url.PathEscape(blueprint),
		url.PathEscape(identifier),
	)
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", path, nil, &raw); err != nil {
		return Scorecard{}, err
	}
	return decodeScorecard(raw)
}

// Create provisions a new scorecard under the provided blueprint.
func (s *Service) Create(ctx context.Context, blueprint string, sc ScorecardDefinition) error {
	if blueprint == "" {
		return fmt.Errorf("scorecards: blueprint identifier required")
	}
	if sc.Identifier == "" {
		return fmt.Errorf("scorecards: definition identifier required")
	}
	path := fmt.Sprintf("/v1/blueprints/%s/scorecards", url.PathEscape(blueprint))
	return s.doer.Do(ctx, "POST", path, sc, nil)
}

// Update replaces the definition of an existing scorecard.
func (s *Service) Update(ctx context.Context, blueprint, identifier string, sc ScorecardDefinition) error {
	if blueprint == "" {
		return fmt.Errorf("scorecards: blueprint identifier required")
	}
	if identifier == "" {
		return fmt.Errorf("scorecards: scorecard identifier required")
	}
	payload := sc
	if payload.Identifier == "" {
		payload.Identifier = identifier
	} else if payload.Identifier != identifier {
		return fmt.Errorf("scorecards: definition identifier %q must match %q", payload.Identifier, identifier)
	}
	path := fmt.Sprintf(
		"/v1/blueprints/%s/scorecards/%s",
		url.PathEscape(blueprint),
		url.PathEscape(identifier),
	)
	return s.doer.Do(ctx, "PUT", path, payload, nil)
}

// Delete removes a scorecard from a blueprint.
func (s *Service) Delete(ctx context.Context, blueprint, identifier string) error {
	if blueprint == "" {
		return fmt.Errorf("scorecards: blueprint identifier required")
	}
	if identifier == "" {
		return fmt.Errorf("scorecards: scorecard identifier required")
	}
	path := fmt.Sprintf(
		"/v1/blueprints/%s/scorecards/%s",
		url.PathEscape(blueprint),
		url.PathEscape(identifier),
	)
	return s.doer.Do(ctx, "DELETE", path, nil, nil)
}

func decodeScorecardList(raw json.RawMessage) ([]Scorecard, error) {
	var wrap struct {
		Scorecards *[]Scorecard `json:"scorecards"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil {
		if wrap.Scorecards != nil {
			return *wrap.Scorecards, nil
		}
	}
	var plain []Scorecard
	if err := json.Unmarshal(raw, &plain); err == nil {
		return plain, nil
	}
	return nil, fmt.Errorf("scorecards: unexpected list response")
}

func decodeScorecard(raw json.RawMessage) (Scorecard, error) {
	var wrap struct {
		Scorecard *Scorecard `json:"scorecard"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil {
		if wrap.Scorecard != nil {
			return *wrap.Scorecard, nil
		}
	}
	var single Scorecard
	if err := json.Unmarshal(raw, &single); err == nil {
		return single, nil
	}
	return Scorecard{}, fmt.Errorf("scorecards: unexpected response")
}
