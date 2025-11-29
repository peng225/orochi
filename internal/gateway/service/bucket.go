package service

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/peng225/orochi/internal/entity"
)

type BucketService struct {
	mu      sync.RWMutex
	buckets map[string]*entity.Bucket
}

func NewBucketService() *BucketService {
	buckets := make(map[string]*entity.Bucket)
	return &BucketService{
		buckets: buckets,
	}
}

func (bs *BucketService) CreateBucket(name, ecConfig string) error {
	slog.Debug("BucketService::CreateBucket called.", "name", name, "ecConfig", ecConfig)
	bs.mu.Lock()
	defer bs.mu.Unlock()
	b, ok := bs.buckets[name]
	if ok {
		if b.ECConfig == ecConfig {
			return nil
		}
		return errors.Join(fmt.Errorf("bucket with the different EC config found: expected=%s, actual=%s",
			ecConfig, b.ECConfig), ErrInvalidParameter)
	}

	bs.buckets[name] = &entity.Bucket{
		Name:     name,
		ECConfig: ecConfig,
		Status:   entity.BucketStatusActive,
	}
	return nil
}

func (bs *BucketService) GetBucket(name string) (*entity.Bucket, error) {
	slog.Debug("BucketService::GetBucket called.", "name", name)
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	b, ok := bs.buckets[name]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

func (bs *BucketService) UpdateBucketStatus(name string, status entity.BucketStatus) error {
	slog.Debug("BucketService::UpdateBucketStatus called.",
		"name", name, "status", status)
	bs.mu.Lock()
	defer bs.mu.Unlock()
	b, ok := bs.buckets[name]
	if !ok {
		return ErrNotFound
	}

	b.Status = status
	return nil
}

func (bs *BucketService) DeleteBucket(name string) {
	slog.Debug("BucketService::DeleteBucket called.", "name", name)
	bs.mu.Lock()
	defer bs.mu.Unlock()

	delete(bs.buckets, name)
}
