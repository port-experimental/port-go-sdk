package blueprints

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

func TestBlueprintPaths(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	if err := svc.Upsert(context.Background(), Blueprint{Identifier: "foo"}); err != nil {
		t.Fatalf("upsert err: %v", err)
	}
	if stub.path != "/v1/blueprints/foo" {
		t.Fatalf("unexpected path: %s", stub.path)
	}
	stub = &stubDoer{}
	svc = New(stub)
	if _, err := svc.Get(context.Background(), "foo"); err != nil {
		t.Fatalf("get err: %v", err)
	}
	if stub.path != "/v1/blueprints/foo" {
		t.Fatalf("bad path: %s", stub.path)
	}
}
