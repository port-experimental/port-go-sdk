package automations

import (
	"context"
	"encoding/json"
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
	StartedAt string         `json:"startedAt,omitempty"`
	Context   map[string]any `json:"context,omitempty"`
	Result    map[string]any `json:"result,omitempty"`
}

type ExecutionRequest struct {
	Context    map[string]any `json:"context,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
	Entity     string         `json:"entity,omitempty"`
	RunAs      string         `json:"run_as,omitempty"`
}

// List returns all automations.
func (s *Service) List(ctx context.Context) ([]Automation, error) {
	params := url.Values{}
	params.Set("trigger_type", "automation")
	params.Set("version", "v2")
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", "/v1/actions?"+params.Encode(), nil, &raw); err != nil {
		return nil, err
	}
	return decodeAutomationList(raw)
}

// Get fetches an automation by identifier.
func (s *Service) Get(ctx context.Context, identifier string) (Automation, error) {
	params := url.Values{}
	params.Set("version", "v2")
	path := fmt.Sprintf("/v1/actions/%s?%s", url.PathEscape(identifier), params.Encode())
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", path, nil, &raw); err != nil {
		return Automation{}, err
	}
	return decodeAutomation(raw)
}

// ListExecutions returns executions for an automation.
func (s *Service) ListExecutions(ctx context.Context, identifier string) ([]Execution, error) {
	params := url.Values{}
	params.Set("action", identifier)
	params.Set("version", "v2")
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", "/v1/actions/runs?"+params.Encode(), nil, &raw); err != nil {
		return nil, err
	}
	return decodeExecutions(raw)
}

// Trigger runs an automation immediately.
func (s *Service) Trigger(ctx context.Context, identifier string, req ExecutionRequest) error {
	path := fmt.Sprintf("/v1/actions/%s/runs", url.PathEscape(identifier))
	vals := url.Values{}
	if req.RunAs != "" {
		vals.Set("run_as", req.RunAs)
	}
	if len(vals) > 0 {
		path = path + "?" + vals.Encode()
	}
	var payload map[string]any
	props := req.Properties
	if len(props) == 0 && len(req.Context) > 0 {
		props = req.Context
	}
	if len(props) > 0 || req.Entity != "" {
		payload = map[string]any{}
		if len(props) > 0 {
			payload["properties"] = props
		}
		if req.Entity != "" {
			payload["entity"] = req.Entity
		}
	}
	return s.doer.Do(ctx, "POST", path, payload, nil)
}

func decodeAutomationList(raw json.RawMessage) ([]Automation, error) {
	var wrap struct {
		Actions     *[]Automation `json:"actions"`
		Automations *[]Automation `json:"automations"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil {
		if wrap.Actions != nil {
			return *wrap.Actions, nil
		}
		if wrap.Automations != nil {
			return *wrap.Automations, nil
		}
	}
	var plain []Automation
	if err := json.Unmarshal(raw, &plain); err == nil {
		return plain, nil
	}
	return nil, fmt.Errorf("automations: unexpected list response")
}

func decodeAutomation(raw json.RawMessage) (Automation, error) {
	var wrap struct {
		Action     *Automation `json:"action"`
		Automation *Automation `json:"automation"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil {
		if wrap.Action != nil {
			return *wrap.Action, nil
		}
		if wrap.Automation != nil {
			return *wrap.Automation, nil
		}
	}
	var single Automation
	if err := json.Unmarshal(raw, &single); err == nil {
		return single, nil
	}
	return Automation{}, fmt.Errorf("automations: unexpected response")
}

func decodeExecutions(raw json.RawMessage) ([]Execution, error) {
	var wrap struct {
		Runs *[]Execution `json:"runs"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil {
		if wrap.Runs != nil {
			return *wrap.Runs, nil
		}
	}
	var plain []Execution
	if err := json.Unmarshal(raw, &plain); err == nil {
		return plain, nil
	}
	return nil, fmt.Errorf("automations: unexpected executions response")
}
