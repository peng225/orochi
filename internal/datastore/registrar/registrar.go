package registrar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/peng225/orochi/internal/manager/api/client"
)

func Register(ctx context.Context, baseURL, baseMgrURL string) error {
	c, err := client.NewClient(baseMgrURL)
	if err != nil {
		return fmt.Errorf("failed to get new client: %w", err)
	}
	req := client.CreateDatastoreRequest{
		BaseURL: &baseURL,
	}
	data, err := json.Marshal(&req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.CreateDatastoreWithBody(ctx, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create datastore: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code for create datastore: %d", resp.StatusCode)
	}
	return nil
}
