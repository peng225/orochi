package service

import (
	"context"

	"github.com/peng225/orochi/internal/manager/entity"
)

type DatastoreService struct {
	dsRepo DatastoreRepository
}

func NewDatastoreService(dsRepo DatastoreRepository) *DatastoreService {
	return &DatastoreService{
		dsRepo: dsRepo,
	}
}

func (dss *DatastoreService) GetDatastore(ctx context.Context, id int64) (*entity.Datastore, error) {
	return dss.dsRepo.GetDatastore(ctx, id)
}

func (dss *DatastoreService) CreateDatastore(ctx context.Context, baseURL string) (int64, error) {
	return dss.dsRepo.CreateDatastore(ctx, &CreateDatastoreRequest{
		BaseURL: baseURL,
	})
}
