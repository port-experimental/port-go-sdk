package main

import (
	"context"
	"fmt"
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
	cli, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	const blueprintID = "example_blueprint"
	if err := cli.Blueprints().Upsert(ctx, blueprints.Blueprint{
		Identifier: blueprintID,
		Title:      "Demo",
		Schema: map[string]any{
			"properties": map[string]any{
				"name":  map[string]any{"type": "string"},
				"owner": map[string]any{"type": "string"},
			},
		},
	}); err != nil {
		log.Fatal(err)
	}
	bp, err := cli.Blueprints().Get(ctx, blueprintID)
	if err != nil {
		log.Fatal(err)
	}
	if schema, ok := bp.Schema["properties"].(map[string]any); ok {
		fmt.Printf("blueprint %s properties:\n", bp.Identifier)
		for name, def := range schema {
			fmt.Printf("  - %s: %v\n", name, def)
		}
		fmt.Printf("total properties: %d\n", len(schema))
	} else {
		fmt.Printf("blueprint %s has no properties\n", bp.Identifier)
	}
}
