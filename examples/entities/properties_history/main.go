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
	apiClient, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := entities.PropertiesHistoryRequest{
		EntityIdentifier:    "demo_service_a",
		BlueprintIdentifier: "example_blueprint",
		PropertyNames:       []string{"name"},
		TimeInterval:        "day",
		TimeRange: &entities.PropertiesHistoryTimeRange{
			Preset: "lastWeek",
		},
	}
	resp, err := apiClient.Entities().PropertiesHistory(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("history min=%f max=%f values=%v\n", resp.Result.MinDate, resp.Result.MaxDate, resp.Result.Data)
}
