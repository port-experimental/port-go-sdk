package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/users"
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
	opts := &users.ListUsersOptions{
		Fields: []string{"email", "firstName", "lastName", "roles.name", "teams"},
	}
	userList, err := apiClient.Users().ListUsers(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range userList {
		fmt.Println("-----")
		fmt.Printf("User: %s (%s)\n", displayName(u.FirstName, u.LastName), u.Email)
		roleNames := collectRoles(u.Roles)
		if len(roleNames) > 0 {
			fmt.Printf("Roles: %s\n", strings.Join(roleNames, ", "))
		} else {
			fmt.Println("Roles: none")
		}
		if len(u.Teams) > 0 {
			fmt.Printf("Teams: %s\n", strings.Join(u.Teams, ", "))
		} else {
			fmt.Println("Teams: none")
		}
	}
}

func displayName(first, last string) string {
	name := strings.TrimSpace(strings.TrimSpace(first) + " " + strings.TrimSpace(last))
	if name == "" {
		return "<unnamed>"
	}
	return name
}

func collectRoles(roles []users.UserRole) []string {
	var out []string
	for _, role := range roles {
		if trimmed := strings.TrimSpace(role.Name); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
