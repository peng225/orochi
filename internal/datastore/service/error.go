package service

import "fmt"

var (
	ErrObjectNotFound error = fmt.Errorf("object not found")
)
