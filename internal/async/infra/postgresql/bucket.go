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
	q *query.Queries
}

func NewBucketRepository(db *sql.DB) *BucketRepository {
	return &BucketRepository{
		q: query.New(db),
	}
}

func (br *BucketRepository) GetBucket(ctx context.Context, id int64) (*entity.Bucket, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := br.q
	if tx != nil {
		q = br.q.WithTx(tx)
	}
	bucket, err := q.SelectBucket(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, process.ErrNotFound
		}
		return nil, fmt.Errorf("failed to select bucket: %w", err)
	}
	return &entity.Bucket{
		ID:     bucket.ID,
		Name:   bucket.Name,
		Status: entity.BucketStatus(bucket.Status),
	}, nil
}

func (br *BucketRepository) DeleteBucket(ctx context.Context, id int64) error {
	tx := psqlutil.TxFromCtx(ctx)
	q := br.q
	if tx != nil {
		q = br.q.WithTx(tx)
	}
	err := q.DeleteBucket(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}
	return nil
}
