package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/manager/api/client"
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

func (c *Client) GetDatastores(ctx context.Context) ([]*entity.Datastore, error) {
	resp, err := c.c.ListDatastores(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ListDatastores failed: %w", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		// Do nothing.
	default:
		return nil, fmt.Errorf("ListDatastores returned unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}
	dss := make([]*entity.Datastore, 0)
	err = json.Unmarshal(data, &dss)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal datastores: %w", err)
	}
	return dss, nil
}
