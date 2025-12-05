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

type User struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Name  string   `json:"name"`
	Role  string   `json:"role"`
	Teams []string `json:"teams,omitempty"`
}

type Team struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type RoleAssignment struct {
	UserID string `json:"userId"`
	RoleID string `json:"roleId"`
}

// InviteRequest describes the user invitation payload.
type InviteRequest struct {
	Email string
	Roles []string
	Teams []string
}

func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
	var resp struct {
		Users []User `json:"users"`
	}
	if err := s.doer.Do(ctx, "GET", "/v1/users", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Users, nil
}

func (s *Service) ListTeams(ctx context.Context) ([]Team, error) {
	var resp struct {
		Teams []Team `json:"teams"`
	}
	if err := s.doer.Do(ctx, "GET", "/v1/teams", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Teams, nil
}

func (s *Service) AssignRole(ctx context.Context, userID, roleID string) error {
	path := fmt.Sprintf("/v1/users/%s/roles", url.PathEscape(userID))
	payload := RoleAssignment{UserID: userID, RoleID: roleID}
	return s.doer.Do(ctx, "POST", path, payload, nil)
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
