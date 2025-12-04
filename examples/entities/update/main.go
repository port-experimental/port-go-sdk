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

	targets := []struct {
		Blueprint string
		Entity    string
		Patch     map[string]any
	}{
		{
			Blueprint: "example_blueprint",
			Entity:    "example_entity",
			Patch: map[string]any{
				"name":        "Demo Entity v2",
				"description": "Updated via SDK",
			},
		},
		{
			Blueprint: "example_feature_blueprint",
			Entity:    "example_feature",
			Patch: map[string]any{
				"description": "Updated feature description v2",
			},
		},
	}
	for _, target := range targets {
		if err := apiClient.Entities().Update(ctx, target.Blueprint, target.Entity, target.Patch); err != nil {
			var perr *porter.Error
			if errors.As(err, &perr) && perr.StatusCode == 404 {
				log.Printf("entity %s not found in %s\n", target.Entity, target.Blueprint)
				continue
			}
			log.Fatal(err)
		}
		log.Printf("entity %s updated in %s\n", target.Entity, target.Blueprint)
	}
}
