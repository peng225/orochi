package datastore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/peng225/orochi/internal/datastore/api/client"
	"github.com/peng225/orochi/internal/gateway/service"
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

func (c *Client) CreateObject(ctx context.Context, object string, data io.Reader) error {
	res, err := c.c.CreateObjectWithBody(ctx, object, "application/octet-stream", data)
	if err != nil {
		return fmt.Errorf("CreateObject failed: %w", err)
	}
	switch res.StatusCode {
	case http.StatusCreated:
		// Do nothing.
	default:
		return fmt.Errorf("CreateObject returned unexpected status code: %d", res.StatusCode)
	}
	return nil
}

func (c *Client) GetObject(ctx context.Context, object string) (io.ReadCloser, error) {
	res, err := c.c.GetObject(ctx, object)
	if err != nil {
		return nil, fmt.Errorf("GetObject failed: %w", err)
	}
	switch res.StatusCode {
	case http.StatusOK:
		// Do nothing.
	case http.StatusNotFound:
		return nil, service.ErrNotFound
	default:
		return nil, fmt.Errorf("GetObject returned unexpected status code: %d", res.StatusCode)
	}
	return res.Body, nil
}

func (c *Client) DeleteObject(ctx context.Context, object string) error {
	res, err := c.c.DeleteObject(ctx, object)
	if err != nil {
		return fmt.Errorf("DeleteObject failed: %w", err)
	}
	switch res.StatusCode {
	case http.StatusNoContent:
		// Do nothing.
	default:
		return fmt.Errorf("DeleteObject returned unexpected status code: %d", res.StatusCode)
	}
	return nil
}
