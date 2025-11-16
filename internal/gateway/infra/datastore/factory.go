package datastore

import (
	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/gateway/service"
)

type ClientFactory struct{}

func NewClientFactory() *ClientFactory {
	return &ClientFactory{}
}

func (cf *ClientFactory) New(ds *entity.Datastore) service.DatastoreClient {
	return NewClient(ds.BaseURL)
}
