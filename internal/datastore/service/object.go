package service

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	// FIXME: parametrize.
	dataRoot    = "/data"
	headerSize  = 8
	maxDataSize = 4*1024*1024*1024 - headerSize
)

type ObjectService interface {
	GetObject(bucket, object string) ([]byte, error)
	CreateObject(bucket, object string, data io.Reader) error
}

type FileObjectStore struct {
}

func NewFileObjectStore() *FileObjectStore {
	return &FileObjectStore{}
}

func (fos *FileObjectStore) GetObject(bucket, object string) ([]byte, error) {
	path := filepath.Join(dataRoot, bucket, object)

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrObjectNotFound
		}
		return nil, fmt.Errorf("failed to open object: %w", err)
	}
	defer f.Close()

	var size uint64
	if err := binary.Read(f, binary.LittleEndian, &size); err != nil {
		return nil, fmt.Errorf("failed to read size header: %w", err)
	}
	if maxDataSize < size {
		return nil, fmt.Errorf("invalid size header: %d", size)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	if uint64(len(data)) != size {
		return nil, fmt.Errorf("data size mismatch: expected %d bytes, got %d bytes", size, len(data))
	}

	return data, nil
}

func (fos *FileObjectStore) CreateObject(bucket, object string, data io.Reader) error {
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

	sizePlaceholder := make([]byte, headerSize)
	if _, err := f.Write(sizePlaceholder); err != nil {
		return fmt.Errorf("failed to write placeholder: %w", err)
	}

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

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek: %w", err)
	}
	sizeBuf := make([]byte, headerSize)
	binary.LittleEndian.PutUint64(sizeBuf, uint64(written))
	if _, err := f.Write(sizeBuf); err != nil {
		return fmt.Errorf("failed to write size header: %w", err)
	}

	if err := os.Rename(tmpPath, objectPath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}
