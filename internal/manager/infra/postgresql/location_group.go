package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/manager/infra/postgresql/sqlc/query"
	"github.com/peng225/orochi/internal/manager/service"
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

func (lgr *LocationGroupRepository) CreateLocationGroup(
	ctx context.Context,
	req *service.CreateLocationGroupRequest,
) (int64, error) {
	tx := psqlutil.TxFromCtx(ctx)
	q := lgr.q
	if tx != nil {
		q = lgr.q.WithTx(tx)
	}
	id, err := q.InsertLocationGroup(ctx, query.InsertLocationGroupParams{
		Datastores: req.Datastores,
		EcConfigID: req.ECConfigID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert location group: %w", err)
	}
	return id, nil
}

func (lgr *LocationGroupRepository) UpdateLocationGroupStatus(
	ctx context.Context, id int64, status entity.LocationGroupStatus,
) error {
	tx := psqlutil.TxFromCtx(ctx)
	q := lgr.q
	if tx != nil {
		q = lgr.q.WithTx(tx)
	}
	err := q.UpdateLocationGroupStatus(ctx, query.UpdateLocationGroupStatusParams{
		ID:     id,
		Status: query.LocationGroupStatus(status),
	})
	if err != nil {
		return fmt.Errorf("failed to update location group status: %w", err)
	}
	return nil
}

// FIXME: should limit the result number. It can be very large.
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
			ID:         lg.ID,
			Datastores: lg.Datastores,
			ECConfigID: lg.EcConfigID,
			Status:     entity.LocationGroupStatus(lg.Status),
		})
	}
	return ret, nil
}
