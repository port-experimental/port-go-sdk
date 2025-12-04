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
	cli, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	users, err := cli.Users().ListUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range users {
		fmt.Println("-----")
		fmt.Printf("User: %s (%s)\n", safeDisplay(u.Name), u.Email)
		fmt.Printf("Role: %s\n", u.Role)
		if len(u.Teams) > 0 {
			fmt.Printf("Teams: %s\n", strings.Join(u.Teams, ", "))
		} else {
			fmt.Println("Teams: none")
		}
	}
}

func safeDisplay(s string) string {
	if strings.TrimSpace(s) == "" {
		return "<unnamed>"
	}
	return s
}
