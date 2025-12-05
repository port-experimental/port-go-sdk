package organization

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Doer matches client.Client for dependency injection.
type Doer interface {
	Do(ctx context.Context, method, path string, body any, out any) error
}

// Service exposes organization-related routes.
type Service struct {
	doer Doer
}

// New constructs an organization service.
func New(doer Doer) *Service {
	return &Service{doer: doer}
}

// Organization models /v1/organization details.
type Organization struct {
	ID           string                    `json:"id,omitempty"`
	Name         string                    `json:"name,omitempty"`
	Settings     *OrganizationSettings     `json:"settings,omitempty"`
	Announcement *OrganizationAnnouncement `json:"announcement,omitempty"`
	Metadata     map[string]any            `json:"metadata,omitempty"`
}

// OrganizationSettings defines organization UI settings.
type OrganizationSettings struct {
	HiddenBlueprints []string `json:"hiddenBlueprints,omitempty"`
	FederatedLogout  bool     `json:"federatedLogout,omitempty"`
	PortalIcon       string   `json:"portalIcon,omitempty"`
	PortalTitle      string   `json:"portalTitle,omitempty"`
}

// OrganizationAnnouncement configures the in-portal banner.
type OrganizationAnnouncement struct {
	Enabled bool    `json:"enabled,omitempty"`
	Content string  `json:"content,omitempty"`
	Link    *string `json:"link,omitempty"`
	Color   string  `json:"color,omitempty"`
}

// UpdateRequest fully replaces organization settings.
type UpdateRequest struct {
	Name         string                    `json:"name"`
	Settings     *OrganizationSettings     `json:"settings,omitempty"`
	Announcement *OrganizationAnnouncement `json:"announcement,omitempty"`
}

// PatchRequest partially updates organization settings.
type PatchRequest struct {
	Name         *string                   `json:"name,omitempty"`
	Settings     *OrganizationSettings     `json:"settings,omitempty"`
	Announcement *OrganizationAnnouncement `json:"announcement,omitempty"`
}

// SecretMetadata describes a stored secret (metadata only).
type SecretMetadata struct {
	SecretName  string `json:"secretName"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

// CreateSecretRequest creates a new organization secret.
type CreateSecretRequest struct {
	SecretName  string `json:"secretName"`
	SecretValue string `json:"secretValue"`
	Description string `json:"description,omitempty"`
}

// UpdateSecretRequest updates secret metadata/value.
type UpdateSecretRequest struct {
	SecretValue string `json:"secretValue,omitempty"`
	Description string `json:"description,omitempty"`
}

// SecretsResponse wraps a list of secrets.
type SecretsResponse struct {
	OK      bool             `json:"ok"`
	Secrets []SecretMetadata `json:"secrets"`
}

// SecretResponse wraps a single secret metadata entry.
type SecretResponse struct {
	OK     bool           `json:"ok"`
	Secret SecretMetadata `json:"secret"`
}

// Get fetches organization details.
func (s *Service) Get(ctx context.Context) (Organization, error) {
	var raw json.RawMessage
	if err := s.doer.Do(ctx, "GET", "/v1/organization", nil, &raw); err != nil {
		return Organization{}, err
	}
	org, err := decodeOrganization(raw)
	if err != nil {
		return Organization{}, err
	}
	return org, nil
}

func decodeOrganization(raw json.RawMessage) (Organization, error) {
	var wrapper struct {
		Organization json.RawMessage `json:"organization"`
	}
	if err := json.Unmarshal(raw, &wrapper); err == nil && len(wrapper.Organization) > 0 {
		var org Organization
		if err := json.Unmarshal(wrapper.Organization, &org); err != nil {
			return Organization{}, err
		}
		return org, nil
	}
	var direct Organization
	if err := json.Unmarshal(raw, &direct); err != nil {
		return Organization{}, err
	}
	return direct, nil
}

// Update replaces the organization name/settings/announcement.
func (s *Service) Update(ctx context.Context, req UpdateRequest) error {
	if req.Name == "" {
		return fmt.Errorf("organization: name is required")
	}
	var resp struct {
		OK bool `json:"ok"`
	}
	return s.doer.Do(ctx, "PUT", "/v1/organization", req, &resp)
}

// Patch partially updates organization details.
func (s *Service) Patch(ctx context.Context, req PatchRequest) error {
	var resp struct {
		OK bool `json:"ok"`
	}
	return s.doer.Do(ctx, "PATCH", "/v1/organization", req, &resp)
}

// ListSecrets fetches organization secret metadata.
func (s *Service) ListSecrets(ctx context.Context) (SecretsResponse, error) {
	var out SecretsResponse
	if err := s.doer.Do(ctx, "GET", "/v1/organization/secrets", nil, &out); err != nil {
		return SecretsResponse{}, err
	}
	return out, nil
}

// CreateSecret stores a new secret and returns its metadata.
func (s *Service) CreateSecret(ctx context.Context, req CreateSecretRequest) (SecretResponse, error) {
	var out SecretResponse
	if err := s.doer.Do(ctx, "POST", "/v1/organization/secrets", req, &out); err != nil {
		return SecretResponse{}, err
	}
	return out, nil
}

// GetSecret fetches metadata for a specific secret.
func (s *Service) GetSecret(ctx context.Context, name string) (SecretResponse, error) {
	var out SecretResponse
	path := fmt.Sprintf("/v1/organization/secrets/%s", url.PathEscape(name))
	if err := s.doer.Do(ctx, "GET", path, nil, &out); err != nil {
		return SecretResponse{}, err
	}
	return out, nil
}

// UpdateSecret updates an existing secret.
func (s *Service) UpdateSecret(ctx context.Context, name string, req UpdateSecretRequest) (SecretResponse, error) {
	var out SecretResponse
	path := fmt.Sprintf("/v1/organization/secrets/%s", url.PathEscape(name))
	if err := s.doer.Do(ctx, "PATCH", path, req, &out); err != nil {
		return SecretResponse{}, err
	}
	return out, nil
}

// DeleteSecret removes a secret.
func (s *Service) DeleteSecret(ctx context.Context, name string) error {
	path := fmt.Sprintf("/v1/organization/secrets/%s", url.PathEscape(name))
	var resp struct {
		OK bool `json:"ok"`
	}
	return s.doer.Do(ctx, "DELETE", path, nil, &resp)
}
