package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
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

	// Replace with blueprint/rule identifiers that exist in your account.
	blueprintID := "service"

	definition := scorecards.ScorecardDefinition{
		Identifier:  "slo-coverage",
		Title:       "SLO Coverage",
		Description: "Ensures every service exposes customer-facing SLO metrics.",
		Levels: []scorecards.ScorecardLevel{
			{Title: "Gold", Color: "gold"},
			{Title: "Silver", Color: "silver"},
			{Title: "Bronze", Color: "bronze"},
		},
		Filter: map[string]any{
			"combinator": "and",
			"conditions": []any{
				map[string]any{
					"property": "tier",
					"operator": "in",
					"value":    []any{"1", "2"},
				},
			},
		},
		Rules: []scorecards.ScorecardRule{
			{
				Identifier: "has-slo",
				Title:      "Has published SLO",
				Level:      "Gold",
				Query: map[string]any{
					"combinator": "and",
					"conditions": []any{
						map[string]any{
							"property": "slo_url",
							"operator": "isNotEmpty",
						},
					},
				},
			},
			{
				Identifier: "alerts-wired",
				Title:      "Alerts wired to pager",
				Level:      "Silver",
				Query: map[string]any{
					"combinator": "and",
					"conditions": []any{
						map[string]any{
							"property": "oncall_rotation",
							"operator": "isNotEmpty",
						},
					},
				},
			},
		},
	}

	if err := apiClient.Scorecards().Create(ctx, blueprintID, definition); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created scorecard %s on blueprint %s\n", definition.Identifier, blueprintID)
}
