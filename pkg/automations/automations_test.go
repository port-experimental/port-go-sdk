package automations

import (
	"context"
	"encoding/json"
	"testing"
)

type stubDoer struct {
	method string
	path   string
	body   any
	resp   []any
}

func (s *stubDoer) Do(ctx context.Context, method, path string, body any, out any) error {
	s.method = method
	s.path = path
	s.body = body
	if out != nil && len(s.resp) > 0 {
		payload := s.resp[0]
		s.resp = s.resp[1:]
		var data []byte
		switch v := payload.(type) {
		case json.RawMessage:
			data = v
		default:
			data, _ = json.Marshal(v)
		}
		switch dst := out.(type) {
		case *json.RawMessage:
			*dst = append((*dst)[:0], data...)
		default:
			_ = json.Unmarshal(data, dst)
		}
	}
	return nil
}

func TestAutomationPaths(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	stub.resp = append(stub.resp, map[string]any{"action": Automation{Identifier: "auto1"}})
	if _, err := svc.Get(context.Background(), "auto1"); err != nil {
		t.Fatalf("get err: %v", err)
	}
	if stub.path != "/v1/actions/auto1?version=v2" {
		t.Fatalf("bad get path %s", stub.path)
	}
	if err := svc.Trigger(context.Background(), "auto1", ExecutionRequest{Context: map[string]any{"foo": "bar"}, RunAs: "user@acme"}); err != nil {
		t.Fatalf("trigger err: %v", err)
	}
	if stub.path != "/v1/actions/auto1/runs?run_as=user%40acme" {
		t.Fatalf("bad trigger path %s", stub.path)
	}
	body, ok := stub.body.(map[string]any)
	if !ok || body["properties"].(map[string]any)["foo"] != "bar" {
		t.Fatalf("unexpected trigger body %#v", stub.body)
	}
	stub.resp = append(stub.resp, map[string]any{"runs": []Execution{{ID: "run1"}}})
	if _, err := svc.ListExecutions(context.Background(), "auto1"); err != nil {
		t.Fatalf("list executions err: %v", err)
	}
	if stub.path != "/v1/actions/runs?action=auto1&version=v2" {
		t.Fatalf("bad exec path %s", stub.path)
	}
	stub.resp = append(stub.resp, map[string]any{"actions": []Automation{{Identifier: "auto1"}}})
	if _, err := svc.List(context.Background()); err != nil {
		t.Fatalf("list err: %v", err)
	}
	if stub.path != "/v1/actions?trigger_type=automation&version=v2" {
		t.Fatalf("bad list path %s", stub.path)
	}

	stub.resp = append(stub.resp,
		map[string]any{"actions": []ActionDefinition{{Identifier: "auto1"}}},
		map[string]any{"action": ActionDefinition{Identifier: "auto1"}},
	)
	if _, err := svc.ListDefinitions(context.Background()); err != nil {
		t.Fatalf("list defs err: %v", err)
	}
	if stub.path != "/v1/actions?trigger_type=automation&version=v2" {
		t.Fatalf("bad list defs path %s", stub.path)
	}
	if _, err := svc.GetActionDefinition(context.Background(), "auto1"); err != nil {
		t.Fatalf("get def err: %v", err)
	}
	if stub.path != "/v1/actions/auto1?version=v2" {
		t.Fatalf("bad get def path %s", stub.path)
	}
	def := ActionDefinition{
		Identifier: "auto1",
		Trigger:    Trigger{Type: "automation", Event: &TriggerEvent{Type: "ENTITY_CREATED", BlueprintIdentifier: "svc"}},
		InvocationMethod: map[string]any{
			"type": "WEBHOOK",
			"url":  "https://example.com",
		},
	}
	if err := svc.CreateAction(context.Background(), def); err != nil {
		t.Fatalf("create action err: %v", err)
	}
	if stub.path != "/v1/actions" {
		t.Fatalf("bad create action path %s", stub.path)
	}
	if err := svc.UpdateAction(context.Background(), "auto1", def); err != nil {
		t.Fatalf("update action err: %v", err)
	}
	if stub.path != "/v1/actions/auto1" || stub.method != "PUT" {
		t.Fatalf("bad update action path %s %s", stub.method, stub.path)
	}
	if err := svc.DeleteAction(context.Background(), "auto1"); err != nil {
		t.Fatalf("delete action err: %v", err)
	}
	if stub.method != "DELETE" || stub.path != "/v1/actions/auto1" {
		t.Fatalf("bad delete action path %s %s", stub.method, stub.path)
	}
}
