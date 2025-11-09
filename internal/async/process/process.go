package process

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/gateway/api/client"
)

type Processor struct {
	period     time.Duration
	jobRepo    JobRepository
	bucketRepo BucketRepository
	gwClient   *client.Client
}

func NewProcessor(
	period time.Duration,
	jobRepo JobRepository,
	bucketRepo BucketRepository,
	gwClient *client.Client,
) *Processor {
	return &Processor{
		period:     period,
		jobRepo:    jobRepo,
		bucketRepo: bucketRepo,
		gwClient:   gwClient,
	}
}

func (p *Processor) Start(ctx context.Context) {
	ticker := time.NewTicker(p.period)
	defer ticker.Stop()

	slog.Info("Processor started.")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Processor stopped.")
			return
		case <-ticker.C:
			jobs, err := p.jobRepo.GetJobs(ctx, 4)
			if err != nil {
				slog.Error("Failed to get jobs.", "err", err)
				break
			}
			// FIXME: parallelize with some locks.
			for _, job := range jobs {
				err := p.processJob(ctx, job)
				if err != nil {
					slog.Error("Failed to process job.",
						"ID", job.ID, "Kind", job.Kind, "err", err)
				}
			}
		}
	}
}

func (p *Processor) processJob(ctx context.Context, job *entity.Job) error {
	slog.Info("Start to process job.", "id", job.ID, "kind", job.Kind)
	var finished bool
	var err error
	switch job.Kind {
	case entity.DeleteAllObjectsInBucket:
		finished, err = p.processDeleteAllObjectsInBucketJob(ctx, job)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected job kind: %s", job.Kind)
	}
	if finished {
		err := p.jobRepo.DeleteJob(ctx, job.ID)
		if err != nil {
			return fmt.Errorf("failed to delete job: %w", err)
		}
	}
	return nil
}

func (p *Processor) processDeleteAllObjectsInBucketJob(ctx context.Context, job *entity.Job) (bool, error) {
	var param entity.DeleteAllObjectsInBucketParam
	err := json.Unmarshal(job.Data, &param)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal job data: %w", err)
	}
	b, err := p.bucketRepo.GetBucket(ctx, param.BucketID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			slog.Info("Bucket not found.")
			return true, nil
		}
		return false, fmt.Errorf("failed to get bucket: %w", err)
	}

	listResp, err := p.gwClient.ListObjects(ctx, b.Name, nil)
	if err != nil {
		return false, fmt.Errorf("failed to list objects: %w", err)
	}
	defer listResp.Body.Close()
	if listResp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", listResp.StatusCode)
	}
	data, err := io.ReadAll(listResp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read body: %w", err)
	}
	objectNames := make([]string, 0)
	err = json.Unmarshal(data, &objectNames)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal object names: %w", err)
	}
	if len(objectNames) == 0 {
		slog.Info("No object found in the bucket. Job finished.", "bucketID", param.BucketID)
		err := p.bucketRepo.DeleteBucket(ctx, param.BucketID)
		if err != nil {
			return false, fmt.Errorf("failed to delete bucket: %w", err)
		}
		return true, nil
	}

	for _, objectName := range objectNames {
		delResp, err := p.gwClient.DeleteObject(ctx, b.Name, objectName)
		if err != nil {
			return false, fmt.Errorf("failed to delete object: %w", err)
		}
		if delResp.StatusCode != http.StatusNoContent {
			return false, fmt.Errorf("unexpected status code: %d", delResp.StatusCode)
		}
	}
	return false, nil
}
