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

	scorecards, err := apiClient.Scorecards().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, sc := range scorecards {
		fmt.Printf("%s (blueprint=%s, rules=%d)\n", sc.Identifier, sc.Blueprint, len(sc.Rules))
	}
	fmt.Printf("total scorecards: %d\n", len(scorecards))
}
