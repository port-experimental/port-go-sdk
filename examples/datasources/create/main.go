package main

import (
	"context"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/datasources"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	apiClient, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	enabled := true
	webhook := datasources.WebhookRequest{
		Identifier:      "webhook_example",
		Title:           "Example Webhook",
		Enabled:         &enabled,
		IntegrationType: "custom",
		Mappings: []datasources.WebhookMapping{
			{
				"blueprint": "example_blueprint",
				"entity": map[string]any{
					"identifier": "{{ event.id }}",
					"title":      "{{ event.title }}",
					"properties": map[string]any{
						"name": "{{ event.name }}",
					},
				},
			},
		},
		Security: &datasources.WebhookSecurity{
			Secret: "replace-me",
		},
	}
	if _, err := apiClient.DataSources().CreateWebhook(ctx, webhook); err != nil {
		log.Fatal(err)
	}
	log.Println("data source created")
}
