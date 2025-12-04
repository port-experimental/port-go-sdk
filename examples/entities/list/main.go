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

	const blueprintID = "example_blueprint"
	resp, err := cli.Entities().List(ctx, blueprintID, &entities.ListOptions{
		Query:   "properties.name:Demo",
		PerPage: 20,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found %d entities\n", len(resp.Data))
}
