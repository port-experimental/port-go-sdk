package users

import (
	"context"
	"fmt"
	"net/url"
	"strings"
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

// User represents a Port user record.
type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"firstName,omitempty"`
	LastName  string     `json:"lastName,omitempty"`
	Type      string     `json:"type,omitempty"`
	Status    string     `json:"status,omitempty"`
	Providers []string   `json:"providers,omitempty"`
	Roles     []UserRole `json:"roles,omitempty"`
	Teams     []string   `json:"teams,omitempty"`
	CreatedAt string     `json:"createdAt,omitempty"`
	UpdatedAt string     `json:"updatedAt,omitempty"`
}

// UserRole reports the name of the role assigned to a user.
type UserRole struct {
	Name string `json:"name"`
}

// Team describes a Port team.
type Team struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Provider    string       `json:"provider,omitempty"`
	Users       []TeamMember `json:"users,omitempty"`
	CreatedAt   string       `json:"createdAt,omitempty"`
	UpdatedAt   string       `json:"updatedAt,omitempty"`
}

// TeamMember contains data for a team member returned by the API.
type TeamMember struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Picture   string `json:"picture,omitempty"`
	Status    string `json:"status,omitempty"`
}

// InviteRequest describes the user invitation payload.
type InviteRequest struct {
	Email string
	Roles []string
	Teams []string
}

// ListUsersOptions controls optional fields returned from the API.
type ListUsersOptions struct {
	Fields []string
}

// ListTeamsOptions controls optional fields returned from the API.
type ListTeamsOptions struct {
	Fields []string
}

// ListUsers returns all users with optional field filtering.
func (s *Service) ListUsers(ctx context.Context, opts *ListUsersOptions) ([]User, error) {
	path := "/v1/users"
	if qs := encodeFieldsQuery(opts); qs != "" {
		path += "?" + qs
	}
	var resp struct {
		Users []User `json:"users"`
	}
	if err := s.doer.Do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Users, nil
}

// ListTeams returns all teams with optional field filtering.
func (s *Service) ListTeams(ctx context.Context, opts *ListTeamsOptions) ([]Team, error) {
	path := "/v1/teams"
	if qs := encodeFieldsQuery(opts); qs != "" {
		path += "?" + qs
	}
	var resp struct {
		Teams []Team `json:"teams"`
	}
	if err := s.doer.Do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Teams, nil
}

// GetUser fetches a single user by email.
func (s *Service) GetUser(ctx context.Context, email string) (User, error) {
	path := fmt.Sprintf("/v1/users/%s", escapeEmailPath(email))
	var resp struct {
		User User `json:"user"`
	}
	if err := s.doer.Do(ctx, "GET", path, nil, &resp); err != nil {
		return User{}, err
	}
	return resp.User, nil
}

// UpdateUserRequest controls user role/team assignment.
type UpdateUserRequest struct {
	Roles []string `json:"roles,omitempty"`
	Teams []string `json:"teams,omitempty"`
}

// UpdateUser mutates the roles/teams assigned to a given user.
func (s *Service) UpdateUser(ctx context.Context, email string, req UpdateUserRequest) (User, error) {
	path := fmt.Sprintf("/v1/users/%s", escapeEmailPath(email))
	var resp struct {
		User User `json:"user"`
	}
	if err := s.doer.Do(ctx, "PATCH", path, req, &resp); err != nil {
		return User{}, err
	}
	return resp.User, nil
}

// AssignRole is a helper that assigns a single role to a user.
func (s *Service) AssignRole(ctx context.Context, userEmail, role string) error {
	role = strings.TrimSpace(role)
	if role == "" {
		return fmt.Errorf("role is required")
	}
	_, err := s.UpdateUser(ctx, userEmail, UpdateUserRequest{Roles: []string{role}})
	return err
}

// DeleteUser removes a user permanently.
func (s *Service) DeleteUser(ctx context.Context, email string) error {
	path := fmt.Sprintf("/v1/users/%s", escapeEmailPath(email))
	return s.doer.Do(ctx, "DELETE", path, nil, nil)
}

// TeamCreateRequest describes the payload for creating teams.
type TeamCreateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Users       []string `json:"users,omitempty"`
}

// TeamPatchRequest updates mutable parts of a team without replacing it.
type TeamPatchRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Users       []string `json:"users,omitempty"`
}

// TeamReplaceRequest fully replaces a team record via PUT.
type TeamReplaceRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Users       []string `json:"users,omitempty"`
}

func (s *Service) CreateTeam(ctx context.Context, req TeamCreateRequest) (Team, error) {
	var resp struct {
		Team Team `json:"team"`
	}
	if err := s.doer.Do(ctx, "POST", "/v1/teams", req, &resp); err != nil {
		return Team{}, err
	}
	return resp.Team, nil
}

func (s *Service) GetTeam(ctx context.Context, name string) (Team, error) {
	path := fmt.Sprintf("/v1/teams/%s", url.PathEscape(name))
	var resp struct {
		Team Team `json:"team"`
	}
	if err := s.doer.Do(ctx, "GET", path, nil, &resp); err != nil {
		return Team{}, err
	}
	return resp.Team, nil
}

func (s *Service) PatchTeam(ctx context.Context, name string, req TeamPatchRequest) (Team, error) {
	path := fmt.Sprintf("/v1/teams/%s", url.PathEscape(name))
	var resp struct {
		Team Team `json:"team"`
	}
	if err := s.doer.Do(ctx, "PATCH", path, req, &resp); err != nil {
		return Team{}, err
	}
	return resp.Team, nil
}

func (s *Service) ReplaceTeam(ctx context.Context, name string, req TeamReplaceRequest) (Team, error) {
	path := fmt.Sprintf("/v1/teams/%s", url.PathEscape(name))
	var resp struct {
		Team Team `json:"team"`
	}
	if err := s.doer.Do(ctx, "PUT", path, req, &resp); err != nil {
		return Team{}, err
	}
	return resp.Team, nil
}

func (s *Service) DeleteTeam(ctx context.Context, name string) error {
	path := fmt.Sprintf("/v1/teams/%s", url.PathEscape(name))
	return s.doer.Do(ctx, "DELETE", path, nil, nil)
}

// RotateCredentials rotates a user's API credentials.
func (s *Service) RotateCredentials(ctx context.Context, email string) error {
	path := fmt.Sprintf("/v1/rotate-credentials/%s", escapeEmailPath(email))
	return s.doer.Do(ctx, "POST", path, nil, nil)
}

// Invite sends an invitation to a user email with optional roles/teams.
func (s *Service) Invite(ctx context.Context, req InviteRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("invite email required")
	}
	invitee := map[string]any{
		"email": req.Email,
	}
	payload := map[string]any{
		"invitee": invitee,
	}
	if len(req.Roles) > 0 {
		invitee["roles"] = req.Roles
	}
	if len(req.Teams) > 0 {
		invitee["teams"] = req.Teams
	}
	return s.doer.Do(ctx, "POST", "/v1/users/invite", payload, nil)
}

func encodeFieldsQuery(opts interface{}) string {
	var fields []string
	switch v := opts.(type) {
	case *ListUsersOptions:
		if v != nil {
			fields = v.Fields
		}
	case *ListTeamsOptions:
		if v != nil {
			fields = v.Fields
		}
	default:
		return ""
	}
	if len(fields) == 0 {
		return ""
	}
	vals := url.Values{}
	for _, f := range fields {
		if strings.TrimSpace(f) == "" {
			continue
		}
		vals.Add("fields", f)
	}
	return vals.Encode()
}

func escapeEmailPath(email string) string {
	escaped := url.PathEscape(email)
	if strings.Contains(escaped, "@") {
		escaped = strings.ReplaceAll(escaped, "@", "%40")
	}
	return escaped
}
