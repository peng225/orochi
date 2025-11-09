package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/peng225/orochi/internal/manager/infra/postgresql/sqlc/query"
	"github.com/peng225/orochi/internal/manager/service"
	"github.com/peng225/orochi/internal/pkg/psqlutil"

	_ "github.com/lib/pq"
)

type JobRepository struct {
	db *sql.DB
	q  *query.Queries
}

func NewJobRepository() *JobRepository {
	db := psqlutil.InitDB()
	return &JobRepository{
		db: db,
		q:  query.New(db),
	}
}

func (br *JobRepository) Close() error {
	return br.db.Close()
}

func (br *JobRepository) CreateJob(ctx context.Context, req *service.CreateJobRequest) (int64, error) {
	id, err := br.q.InsertJob(ctx, query.InsertJobParams{
		Kind: req.Kind,
		Data: req.Data,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert job: %w", err)
	}
	return id, nil
}
