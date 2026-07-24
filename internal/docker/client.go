package docker

import (
	"context"
	"github.com/docker/docker/client"
)

type Client struct {
	Cli *client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{Cli: cli}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.Cli.Ping(ctx)
	return err
}
