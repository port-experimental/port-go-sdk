package main

import (
	"context"
	"fmt"
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
	cli, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	blueprints := []string{"example_blueprint", "example_feature_blueprint"}
	for _, bp := range blueprints {
		resp, err := cli.Entities().List(ctx, bp, &entities.ListOptions{
			Limit: 20,
		})
		if err != nil {
			log.Fatalf("list %s: %v", bp, err)
		}
		fmt.Printf("[%s] found %d entities\n", bp, len(resp.Entities))
		for _, ent := range resp.Entities {
			fmt.Printf("  - %s (%s)\n", ent.Identifier, ent.Blueprint)
		}
	}
}
