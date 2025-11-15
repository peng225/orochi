package datastore

import (
	"context"
	"fmt"
	"net/http"

	"github.com/peng225/orochi/internal/datastore/api/client"
)

type Client struct {
	c *client.Client
}

func NewClient(baseURL string) *Client {
	c, err := client.NewClient(baseURL)
	if err != nil {
		panic(err)
	}
	return &Client{
		c: c,
	}
}

func (c *Client) CheckHealthStatus(ctx context.Context) error {
	res, err := c.c.CheckHealthStatus(ctx)
	if err != nil {
		return fmt.Errorf("CheckHealthStatus failed: %w", err)
	}
	switch res.StatusCode {
	case http.StatusOK:
		// Do nothing.
	default:
		return fmt.Errorf("CheckHealthStatus returned unexpected status code: %d", res.StatusCode)
	}
	return nil
}
