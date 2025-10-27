package service

import (
	"context"
	"regexp"

	"github.com/peng225/orochi/internal/entity"
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
	if !isValidURL(baseURL) {
		return 0, ErrInvalidParameter
	}
	return dss.dsRepo.CreateDatastore(ctx, &CreateDatastoreRequest{
		BaseURL: baseURL,
	})
}

func isValidURL(s string) bool {
	re := regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+(:[0-9]+)?$`)
	return re.MatchString(s)
}
