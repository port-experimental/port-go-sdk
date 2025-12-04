package automations

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

func TestAutomationPaths(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	if _, err := svc.Get(context.Background(), "auto1"); err != nil {
		t.Fatalf("get err: %v", err)
	}
	if stub.path != "/v1/automations/auto1" {
		t.Fatalf("bad get path %s", stub.path)
	}
	if err := svc.Trigger(context.Background(), "auto1", ExecutionRequest{Context: map[string]any{"foo": "bar"}}); err != nil {
		t.Fatalf("trigger err: %v", err)
	}
	if stub.path != "/v1/automations/auto1/trigger" {
		t.Fatalf("bad trigger path %s", stub.path)
	}
	if _, err := svc.ListExecutions(context.Background(), "auto1"); err != nil {
		t.Fatalf("list executions err: %v", err)
	}
	if stub.path != "/v1/automations/auto1/executions" {
		t.Fatalf("bad exec path %s", stub.path)
	}
}
