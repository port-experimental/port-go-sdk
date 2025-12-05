package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/client"
	"github.com/port-experimental/port-go-sdk/pkg/config"
)

func main() {
	targetEmail := strings.TrimSpace(os.Getenv("PORT_INVITE_EMAIL"))
	if targetEmail == "" {
		log.Fatal("PORT_INVITE_EMAIL must be set to the user email to rotate credentials for")
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

	if err := apiClient.Auth().RotateCredentials(ctx, targetEmail); err != nil {
		log.Fatal(err)
	}
	log.Printf("rotated credentials for %s", targetEmail)
}
