package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/entities"
	"github.com/port-experimental/port-go-sdk/pkg/porter"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	blueprints := []struct {
		ID     string
		Entity entities.Entity
	}{
		{
			ID: "example_blueprint",
			Entity: entities.Entity{
				Identifier: "example_entity",
				Properties: map[string]any{
					"name":  "Demo Entity",
					"owner": "team@example.com",
				},
			},
		},
		{
			ID: "example_feature_blueprint",
			Entity: entities.Entity{
				Identifier: "example_feature",
				Properties: map[string]any{
					"name":        "AI Feature",
					"description": "Feature entity used for relation demos",
				},
			},
		},
	}

	for _, bp := range blueprints {
		if err := createEntity(ctx, c, bp.ID, bp.Entity); err != nil {
			log.Fatal(err)
		}
	}
}

func createEntity(ctx context.Context, cli *client.Client, blueprintID string, ent entities.Entity) error {
	err := cli.Entities().Create(ctx, blueprintID, ent)
	if err != nil {
		var perr *porter.Error
		if errors.As(err, &perr) && perr.StatusCode == 409 {
			log.Printf("entity %s already exists in %s\n", ent.Identifier, blueprintID)
			return nil
		}
		return err
	}
	log.Printf("entity %s created in blueprint %s\n", ent.Identifier, blueprintID)
	return nil
}
