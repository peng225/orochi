package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/gateway/infra/postgresql/sqlc/query"
	"github.com/peng225/orochi/internal/gateway/service"
	"github.com/peng225/orochi/internal/pkg/psqlutil"

	_ "github.com/lib/pq"
)

type ObjectMetadataRepository struct {
	q *query.Queries
}

func NewObjectMetadataRepository(db *sql.DB) *ObjectMetadataRepository {
	return &ObjectMetadataRepository{
		q: query.New(db),
	}
}

func (omr *ObjectMetadataRepository) CreateObjectMetadata(
	ctx context.Context, req *service.CreateObjectMetadataRequest,
) (int64, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := omr.q
	if tx != nil {
		q = omr.q.WithTx(tx)
	}
	id, err := q.InsertObjectMetadata(ctx, query.InsertObjectMetadataParams{
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
	tx := psqlutil.TxFromCtx(ctx)
	q := omr.q
	if tx != nil {
		q = omr.q.WithTx(tx)
	}
	om, err := q.SelectObjectMetadataByName(ctx, query.SelectObjectMetadataByNameParams{
		Name:     name,
		BucketID: bucketID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to select object metadata: %w", err)
	}
	if len(om) == 0 {
		return nil, service.ErrNotFound
	}
	if len(om) > 1 {
		// FIXME: need to fix when snapshot is supported.
		return nil, fmt.Errorf("unsupported situation")
	}
	return &entity.ObjectMetadata{
		ID:              om[0].ID,
		Name:            om[0].Name,
		Status:          entity.ObjectStatus(om[0].Status),
		BucketID:        om[0].BucketID,
		LocationGroupID: om[0].LocationGroupID,
	}, nil
}

func (omr *ObjectMetadataRepository) ChangeObjectStatus(
	ctx context.Context, id int64, status entity.ObjectStatus,
) error {
	tx := psqlutil.TxFromCtx(ctx)
	q := omr.q
	if tx != nil {
		q = omr.q.WithTx(tx)
	}
	err := q.UpdateObjectMetadataStatus(ctx, query.UpdateObjectMetadataStatusParams{
		ID:     id,
		Status: query.ObjectStatus(status),
	})
	if err != nil {
		return fmt.Errorf("failed to update object metadata: %w", err)
	}
	return nil
}

func (omr *ObjectMetadataRepository) DeleteObjectMetadata(ctx context.Context, id int64) error {
	tx := psqlutil.TxFromCtx(ctx)
	q := omr.q
	if tx != nil {
		q = omr.q.WithTx(tx)
	}
	err := q.DeleteObjectMetadata(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete object metadata: %w", err)
	}
	return nil
}

func (omr *ObjectMetadataRepository) GetObjectMetadatas(
	ctx context.Context, req *service.GetObjectMetadatasRequest,
) ([]*entity.ObjectMetadata, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := omr.q
	if tx != nil {
		q = omr.q.WithTx(tx)
	}
	ret, err := q.SelectObjectMetadatas(ctx, query.SelectObjectMetadatasParams{
		ID:       req.StartFrom,
		BucketID: req.BucketID,
		Limit:    int32(req.Limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object metadatas: %w", err)
	}

	oms := make([]*entity.ObjectMetadata, 0, len(ret))
	for _, v := range ret {
		oms = append(oms, &entity.ObjectMetadata{
			ID:              v.ID,
			Name:            v.Name,
			BucketID:        v.BucketID,
			LocationGroupID: v.LocationGroupID,
		})
	}
	return oms, nil
}
