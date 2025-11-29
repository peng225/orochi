package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/manager/infra/postgresql/sqlc/query"
	"github.com/peng225/orochi/internal/manager/service"
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

func (br *BucketRepository) CreateBucket(ctx context.Context, req *service.CreateBucketRequest) (int64, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := br.q
	if tx != nil {
		q = br.q.WithTx(tx)
	}
	id, err := q.InsertBucket(ctx, query.InsertBucketParams{
		Name:     req.Name,
		EcConfig: req.ECConfig,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert bucket: %w", err)
	}
	return id, nil
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
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("failed to select bucket: %w", err)
	}
	return &entity.Bucket{
		ID:       bucket.ID,
		Name:     bucket.Name,
		ECConfig: bucket.EcConfig,
		Status:   entity.BucketStatus(bucket.Status),
	}, nil
}

func (br *BucketRepository) GetBucketByName(ctx context.Context, name string) (*entity.Bucket, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := br.q
	if tx != nil {
		q = br.q.WithTx(tx)
	}
	bucket, err := q.SelectBucketByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("failed to select bucket by name: %w", err)
	}
	return &entity.Bucket{
		ID:       bucket.ID,
		Name:     bucket.Name,
		ECConfig: bucket.EcConfig,
		Status:   entity.BucketStatus(bucket.Status),
	}, nil
}

func (br *BucketRepository) ChangeBucketStatus(ctx context.Context, id int64, status entity.BucketStatus) error {
	tx := psqlutil.TxFromCtx(ctx)
	q := br.q
	if tx != nil {
		q = br.q.WithTx(tx)
	}
	err := q.UpdateBucketStatus(ctx, query.UpdateBucketStatusParams{
		ID:     id,
		Status: query.BucketStatus(status),
	})
	if err != nil {
		return fmt.Errorf("failed to update bucket status: %w", err)
	}
	return nil
}
