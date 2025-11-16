package process

import (
	"context"

	"github.com/peng225/orochi/internal/entity"
)

type Transaction interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type JobRepository interface {
	GetJobs(ctx context.Context, limit int) ([]*entity.Job, error)
	DeleteJob(ctx context.Context, id int64) error
}

type BucketRepository interface {
	GetBucket(ctx context.Context, id int64) (*entity.Bucket, error)
	DeleteBucket(ctx context.Context, id int64) error
}

type DatastoreRepository interface {
	GetDatastores(ctx context.Context) ([]*entity.Datastore, error)
	ChangeDatastoreStatus(ctx context.Context, id int64, status entity.DatastoreStatus) error
}

type GatewayClient interface {
	DeleteObject(ctx context.Context, bucket, object string) error
	ListObjectNames(ctx context.Context, bucket string) ([]string, error)
}

type DatastoreClient interface {
	CheckHealthStatus(ctx context.Context) error
}

type DatastoreClientFactory interface {
	New(ds *entity.Datastore) DatastoreClient
}
