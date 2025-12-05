package entities

import (
	"context"
	"net/url"
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
	case *BulkEntitiesResponse:
		*v = BulkEntitiesResponse{OK: true}
	case *BulkDeleteResponse:
		*v = BulkDeleteResponse{OK: true, DeletedEntities: []string{"a"}}
	case *AggregateResponse:
		*v = AggregateResponse{OK: true}
	case *AggregateOverTimeResponse:
		*v = AggregateOverTimeResponse{OK: true}
	case *PropertiesHistoryResponse:
		*v = PropertiesHistoryResponse{OK: true}
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
	opts := &ListOptions{
		Include:                                 []string{"identifier", "properties.name"},
		Exclude:                                 []string{"properties.large_payload"},
		ExcludeCalculatedProperties:             true,
		AttachTitleToRelation:                   true,
		AttachIdentifierToTitleMirrorProperties: true,
		AllowPartialResults:                     true,
	}
	resp, err := svc.List(context.Background(), "bp", opts)
	if err != nil || len(resp.Entities) != 1 || !resp.OK {
		t.Fatalf("list err %v resp %+v", err, resp)
	}
	if stub.method != "GET" {
		t.Fatalf("expected GET, got %s", stub.method)
	}
	u, err := url.Parse(stub.path)
	if err != nil {
		t.Fatalf("parse path: %v", err)
	}
	values := u.Query()
	if got := values["include"]; len(got) != 2 || got[0] != "identifier" || got[1] != "properties.name" {
		t.Fatalf("include mismatch: %#v", got)
	}
	if got := values["exclude"]; len(got) != 1 || got[0] != "properties.large_payload" {
		t.Fatalf("exclude mismatch: %#v", got)
	}
	if values.Get("exclude_calculated_properties") != "true" {
		t.Fatalf("missing exclude_calculated_properties flag: %s", values.Encode())
	}
	if values.Get("attach_title_to_relation") != "true" || values.Get("attach_identifier_to_title_mirror_properties") != "true" {
		t.Fatalf("missing attach flags: %s", values.Encode())
	}
	if values.Get("allow_partial_results") != "true" {
		t.Fatalf("missing allow_partial_results: %s", values.Encode())
	}
	if stub.body != nil {
		t.Fatalf("expected nil body, got %#v", stub.body)
	}
}

func TestListRejectsSearchOptions(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	_, err := svc.List(context.Background(), "bp", &ListOptions{
		Query: map[string]any{"foo": "bar"},
	})
	if err == nil {
		t.Fatalf("expected error when using query in List")
	}
	if stub.method != "" {
		t.Fatalf("expected request not sent")
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

func TestBulkUpsert(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	ents := []Entity{
		{Identifier: "a", Properties: map[string]any{"name": "A"}},
		{Identifier: "b", Title: "B"},
	}
	resp, err := svc.BulkUpsert(context.Background(), "bp", ents)
	if err != nil || !resp.OK {
		t.Fatalf("bulk upsert err %v resp %+v", err, resp)
	}
	if stub.method != "POST" || stub.path != "/v1/blueprints/bp/entities/bulk" {
		t.Fatalf("unexpected path %s %s", stub.method, stub.path)
	}
	body, ok := stub.body.(map[string]any)
	if !ok {
		t.Fatalf("unexpected payload %#v", stub.body)
	}
	items, ok := body["entities"].([]map[string]any)
	if !ok || len(items) != 2 {
		t.Fatalf("entities payload mismatch %#v", body["entities"])
	}
	if items[0]["identifier"] != "a" || items[1]["identifier"] != "b" {
		t.Fatalf("bad identifiers %#v", items)
	}
}

func TestBulkDelete(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	opts := &BulkDeleteOptions{DeleteDependents: true, RunID: "abc"}
	resp, err := svc.BulkDelete(context.Background(), "bp", []string{"a", "b"}, opts)
	if err != nil || !resp.OK {
		t.Fatalf("bulk delete err %v resp %+v", err, resp)
	}
	if stub.method != "POST" {
		t.Fatalf("expected POST, got %s", stub.method)
	}
	if stub.path != "/v1/blueprints/bp/bulk/entities/delete?delete_dependents=true&run_id=abc" {
		t.Fatalf("unexpected path %s", stub.path)
	}
	body, ok := stub.body.(map[string]any)
	if !ok {
		t.Fatalf("unexpected payload %#v", stub.body)
	}
	ids := body["entities"].([]string)
	if len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Fatalf("bad ids %#v", ids)
	}
}

func TestSearchBlueprint(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	resp, err := svc.SearchBlueprint(context.Background(), "bp", SearchOptions{
		Query: map[string]any{
			"combinator": "and",
			"rules": []map[string]any{
				{"property": "name", "operator": "=", "value": "demo"},
			},
		},
		Include: []string{"identifier"},
		Limit:   5,
	})
	if err != nil || !resp.OK {
		t.Fatalf("search err %v resp %+v", err, resp)
	}
	if stub.method != "POST" || stub.path != "/v1/blueprints/bp/entities/search" {
		t.Fatalf("unexpected search path %s %s", stub.method, stub.path)
	}
	body, ok := stub.body.(map[string]any)
	if !ok || body["limit"] != 5 {
		t.Fatalf("missing limit %#v", stub.body)
	}
	query := body["query"].(map[string]any)
	if query["combinator"] != "and" {
		t.Fatalf("bad query %#v", query)
	}
}

func TestAggregate(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	req := AggregateRequest{
		"func":  "count",
		"query": map[string]any{"combinator": "and", "rules": []map[string]any{}},
	}
	resp, err := svc.Aggregate(context.Background(), req)
	if err != nil || !resp.OK {
		t.Fatalf("aggregate err %v resp %+v", err, resp)
	}
	if stub.method != "POST" || stub.path != "/v1/entities/aggregate" {
		t.Fatalf("bad aggregate call %s %s", stub.method, stub.path)
	}
	switch body := stub.body.(type) {
	case AggregateRequest:
		if body["func"] != "count" {
			t.Fatalf("aggregate func mismatch %#v", body)
		}
	case map[string]any:
		if body["func"] != "count" {
			t.Fatalf("aggregate func mismatch %#v", body)
		}
	default:
		t.Fatalf("unexpected body %#v", stub.body)
	}
}

func TestAggregateOverTime(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	req := AggregateOverTimeRequest{
		Blueprint:       "bp",
		TimeRange:       AggregateTimeRange{Preset: "lastWeek"},
		TimeInterval:    "day",
		Query:           map[string]any{"combinator": "and", "rules": []map[string]any{}},
		MeasureTimeBy:   "createdAt",
		AggregationType: "countEntities",
		Func:            "count",
	}
	resp, err := svc.AggregateOverTime(context.Background(), req)
	if err != nil || !resp.OK {
		t.Fatalf("aggregate over time err %v resp %+v", err, resp)
	}
	if stub.method != "POST" || stub.path != "/v1/entities/aggregate-over-time" {
		t.Fatalf("bad aggregate over time call %s %s", stub.method, stub.path)
	}
	if !reflect.DeepEqual(stub.body, req) {
		t.Fatalf("request mismatch %#v", stub.body)
	}
}

func TestPropertiesHistory(t *testing.T) {
	stub := &stubDoer{}
	svc := New(stub)
	req := PropertiesHistoryRequest{
		EntityIdentifier:    "ent",
		BlueprintIdentifier: "bp",
		PropertyNames:       []string{"p1", "p2"},
		TimeInterval:        "day",
		TimeRange:           &PropertiesHistoryTimeRange{Preset: "lastWeek"},
	}
	resp, err := svc.PropertiesHistory(context.Background(), req)
	if err != nil || !resp.OK {
		t.Fatalf("properties history err %v resp %+v", err, resp)
	}
	if stub.method != "POST" || stub.path != "/v1/entities/properties-history" {
		t.Fatalf("bad properties history call %s %s", stub.method, stub.path)
	}
	if !reflect.DeepEqual(stub.body, req) {
		t.Fatalf("request mismatch %#v", stub.body)
	}
}
