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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	secrets, err := apiClient.Organization().ListSecrets(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("found %d secrets\n", len(secrets.Secrets))

	// Optionally create/update/delete a secret when env vars provided.
	name := strings.TrimSpace(os.Getenv("PORT_SECRET_NAME"))
	value := strings.TrimSpace(os.Getenv("PORT_SECRET_VALUE"))
	if name == "" || value == "" {
		return
	}

	createReq := organization.CreateSecretRequest{
		SecretName: name,
		// Never log or check the value into version control; this environment
		// variable is for demo purposes only.
		SecretValue: value,
		Description: "example secret managed by sdk",
	}
	secretResp, err := apiClient.Organization().CreateSecret(ctx, createReq)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("created secret %s", secretResp.Secret.SecretName)

	updateReq := organization.UpdateSecretRequest{
		Description: "updated via example",
	}
	if _, err := apiClient.Organization().UpdateSecret(ctx, name, updateReq); err != nil {
		log.Fatal(err)
	}
	log.Printf("updated secret %s description", name)

	if err := apiClient.Organization().DeleteSecret(ctx, name); err != nil {
		log.Fatal(err)
	}
	log.Printf("deleted secret %s", name)
}
