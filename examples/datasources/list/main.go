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
	sources, err := apiClient.DataSources().ListIntegrations(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found %d integrations\n", len(sources))
	webhooks, err := apiClient.DataSources().ListWebhooks(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found %d webhooks\n", len(webhooks))
}
