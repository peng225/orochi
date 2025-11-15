package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/peng225/orochi/internal/async/infra/postgresql/sqlc/query"
	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/pkg/psqlutil"

	_ "github.com/lib/pq"
)

type DatastoreRepository struct {
	q *query.Queries
}

func NewDatastoreRepository(db *sql.DB) *DatastoreRepository {
	return &DatastoreRepository{
		q: query.New(db),
	}
}

func (dr *DatastoreRepository) GetDatastores(ctx context.Context) ([]*entity.Datastore, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := dr.q
	if tx != nil {
		q = dr.q.WithTx(tx)
	}
	dss, err := q.SelectDatastores(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to select datastores: %w", err)
	}
	datastores := make([]*entity.Datastore, len(dss))
	for i, ds := range dss {
		datastores[i] = &entity.Datastore{
			ID:      ds.ID,
			BaseURL: ds.BaseUrl,
			Status:  entity.DatastoreStatus(ds.Status),
		}
	}
	return datastores, nil
}

func (dr *DatastoreRepository) ChangeDatastoreStatus(
	ctx context.Context, id int64, status entity.DatastoreStatus,
) error {
	tx := psqlutil.TxFromCtx(ctx)
	q := dr.q
	if tx != nil {
		q = dr.q.WithTx(tx)
	}
	err := q.UpdateDatastoreStatus(ctx, query.UpdateDatastoreStatusParams{
		ID:     id,
		Status: query.DatastoreStatus(status),
	})
	if err != nil {
		return fmt.Errorf("failed to update datastore status: %w", err)
	}
	return nil
}
