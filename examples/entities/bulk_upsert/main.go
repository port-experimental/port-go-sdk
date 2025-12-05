package main

import (
	"context"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/entities"
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

	batches := map[string][]entities.Entity{
		"example_blueprint": {
			{
				Identifier: "demo_service_a",
				Title:      "Demo Service A",
				Properties: map[string]any{
					"name":        "Demo Service A",
					"environment": "production",
				},
			},
			{
				Identifier: "demo_service_b",
				Properties: map[string]any{
					"name":        "Demo Service B",
					"environment": "staging",
				},
			},
			{
				Identifier: "demo_service_c",
				Properties: map[string]any{
					"name":        "Demo Service C",
					"environment": "development",
				},
			},
		},
		"example_feature_blueprint": {
			{
				Identifier: "example_feature_rollout",
				Properties: map[string]any{
					"name":        "Rollout Guard",
					"description": "Feature flag entity controlling guarded rollouts",
				},
			},
			{
				Identifier: "example_feature_canary",
				Properties: map[string]any{
					"name":        "Canary Slots",
					"description": "Canary capacity allocation",
				},
			},
			{
				Identifier: "example_feature_payments",
				Properties: map[string]any{
					"name":        "Payments Feature",
					"description": "Handles card capture mutations",
				},
			},
		},
	}

	for blueprint, batch := range batches {
		resp, err := apiClient.Entities().BulkUpsert(ctx, blueprint, batch)
		if err != nil {
			log.Fatalf("bulk upsert %s: %v", blueprint, err)
		}
		for _, res := range resp.Entities {
			log.Printf("[%s] entity %s created=%t", blueprint, res.Identifier, res.Created)
		}
	}
}
