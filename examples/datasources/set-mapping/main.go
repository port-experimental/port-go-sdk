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
	mappingConfig := datasources.IntegrationConfig{
		"resources": []any{
			map[string]any{
				"kind": "services",
				"selector": map[string]any{
					"query": "true",
				},
				"port": map[string]any{
					"entity": map[string]any{
						"blueprint":  "example_blueprint",
						"identifier": "{{ item.id }}",
						"title":      "{{ item.name }}",
						"properties": map[string]any{
							"name": "{{ item.name }}",
						},
					},
				},
			},
		},
	}
	req := datasources.IntegrationConfigRequest{Config: mappingConfig}
	if err := apiClient.DataSources().UpdateIntegrationConfig(ctx, "integration_id", req); err != nil {
		log.Fatal(err)
	}
	log.Println("mapping uploaded")
}
