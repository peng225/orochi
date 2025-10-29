package service

import (
	"context"
	"fmt"
	"regexp"
)

type BucketService struct {
	bucketRepo BucketRepository
}

func NewBucketService(bucketRepo BucketRepository) *BucketService {
	return &BucketService{
		bucketRepo: bucketRepo,
	}
}

func (bs *BucketService) CreateBucket(ctx context.Context, name string) (int64, error) {
	if !isValidBucketName(name) {
		return 0, ErrInvalidParameter
	}
	id, err := bs.bucketRepo.CreateBucket(ctx, &CreateBucketRequest{
		Name: name,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create bucket: %w", err)
	}
	return id, nil
}

func isValidBucketName(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	return re.MatchString(s)
}
