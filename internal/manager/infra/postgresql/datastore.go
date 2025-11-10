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

type DatastoreRepository struct {
	q *query.Queries
}

func NewDatastoreRepository(db *sql.DB) *DatastoreRepository {
	return &DatastoreRepository{
		q: query.New(db),
	}
}

func (dr *DatastoreRepository) CreateDatastore(ctx context.Context, req *service.CreateDatastoreRequest) (int64, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := dr.q
	if tx != nil {
		q = dr.q.WithTx(tx)
	}
	id, err := q.InsertDatastore(ctx, req.BaseURL)
	if err != nil {
		return 0, fmt.Errorf("failed to insert datastore: %w", err)
	}
	return id, nil
}

func (dr *DatastoreRepository) GetDatastore(ctx context.Context, id int64) (*entity.Datastore, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := dr.q
	if tx != nil {
		q = dr.q.WithTx(tx)
	}
	ds, err := q.SelectDatastore(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("failed to select datastore: %w", err)
	}
	return &entity.Datastore{
		ID:      ds.ID,
		BaseURL: ds.BaseUrl,
	}, nil
}

func (dr *DatastoreRepository) GetDatastoreIDs(ctx context.Context) ([]int64, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := dr.q
	if tx != nil {
		q = dr.q.WithTx(tx)
	}
	dsIDs, err := q.SelectDatastoreIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to select datastore IDs: %w", err)
	}
	return dsIDs, nil
}
