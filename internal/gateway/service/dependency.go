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
	CreateObject(ctx context.Context, bucket, object string, data io.Reader) error
	GetObject(ctx context.Context, bucket, object string) (io.ReadCloser, error)
	DeleteObject(ctx context.Context, bucket, object string) error
}

type ChunkRepositoryFactory interface {
	New(ds *entity.Datastore) ChunkRepository
}

type CreateObjectMetadataRequest struct {
	Name            string
	BucketID        int64
	LocationGroupID int64
}

type ObjectMetadataRepository interface {
	CreateObjectMetadata(ctx context.Context, req *CreateObjectMetadataRequest) (int64, error)
	GetObjectMetadataByName(ctx context.Context, name string, bucketID int64) (*entity.ObjectMetadata, error)
	DeleteObjectMetadata(ctx context.Context, id int64) error
}

type BucketRepository interface {
	GetBucketByName(ctx context.Context, name string) (*entity.Bucket, error)
}

type LocationGroupRepository interface {
	GetLocationGroup(ctx context.Context, id int64) (*entity.LocationGroup, error)
	GetLocationGroups(ctx context.Context) ([]*entity.LocationGroup, error)
}
