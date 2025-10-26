package service

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/peng225/orochi/internal/gateway/api/client"
)

type ObjectService interface {
	GetObject(ctx context.Context, bucket, object string) (io.ReadCloser, error)
	CreateObject(ctx context.Context, bucket, object string, data io.Reader) error
}

type DataStoreObjectStore struct {
	c *client.Client
}

func NewDataStoreObjectStore(baseURL string) *DataStoreObjectStore {
	c, err := client.NewClient(baseURL)
	if err != nil {
		panic(err)
	}
	return &DataStoreObjectStore{
		c: c,
	}
}

func (dsos *DataStoreObjectStore) GetObject(ctx context.Context, bucket, object string) (io.ReadCloser, error) {
	res, err := dsos.c.GetObject(ctx, bucket, object)
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

func (dsos *DataStoreObjectStore) CreateObject(ctx context.Context, bucket, object string, data io.Reader) error {
	res, err := dsos.c.CreateObjectWithBody(ctx, bucket, object, "application/octet-stream", data)
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
