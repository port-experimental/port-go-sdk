package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/porter"
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// These identifiers line up with the blueprint created by examples/blueprints/create.
	const (
		blueprintID = "example_blueprint"
		oldProperty = "name"
		newProperty = "display_name"
	)
	if err := apiClient.Blueprints().RenameProperty(ctx, blueprintID, oldProperty, newProperty); err != nil {
		var perr *porter.Error
		if errors.As(err, &perr) && perr.StatusCode == 404 {
			if alreadyRenamed(ctx, apiClient, blueprintID, newProperty) {
				log.Printf("property already renamed to %s on blueprint %s", newProperty, blueprintID)
				return
			}
		}
		log.Fatal(err)
	}
	log.Printf("renamed property %s -> %s on blueprint %s", oldProperty, newProperty, blueprintID)
}

func alreadyRenamed(ctx context.Context, apiClient *client.Client, blueprintID, newProperty string) bool {
	bp, err := apiClient.Blueprints().Get(ctx, blueprintID)
	if err != nil {
		return false
	}
	props, ok := bp.Schema["properties"].(map[string]any)
	if !ok {
		return false
	}
	_, exists := props[newProperty]
	return exists
}
