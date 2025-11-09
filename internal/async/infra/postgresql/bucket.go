package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/peng225/orochi/internal/async/infra/postgresql/sqlc/query"
	"github.com/peng225/orochi/internal/async/process"
	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/pkg/psqlutil"

	_ "github.com/lib/pq"
)

type BucketRepository struct {
	db *sql.DB
	q  *query.Queries
}

func NewBucketRepository() *BucketRepository {
	db := psqlutil.InitDB()
	return &BucketRepository{
		db: db,
		q:  query.New(db),
	}
}

func (br *BucketRepository) Close() error {
	return br.db.Close()
}

func (br *BucketRepository) GetBucket(ctx context.Context, id int64) (*entity.Bucket, error) {
	bucket, err := br.q.SelectBucket(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, process.ErrNotFound
		}
		return nil, fmt.Errorf("failed to select bucket: %w", err)
	}
	return &entity.Bucket{
		ID:     bucket.ID,
		Name:   bucket.Name,
		Status: string(bucket.Status),
	}, nil
}

func (br *BucketRepository) DeleteBucket(ctx context.Context, id int64) error {
	err := br.q.DeleteBucket(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}
	return nil
}
