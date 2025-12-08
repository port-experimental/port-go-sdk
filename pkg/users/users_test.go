package users

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

type stubDoer struct {
	method    string
	path      string
	body      any
	responses []any
}

func (s *stubDoer) push(resp any) {
	s.responses = append(s.responses, resp)
}

func (s *stubDoer) Do(ctx context.Context, method, path string, body any, out any) error {
	s.method = method
	s.path = path
	s.body = body
	if out != nil && len(s.responses) > 0 {
		payload := s.responses[0]
		s.responses = s.responses[1:]
		data, _ := json.Marshal(payload)
		_ = json.Unmarshal(data, out)
	}
	return nil
}

func TestUserAndTeamPaths(t *testing.T) {
	ctx := context.Background()
	stub := &stubDoer{}
	svc := New(stub)

	stub.push(map[string]any{"users": []User{}})
	if _, err := svc.ListUsers(ctx, &ListUsersOptions{Fields: []string{"email", "roles.name"}}); err != nil {
		t.Fatalf("list users err: %v", err)
	}
	wantUsersPath := "/v1/users?fields=email&fields=roles.name"
	if stub.path != wantUsersPath {
		t.Fatalf("bad users path: %s", stub.path)
	}

	stub.push(map[string]any{"teams": []Team{}})
	if _, err := svc.ListTeams(ctx, &ListTeamsOptions{Fields: []string{"name"}}); err != nil {
		t.Fatalf("list teams err: %v", err)
	}
	if stub.path != "/v1/teams?fields=name" {
		t.Fatalf("bad teams path: %s", stub.path)
	}

	stub.push(map[string]any{"user": User{Email: "alice@example.com"}})
	if _, err := svc.GetUser(ctx, "alice@example.com"); err != nil {
		t.Fatalf("get user err: %v", err)
	}
	if stub.path != "/v1/users/alice%40example.com" {
		t.Fatalf("bad get user path: %s", stub.path)
	}

	updateReq := UpdateUserRequest{Roles: []string{"admin"}}
	stub.push(map[string]any{"user": User{Email: "alice@example.com", Roles: []UserRole{{Name: "admin"}}}})
	if _, err := svc.UpdateUser(ctx, "alice@example.com", updateReq); err != nil {
		t.Fatalf("update user err: %v", err)
	}
	if stub.method != "PATCH" || stub.path != "/v1/users/alice%40example.com" {
		t.Fatalf("bad update user path: %s %s", stub.method, stub.path)
	}
	if !reflect.DeepEqual(stub.body, updateReq) {
		t.Fatalf("unexpected update payload %#v", stub.body)
	}

	stub.push(map[string]any{"user": User{Email: "alice@example.com"}})
	if err := svc.AssignRole(ctx, "alice@example.com", "admin"); err != nil {
		t.Fatalf("assign role err: %v", err)
	}
	assignReq, _ := stub.body.(UpdateUserRequest)
	if len(assignReq.Roles) != 1 || assignReq.Roles[0] != "admin" {
		t.Fatalf("bad assign payload %#v", stub.body)
	}

	if err := svc.DeleteUser(ctx, "alice@example.com"); err != nil {
		t.Fatalf("delete user err: %v", err)
	}
	if stub.method != "DELETE" || stub.path != "/v1/users/alice%40example.com" {
		t.Fatalf("bad delete user path %s %s", stub.method, stub.path)
	}

	if err := svc.RotateCredentials(ctx, "alice@example.com"); err != nil {
		t.Fatalf("rotate creds err: %v", err)
	}
	if stub.method != "POST" || stub.path != "/v1/rotate-credentials/alice%40example.com" {
		t.Fatalf("bad rotate path %s %s", stub.method, stub.path)
	}

	createReq := TeamCreateRequest{Name: "Dev", Description: "Dev team"}
	stub.push(map[string]any{"team": Team{Name: "Dev"}})
	if _, err := svc.CreateTeam(ctx, createReq); err != nil {
		t.Fatalf("create team err: %v", err)
	}
	if stub.method != "POST" || stub.path != "/v1/teams" {
		t.Fatalf("bad create team path %s %s", stub.method, stub.path)
	}
	if !reflect.DeepEqual(stub.body, createReq) {
		t.Fatalf("unexpected team create payload %#v", stub.body)
	}

	stub.push(map[string]any{"team": Team{Name: "Dev"}})
	if _, err := svc.GetTeam(ctx, "Dev"); err != nil {
		t.Fatalf("get team err: %v", err)
	}
	if stub.path != "/v1/teams/Dev" {
		t.Fatalf("bad get team path %s", stub.path)
	}

	desc := "updated"
	patchReq := TeamPatchRequest{Description: &desc}
	stub.push(map[string]any{"team": Team{Name: "Dev", Description: "updated"}})
	if _, err := svc.PatchTeam(ctx, "Dev", patchReq); err != nil {
		t.Fatalf("patch team err: %v", err)
	}
	if stub.method != "PATCH" || stub.path != "/v1/teams/Dev" {
		t.Fatalf("bad patch team path %s %s", stub.method, stub.path)
	}
	if !reflect.DeepEqual(stub.body, patchReq) {
		t.Fatalf("unexpected patch payload %#v", stub.body)
	}

	replaceReq := TeamReplaceRequest{Name: "Backend"}
	stub.push(map[string]any{"team": Team{Name: "Backend"}})
	if _, err := svc.ReplaceTeam(ctx, "Dev", replaceReq); err != nil {
		t.Fatalf("replace team err: %v", err)
	}
	if stub.method != "PUT" || stub.path != "/v1/teams/Dev" {
		t.Fatalf("bad replace team path %s %s", stub.method, stub.path)
	}
	if !reflect.DeepEqual(stub.body, replaceReq) {
		t.Fatalf("unexpected replace payload %#v", stub.body)
	}

	if err := svc.DeleteTeam(ctx, "Backend"); err != nil {
		t.Fatalf("delete team err: %v", err)
	}
	if stub.method != "DELETE" || stub.path != "/v1/teams/Backend" {
		t.Fatalf("bad delete team path %s %s", stub.method, stub.path)
	}
}
