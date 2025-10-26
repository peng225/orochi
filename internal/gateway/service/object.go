package service

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/peng225/orochi/internal/gateway/api/client"
)

type ObjectService struct {
	// FIXME: hide the implementation detail by adding repository layer.
	c *client.Client
}

func NewObjectStore(baseURL string) *ObjectService {
	c, err := client.NewClient(baseURL)
	if err != nil {
		panic(err)
	}
	return &ObjectService{
		c: c,
	}
}

func (osvc *ObjectService) GetObject(ctx context.Context, bucket, object string) (io.ReadCloser, error) {
	res, err := osvc.c.GetObject(ctx, bucket, object)
	if err != nil {
		return nil, fmt.Errorf("GetObject failed: %w", err)
	}
	switch res.StatusCode {
	case http.StatusOK:
		// 何もしない。
	case http.StatusNotFound:
		return nil, ErrObjectNotFound
	default:
		return nil, fmt.Errorf("GetObject returned unexpected status code: %d", res.StatusCode)
	}
	return res.Body, nil
}

func (osvc *ObjectService) CreateObject(ctx context.Context, bucket, object string, data io.Reader) error {
	res, err := osvc.c.CreateObjectWithBody(ctx, bucket, object, "application/octet-stream", data)
	if err != nil {
		return fmt.Errorf("CreateObject failed: %w", err)
	}
	switch res.StatusCode {
	case http.StatusOK:
		// 何もしない。
	default:
		return fmt.Errorf("CreateObject returned unexpected status code: %d", res.StatusCode)
	}
	return nil
}
