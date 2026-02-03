package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/examples/scorecards/internal/setup"
	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/porter"
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

	if err := apiClient.Scorecards().Delete(ctx, blueprintID, definition.Identifier); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted scorecard %s from blueprint %s\n", definition.Identifier, blueprintID)

	// Re-create the sample so the other examples continue to work.
	if err := setup.EnsureScorecard(ctx, apiClient, blueprintID, definition); err != nil {
		var perr *porter.Error
		if errors.As(err, &perr) {
			log.Fatalf("failed to restore scorecard: %s (%s)", perr.Message, perr.Body)
		}
		log.Fatalf("failed to restore scorecard: %v", err)
	}
}
