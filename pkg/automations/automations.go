package automations

import (
	"context"
	"fmt"
	"net/url"
)

type Doer interface {
	Do(ctx context.Context, method, path string, body any, out any) error
}

type Service struct {
	doer Doer
}

func New(doer Doer) *Service {
	return &Service{doer: doer}
}

type Automation struct {
	Identifier  string `json:"identifier"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

type Execution struct {
	ID        string         `json:"id"`
	Status    string         `json:"status"`
	StartedAt string         `json:"startedAt"`
	Context   map[string]any `json:"context,omitempty"`
	Result    map[string]any `json:"result,omitempty"`
}

type ExecutionRequest struct {
	Context map[string]any `json:"context,omitempty"`
}

// List returns all automations.
func (s *Service) List(ctx context.Context) ([]Automation, error) {
	var out []Automation
	err := s.doer.Do(ctx, "GET", "/v1/automations", nil, &out)
	return out, err
}

// Get fetches an automation by identifier.
func (s *Service) Get(ctx context.Context, identifier string) (Automation, error) {
	var out Automation
	path := fmt.Sprintf("/v1/automations/%s", url.PathEscape(identifier))
	err := s.doer.Do(ctx, "GET", path, nil, &out)
	return out, err
}

// ListExecutions returns executions for an automation.
func (s *Service) ListExecutions(ctx context.Context, identifier string) ([]Execution, error) {
	var out []Execution
	path := fmt.Sprintf("/v1/automations/%s/executions", url.PathEscape(identifier))
	err := s.doer.Do(ctx, "GET", path, nil, &out)
	return out, err
}

// Trigger runs an automation immediately.
func (s *Service) Trigger(ctx context.Context, identifier string, req ExecutionRequest) error {
	path := fmt.Sprintf("/v1/automations/%s/trigger", url.PathEscape(identifier))
	return s.doer.Do(ctx, "POST", path, req, nil)
}
