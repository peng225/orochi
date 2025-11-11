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
)

type LocationGroupRepository struct {
	q *query.Queries
}

func NewLocationGroupRepository(db *sql.DB) *LocationGroupRepository {
	return &LocationGroupRepository{
		q: query.New(db),
	}
}

func (lgr *LocationGroupRepository) GetLocationGroup(
	ctx context.Context, id int64,
) (*entity.LocationGroup, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := lgr.q
	if tx != nil {
		q = lgr.q.WithTx(tx)
	}
	lg, err := q.SelectLocationGroup(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("failed to select location group: %w", err)
	}
	return &entity.LocationGroup{
		ID:                lg.ID,
		CurrentDatastores: lg.CurrentDatastores,
		DesiredDatastores: lg.DesiredDatastores,
		ECConfigID:        lg.EcConfigID,
	}, nil
}

func (lgr *LocationGroupRepository) GetLocationGroupsByECConfigID(
	ctx context.Context, ecConfigID int64,
) ([]*entity.LocationGroup, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := lgr.q
	if tx != nil {
		q = lgr.q.WithTx(tx)
	}
	lgs, err := q.SelectLocationGroupsByECConfigID(ctx, ecConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to select location group: %w", err)
	}
	ret := make([]*entity.LocationGroup, 0, len(lgs))
	for _, lg := range lgs {
		ret = append(ret, &entity.LocationGroup{
			ID:                lg.ID,
			CurrentDatastores: lg.CurrentDatastores,
			DesiredDatastores: lg.DesiredDatastores,
			ECConfigID:        lg.EcConfigID,
		})
	}
	return ret, nil
}
