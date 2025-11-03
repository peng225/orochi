package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// FIXME: parametrize.
	dataRoot    = "/data"
	maxDataSize = 4 * 1024 * 1024 * 1024
)

type ObjectService struct {
}

func NewObjectStore() *ObjectService {
	return &ObjectService{}
}

func (osvc *ObjectService) CreateObject(bucket, object string, data io.Reader) error {
	err := checkBucketFormat(bucket)
	if err != nil {
		return err
	}
	err = checkObjectFormat(object)
	if err != nil {
		return err
	}

	bucketPath := filepath.Join(dataRoot, bucket)
	objectPath := filepath.Join(bucketPath, object)

	if err := os.MkdirAll(bucketPath, 0755); err != nil {
		return fmt.Errorf("failed to create bucket dir: %w", err)
	}

	if err := os.Remove(objectPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to remove existing file: %w", err)
	}

	tmpPath := objectPath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer f.Close()

	var written int64
	buf := make([]byte, maxDataSize)
	for {
		n, err := data.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				return fmt.Errorf("failed to write data: %w", werr)
			}
			written += int64(n)
		}
		if written > maxDataSize {
			return fmt.Errorf("received too large data")
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read data: %w", err)
		}
	}
	// FIXME: synchronous fsync call may lead to the performance degradation.
	if err := f.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	if err := os.Rename(tmpPath, objectPath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}

func (osvc *ObjectService) GetObject(bucket, object string) ([]byte, error) {
	err := checkBucketFormat(bucket)
	if err != nil {
		return nil, err
	}
	err = checkObjectFormat(object)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dataRoot, bucket, object)

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrObjectNotFound
		}
		return nil, fmt.Errorf("failed to open object: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	return data, nil
}

func (osvc *ObjectService) DeleteObject(bucket, object string) error {
	err := checkBucketFormat(bucket)
	if err != nil {
		return err
	}
	err = checkObjectFormat(object)
	if err != nil {
		return err
	}
	path := filepath.Join(dataRoot, bucket, object)
	err = os.Remove(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to remove file: %w", err)
	}

	// FIXME: object may include/the/intermediate/directories.
	// In that case, those directory may be empty, and should be deleted.

	return nil
}

func checkBucketFormat(bucket string) error {
	if strings.Contains(bucket, "/") {
		return fmt.Errorf("invalid bucket format")
	}
	return nil
}

func checkObjectFormat(object string) error {
	items := strings.Split(object, "/")
	for _, item := range items {
		if len(item) == 0 {
			return fmt.Errorf("invalid object format")
		}
	}
	return nil
}
