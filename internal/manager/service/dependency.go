package service

import (
	"context"

	"github.com/peng225/orochi/internal/entity"
)

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
}

type LocationGroupRepository interface {
	CreateLocationGroup(ctx context.Context, req *CreateLocationGroupRequest) (int64, error)
	UpdateDesiredDatastores(ctx context.Context, id int64, desiredDatastores []int64) error
	GetLocationGroups(ctx context.Context) ([]*entity.LocationGroup, error)
}

type CreateBucketRequest struct {
	Name string
}

type BucketRepository interface {
	CreateBucket(ctx context.Context, req *CreateBucketRequest) (int64, error)
	GetBucketByName(ctx context.Context, name string) (*entity.Bucket, error)
	ChangeBucketStatus(ctx context.Context, id int64, status string) error
}

type CreateJobRequest struct {
	Name string
	Data []byte
}

type JobRepository interface {
	CreateJob(ctx context.Context, req *CreateJobRequest) (int64, error)
}
