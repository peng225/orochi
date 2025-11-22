package service

import (
	"context"

	"github.com/peng225/orochi/internal/entity"
)

type Transaction interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type CreateDatastoreRequest struct {
	BaseURL string
}

type DatastoreRepository interface {
	CreateDatastore(ctx context.Context, req *CreateDatastoreRequest) (int64, error)
	GetDatastore(ctx context.Context, id int64) (*entity.Datastore, error)
	GetDatastoreIDs(ctx context.Context) ([]int64, error)
}

type CreateLocationGroupRequest struct {
	Datastores []int64
	ECConfigID int64
}

type LocationGroupRepository interface {
	CreateLocationGroup(ctx context.Context, req *CreateLocationGroupRequest) (int64, error)
	GetLocationGroupsByECConfigID(ctx context.Context, ecConfigID int64) ([]*entity.LocationGroup, error)
	UpdateLocationGroupStatus(ctx context.Context, id int64, status entity.LocationGroupStatus) error
}

type CreateBucketRequest struct {
	Name       string
	ECConfigID int64
}

type BucketRepository interface {
	CreateBucket(ctx context.Context, req *CreateBucketRequest) (int64, error)
	GetBucket(ctx context.Context, id int64) (*entity.Bucket, error)
	GetBucketByName(ctx context.Context, name string) (*entity.Bucket, error)
	ChangeBucketStatus(ctx context.Context, id int64, status entity.BucketStatus) error
}

type CreateJobRequest struct {
	Kind string
	Data []byte
}

type JobRepository interface {
	CreateJob(ctx context.Context, req *CreateJobRequest) (int64, error)
}

type CreateECConfigRequest struct {
	NumData   int
	NumParity int
}

type ECConfigRepository interface {
	CreateECConfig(ctx context.Context, req *CreateECConfigRequest) (int64, error)
	GetECConfigByNumbers(ctx context.Context, numData, numParity int) (*entity.ECConfig, error)
	GetECConfigs(ctx context.Context) ([]*entity.ECConfig, error)
}
