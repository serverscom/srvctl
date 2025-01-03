package client

import (
	"context"
	"fmt"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

// Client wraps serverscom API client
type Client struct {
	scClient *serverscom.Client
}

func NewClient(token string, endpoint string) *Client {
	if endpoint != "" {
		return &Client{
			scClient: serverscom.NewClientWithEndpoint(token, endpoint),
		}
	}
	return &Client{
		scClient: serverscom.NewClient(token),
	}
}

// VerifyCredentials checks that token is valid by executing /hosts
func (c *Client) VerifyCredentials(ctx context.Context) error {
	_, err := c.scClient.Hosts.Collection().List(ctx)
	if err != nil {
		return fmt.Errorf("invalid API token: %w", err)
	}
	return nil
}

func (c *Client) GetScClient() *serverscom.Client {
	return c.scClient
}

// SetVerbose reads verbose from config or cmd flag a
func (c *Client) SetVerbose(verbose bool) *Client {
	c.scClient.SetVerbose(verbose)
	return c
}
