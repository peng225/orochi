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

func (lgr *LocationGroupRepository) GetLocationGroup(
	ctx context.Context, id int64,
) (*entity.LocationGroup, error) {
	lg, err := lgr.q.SelectLocationGroup(ctx, id)
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
	}, nil
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
