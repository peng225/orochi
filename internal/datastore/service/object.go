package service

import (
	"io"
)

type ObjectService interface {
	GetObject(bucket, object string) (io.Reader, error)
	CreateObject(bucket, object string, data io.Reader) error
}

type FileObjectStore struct {
}

func NewFileObjectStore() *FileObjectStore {
	return &FileObjectStore{}
}

func (os *FileObjectStore) GetObject(bucket, object string) (io.Reader, error) {
	return nil, nil
}

func (os *FileObjectStore) CreateObject(bucket, object string, data io.Reader) error {
	return nil
}
