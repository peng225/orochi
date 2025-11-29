package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/peng225/orochi/internal/gateway/api/client"
	pclient "github.com/peng225/orochi/internal/gateway/api/private/client"
)

type Client struct {
	c  *client.Client
	pc *pclient.Client
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
	pc, err := pclient.NewClient(baseURL, pclient.WithHTTPClient(
		&http.Client{
			Timeout: 2 * time.Second,
		},
	))
	if err != nil {
		panic(err)
	}
	return &Client{
		c:  c,
		pc: pc,
	}
}

func (c *Client) DeleteObject(ctx context.Context, bucket, object string) error {
	resp, err := c.c.DeleteObject(ctx, bucket, object)
	if err != nil {
		return fmt.Errorf("DeleteObject failed: %w", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) ListObjectNames(ctx context.Context, bucket string) ([]string, error) {
	resp, err := c.c.ListObjects(ctx, bucket, nil)
	if err != nil {
		return nil, fmt.Errorf("ListObjects failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	objectNames := make([]string, 0)
	err = json.Unmarshal(data, &objectNames)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal object names: %w", err)
	}
	return objectNames, nil
}

func (c *Client) DeleteBucket(ctx context.Context, name string) error {
	resp, err := c.pc.DeleteBucket(ctx, name)
	if err != nil {
		return fmt.Errorf("DeleteBucket failed: %w", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
