package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/peng225/orochi/internal/gateway/infra/postgresql/sqlc/query"
	"github.com/peng225/orochi/internal/pkg/psqlutil"

	_ "github.com/lib/pq"
)

type ObjectVersionRepository struct {
	q *query.Queries
}

func NewObjectVersionRepository(db *sql.DB) *ObjectVersionRepository {
	return &ObjectVersionRepository{
		q: query.New(db),
	}
}

func (ovr *ObjectVersionRepository) CreateObjectVersion(ctx context.Context, objectID int64) (int64, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := ovr.q
	if tx != nil {
		q = ovr.q.WithTx(tx)
	}
	id, err := q.InsertObjectVersion(ctx, query.InsertObjectVersionParams{
		ObjectID:   objectID,
		UpdateTime: time.Now(),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert object version: %w", err)
	}
	return id, nil
}

func (ovr *ObjectVersionRepository) DeleteObjectVersionsByObjectID(ctx context.Context, objectID int64) error {
	tx := psqlutil.TxFromCtx(ctx)
	q := ovr.q
	if tx != nil {
		q = ovr.q.WithTx(tx)
	}
	err := q.DeleteObjectVersionsByObjectID(ctx, objectID)
	if err != nil {
		return fmt.Errorf("failed to delete object versions: %w", err)
	}
	return nil
}
