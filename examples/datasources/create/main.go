package main

import (
	"context"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/datasources"
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
	ds := datasources.DataSource{
		Identifier: "webhook_example",
		Title:      "Example Webhook",
		Type:       "webhook",
		Config: map[string]any{
			"url": "https://ingest.getport.io/example",
		},
	}
	if err := apiClient.DataSources().Create(ctx, ds); err != nil {
		log.Fatal(err)
	}
	log.Println("data source created")
}
