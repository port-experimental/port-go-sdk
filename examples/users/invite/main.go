package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
	"github.com/port-experimental/port-go-sdk/pkg/users"
)

func main() {
	targetEmail := strings.TrimSpace(os.Getenv("PORT_INVITE_EMAIL"))
	if targetEmail == "" {
		log.Fatal("set PORT_INVITE_EMAIL to the address you want to invite")
	}

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

	req := users.InviteRequest{
		Email: targetEmail,
		Roles: []string{"Member"},
	}
	if err := apiClient.Users().Invite(ctx, req); err != nil {
		log.Fatal(err)
	}
	log.Printf("invited %s\n", targetEmail)
}
