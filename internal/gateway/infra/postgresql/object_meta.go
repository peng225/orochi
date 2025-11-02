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

type ObjectMetadataRepository struct {
	db *sql.DB
	q  *query.Queries
}

func NewObjectMetadataRepository() *ObjectMetadataRepository {
	db := psqlutil.InitDB()
	return &ObjectMetadataRepository{
		db: db,
		q:  query.New(db),
	}
}

func (omr *ObjectMetadataRepository) Close() error {
	return omr.db.Close()
}

func (omr *ObjectMetadataRepository) CreateObjectMetadata(ctx context.Context, req *service.CreateObjectMetadataRequest) (int64, error) {
	id, err := omr.q.CreateObjectMetadata(ctx, query.CreateObjectMetadataParams{
		Name:            req.Name,
		BucketID:        req.BucketID,
		LocationGroupID: req.LocationGroupID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert object metadata: %w", err)
	}
	return id, nil
}

func (omr *ObjectMetadataRepository) GetObjectMetadataByName(
	ctx context.Context,
	name string,
	bucketID int64,
) (*entity.ObjectMetadata, error) {
	om, err := omr.q.SelectObjectMetadataByName(ctx, query.SelectObjectMetadataByNameParams{
		Name:     name,
		BucketID: bucketID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("failed to select object metadata: %w", err)
	}
	if len(om) != 1 {
		// FIXME: need to fix when snapshot is supported.
		return nil, fmt.Errorf("unsupported situation")
	}
	return &entity.ObjectMetadata{
		ID:              om[0].ID,
		Name:            om[0].Name,
		BucketID:        om[0].BucketID,
		LocationGroupID: om[0].LocationGroupID,
	}, nil
}
