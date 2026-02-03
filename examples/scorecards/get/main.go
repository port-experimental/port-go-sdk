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

	// Replace with identifiers that exist in your account.
	blueprintID := "service"
	scorecardID := "slo-coverage"

	scorecard, err := apiClient.Scorecards().Get(ctx, blueprintID, scorecardID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s (%s) has %d rules\n", scorecard.Title, scorecard.Identifier, len(scorecard.Rules))
	for _, rule := range scorecard.Rules {
		fmt.Printf("- [%s] %s (level=%s)\n", rule.Identifier, rule.Title, rule.Level)
	}
}
