package users

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

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
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

func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
	var out []User
	err := s.doer.Do(ctx, "GET", "/v1/users", nil, &out)
	return out, err
}

func (s *Service) ListTeams(ctx context.Context) ([]Team, error) {
	var out []Team
	err := s.doer.Do(ctx, "GET", "/v1/teams", nil, &out)
	return out, err
}

func (s *Service) AssignRole(ctx context.Context, userID, roleID string) error {
	path := fmt.Sprintf("/v1/users/%s/roles", url.PathEscape(userID))
	payload := RoleAssignment{UserID: userID, RoleID: roleID}
	return s.doer.Do(ctx, "POST", path, payload, nil)
}
