package process

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/peng225/orochi/internal/entity"
	"golang.org/x/sync/errgroup"
)

// FIXME: what to do if gateway is deleted or newly created?
type Processor struct {
	period        time.Duration
	gwClientIndex int
	tx            Transaction
	jobRepo       JobRepository
	bucketRepo    BucketRepository
	dsRepo        DatastoreRepository
	gwClients     []GatewayClient
	dscFactory    DatastoreClientFactory
	dsDownCount   map[int64]int
}

func NewProcessor(
	period time.Duration,
	tx Transaction,
	jobRepo JobRepository,
	bucketRepo BucketRepository,
	dsRepo DatastoreRepository,
	gwClients []GatewayClient,
	dscFactory DatastoreClientFactory,
) *Processor {
	return &Processor{
		period:        period,
		gwClientIndex: 0,
		tx:            tx,
		jobRepo:       jobRepo,
		bucketRepo:    bucketRepo,
		dsRepo:        dsRepo,
		gwClients:     gwClients,
		dscFactory:    dscFactory,
		dsDownCount:   make(map[int64]int),
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
			err = p.checkDatastoreHealthStatus(ctx)
			if err != nil {
				slog.Error("Failed to check datastore health status.", "err", err)
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

func (p *Processor) selectGWClient() GatewayClient {
	defer func() {
		p.gwClientIndex = (p.gwClientIndex + 1) % len(p.gwClients)
	}()
	return p.gwClients[p.gwClientIndex]
}

func (p *Processor) processDeleteAllObjectsInBucketJob(ctx context.Context, job *entity.Job) (bool, error) {
	var param entity.DeleteAllObjectsInBucketParam
	err := json.Unmarshal(job.Data, &param)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal job data: %w", err)
	}

	finished := false
	err = p.tx.Do(ctx, func(ctx context.Context) error {
		b, err := p.bucketRepo.GetBucket(ctx, param.BucketID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				slog.Info("Bucket not found.")
				finished = true
				return nil
			}
			return fmt.Errorf("failed to get bucket: %w", err)
		}

		objectNames, err := p.selectGWClient().ListObjectNames(ctx, b.Name)
		if err != nil {
			return fmt.Errorf("failed to list objects: %w", err)
		}
		if len(objectNames) == 0 {
			slog.Info("No object found in the bucket. Job finished.", "bucketID", param.BucketID)
			err := p.bucketRepo.DeleteBucket(ctx, param.BucketID)
			if err != nil {
				return fmt.Errorf("failed to delete bucket: %w", err)
			}
			finished = true
			return nil
		}

		for _, objectName := range objectNames {
			err := p.selectGWClient().DeleteObject(ctx, b.Name, objectName)
			if err != nil {
				return fmt.Errorf("failed to delete object: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("transaction failed: %w", err)
	}
	return finished, nil
}

func (p *Processor) checkDatastoreHealthStatus(ctx context.Context) error {
	var dss []*entity.Datastore
	err := p.tx.Do(ctx, func(ctx context.Context) error {
		var err error
		dss, err = p.dsRepo.GetDatastores(ctx)
		if err != nil {
			return fmt.Errorf("failed to get datastores: %w", err)
		}

		eg := new(errgroup.Group)
		for _, ds := range dss {
			eg.Go(func() error {
				dsClient := p.dscFactory.New(ds)
				err := dsClient.CheckHealthStatus(ctx)
				if err != nil {
					if ds.Status == entity.DatastoreStatusDown {
						return nil
					}
					p.dsDownCount[ds.ID]++
					if p.dsDownCount[ds.ID] >= 2 {
						err := p.dsRepo.ChangeDatastoreStatus(ctx, ds.ID, entity.DatastoreStatusDown)
						if err != nil {
							return fmt.Errorf("failed to update datastore status: %w", err)
						}
						slog.Info("Datastore outage detected.", "id", ds.ID)
					}
				}
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return fmt.Errorf("failed to check health status of datastores: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	shouldBeDeletedIDs := make([]int64, 0)
OUTER:
	for dsID := range p.dsDownCount {
		for _, ds := range dss {
			if ds.ID == dsID {
				continue OUTER
			}
		}
		shouldBeDeletedIDs = append(shouldBeDeletedIDs, dsID)
	}
	for _, id := range shouldBeDeletedIDs {
		delete(p.dsDownCount, id)
	}

	return nil
}
