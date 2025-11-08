package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/gateway/infra/postgresql/sqlc/query"
	"github.com/peng225/orochi/internal/gateway/service"
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

func (br *BucketRepository) GetBucketByName(ctx context.Context, name string) (*entity.Bucket, error) {
	b, err := br.q.SelectBucketByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("failed to select bucket: %w", err)
	}
	return &entity.Bucket{
		ID:     b.ID,
		Name:   b.Name,
		Status: string(b.Status),
	}, nil
}
