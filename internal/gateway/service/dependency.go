package service

import (
	"context"
	"io"

	"github.com/peng225/orochi/internal/entity"
)

type CreateDatastoreRequest struct {
	BaseURL string
}

type DatastoreRepository interface {
	GetDatastores(ctx context.Context) ([]*entity.Datastore, error)
}

type ChunkRepository interface {
	GetObject(ctx context.Context, bucket, object string) (io.ReadCloser, error)
	CreateObject(ctx context.Context, bucket, object string, data io.Reader) error
}

type ChunkRepositoryFactory interface {
	New(ds *entity.Datastore) ChunkRepository
}
