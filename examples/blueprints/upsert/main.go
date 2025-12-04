package main

import (
	"context"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/blueprints"
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
	const blueprintID = "example_blueprint"
	bp := blueprints.Blueprint{
		Identifier: blueprintID,
		Title:      "Demo Blueprint",
		Schema: map[string]any{
			"properties": map[string]any{
				"name":   map[string]any{"type": "string"},
				"x1name": map[string]any{"type": "string"},
			},
		},
	}
	if err := apiClient.Blueprints().Upsert(ctx, bp); err != nil {
		log.Fatal(err)
	}
	log.Printf("blueprint %s upserted\n", blueprintID)
}
