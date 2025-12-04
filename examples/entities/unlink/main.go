package main

import (
	"context"
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
	cli, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := cli.Entities().UnlinkRelation(ctx, "demo_blueprint", "demo", "owner", []string{"user:123"}); err != nil {
		log.Fatal(err)
	}
	log.Println("relation unlinked")
}
