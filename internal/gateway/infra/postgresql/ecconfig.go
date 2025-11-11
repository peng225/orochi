package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/gateway/infra/postgresql/sqlc/query"
	"github.com/peng225/orochi/internal/gateway/service"
	"github.com/peng225/orochi/internal/pkg/psqlutil"

	_ "github.com/lib/pq"
)

type ECConfigRepository struct {
	q *query.Queries
}

func NewECConfigRepository(db *sql.DB) *ECConfigRepository {
	return &ECConfigRepository{
		q: query.New(db),
	}
}

func (ecr *ECConfigRepository) GetECConfig(ctx context.Context, id int64) (*entity.ECConfig, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := ecr.q
	if tx != nil {
		q = ecr.q.WithTx(tx)
	}
	ecc, err := q.SelectECConfig(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("failed to select EC config: %w", err)
	}
	return &entity.ECConfig{
		ID:        ecc.ID,
		NumData:   int(ecc.NumData),
		NumParity: int(ecc.NumParity),
	}, nil
}
