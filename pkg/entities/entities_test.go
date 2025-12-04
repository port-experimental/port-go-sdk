package entities

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
	switch v := out.(type) {
	case *ListResponse:
		*v = ListResponse{Entities: []Entity{{Identifier: "a"}}, OK: true}
	case *Entity:
		*v = Entity{Identifier: "demo"}
	}
	return nil
}

func TestUpsert(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	ent := Entity{Identifier: "demo", Properties: map[string]any{"name": "Demo"}}
	if err := svc.Upsert(context.Background(), "my_blueprint", ent); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if stub.method != "POST" || stub.path != "/v1/blueprints/my_blueprint/entities?upsert=true&merge=true" {
		t.Fatalf("unexpected path %s %s", stub.method, stub.path)
	}
	payload := stub.body.(map[string]any)
	if payload["identifier"] != "demo" {
		t.Fatalf("payload mismatch")
	}
}

func TestList(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	resp, err := svc.List(context.Background(), "bp", &ListOptions{
		Query: map[string]any{
			"combinator": "and",
			"rules": []map[string]any{
				{"property": "name", "operator": "=", "value": "demo"},
			},
		},
		Include: []string{"identifier"},
		Exclude: []string{"properties.large_payload"},
		Limit:   5,
	})
	if err != nil || len(resp.Entities) != 1 || !resp.OK {
		t.Fatalf("list err %v resp %+v", err, resp)
	}
	if stub.method != "POST" {
		t.Fatalf("expected POST, got %s", stub.method)
	}
	if stub.path != "/v1/blueprints/bp/entities/search" {
		t.Fatalf("unexpected path: %s", stub.path)
	}
	body, ok := stub.body.(map[string]any)
	if !ok {
		t.Fatalf("unexpected body type: %#v", stub.body)
	}
	if body["limit"] != 5 {
		t.Fatalf("limit missing: %#v", body)
	}
	if body["query"].(map[string]any)["combinator"] != "and" {
		t.Fatalf("query mismatch: %#v", body["query"])
	}
}

func TestUpdate(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	props := map[string]any{"foo": "bar"}
	if err := svc.Update(context.Background(), "bp", "id", props); err != nil {
		t.Fatalf("update err: %v", err)
	}
	want := map[string]any{"properties": props}
	if !reflect.DeepEqual(stub.body, want) {
		t.Fatalf("bad payload: %#v", stub.body)
	}
}
