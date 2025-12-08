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
	ds, err := apiClient.DataSources().GetIntegration(ctx, "integration_id", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("integration %s app=%s\n", ds.Identifier, ds.InstallationAppType)
}
