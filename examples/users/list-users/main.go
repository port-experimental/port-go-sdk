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
	userList, err := apiClient.Users().ListUsers(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	teamCache := map[string]*users.Team{}
	getTeam := func(name string) *users.Team {
		if cached, ok := teamCache[name]; ok {
			return cached
		}
		team, err := apiClient.Users().GetTeam(ctx, name)
		if err != nil {
			log.Printf("warn: failed to fetch team %s: %v\n", name, err)
			return nil
		}
		teamCopy := team
		teamCache[name] = &teamCopy
		return &teamCopy
	}

	for _, u := range userList {
		fmt.Println("-----")
		fmt.Printf("User: %s (%s)\n", displayName(u.FirstName, u.LastName), u.Email)
		detail, err := apiClient.Users().GetUser(ctx, u.Email)
		if err != nil {
			log.Printf("warn: failed to fetch details for %s: %v\n", u.Email, err)
			detail = u
		}
		roleNames := collectRoles(detail.Roles)
		if len(roleNames) > 0 {
			fmt.Printf("Roles: %s\n", strings.Join(roleNames, ", "))
		} else {
			fmt.Println("Roles: none")
		}
		printTeams(detail.Teams, getTeam)
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

func printTeams(teamNames []string, resolve func(string) *users.Team) {
	if len(teamNames) == 0 {
		fmt.Println("Teams: none")
		return
	}
	fmt.Println("Teams:")
	for _, name := range teamNames {
		line := fmt.Sprintf("  - %s", name)
		if team := resolve(name); team != nil {
			members := teamMemberEmails(team.Users)
			if len(members) > 0 {
				line = fmt.Sprintf("%s (members: %s)", line, strings.Join(members, ", "))
			}
		}
		fmt.Println(line)
	}
}

func teamMemberEmails(users []users.TeamMember) []string {
	var out []string
	for _, member := range users {
		if email := strings.TrimSpace(member.Email); email != "" {
			out = append(out, email)
		}
	}
	return out
}
