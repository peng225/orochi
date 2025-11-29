package gateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/peng225/orochi/internal/gateway/api/private/client"
)

type Client struct {
	c *client.Client
}

func NewClient(baseURL string) *Client {
	c, err := client.NewClient(baseURL, client.WithHTTPClient(
		&http.Client{
			Timeout: 2 * time.Second,
		},
	))
	if err != nil {
		panic(err)
	}
	return &Client{
		c: c,
	}
}

func (c *Client) CreateBucket(ctx context.Context, name, ecConfig string) error {
	resp, err := c.c.CreateBucket(ctx, client.CreateBucketRequest{
		Name:     &name,
		EcConfig: &ecConfig,
	})
	if err != nil {
		return fmt.Errorf("CreateBucket failed: %w", err)
	}
	switch resp.StatusCode {
	case http.StatusCreated:
		// Do nothing.
	default:
		return fmt.Errorf("CreateBucket returned unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) ChangeBucketStatusToDeleting(ctx context.Context, name string) error {
	resp, err := c.c.UpdateBucketStatusToDeleting(ctx, name)
	if err != nil {
		return fmt.Errorf("UpdateBucketStatusToDeleting failed: %w", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		// Do nothing.
	default:
		return fmt.Errorf("UpdateBucketStatusToDeleting returned unexpected status code: %d",
			resp.StatusCode)
	}

	return nil
}

func (c *Client) DeleteBucket(ctx context.Context, name string) error {
	resp, err := c.c.DeleteBucket(ctx, name)
	if err != nil {
		return fmt.Errorf("DeleteBucket failed: %w", err)
	}
	switch resp.StatusCode {
	case http.StatusNoContent:
		// Do nothing.
	default:
		return fmt.Errorf("DeleteBucket returned unexpected status code: %d",
			resp.StatusCode)
	}

	return nil
}
