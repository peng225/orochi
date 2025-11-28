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
	lgService  *LocationGroupService
	bucketRepo BucketRepository
	jobRepo    JobRepository
	eccRepo    ECConfigRepository
}

func NewBucketService(
	tx Transaction,
	lgService *LocationGroupService,
	bucketRepo BucketRepository,
	jobRepo JobRepository,
	eccRepo ECConfigRepository,
) *BucketService {
	return &BucketService{
		tx:         tx,
		lgService:  lgService,
		bucketRepo: bucketRepo,
		jobRepo:    jobRepo,
		eccRepo:    eccRepo,
	}
}

func (bs *BucketService) CreateBucket(ctx context.Context, name, ecConfigStr string) (int64, error) {
	if !validBucketName.MatchString(name) {
		return 0, errors.Join(fmt.Errorf("invalid bucket name: %s", name), ErrInvalidParameter)
	}
	numData, numParity, err := entity.ParseECConfig(ecConfigStr)
	if err != nil {
		return 0, errors.Join(fmt.Errorf("failed to parse EC config: %w", err), ErrInvalidParameter)
	}

	var id int64
	err = bs.tx.Do(ctx, func(ctx context.Context) error {
		ecConfig, err := bs.eccRepo.GetECConfigByNumbers(ctx, numData, numParity)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return fmt.Errorf("failed to get EC config by numbers: %w", err)
			}
			_, err := bs.eccRepo.CreateECConfig(ctx, &CreateECConfigRequest{
				NumData:   numData,
				NumParity: numParity,
			})
			if err != nil {
				return fmt.Errorf("failed to create EC config: %w", err)
			}
			ecConfig, err = bs.eccRepo.GetECConfigByNumbers(ctx, numData, numParity)
			if err != nil {
				return fmt.Errorf("failed to get EC config by numbers: %w", err)
			}
		}

		bucket, err := bs.bucketRepo.GetBucketByName(ctx, name)
		if err == nil {
			if bucket.ECConfigID != ecConfig.ID {
				return errors.Join(fmt.Errorf("bucket with different EC config ID found: expected=%d, actual=%d",
					ecConfig.ID, bucket.ECConfigID), ErrConflict)
			} else if bucket.Status != entity.BucketStatusActive {
				return errors.Join(fmt.Errorf("bucket with different status found: %s", bucket.Status), ErrConflict)
			}
			id = bucket.ID
			return nil
		} else if !errors.Is(err, ErrNotFound) {
			return fmt.Errorf("failed to get bucket by name: %w", err)
		}

		id, err = bs.bucketRepo.CreateBucket(ctx, &CreateBucketRequest{
			Name:       name,
			ECConfigID: ecConfig.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}

		err = bs.lgService.ReconstructLocationGroups(ctx, ecConfig)
		if err != nil {
			return fmt.Errorf("failed to reconstruct location groups: %w", err)
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("transaction failed: %w", err)
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
		err := bs.bucketRepo.ChangeBucketStatus(ctx, id, entity.BucketStatusDeleting)
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
