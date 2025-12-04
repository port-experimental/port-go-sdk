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
	const (
		blueprintID = "example_blueprint"
		entityID    = "example_entity"
	)
	ent := entities.Entity{
		Identifier: entityID,
		Properties: map[string]any{
			"name":  "Demo Entity",
			"owner": "team@example.com",
		},
	}
	err = c.Entities().Create(ctx, blueprintID, ent)
	if err != nil {
		var perr *porter.Error
		if errors.As(err, &perr) && perr.StatusCode == 409 {
			log.Printf("entity %s already exists in %s\n", entityID, blueprintID)
			return
		}
		log.Fatal(err)
	}
	log.Printf("entity %s created in blueprint %s\n", entityID, blueprintID)
}
