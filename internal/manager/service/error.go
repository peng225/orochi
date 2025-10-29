package service

import "fmt"

var (
	ErrInvalidParameter error = fmt.Errorf("invalid parameter")
	ErrNotFound         error = fmt.Errorf("not found")
)
