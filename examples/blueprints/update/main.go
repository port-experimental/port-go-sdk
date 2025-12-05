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
		blueprintID       = "example_blueprint"
		oldProperty       = "name"
		newProperty       = "display_name"
		oldMirrorProperty = "feature_names"
		newMirrorProperty = "renamed_feature_names"
		oldRelation       = "features"
		newRelation       = "linked_features"
	)
	if err := apiClient.Blueprints().RenameProperty(ctx, blueprintID, oldProperty, newProperty); err != nil {
		handleRenameError(ctx, apiClient, blueprintID, newProperty, "property", propertyExists, err)
	} else {
		log.Printf("renamed property %s -> %s on blueprint %s", oldProperty, newProperty, blueprintID)
	}
	if err := apiClient.Blueprints().RenameMirrorProperty(ctx, blueprintID, oldMirrorProperty, newMirrorProperty); err != nil {
		handleRenameError(ctx, apiClient, blueprintID, newMirrorProperty, "mirror property", mirrorPropertyExists, err)
	} else {
		log.Printf("renamed mirror property %s -> %s on blueprint %s", oldMirrorProperty, newMirrorProperty, blueprintID)
	}
	if err := apiClient.Blueprints().RenameRelation(ctx, blueprintID, oldRelation, newRelation); err != nil {
		handleRenameError(ctx, apiClient, blueprintID, newRelation, "relation", relationExists, err)
	} else {
		log.Printf("renamed relation %s -> %s on blueprint %s", oldRelation, newRelation, blueprintID)
	}
}

type existenceCheck func(ctx context.Context, apiClient *client.Client, blueprintID, identifier string) bool

func handleRenameError(ctx context.Context, apiClient *client.Client, blueprintID, identifier, label string, checker existenceCheck, err error) {
	var perr *porter.Error
	if errors.As(err, &perr) && perr.StatusCode == 404 {
		if checker != nil && checker(ctx, apiClient, blueprintID, identifier) {
			log.Printf("%s already renamed to %s on blueprint %s", label, identifier, blueprintID)
			return
		}
		log.Printf("skipping %s rename on blueprint %s: %s", label, blueprintID, perr.Message)
		return
	}
	log.Fatal(err)
}

func propertyExists(ctx context.Context, apiClient *client.Client, blueprintID, property string) bool {
	bp, err := apiClient.Blueprints().Get(ctx, blueprintID)
	if err != nil {
		return false
	}
	props, ok := bp.Schema["properties"].(map[string]any)
	if !ok {
		return false
	}
	_, exists := props[property]
	return exists
}

func mirrorPropertyExists(ctx context.Context, apiClient *client.Client, blueprintID, mirror string) bool {
	bp, err := apiClient.Blueprints().Get(ctx, blueprintID)
	if err != nil {
		return false
	}
	mirrors, ok := bp.Schema["mirrorProperties"].(map[string]any)
	if !ok {
		return false
	}
	_, exists := mirrors[mirror]
	return exists
}

func relationExists(ctx context.Context, apiClient *client.Client, blueprintID, relation string) bool {
	bp, err := apiClient.Blueprints().Get(ctx, blueprintID)
	if err != nil {
		return false
	}
	if bp.Relations == nil {
		return false
	}
	_, exists := bp.Relations[relation]
	return exists
}
