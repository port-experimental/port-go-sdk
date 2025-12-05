package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/blueprints"
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

	const blueprintID = "example_blueprint"

	perms, err := apiClient.Blueprints().GetPermissions(ctx, blueprintID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("current entity permissions: %+v\n", perms.Entities)

	update := blueprints.BlueprintPermissions{
		Entities: &blueprints.BlueprintEntityPermissions{
			Read: &blueprints.BlueprintPermissionRule{
				Roles: []string{"Admin"},
			},
			UpdateProperties: map[string]blueprints.BlueprintPermissionRule{
				"owner": {
					Roles: []string{"Admin"},
				},
			},
		},
	}

	if err := apiClient.Blueprints().UpdatePermissions(ctx, blueprintID, update); err != nil {
		log.Fatal(err)
	}
	log.Printf("updated permissions for blueprint %s", blueprintID)
}
