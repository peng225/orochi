package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/peng225/orochi/internal/entity"
)

var (
	validBucketName = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
)

type BucketService struct {
	tx         Transaction
	bucketRepo BucketRepository
	jobRepo    JobRepository
}

func NewBucketService(tx Transaction, bucketRepo BucketRepository, jobRepo JobRepository) *BucketService {
	return &BucketService{
		tx:         tx,
		bucketRepo: bucketRepo,
		jobRepo:    jobRepo,
	}
}

func (bs *BucketService) CreateBucket(ctx context.Context, name string) (int64, error) {
	if !validBucketName.MatchString(name) {
		return 0, errors.Join(fmt.Errorf("invalid bucket name: %s", name), ErrInvalidParameter)
	}
	id, err := bs.bucketRepo.CreateBucket(ctx, &CreateBucketRequest{
		Name: name,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create bucket: %w", err)
	}
	return id, nil
}

func (bs *BucketService) GetBucket(ctx context.Context, id int64) (*entity.Bucket, error) {
	bucket, err := bs.bucketRepo.GetBucket(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}
	return bucket, nil
}

func (bs *BucketService) DeleteBucket(ctx context.Context, id int64) error {
	var jobID int64
	err := bs.tx.Do(ctx, func(ctx context.Context) error {
		err := bs.bucketRepo.ChangeBucketStatus(ctx, id, "deleted")
		if err != nil {
			return fmt.Errorf("failed to change bucket status: %w", err)
		}
		param := entity.DeleteAllObjectsInBucketParam{
			BucketID: id,
		}
		data, err := json.Marshal(&param)
		if err != nil {
			return fmt.Errorf("failed to marshal json: %w", err)
		}
		jobID, err = bs.jobRepo.CreateJob(ctx, &CreateJobRequest{
			Kind: entity.DeleteAllObjectsInBucket,
			Data: data,
		})
		if err != nil {
			return fmt.Errorf("failed create job: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	slog.Info("Created a job to delete all objects in a bucket.",
		"bucketID", id, "jobID", jobID)
	return nil
}
