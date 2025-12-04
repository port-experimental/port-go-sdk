package main

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const blueprintID = "example_blueprint"
	bp := blueprints.Blueprint{
		Identifier:  blueprintID,
		Title:       "Example Blueprint",
		Description: "Created via Go SDK",
		Schema: map[string]any{
			"properties": map[string]any{
				"name": map[string]any{
					"type": "string",
				},
				"owner": map[string]any{
					"type": "string",
				},
			},
		},
		Icon: "Cube",
	}
	if err := cli.Blueprints().Upsert(ctx, bp); err != nil {
		log.Fatal(err)
	}
	log.Printf("blueprint %s created\n", blueprintID)
}
