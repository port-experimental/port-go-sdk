package organization

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

type stubDoer struct {
	method string
	path   string
	body   any
	resp   any
}

func (s *stubDoer) Do(ctx context.Context, method, path string, body any, out any) error {
	s.method = method
	s.path = path
	s.body = body
	if out == nil {
		return nil
	}
	payload := s.resp
	if payload == nil {
		payload = map[string]any{"ok": true}
	}
	b, _ := json.Marshal(payload)
	switch v := out.(type) {
	case *json.RawMessage:
		*v = append((*v)[:0], b...)
	default:
		_ = json.Unmarshal(b, v)
	}
	return nil
}

func TestGetDirect(t *testing.T) {
	stub := &stubDoer{
		resp: Organization{Name: "Demo Org"},
	}
	svc := New(stub)
	org, err := svc.Get(context.Background())
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
	if org.Name != "Demo Org" {
		t.Fatalf("unexpected org %+v", org)
	}
	if stub.method != "GET" || stub.path != "/v1/organization" {
		t.Fatalf("bad call %s %s", stub.method, stub.path)
	}
}

func TestGetWrapped(t *testing.T) {
	stub := &stubDoer{
		resp: map[string]any{
			"organization": map[string]any{"name": "Wrapped Org"},
		},
	}
	svc := New(stub)
	org, err := svc.Get(context.Background())
	if err != nil {
		t.Fatalf("get wrapped err: %v", err)
	}
	if org.Name != "Wrapped Org" {
		t.Fatalf("unexpected org %+v", org)
	}
}

func TestUpdate(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	req := UpdateRequest{Name: "New Name"}
	if err := svc.Update(context.Background(), req); err != nil {
		t.Fatalf("update err: %v", err)
	}
	if stub.method != "PUT" || stub.path != "/v1/organization" {
		t.Fatalf("bad path %s %s", stub.method, stub.path)
	}
	if !reflect.DeepEqual(stub.body, req) {
		t.Fatalf("bad body %#v", stub.body)
	}
}

func TestPatch(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	title := "Portal"
	req := PatchRequest{Name: &title}
	if err := svc.Patch(context.Background(), req); err != nil {
		t.Fatalf("patch err: %v", err)
	}
	if stub.method != "PATCH" {
		t.Fatalf("expected PATCH got %s", stub.method)
	}
}

func TestListSecrets(t *testing.T) {
	stub := &stubDoer{
		resp: SecretsResponse{
			OK: true,
			Secrets: []SecretMetadata{
				{SecretName: "demo"},
			},
		},
	}
	svc := New(stub)
	secrets, err := svc.ListSecrets(context.Background())
	if err != nil {
		t.Fatalf("list secrets err: %v", err)
	}
	if len(secrets.Secrets) != 1 {
		t.Fatalf("unexpected secrets %+v", secrets)
	}
}

func TestCreateSecret(t *testing.T) {
	stub := &stubDoer{
		resp: SecretResponse{
			OK: true,
			Secret: SecretMetadata{
				SecretName: "demo",
			},
		},
	}
	svc := New(stub)
	req := CreateSecretRequest{SecretName: "demo", SecretValue: "value"}
	resp, err := svc.CreateSecret(context.Background(), req)
	if err != nil {
		t.Fatalf("create secret err: %v", err)
	}
	if resp.Secret.SecretName != "demo" {
		t.Fatalf("unexpected resp %+v", resp)
	}
}

func TestGetSecret(t *testing.T) {
	stub := &stubDoer{
		resp: SecretResponse{
			OK: true,
			Secret: SecretMetadata{
				SecretName: "demo",
			},
		},
	}
	svc := New(stub)
	resp, err := svc.GetSecret(context.Background(), "demo")
	if err != nil {
		t.Fatalf("get secret err: %v", err)
	}
	if resp.Secret.SecretName != "demo" {
		t.Fatalf("unexpected resp %+v", resp)
	}
	if stub.path != "/v1/organization/secrets/demo" {
		t.Fatalf("unexpected path %s", stub.path)
	}
}

func TestUpdateSecret(t *testing.T) {
	stub := &stubDoer{
		resp: SecretResponse{
			OK:     true,
			Secret: SecretMetadata{SecretName: "demo"},
		},
	}
	svc := New(stub)
	_, err := svc.UpdateSecret(context.Background(), "demo", UpdateSecretRequest{Description: "updated"})
	if err != nil {
		t.Fatalf("update secret err: %v", err)
	}
	if stub.method != "PATCH" {
		t.Fatalf("expected PATCH, got %s", stub.method)
	}
}

func TestDeleteSecret(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	if err := svc.DeleteSecret(context.Background(), "demo"); err != nil {
		t.Fatalf("delete secret err: %v", err)
	}
	if stub.method != "DELETE" {
		t.Fatalf("expected DELETE got %s", stub.method)
	}
}
