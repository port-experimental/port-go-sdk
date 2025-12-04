package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
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

	const (
		blueprintID = "example_blueprint"
		entityID    = "example_entity"
	)
	ent, err := apiClient.Entities().Get(ctx, blueprintID, entityID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("entity %s %+v\n", ent.Identifier, ent.Properties)
}
