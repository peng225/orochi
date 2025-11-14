package service

import "fmt"

var (
	ErrObjectNotFound   error = fmt.Errorf("object not found")
	ErrInvalidParameter error = fmt.Errorf("invalid parameter")
	ErrTooLargeObject   error = fmt.Errorf("too large object")
)
