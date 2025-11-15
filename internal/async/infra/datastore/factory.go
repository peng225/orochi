package datastore

import (
	"github.com/peng225/orochi/internal/async/process"
	"github.com/peng225/orochi/internal/entity"
)

type ClientFactory struct{}

func NewClientFactory() *ClientFactory {
	return &ClientFactory{}
}

func (cf *ClientFactory) New(ds *entity.Datastore) process.DatastoreClient {
	return NewClient(ds.BaseURL)
}
