package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/blueprints"
	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/porter"
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const (
		ownerBlueprint     = "example_blueprint"
		ownerDisplayTitle  = "Example Blueprint"
		dependentBlueprint = "example_feature_blueprint"
		dependentTitle     = "Example Feature Blueprint"
	)
	ensureBlueprint(ctx, cli, dependentBlueprint, dependentTitle, map[string]any{
		"name":        map[string]any{"type": "string"},
		"description": map[string]any{"type": "string"},
	}, nil)

	ensureBlueprint(ctx, cli, ownerBlueprint, ownerDisplayTitle, map[string]any{
		"name":  map[string]any{"type": "string"},
		"owner": map[string]any{"type": "string"},
	}, map[string]blueprints.Relation{
		"features": {
			Title:  "Features",
			Target: dependentBlueprint,
			Many:   true,
		},
	})

	fmt.Println("Blueprint scaffolding complete.")
}

func ensureBlueprint(ctx context.Context, cli *client.Client, id, title string, properties map[string]any, relations map[string]blueprints.Relation) {
	existed := false
	if _, err := cli.Blueprints().Get(ctx, id); err == nil {
		existed = true
	} else {
		var perr *porter.Error
		if !errors.As(err, &perr) || perr.StatusCode != 404 {
			log.Fatalf("failed to check blueprint: %v", err)
		}
	}
	bp := blueprints.Blueprint{
		Identifier: id,
		Title:      title,
		Schema:     map[string]any{"properties": properties},
		Relations:  relations,
		Icon:       "Cube",
	}
	if err := cli.Blueprints().Upsert(ctx, bp); err != nil {
		log.Fatalf("failed to upsert blueprint %s: %v", id, err)
	}
	if existed {
		log.Printf("blueprint %s updated", id)
	} else {
		log.Printf("blueprint %s created", id)
	}
}
