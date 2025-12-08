package datasources

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
		switch dst := out.(type) {
		case *json.RawMessage:
			*dst = append((*dst)[:0], data...)
		default:
			_ = json.Unmarshal(data, dst)
		}
	}
	return nil
}

func TestIntegrationPaths(t *testing.T) {
	ctx := context.Background()
	stub := &stubDoer{}
	svc := New(stub)
	enabled := true
	stub.push([]Integration{})
	if _, err := svc.ListIntegrations(ctx, &ListIntegrationsOptions{ActionsProcessingEnabled: &enabled}); err != nil {
		t.Fatalf("list integrations: %v", err)
	}
	if stub.method != "GET" || stub.path != "/v1/integration?actionsProcessingEnabled=true" {
		t.Fatalf("unexpected list path %s %s", stub.method, stub.path)
	}

	stub.push(Integration{Identifier: "hook"})
	if _, err := svc.GetIntegration(ctx, "hook", &GetIntegrationOptions{ByField: "logIngestId"}); err != nil {
		t.Fatalf("get integration: %v", err)
	}
	if stub.path != "/v1/integration/hook?byField=logIngestId" {
		t.Fatalf("bad get path %s", stub.path)
	}

	update := IntegrationUpdateRequest{Title: "test"}
	if err := svc.UpdateIntegration(ctx, "hook", update); err != nil {
		t.Fatalf("update integration: %v", err)
	}
	if stub.method != "PATCH" || stub.path != "/v1/integration/hook" {
		t.Fatalf("bad update path %s %s", stub.method, stub.path)
	}
	if !reflect.DeepEqual(stub.body, update) {
		t.Fatalf("update body mismatch %#v", stub.body)
	}

	if err := svc.DeleteIntegration(ctx, "hook"); err != nil {
		t.Fatalf("delete integration: %v", err)
	}
	if stub.method != "DELETE" || stub.path != "/v1/integration/hook" {
		t.Fatalf("bad delete path %s %s", stub.method, stub.path)
	}

	cfg := IntegrationConfigRequest{Config: IntegrationConfig{"foo": "bar"}}
	if err := svc.UpdateIntegrationConfig(ctx, "hook", cfg); err != nil {
		t.Fatalf("update config: %v", err)
	}
	if stub.method != "PATCH" || stub.path != "/v1/integration/hook/config" {
		t.Fatalf("bad config path %s %s", stub.method, stub.path)
	}
	if !reflect.DeepEqual(stub.body, cfg) {
		t.Fatalf("config body mismatch %#v", stub.body)
	}

	stub.push(IntegrationLogs{Logs: []IntegrationLogEntry{}})
	opts := &ListIntegrationLogsOptions{Limit: 50, Direction: "up", LogID: "abc", Timestamp: "ts"}
	if _, err := svc.ListIntegrationLogs(ctx, "hook", opts); err != nil {
		t.Fatalf("list logs: %v", err)
	}
	wantPath := "/v1/integration/hook/logs?direction=up&limit=50&log_id=abc&timestamp=ts"
	if stub.path != wantPath {
		t.Fatalf("bad logs path %s", stub.path)
	}
}

func TestWebhookPaths(t *testing.T) {
	ctx := context.Background()
	stub := &stubDoer{}
	svc := New(stub)
	stub.push([]Webhook{})
	if _, err := svc.ListWebhooks(ctx); err != nil {
		t.Fatalf("list webhooks: %v", err)
	}
	if stub.method != "GET" || stub.path != "/v1/webhooks" {
		t.Fatalf("bad list path %s %s", stub.method, stub.path)
	}

	stub.push(Webhook{Identifier: "hook"})
	if _, err := svc.GetWebhook(ctx, "hook"); err != nil {
		t.Fatalf("get webhook: %v", err)
	}
	if stub.path != "/v1/webhooks/hook" {
		t.Fatalf("bad get path %s", stub.path)
	}

	stub.push(Webhook{Identifier: "hook"})
	enabled := true
	req := WebhookRequest{Identifier: "hook", Title: "Hook", Enabled: &enabled}
	if _, err := svc.CreateWebhook(ctx, req); err != nil {
		t.Fatalf("create webhook: %v", err)
	}
	if stub.method != "POST" || stub.path != "/v1/webhooks" {
		t.Fatalf("bad create path %s %s", stub.method, stub.path)
	}
	if !reflect.DeepEqual(stub.body, req) {
		t.Fatalf("create body mismatch %#v", stub.body)
	}

	stub.push(Webhook{Identifier: "hook"})
	update := WebhookRequest{Title: "New Title"}
	if _, err := svc.UpdateWebhook(ctx, "hook", update); err != nil {
		t.Fatalf("update webhook: %v", err)
	}
	if stub.method != "PATCH" || stub.path != "/v1/webhooks/hook" {
		t.Fatalf("bad update path %s %s", stub.method, stub.path)
	}
	if !reflect.DeepEqual(stub.body, update) {
		t.Fatalf("update body mismatch %#v", stub.body)
	}

	if err := svc.DeleteWebhook(ctx, "hook"); err != nil {
		t.Fatalf("delete webhook: %v", err)
	}
	if stub.method != "DELETE" || stub.path != "/v1/webhooks/hook" {
		t.Fatalf("bad delete path %s %s", stub.method, stub.path)
	}
}

func TestRotateSecretPath(t *testing.T) {
	ctx := context.Background()
	stub := &stubDoer{}
	svc := New(stub)
	stub.push(map[string]any{"app": App{ID: "app1"}})
	if _, err := svc.RotateAppSecret(ctx, "app1"); err != nil {
		t.Fatalf("rotate secret: %v", err)
	}
	if stub.method != "POST" || stub.path != "/v1/apps/app1/rotate-secret" {
		t.Fatalf("bad rotate path %s %s", stub.method, stub.path)
	}
}
