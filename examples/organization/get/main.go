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

	org, err := apiClient.Organization().Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	hiddenCount := 0
	if org.Settings != nil {
		hiddenCount = len(org.Settings.HiddenBlueprints)
	}
	fmt.Printf("Organization: %s (hidden blueprints: %d)\n", org.Name, hiddenCount)
	if org.Announcement != nil && org.Announcement.Enabled {
		fmt.Printf("Announcement: %s\n", org.Announcement.Content)
	}
}
