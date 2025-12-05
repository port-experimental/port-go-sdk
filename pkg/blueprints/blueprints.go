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

// RenameProperty changes the identifier of a property within a blueprint schema.
func (s *Service) RenameProperty(ctx context.Context, blueprintID, propertyID, newName string) error {
	if blueprintID == "" || propertyID == "" || newName == "" {
		return fmt.Errorf("blueprint, property and new names are required")
	}
	body := map[string]string{"newPropertyName": newName}
	path := fmt.Sprintf(
		"/v1/blueprints/%s/properties/%s/rename",
		url.PathEscape(blueprintID),
		url.PathEscape(propertyID),
	)
	return s.doer.Do(ctx, "PATCH", path, body, nil)
}

// RenameMirrorProperty changes the identifier of a mirror property.
func (s *Service) RenameMirrorProperty(ctx context.Context, blueprintID, mirrorID, newName string) error {
	if blueprintID == "" || mirrorID == "" || newName == "" {
		return fmt.Errorf("blueprint, mirror and new names are required")
	}
	body := map[string]string{"newMirrorName": newName}
	path := fmt.Sprintf(
		"/v1/blueprints/%s/mirror/%s/rename",
		url.PathEscape(blueprintID),
		url.PathEscape(mirrorID),
	)
	return s.doer.Do(ctx, "PATCH", path, body, nil)
}

// RenameRelation changes the identifier of a relation on a blueprint.
func (s *Service) RenameRelation(ctx context.Context, blueprintID, relationID, newName string) error {
	if blueprintID == "" || relationID == "" || newName == "" {
		return fmt.Errorf("blueprint, relation and new names are required")
	}
	body := map[string]string{"newRelationIdentifier": newName}
	path := fmt.Sprintf(
		"/v1/blueprints/%s/relations/%s/rename",
		url.PathEscape(blueprintID),
		url.PathEscape(relationID),
	)
	return s.doer.Do(ctx, "PATCH", path, body, nil)
}

// Relation defines a blueprint relation.
type Relation struct {
	Title    string `json:"title"`
	Target   string `json:"target"`
	Many     bool   `json:"many"`
	Required bool   `json:"required,omitempty"`
}

// BlueprintPermissions represents RBAC rules applied to a blueprint.
type BlueprintPermissions struct {
	Entities *BlueprintEntityPermissions `json:"entities,omitempty"`
}

// BlueprintEntityPermissions controls access to blueprint entities.
type BlueprintEntityPermissions struct {
	Read             *BlueprintPermissionRule           `json:"read,omitempty"`
	Register         *BlueprintPermissionRule           `json:"register,omitempty"`
	Update           *BlueprintPermissionRule           `json:"update,omitempty"`
	Unregister       *BlueprintPermissionRule           `json:"unregister,omitempty"`
	UpdateProperties map[string]BlueprintPermissionRule `json:"updateProperties,omitempty"`
	UpdateRelations  map[string]BlueprintPermissionRule `json:"updateRelations,omitempty"`
}

// BlueprintPermissionRule describes who can perform a given action.
type BlueprintPermissionRule struct {
	Users       []string       `json:"users,omitempty"`
	Teams       []string       `json:"teams,omitempty"`
	Roles       []string       `json:"roles,omitempty"`
	OwnedByTeam bool           `json:"ownedByTeam,omitempty"`
	Policy      map[string]any `json:"policy,omitempty"`
}

// GetPermissions fetches the permissions configured for a blueprint.
func (s *Service) GetPermissions(ctx context.Context, blueprintID string) (BlueprintPermissions, error) {
	path := fmt.Sprintf("/v1/blueprints/%s/permissions", url.PathEscape(blueprintID))
	var resp struct {
		Permissions BlueprintPermissions `json:"permissions"`
	}
	if err := s.doer.Do(ctx, "GET", path, nil, &resp); err != nil {
		return BlueprintPermissions{}, err
	}
	return resp.Permissions, nil
}

// UpdatePermissions replaces the permissions configuration of a blueprint.
func (s *Service) UpdatePermissions(ctx context.Context, blueprintID string, perms BlueprintPermissions) error {
	path := fmt.Sprintf("/v1/blueprints/%s/permissions", url.PathEscape(blueprintID))
	return s.doer.Do(ctx, "PATCH", path, perms, nil)
}
