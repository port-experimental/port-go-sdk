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

	req := entities.AggregateOverTimeRequest{
		Blueprint:       "example_blueprint",
		TimeRange:       entities.AggregateTimeRange{Preset: "lastWeek"},
		TimeInterval:    "day",
		Query:           map[string]any{"combinator": "and", "rules": []map[string]any{}},
		MeasureTimeBy:   "createdAt",
		AggregationType: "countEntities",
		Func:            "count",
	}
	resp, err := apiClient.Entities().AggregateOverTime(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("min=%f max=%f\n", resp.Result.MinDate, resp.Result.MaxDate)
	for _, row := range resp.Result.Data {
		fmt.Println(row)
	}
}
