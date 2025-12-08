package main

import (
	"context"
	"log"
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
	// Replace the email and role name when executing.
	if err := apiClient.Users().AssignRole(ctx, "user@example.com", "Port Admin"); err != nil {
		log.Fatal(err)
	}
	log.Println("role assigned")
}
