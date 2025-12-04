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
	cli, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	const (
		blueprintID = "example_blueprint"
		entityID    = "example_entity"
	)
	patch := map[string]any{
		"name":        "Demo Entity v2",
		"description": "Updated via SDK",
	}
	if err := cli.Entities().Update(ctx, blueprintID, entityID, patch); err != nil {
		var perr *porter.Error
		if errors.As(err, &perr) && perr.StatusCode == 404 {
			log.Printf("entity %s not found in %s\n", entityID, blueprintID)
			return
		}
		log.Fatal(err)
	}
	log.Println("entity updated")
}
