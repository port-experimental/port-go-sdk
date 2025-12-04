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
		fmt.Printf("%s - %s\n", u.Email, u.Role)
	}
}
