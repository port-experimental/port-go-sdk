package main

import (
	"context"
	"fmt"
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

	req := entities.AggregateRequest{
		"func": "count",
		"query": map[string]any{
			"combinator": "and",
			"rules": []map[string]any{
				{
					"property": "$identifier",
					"operator": "contains",
					"value":    "example_feature",
				},
			},
		},
	}
	resp, err := apiClient.Entities().Aggregate(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("aggregate matched %d entities across %d blueprints\n", len(resp.Entities), len(resp.MatchingBlueprints))
	for _, ent := range resp.Entities {
		fmt.Printf("- %s (%s)\n", ent.Identifier, ent.Blueprint)
	}
}
