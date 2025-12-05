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

	query := map[string]any{
		"combinator": "and",
		"rules": []map[string]any{
			{
				"property": "$identifier",
				"operator": "contains",
				"value":    "demo",
			},
		},
	}
	resp, err := apiClient.Entities().SearchBlueprint(ctx, "example_blueprint", entities.SearchOptions{
		Query: query,
		Limit: 10,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found %d entities matching identifier filter\n", len(resp.Entities))
	for _, ent := range resp.Entities {
		id := ent.Identifier
		if id == "" {
			id = "<unknown>"
		}
		title := ent.Title
		if title == "" {
			title = "<no title>"
		}
		bp := ent.Blueprint
		if bp == "" {
			bp = "<unknown>"
		}
		fmt.Printf("- id=%s title=%s blueprint=%s\n", id, title, bp)
	}
}
