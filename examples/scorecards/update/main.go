package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/examples/scorecards/internal/setup"
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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	blueprintID := setup.BlueprintID
	definition := setup.BuildScorecardDefinition()

	if err := setup.EnsureBlueprint(ctx, apiClient, blueprintID); err != nil {
		log.Fatalf("ensure blueprint: %v", err)
	}
	if err := setup.EnsureScorecard(ctx, apiClient, blueprintID, definition); err != nil {
		log.Fatalf("ensure scorecard: %v", err)
	}

	// Mutate the definition before sending an update.
	definition.Title = "SLO Coverage (updated)"
	if len(definition.Rules) > 0 {
		definition.Rules[0].Level = "Silver"
	}

	if err := apiClient.Scorecards().Update(ctx, blueprintID, definition.Identifier, definition); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("updated scorecard %s on blueprint %s\n", definition.Identifier, blueprintID)
}
