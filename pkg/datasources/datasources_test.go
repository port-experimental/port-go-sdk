package datasources

import (
	"context"
	"reflect"
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

func TestPaths(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	ds := DataSource{Identifier: "hook", Title: "Hook", Type: "webhook"}
	if err := svc.Create(context.Background(), ds); err != nil {
		t.Fatalf("create err: %v", err)
	}
	if stub.path != "/v1/data_sources" {
		t.Fatalf("bad create path %s", stub.path)
	}
	if err := svc.Delete(context.Background(), "hook"); err != nil {
		t.Fatalf("delete err: %v", err)
	}
	if stub.path != "/v1/data_sources/hook" {
		t.Fatalf("bad delete path %s", stub.path)
	}
	if err := svc.RotateSecret(context.Background(), "hook"); err != nil {
		t.Fatalf("rotate err: %v", err)
	}
	if stub.path != "/v1/data_sources/hook/rotate_secret" {
		t.Fatalf("rotate path %s", stub.path)
	}
	mapping := map[string]any{"a": "b"}
	if err := svc.SetMapping(context.Background(), "hook", mapping); err != nil {
		t.Fatalf("mapping err: %v", err)
	}
	if stub.path != "/v1/data_sources/hook/mapping" {
		t.Fatalf("mapping path %s", stub.path)
	}
	if !reflect.DeepEqual(stub.body, mapping) {
		t.Fatalf("mapping body mismatch")
	}
}
