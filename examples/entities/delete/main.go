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
	targets := []struct {
		Blueprint string
		Entity    string
	}{
		{"example_blueprint", "example_entity"},
		{"example_feature_blueprint", "example_feature"},
	}
	for _, target := range targets {
		if err := cli.Entities().Delete(ctx, target.Blueprint, target.Entity); err != nil {
			var perr *porter.Error
			if errors.As(err, &perr) && perr.StatusCode == 404 {
				log.Printf("entity %s already absent from %s\n", target.Entity, target.Blueprint)
				continue
			}
			log.Fatal(err)
		}
		log.Printf("entity %s deleted from %s\n", target.Entity, target.Blueprint)
	}
}
