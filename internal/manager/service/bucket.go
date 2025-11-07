package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/job"
)

type BucketService struct {
	bucketRepo BucketRepository
	jobRepo    JobRepository
}

func NewBucketService(bucketRepo BucketRepository, jobRepo JobRepository) *BucketService {
	return &BucketService{
		bucketRepo: bucketRepo,
		jobRepo:    jobRepo,
	}
}

func (bs *BucketService) CreateBucket(ctx context.Context, name string) (int64, error) {
	if !isValidBucketName(name) {
		return 0, ErrInvalidParameter
	}
	buckets, err := bs.bucketRepo.GetBucketsByName(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("failed to get buckets by name: %w", err)
	}
	bucketCount := 0
	var bucket *entity.Bucket
	for _, b := range buckets {
		if b.Status == "deleted" {
			continue
		}
		bucketCount++
		bucket = b
	}
	switch bucketCount {
	case 0:
		// Bucket not found. Should create a new bucket.
	case 1:
		return bucket.ID, nil
	default:
		return 0, fmt.Errorf("unexpected number of buckets found: %d", bucketCount)
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

func (bs *BucketService) GetBucket(ctx context.Context, id int64) (*entity.Bucket, error) {
	bucket, err := bs.bucketRepo.GetBucket(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}
	return bucket, nil
}

func (bs *BucketService) DeleteBucket(ctx context.Context, id int64) error {
	err := bs.bucketRepo.ChangeBucketStatus(ctx, id, "deleted")
	if err != nil {
		return fmt.Errorf("failed to change bucket status: %w", err)
	}
	param := job.DeleteAllObjectsInBUcketParam{
		BucketID: id,
	}
	data, err := json.Marshal(&param)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	jobID, err := bs.jobRepo.CreateJob(ctx, &CreateJobRequest{
		Name: "DeleteAllObjectsInBucket",
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("failed create job: %w", err)
	}
	slog.Info("Created a job to delete all objects in a bucket.",
		"bucketID", id, "jobID", jobID)
	return nil
}
