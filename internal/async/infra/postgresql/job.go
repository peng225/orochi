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

func (br *JobRepository) GetJobs(ctx context.Context, limit int) ([]*entity.Job, error) {
	jobs, err := br.q.SelectJobs(ctx, int32(limit))
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

func (br *JobRepository) DeleteJob(ctx context.Context, id int64) error {
	err := br.q.DeleteJob(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}
	return nil
}
