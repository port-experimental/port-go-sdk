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
	if err := cli.Blueprints().Upsert(ctx, blueprints.Blueprint{
		Identifier: "demo_blueprint",
		Title:      "Demo",
		Schema:     map[string]any{"properties": map[string]any{"name": map[string]any{"type": "string"}}},
	}); err != nil {
		log.Fatal(err)
	}
	bp, err := cli.Blueprints().Get(ctx, "demo_blueprint")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("blueprint %s has %d properties\n", bp.Identifier, len(bp.Schema))
}
