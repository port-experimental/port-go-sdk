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
	"github.com/port-experimental/port-go-sdk/pkg/scorecards"
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

	// Use a dedicated example blueprint so we don't mutate an existing catalog entry.
	blueprintID := setup.BlueprintID
	if err := setup.EnsureBlueprint(ctx, apiClient, blueprintID); err != nil {
		log.Fatalf("ensure blueprint: %v", err)
	}

	definition := setup.BuildScorecardDefinition()

	if err := createOrUpdateScorecard(ctx, apiClient, blueprintID, definition); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("scorecard %s is ready on blueprint %s\n", definition.Identifier, blueprintID)
}

func createOrUpdateScorecard(ctx context.Context, cli *client.Client, blueprintID string, def scorecards.ScorecardDefinition) error {
	if err := cli.Scorecards().Create(ctx, blueprintID, def); err != nil {
		var perr *porter.Error
		if errors.As(err, &perr) && (perr.StatusCode == 409 || perr.StatusCode == 422) {
			if err := cli.Scorecards().Update(ctx, blueprintID, def.Identifier, def); err != nil {
				return enrichAPIError(err)
			}
			return nil
		}
		return enrichAPIError(err)
	}
	return nil
}

func enrichAPIError(err error) error {
	var perr *porter.Error
	if errors.As(err, &perr) && len(perr.Body) > 0 {
		return fmt.Errorf("%w: %s", err, perr.Body)
	}
	return err
}
