package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/gateway/infra/postgresql/sqlc/query"
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

func (br *BucketRepository) GetBucketsByName(ctx context.Context, name string) ([]*entity.Bucket, error) {
	buckets, err := br.q.SelectBucketsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to select buckets by name: %w", err)
	}
	res := make([]*entity.Bucket, 0, len(buckets))
	for _, b := range buckets {
		res = append(res, &entity.Bucket{
			ID:     b.ID,
			Name:   b.Name,
			Status: string(b.Status),
		})
	}
	return res, nil
}
