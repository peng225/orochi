package process

import (
	"context"

	"github.com/peng225/orochi/internal/entity"
)

type JobRepository interface {
	GetJobs(ctx context.Context, limit int) ([]*entity.Job, error)
	DeleteJob(ctx context.Context, id int64) error
}

type BucketRepository interface {
	GetBucket(ctx context.Context, id int64) (*entity.Bucket, error)
	DeleteBucket(ctx context.Context, id int64) error
}
