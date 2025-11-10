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
	q *query.Queries
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{
		q: query.New(db),
	}
}

func (jr *JobRepository) CreateJob(ctx context.Context, req *service.CreateJobRequest) (int64, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := jr.q
	if tx != nil {
		q = jr.q.WithTx(tx)
	}
	id, err := q.InsertJob(ctx, query.InsertJobParams{
		Kind: req.Kind,
		Data: req.Data,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert job: %w", err)
	}
	return id, nil
}
