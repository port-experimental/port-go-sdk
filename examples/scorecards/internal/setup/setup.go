package setup

import (
	"context"
	"errors"
	"fmt"

	"github.com/port-experimental/port-go-sdk/pkg/blueprints"
	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/porter"
	"github.com/port-experimental/port-go-sdk/pkg/scorecards"
)

// BlueprintID is the sandbox blueprint all scorecard examples operate on.
const BlueprintID = "scorecard-example-service"

// BuildScorecardDefinition returns the shared definition used across examples.
func BuildScorecardDefinition() scorecards.ScorecardDefinition {
	return scorecards.ScorecardDefinition{
		Identifier: "slo-coverage",
		Title:      "SLO Coverage",
		Levels: []scorecards.ScorecardLevel{
			{Title: "Basic", Color: "paleBlue"},
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
}

// EnsureBlueprint upserts the shared blueprint so examples always have the properties they rely on.
func EnsureBlueprint(ctx context.Context, cli *client.Client, id string) error {
	err := cli.Blueprints().Upsert(ctx, blueprints.Blueprint{
		Identifier:  id,
		Title:       "Scorecard Example Service",
		Description: "Sandbox blueprint used by the scorecard examples.",
		Icon:        "Cube",
		Schema: map[string]any{
			"properties": map[string]any{
				"tier": map[string]any{
					"type":  "string",
					"title": "Tier",
					"enum":  []string{"1", "2", "3"},
				},
				"slo_url": map[string]any{
					"type":  "string",
					"title": "SLO URL",
				},
				"oncall_rotation": map[string]any{
					"type":  "string",
					"title": "On-call Rotation",
				},
			},
		},
	})
	return enrichAPIError(err)
}

// EnsureScorecard guarantees the shared scorecard exists by creating or updating it.
func EnsureScorecard(ctx context.Context, cli *client.Client, blueprintID string, def scorecards.ScorecardDefinition) error {
	_, err := cli.Scorecards().Get(ctx, blueprintID, def.Identifier)
	if err != nil {
		var perr *porter.Error
		if errors.As(err, &perr) && perr.StatusCode == 404 {
			return enrichAPIError(cli.Scorecards().Create(ctx, blueprintID, def))
		}
		return enrichAPIError(err)
	}
	return enrichAPIError(cli.Scorecards().Update(ctx, blueprintID, def.Identifier, def))
}

func enrichAPIError(err error) error {
	if err == nil {
		return nil
	}
	var perr *porter.Error
	if errors.As(err, &perr) && len(perr.Body) > 0 {
		return fmt.Errorf("%w: %s", err, perr.Body)
	}
	return err
}
