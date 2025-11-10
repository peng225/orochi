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

type JobRepository struct {
	q *query.Queries
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{
		q: query.New(db),
	}
}

func (jr *JobRepository) GetJobs(ctx context.Context, limit int) ([]*entity.Job, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := jr.q
	if tx != nil {
		q = jr.q.WithTx(tx)
	}
	jobs, err := q.SelectJobs(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to select jobs: %w", err)
	}
	res := make([]*entity.Job, 0, len(jobs))
	for _, job := range jobs {
		res = append(res, &entity.Job{
			ID:   job.ID,
			Kind: job.Kind,
			Data: job.Data,
		})
	}
	return res, nil
}

func (jr *JobRepository) DeleteJob(ctx context.Context, id int64) error {
	tx := psqlutil.TxFromCtx(ctx)
	q := jr.q
	if tx != nil {
		q = jr.q.WithTx(tx)
	}
	err := q.DeleteJob(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}
	return nil
}
