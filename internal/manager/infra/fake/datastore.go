package fake

import (
	"context"
	"slices"
	"sync"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/manager/service"

	_ "github.com/lib/pq"
)

type FakeDatastoreRepository struct {
	mu         sync.Mutex
	nextID     int64
	datastores map[int64]*entity.Datastore
}

func NewFakeDatastoreRepository() *FakeDatastoreRepository {
	return &FakeDatastoreRepository{
		nextID:     1,
		datastores: make(map[int64]*entity.Datastore),
	}
}

func (dr *FakeDatastoreRepository) CreateDatastore(ctx context.Context, req *service.CreateDatastoreRequest) (int64, error) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	id := dr.nextID
	dr.datastores[dr.nextID] = &entity.Datastore{
		ID:      id,
		BaseURL: req.BaseURL,
	}
	dr.nextID++
	return id, nil
}

func (dr *FakeDatastoreRepository) GetDatastore(ctx context.Context, id int64) (*entity.Datastore, error) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	if _, ok := dr.datastores[id]; !ok {
		return nil, service.ErrNotFound
	}
	return dr.datastores[id], nil
}

func (dr *FakeDatastoreRepository) GetDatastoreByBaseURL(
	ctx context.Context, baseURL string,
) (*entity.Datastore, error) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	for _, ds := range dr.datastores {
		if ds.BaseURL == baseURL {
			return ds, nil
		}
	}
	return nil, service.ErrNotFound
}

func (dr *FakeDatastoreRepository) GetDatastores(
	ctx context.Context, req *service.GetDatastoresRequest,
) ([]*entity.Datastore, error) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	dss := make([]*entity.Datastore, 0, len(dr.datastores))
	for _, ds := range dr.datastores {
		if ds.ID >= req.StartFrom {
			dss = append(dss, ds)
		}
	}
	slices.SortFunc(dss, func(a, b *entity.Datastore) int {
		if a.ID < b.ID {
			return -1
		} else if a.ID == b.ID {
			return 0
		}
		return 1
	})
	return dss[:min(len(dss), req.Limit)], nil
}

func (dr *FakeDatastoreRepository) GetDatastoreIDs(ctx context.Context) ([]int64, error) {
	dr.mu.Lock()
	defer dr.mu.Unlock()
	dsIDs := make([]int64, 0, len(dr.datastores))
	for _, ds := range dr.datastores {
		dsIDs = append(dsIDs, ds.ID)
	}
	return dsIDs, nil
}
