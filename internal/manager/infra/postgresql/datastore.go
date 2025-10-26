package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/manager/infra/postgresql/sqlc/query"
	"github.com/peng225/orochi/internal/manager/service"

	_ "github.com/lib/pq"
)

type PSQLDatastoreRepository struct {
	db *sql.DB
	q  *query.Queries
}

func NewPSQLDatastoreRepository() *PSQLDatastoreRepository {
	dsn := os.Getenv("DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return &PSQLDatastoreRepository{
		db: db,
		q:  query.New(db),
	}
}

func (pdr *PSQLDatastoreRepository) Close() error {
	return pdr.db.Close()
}

func (pdr *PSQLDatastoreRepository) GetDatastore(ctx context.Context, id int64) (*entity.Datastore, error) {
	ds, err := pdr.q.SelectDatastore(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrDatastoreNotFound
		}
		return nil, fmt.Errorf("failed to select datastore: %w", err)
	}
	return &entity.Datastore{
		ID:      ds.ID,
		BaseURL: ds.BaseUrl,
	}, nil
}

func (pdr *PSQLDatastoreRepository) CreateDatastore(ctx context.Context, req *service.CreateDatastoreRequest) (int64, error) {
	id, err := pdr.q.InsertDatastore(ctx, req.BaseURL)
	if err != nil {
		return 0, fmt.Errorf("failed to insert datastore: %w", err)
	}
	return id, nil
}
