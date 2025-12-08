package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// Doer matches client.Client for dependency injection.
type Doer interface {
	Do(ctx context.Context, method, path string, body any, out any) error
}

// Service manages Port integrations (data sources) and webhooks.
type Service struct {
	doer Doer
}

// New creates a datasources service.
func New(doer Doer) *Service {
	return &Service{doer: doer}
}

// Integration represents a Port data source.
type Integration struct {
	Identifier               string            `json:"identifier"`
	Title                    string            `json:"title,omitempty"`
	InstallationAppType      string            `json:"installationAppType,omitempty"`
	ActionsProcessingEnabled bool              `json:"actionsProcessingEnabled,omitempty"`
	Spec                     map[string]any    `json:"spec,omitempty"`
	Config                   IntegrationConfig `json:"config,omitempty"`
	ChangelogDestination     map[string]any    `json:"changelogDestination,omitempty"`
	Version                  string            `json:"version,omitempty"`
}

// IntegrationConfig is the raw configuration payload returned by the API.
type IntegrationConfig map[string]any

// ListIntegrationsOptions controls filtering on ListIntegrations.
type ListIntegrationsOptions struct {
	ActionsProcessingEnabled *bool
}

// GetIntegrationOptions allows selecting the identifier field.
type GetIntegrationOptions struct {
	ByField string
}

// IntegrationUpdateRequest updates mutable integration fields.
type IntegrationUpdateRequest struct {
	Title                    string         `json:"title,omitempty"`
	InstallationAppType      string         `json:"installationAppType,omitempty"`
	ActionsProcessingEnabled *bool          `json:"actionsProcessingEnabled,omitempty"`
	Spec                     map[string]any `json:"spec,omitempty"`
	ChangelogDestination     map[string]any `json:"changelogDestination,omitempty"`
	Version                  string         `json:"version,omitempty"`
}

// IntegrationConfigRequest uploads a new mapping/config definition.
type IntegrationConfigRequest struct {
	Config IntegrationConfig `json:"config"`
}

// IntegrationLogs contains audit log entries for an integration.
type IntegrationLogs struct {
	Logs     []IntegrationLogEntry `json:"logs"`
	Next     string                `json:"next,omitempty"`
	Previous string                `json:"previous,omitempty"`
	Metadata map[string]any        `json:"metadata,omitempty"`
}

// IntegrationLogEntry represents a single log; payload is intentionally loose.
type IntegrationLogEntry map[string]any

// ListIntegrationLogsOptions controls pagination/filtering.
type ListIntegrationLogsOptions struct {
	Limit     int
	Timestamp string
	LogID     string
	Direction string
}

// Webhook represents an inbound webhook data source definition.
type Webhook struct {
	Identifier      string           `json:"identifier"`
	Title           string           `json:"title"`
	Description     string           `json:"description,omitempty"`
	Icon            string           `json:"icon,omitempty"`
	Enabled         bool             `json:"enabled"`
	IntegrationType string           `json:"integrationType,omitempty"`
	Mappings        []WebhookMapping `json:"mappings,omitempty"`
	Security        *WebhookSecurity `json:"security,omitempty"`
}

// WebhookRequest is used to create/update webhook definitions.
type WebhookRequest struct {
	Identifier      string           `json:"identifier,omitempty"`
	Title           string           `json:"title,omitempty"`
	Description     string           `json:"description,omitempty"`
	Icon            string           `json:"icon,omitempty"`
	Enabled         *bool            `json:"enabled,omitempty"`
	IntegrationType string           `json:"integrationType,omitempty"`
	Mappings        []WebhookMapping `json:"mappings,omitempty"`
	Security        *WebhookSecurity `json:"security,omitempty"`
}

// WebhookMapping mirrors the flexible structure accepted by the API.
type WebhookMapping map[string]any

// WebhookSecurity configures request validation.
type WebhookSecurity struct {
	RequestIdentifierPath string `json:"requestIdentifierPath,omitempty"`
	Secret                string `json:"secret,omitempty"`
	SignatureAlgorithm    string `json:"signatureAlgorithm,omitempty"`
	SignatureHeaderName   string `json:"signatureHeaderName,omitempty"`
	SignaturePrefix       string `json:"signaturePrefix,omitempty"`
}

// App represents a Port credential set (used for rotating webhook secrets).
type App struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Secret    string `json:"secret,omitempty"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// ListIntegrations returns all integrations with optional filtering.
func (s *Service) ListIntegrations(ctx context.Context, opts *ListIntegrationsOptions) ([]Integration, error) {
	path := "/v1/integration"
	if opts != nil && opts.ActionsProcessingEnabled != nil {
		q := url.Values{}
		q.Set("actionsProcessingEnabled", strconv.FormatBool(*opts.ActionsProcessingEnabled))
		path += "?" + q.Encode()
	}
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", path, nil, &raw); err != nil {
		return nil, err
	}
	return decodeIntegrationList(raw)
}

// GetIntegration fetches an integration by identifier.
func (s *Service) GetIntegration(ctx context.Context, identifier string, opts *GetIntegrationOptions) (Integration, error) {
	path := fmt.Sprintf("/v1/integration/%s", url.PathEscape(identifier))
	if opts != nil && opts.ByField != "" {
		q := url.Values{}
		q.Set("byField", opts.ByField)
		path += "?" + q.Encode()
	}
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", path, nil, &raw); err != nil {
		return Integration{}, err
	}
	return decodeIntegration(raw)
}

// UpdateIntegration applies partial updates to an integration.
func (s *Service) UpdateIntegration(ctx context.Context, identifier string, req IntegrationUpdateRequest) error {
	path := fmt.Sprintf("/v1/integration/%s", url.PathEscape(identifier))
	return s.doer.Do(ctx, "PATCH", path, req, nil)
}

// DeleteIntegration removes a data source.
func (s *Service) DeleteIntegration(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("/v1/integration/%s", url.PathEscape(identifier))
	return s.doer.Do(ctx, "DELETE", path, nil, nil)
}

// UpdateIntegrationConfig uploads a new mapping/configuration definition.
func (s *Service) UpdateIntegrationConfig(ctx context.Context, identifier string, cfg IntegrationConfigRequest) error {
	path := fmt.Sprintf("/v1/integration/%s/config", url.PathEscape(identifier))
	return s.doer.Do(ctx, "PATCH", path, cfg, nil)
}

// ListIntegrationLogs returns the audit logs for an integration.
func (s *Service) ListIntegrationLogs(ctx context.Context, identifier string, opts *ListIntegrationLogsOptions) (IntegrationLogs, error) {
	values := url.Values{}
	if opts != nil {
		if opts.Limit > 0 {
			values.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Timestamp != "" {
			values.Set("timestamp", opts.Timestamp)
		}
		if opts.LogID != "" {
			values.Set("log_id", opts.LogID)
		}
		if opts.Direction != "" {
			values.Set("direction", opts.Direction)
		}
	}
	path := fmt.Sprintf("/v1/integration/%s/logs", url.PathEscape(identifier))
	if qs := values.Encode(); qs != "" {
		path += "?" + qs
	}
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", path, nil, &raw); err != nil {
		return IntegrationLogs{}, err
	}
	return decodeIntegrationLogs(raw)
}

// ListWebhooks returns all webhook definitions.
func (s *Service) ListWebhooks(ctx context.Context) ([]Webhook, error) {
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", "/v1/webhooks", nil, &raw); err != nil {
		return nil, err
	}
	return decodeWebhookList(raw)
}

// GetWebhook fetches a webhook definition.
func (s *Service) GetWebhook(ctx context.Context, identifier string) (Webhook, error) {
	path := fmt.Sprintf("/v1/webhooks/%s", url.PathEscape(identifier))
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", path, nil, &raw); err != nil {
		return Webhook{}, err
	}
	return decodeWebhook(raw)
}

// CreateWebhook creates a webhook data source.
func (s *Service) CreateWebhook(ctx context.Context, hook WebhookRequest) (Webhook, error) {
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "POST", "/v1/webhooks", hook, &raw); err != nil {
		return Webhook{}, err
	}
	return decodeWebhook(raw)
}

// UpdateWebhook updates a webhook definition.
func (s *Service) UpdateWebhook(ctx context.Context, identifier string, hook WebhookRequest) (Webhook, error) {
	path := fmt.Sprintf("/v1/webhooks/%s", url.PathEscape(identifier))
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "PATCH", path, hook, &raw); err != nil {
		return Webhook{}, err
	}
	return decodeWebhook(raw)
}

// DeleteWebhook removes a webhook.
func (s *Service) DeleteWebhook(ctx context.Context, identifier string) error {
	path := fmt.Sprintf("/v1/webhooks/%s", url.PathEscape(identifier))
	return s.doer.Do(ctx, "DELETE", path, nil, nil)
}

// RotateAppSecret rotates the secret for a credential set.
func (s *Service) RotateAppSecret(ctx context.Context, id string) (App, error) {
	path := fmt.Sprintf("/v1/apps/%s/rotate-secret", url.PathEscape(id))
	var out struct {
		App App `json:"app"`
	}
	err := s.doer.Do(ctx, "POST", path, nil, &out)
	return out.App, err
}

func decodeIntegrationList(raw json.RawMessage) ([]Integration, error) {
	var wrap struct {
		Integrations *[]Integration `json:"integrations"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil {
		if wrap.Integrations != nil {
			return *wrap.Integrations, nil
		}
	}
	var plain []Integration
	if err := json.Unmarshal(raw, &plain); err == nil {
		return plain, nil
	}
	return nil, fmt.Errorf("datasources: unexpected integration list response")
}

func decodeIntegration(raw json.RawMessage) (Integration, error) {
	var wrap struct {
		Integration *Integration `json:"integration"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil {
		if wrap.Integration != nil {
			return *wrap.Integration, nil
		}
	}
	var single Integration
	if err := json.Unmarshal(raw, &single); err == nil {
		return single, nil
	}
	return Integration{}, fmt.Errorf("datasources: unexpected integration response")
}

func decodeIntegrationLogs(raw json.RawMessage) (IntegrationLogs, error) {
	var logs IntegrationLogs
	if err := json.Unmarshal(raw, &logs); err == nil {
		if logs.Logs != nil || logs.Next != "" || logs.Previous != "" || len(logs.Metadata) > 0 {
			return logs, nil
		}
	}
	var wrap struct {
		Logs []IntegrationLogEntry `json:"logs"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil && wrap.Logs != nil {
		return IntegrationLogs{Logs: wrap.Logs}, nil
	}
	var arr []IntegrationLogEntry
	if err := json.Unmarshal(raw, &arr); err == nil {
		return IntegrationLogs{Logs: arr}, nil
	}
	return IntegrationLogs{}, fmt.Errorf("datasources: unexpected integration logs response")
}

func decodeWebhookList(raw json.RawMessage) ([]Webhook, error) {
	type top struct {
		Webhooks []Webhook `json:"webhooks"`
		Items    []Webhook `json:"items"`
	}
	var wrap top
	if err := json.Unmarshal(raw, &wrap); err == nil {
		switch {
		case wrap.Webhooks != nil:
			return wrap.Webhooks, nil
		case wrap.Items != nil:
			return wrap.Items, nil
		}
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err == nil {
		if wraw, ok := envelope["webhooks"]; ok {
			var simple []Webhook
			if err := json.Unmarshal(wraw, &simple); err == nil {
				return simple, nil
			}
			var items struct {
				Items []Webhook `json:"items"`
			}
			if err := json.Unmarshal(wraw, &items); err == nil && items.Items != nil {
				return items.Items, nil
			}
		}
		if itemsRaw, ok := envelope["items"]; ok {
			var items []Webhook
			if err := json.Unmarshal(itemsRaw, &items); err == nil {
				return items, nil
			}
		}
	}
	var plain []Webhook
	if err := json.Unmarshal(raw, &plain); err == nil {
		return plain, nil
	}
	// Fallback to empty slice instead of propagating an error; safer for clients.
	return []Webhook{}, nil
}

func decodeWebhook(raw json.RawMessage) (Webhook, error) {
	var wrap struct {
		Webhook *Webhook `json:"webhook"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil {
		if wrap.Webhook != nil {
			return *wrap.Webhook, nil
		}
	}
	var single Webhook
	if err := json.Unmarshal(raw, &single); err == nil {
		return single, nil
	}
	return Webhook{}, fmt.Errorf("datasources: unexpected webhook response")
}
