package service

import "fmt"

var (
	ErrDatastoreNotFound error = fmt.Errorf("datastore not found")
)
