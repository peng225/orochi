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
)

type LocationGroupRepository struct {
	db *sql.DB
	q  *query.Queries
}

func NewLocationGroupRepository() *LocationGroupRepository {
	db := psqlutil.InitDB()
	return &LocationGroupRepository{
		db: db,
		q:  query.New(db),
	}
}

func (lgr *LocationGroupRepository) Close() error {
	return lgr.db.Close()
}

func (lgr *LocationGroupRepository) CreateLocationGroup(
	ctx context.Context,
	req *service.CreateLocationGroupRequest,
) (int64, error) {
	id, err := lgr.q.InsertLocationGroup(ctx, query.InsertLocationGroupParams{
		CurrentDatastores: req.Datastores,
		DesiredDatastores: req.Datastores,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to insert location group: %w", err)
	}
	return id, nil
}

func (lgr *LocationGroupRepository) UpdateDesiredDatastores(
	ctx context.Context,
	id int64,
	desiredDatastores []int64,
) error {
	err := lgr.q.UpdateDesiredDatastores(ctx, query.UpdateDesiredDatastoresParams{
		ID:                id,
		DesiredDatastores: desiredDatastores,
	})
	if err != nil {
		return fmt.Errorf("failed to update desired datastores: %w", err)
	}
	return nil
}

func (lgr *LocationGroupRepository) GetLocationGroups(ctx context.Context) ([]*entity.LocationGroup, error) {
	lgs, err := lgr.q.SelectLocationGroups(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("failed to select location group: %w", err)
	}
	ret := make([]*entity.LocationGroup, 0, len(lgs))
	for _, lg := range lgs {
		ret = append(ret, &entity.LocationGroup{
			ID:                lg.ID,
			CurrentDatastores: lg.CurrentDatastores,
			DesiredDatastores: lg.DesiredDatastores,
		})
	}
	return ret, nil
}
