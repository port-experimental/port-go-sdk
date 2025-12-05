package main

import (
	"context"
	"log"
	"time"

	"github.com/port-experimental/port-go-sdk/pkg/auth"
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

	resp, err := apiClient.Auth().RequestAccessToken(ctx, auth.AccessTokenRequest{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	})
	if err != nil {
		log.Fatal(err)
	}
	prefix := resp.AccessToken
	if len(prefix) > 8 {
		prefix = prefix[:8]
	}
	log.Printf("token starts with %s..., expires in %d seconds", prefix, resp.ExpiresIn)
}
