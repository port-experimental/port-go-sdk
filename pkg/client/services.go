package client

import (
	"github.com/port-experimental/port-go-sdk/pkg/automations"
	"github.com/port-experimental/port-go-sdk/pkg/blueprints"
	"github.com/port-experimental/port-go-sdk/pkg/datasources"
	"github.com/port-experimental/port-go-sdk/pkg/entities"
	"github.com/port-experimental/port-go-sdk/pkg/users"
)

// Entities exposes entity endpoints.
func (c *Client) Entities() *entities.Service {
	return entities.New(c)
}

// Blueprints exposes blueprint endpoints.
func (c *Client) Blueprints() *blueprints.Service {
	return blueprints.New(c)
}

// DataSources exposes data source endpoints.
func (c *Client) DataSources() *datasources.Service {
	return datasources.New(c)
}

// Automations exposes automation endpoints.
func (c *Client) Automations() *automations.Service {
	return automations.New(c)
}

// Users exposes user/team endpoints.
func (c *Client) Users() *users.Service {
	return users.New(c)
}
