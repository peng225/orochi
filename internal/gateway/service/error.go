package service

import "fmt"

var (
	ErrNotFound              error = fmt.Errorf("not found")
	ErrInvalidParameter      error = fmt.Errorf("invalid parameter")
	ErrLocationGroupNotFound error = fmt.Errorf("location group not found")
	ErrObjectNotActive       error = fmt.Errorf("object not active")
)
