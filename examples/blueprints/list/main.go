package main

import (
	"context"
	"fmt"
	"log"
	"strings"
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
	bps, err := apiClient.Blueprints().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found %d blueprints\n", len(bps))
	for _, bp := range bps {
		fmt.Println("-----")
		fmt.Printf("Blueprint: %s (%s)\n", bp.Identifier, bp.Title)
		if schema, ok := bp.Schema["properties"].(map[string]any); ok && len(schema) > 0 {
			fmt.Println("Properties:")
			for name, value := range schema {
				fmt.Printf("  - %s: %s\n", name, describeProperty(value))
			}
		} else {
			fmt.Println("Properties: none")
		}
	}
}

func describeProperty(raw any) string {
	switch v := raw.(type) {
	case map[string]any:
		if t, ok := v["type"].(string); ok {
			var extras []string
			for key, val := range v {
				if key == "type" {
					continue
				}
				extras = append(extras, fmt.Sprintf("%s=%v", key, val))
			}
			if len(extras) > 0 {
				return fmt.Sprintf("type=%s (%s)", t, strings.Join(extras, ", "))
			}
			return fmt.Sprintf("type=%s", t)
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
