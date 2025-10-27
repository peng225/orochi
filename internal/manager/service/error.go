package service

import "fmt"

var (
	ErrInvalidParameter  error = fmt.Errorf("invalid parameter")
	ErrDatastoreNotFound error = fmt.Errorf("datastore not found")
)
