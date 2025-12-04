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
		*v = ListResponse{Data: []Entity{{Identifier: "a"}}}
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
	resp, err := svc.List(context.Background(), "bp", &ListOptions{Query: "foo", Page: 2, PerPage: 5})
	if err != nil || len(resp.Data) != 1 {
		t.Fatalf("list err %v resp %+v", err, resp)
	}
	if stub.path != "/v1/blueprints/bp/entities?page=2&per_page=5&query=foo" {
		t.Fatalf("unexpected query path: %s", stub.path)
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
