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

type ECConfigRepository struct {
	q *query.Queries
}

func NewECConfigRepository(db *sql.DB) *ECConfigRepository {
	return &ECConfigRepository{
		q: query.New(db),
	}
}

func (ecr *ECConfigRepository) CreateECConfig(ctx context.Context, req *service.CreateECConfigRequest) (int64, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := ecr.q
	if tx != nil {
		q = ecr.q.WithTx(tx)
	}
	id, err := q.InsertECConfig(ctx, query.InsertECConfigParams{
		NumData:   int32(req.NumData),
		NumParity: int32(req.NumParity),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert EC config: %w", err)
	}
	return id, nil
}

func (ecr *ECConfigRepository) GetECConfigByNumbers(ctx context.Context, numData, numParity int) (*entity.ECConfig, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := ecr.q
	if tx != nil {
		q = ecr.q.WithTx(tx)
	}
	ecc, err := q.SelectECConfigByNumbers(ctx, query.SelectECConfigByNumbersParams{
		NumData:   int32(numData),
		NumParity: int32(numParity),
	})
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

func (ecr *ECConfigRepository) GetECConfigs(ctx context.Context) ([]*entity.ECConfig, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := ecr.q
	if tx != nil {
		q = ecr.q.WithTx(tx)
	}
	eccs, err := q.SelectECConfigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to select EC configs: %w", err)
	}
	ret := make([]*entity.ECConfig, 0, len(eccs))
	for _, ecc := range eccs {
		ret = append(ret, &entity.ECConfig{
			ID:        ecc.ID,
			NumData:   int(ecc.NumData),
			NumParity: int(ecc.NumParity),
		})
	}
	return ret, nil
}
