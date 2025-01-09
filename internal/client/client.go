package client

import (
	"context"

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

func NewWithClient(client *serverscom.Client) *Client {
	return &Client{
		scClient: client,
	}
}

// VerifyCredentials checks that token is valid by executing /hosts
func (c *Client) VerifyCredentials(ctx context.Context) error {
	_, err := c.scClient.Hosts.Collection().List(ctx)
	return err
}

func (c *Client) GetScClient() *serverscom.Client {
	return c.scClient
}

// SetVerbose reads verbose from config or cmd flag a
func (c *Client) SetVerbose(verbose bool) *Client {
	c.scClient.SetVerbose(verbose)
	return c
}
