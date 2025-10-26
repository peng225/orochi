package service

import (
	"context"

	"github.com/peng225/orochi/internal/manager/entity"
)

type CreateDatastoreRequest struct {
	BaseURL string
}

type DatastoreRepository interface {
	GetDatastore(ctx context.Context, id int64) (*entity.Datastore, error)
	CreateDatastore(ctx context.Context, req *CreateDatastoreRequest) (int64, error)
}
