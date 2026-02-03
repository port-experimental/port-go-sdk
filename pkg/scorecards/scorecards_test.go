package scorecards

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

func TestServicePaths(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)

	stub.resp = append(stub.resp, map[string]any{"scorecards": []Scorecard{{Identifier: "sc1"}}})
	if _, err := svc.List(context.Background()); err != nil {
		t.Fatalf("list err: %v", err)
	}
	if stub.path != "/v1/scorecards" || stub.method != "GET" {
		t.Fatalf("unexpected list path: %s %s", stub.method, stub.path)
	}

	stub.resp = append(stub.resp, map[string]any{"scorecards": []Scorecard{{Identifier: "sc2"}}})
	if _, err := svc.ListBlueprint(context.Background(), "bp"); err != nil {
		t.Fatalf("list blueprint err: %v", err)
	}
	if stub.path != "/v1/blueprints/bp/scorecards" || stub.method != "GET" {
		t.Fatalf("unexpected blueprint list path: %s %s", stub.method, stub.path)
	}

	stub.resp = append(stub.resp, map[string]any{"scorecard": Scorecard{Identifier: "sc2"}})
	if _, err := svc.Get(context.Background(), "bp", "sc2"); err != nil {
		t.Fatalf("get err: %v", err)
	}
	if stub.path != "/v1/blueprints/bp/scorecards/sc2" || stub.method != "GET" {
		t.Fatalf("unexpected get path: %s %s", stub.method, stub.path)
	}

	def := ScorecardDefinition{
		Identifier: "sc-create",
		Title:      "Created",
		Rules: []ScorecardRule{
			{Identifier: "rule1", Title: "Has owner", Level: "Gold", Query: map[string]any{"conditions": []any{}}},
		},
	}
	if err := svc.Create(context.Background(), "bp", def); err != nil {
		t.Fatalf("create err: %v", err)
	}
	if stub.path != "/v1/blueprints/bp/scorecards" || stub.method != "POST" {
		t.Fatalf("unexpected create path: %s %s", stub.method, stub.path)
	}

	update := ScorecardDefinition{
		Title: "Updated",
		Rules: []ScorecardRule{
			{Identifier: "rule1", Title: "Has owner", Level: "Gold", Query: map[string]any{"conditions": []any{}}},
		},
	}
	if err := svc.Update(context.Background(), "bp", "sc-create", update); err != nil {
		t.Fatalf("update err: %v", err)
	}
	if stub.path != "/v1/blueprints/bp/scorecards/sc-create" || stub.method != "PUT" {
		t.Fatalf("unexpected update path: %s %s", stub.method, stub.path)
	}
	body, ok := stub.body.(ScorecardDefinition)
	if !ok || body.Identifier != "sc-create" {
		t.Fatalf("identifier should be injected into payload: %#v", stub.body)
	}

	if err := svc.Delete(context.Background(), "bp", "sc-create"); err != nil {
		t.Fatalf("delete err: %v", err)
	}
	if stub.path != "/v1/blueprints/bp/scorecards/sc-create" || stub.method != "DELETE" {
		t.Fatalf("unexpected delete path: %s %s", stub.method, stub.path)
	}
}

func TestUpdateIdentifierMismatch(t *testing.T) {
	svc := New(&stubDoer{})
	update := ScorecardDefinition{
		Identifier: "foo",
		Rules: []ScorecardRule{
			{Identifier: "rule1", Title: "rule", Level: "Gold", Query: map[string]any{"conditions": []any{}}},
		},
	}
	if err := svc.Update(context.Background(), "bp", "bar", update); err == nil {
		t.Fatalf("expected identifier mismatch error")
	}
}
