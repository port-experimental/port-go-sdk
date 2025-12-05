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
	apiClient, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ids := []string{"demo_service_a", "demo_service_b"}
	resp, err := apiClient.Entities().BulkDelete(ctx, "example_blueprint", ids, &entities.BulkDeleteOptions{
		DeleteDependents: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("deleted %d entities", len(resp.DeletedEntities))
}
