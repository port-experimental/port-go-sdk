package users

import (
	"context"
	"testing"
)

type stubDoer struct {
	method string
	path   string
	body   any
}

func (s *stubDoer) Do(ctx context.Context, method, path string, body any, out any) error {
	s.method = method
	s.path = path
	s.body = body
	return nil
}

func TestUserPaths(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	if _, err := svc.ListUsers(context.Background()); err != nil {
		t.Fatalf("list users err: %v", err)
	}
	if stub.path != "/v1/users" {
		t.Fatalf("bad users path: %s", stub.path)
	}
	if _, err := svc.ListTeams(context.Background()); err != nil {
		t.Fatalf("list teams err: %v", err)
	}
	if stub.path != "/v1/teams" {
		t.Fatalf("bad teams path: %s", stub.path)
	}
	if err := svc.AssignRole(context.Background(), "user1", "role1"); err != nil {
		t.Fatalf("assign role err: %v", err)
	}
	if stub.path != "/v1/users/user1/roles" {
		t.Fatalf("bad assign path %s", stub.path)
	}
	if ra, ok := stub.body.(RoleAssignment); !ok || ra.RoleID != "role1" {
		t.Fatalf("bad payload %#v", stub.body)
	}
}
