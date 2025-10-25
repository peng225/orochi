package service

import "fmt"

var (
	ErrBucketNotFound error = fmt.Errorf("bucket not found")
	ErrObjectNotFound error = fmt.Errorf("object not found")
)
