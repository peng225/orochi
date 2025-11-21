package service

import (
	"context"
	"io"

	"github.com/peng225/orochi/internal/entity"
)

type Transaction interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type CreateDatastoreRequest struct {
	BaseURL string
}

type DatastoreRepository interface {
	GetDatastores(ctx context.Context) ([]*entity.Datastore, error)
}

type DatastoreClient interface {
	CreateObject(ctx context.Context, object string, data io.Reader) error
	GetObject(ctx context.Context, object string) (io.ReadCloser, error)
	DeleteObject(ctx context.Context, object string) error
}

type DatastoreClientFactory interface {
	New(ds *entity.Datastore) DatastoreClient
}

type CreateObjectMetadataRequest struct {
	Name            string
	BucketID        int64
	LocationGroupID int64
}

type GetObjectMetadatasRequest struct {
	BucketID  int64
	StartFrom int64
	Limit     int
}

type ObjectMetadataRepository interface {
	CreateObjectMetadata(ctx context.Context, req *CreateObjectMetadataRequest) (int64, error)
	GetObjectMetadataByName(ctx context.Context, name string, bucketID int64) (*entity.ObjectMetadata, error)
	GetObjectMetadatas(ctx context.Context, req *GetObjectMetadatasRequest) ([]*entity.ObjectMetadata, error)
	ChangeObjectStatus(ctx context.Context, id int64, status entity.ObjectStatus) error
	DeleteObjectMetadata(ctx context.Context, id int64) error
}

type ObjectVersionRepository interface {
	CreateObjectVersion(ctx context.Context, objectID int64) (int64, error)
	DeleteObjectVersionsByObjectID(ctx context.Context, objectID int64) error
}

type BucketRepository interface {
	GetBucketByName(ctx context.Context, name string) (*entity.Bucket, error)
}

type LocationGroupRepository interface {
	GetLocationGroup(ctx context.Context, id int64) (*entity.LocationGroup, error)
	GetLocationGroupsByECConfigID(ctx context.Context, ecConfigID int64) ([]*entity.LocationGroup, error)
}

type ECConfigRepository interface {
	GetECConfig(ctx context.Context, id int64) (*entity.ECConfig, error)
}
