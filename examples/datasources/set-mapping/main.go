package main

import (
	"context"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
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
	mapping := map[string]any{
		"blueprint": "example_blueprint",
		"mapping": map[string]any{
			"identifier": "{{ event.id }}",
			"properties": map[string]any{
				"name": "{{ event.name }}",
			},
		},
	}
	if err := apiClient.DataSources().SetMapping(ctx, "webhook_example", mapping); err != nil {
		log.Fatal(err)
	}
	log.Println("mapping uploaded")
}
