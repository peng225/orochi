package fake

import (
	"context"
	"sync"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/manager/service"
)

type FakeLocationGroupRepository struct {
	mu             sync.Mutex
	nextID         int64
	locationGroups map[int64]*entity.LocationGroup
}

func NewFakeLocationGroupRepository() *FakeLocationGroupRepository {
	return &FakeLocationGroupRepository{
		nextID:         1,
		locationGroups: make(map[int64]*entity.LocationGroup),
	}
}

func (lgr *FakeLocationGroupRepository) CreateLocationGroup(
	ctx context.Context,
	req *service.CreateLocationGroupRequest,
) (int64, error) {
	lgr.mu.Lock()
	defer lgr.mu.Unlock()
	id := lgr.nextID
	lgr.locationGroups[lgr.nextID] = &entity.LocationGroup{
		ID:         id,
		Datastores: req.Datastores,
		ECConfigID: req.ECConfigID,
	}
	lgr.nextID++
	return id, nil
}

func (lgr *FakeLocationGroupRepository) UpdateLocationGroupStatus(
	ctx context.Context,
	id int64,
	status entity.LocationGroupStatus,
) error {
	lgr.mu.Lock()
	defer lgr.mu.Unlock()
	if _, ok := lgr.locationGroups[id]; !ok {
		return service.ErrNotFound
	}
	lgr.locationGroups[id].Status = status
	return nil
}

func (lgr *FakeLocationGroupRepository) GetLocationGroupsByECConfigID(
	ctx context.Context, ecConfigID int64,
) ([]*entity.LocationGroup, error) {
	lgr.mu.Lock()
	defer lgr.mu.Unlock()
	ret := make([]*entity.LocationGroup, 0)
	for _, lg := range lgr.locationGroups {
		if lg.ECConfigID == ecConfigID {
			ret = append(ret, lg)
		}
	}
	return ret, nil
}
