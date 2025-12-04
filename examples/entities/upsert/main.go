package main

import (
	"context"
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
			"name": "Demo Entity",
		},
	}
	if err := c.Entities().Upsert(ctx, blueprintID, ent); err != nil {
		log.Fatal(err)
	}
	log.Println("entity upsert succeeded")
}
