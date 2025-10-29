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

type DatastoreRepository struct {
	db *sql.DB
	q  *query.Queries
}

func NewDatastoreRepository() *DatastoreRepository {
	db := psqlutil.InitDB()
	return &DatastoreRepository{
		db: db,
		q:  query.New(db),
	}
}

func (dr *DatastoreRepository) Close() error {
	return dr.db.Close()
}

func (dr *DatastoreRepository) GetDatastores(ctx context.Context) ([]*entity.Datastore, error) {
	dss, err := dr.q.SelectDatastores(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to select datastores: %w", err)
	}
	datastores := make([]*entity.Datastore, len(dss))
	for i, ds := range dss {
		datastores[i] = &entity.Datastore{
			ID:      ds.ID,
			BaseURL: ds.BaseUrl,
		}
	}
	return datastores, nil
}
