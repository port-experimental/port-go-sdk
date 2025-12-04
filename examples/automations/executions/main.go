package main

import (
	"context"
	"fmt"
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
	execs, err := apiClient.Automations().ListExecutions(ctx, "automation_id")
	if err != nil {
		log.Fatal(err)
	}
	for _, ex := range execs {
		fmt.Printf("%s status=%s\n", ex.ID, ex.Status)
	}
}
