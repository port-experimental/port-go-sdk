package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/organization"
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

	current, err := apiClient.Organization().Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	title := strings.TrimSpace(os.Getenv("PORT_PORTAL_TITLE"))
	if title == "" {
		title = "Developer Portal"
	}
	name := current.Name
	if envName := strings.TrimSpace(os.Getenv("PORT_ORG_NAME")); envName != "" {
		name = envName
	}
	if name == "" {
		log.Fatal("organization name is required (either existing or via PORT_ORG_NAME)")
	}
	req := organization.PatchRequest{
		Name: &name,
		Settings: &organization.OrganizationSettings{
			PortalTitle: title,
		},
	}
	if err := apiClient.Organization().Patch(ctx, req); err != nil {
		log.Fatal(err)
	}
	log.Printf("updated portal title to %q", title)
}
