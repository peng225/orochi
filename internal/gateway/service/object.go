package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	randv2 "math/rand/v2"
	"slices"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/pkg/ec"
	"golang.org/x/sync/errgroup"
)

const (
	minECChunkSizeInByte = 4 * 1024
)

type ObjectService struct {
	chunkRepos map[int64]ChunkRepository
	crFactory  ChunkRepositoryFactory
	dsRepo     DatastoreRepository
	omRepo     ObjectMetadataRepository
	bucketRepo BucketRepository
	lgRepo     LocationGroupRepository
}

func NewObjectStore(
	chunkRepos map[int64]ChunkRepository,
	crFactory ChunkRepositoryFactory,
	dsRepo DatastoreRepository,
	omRepo ObjectMetadataRepository,
	bucketRepo BucketRepository,
	lgRepo LocationGroupRepository,
) *ObjectService {
	if chunkRepos == nil {
		chunkRepos = make(map[int64]ChunkRepository)
	}
	return &ObjectService{
		chunkRepos: chunkRepos,
		crFactory:  crFactory,
		dsRepo:     dsRepo,
		omRepo:     omRepo,
		bucketRepo: bucketRepo,
		lgRepo:     lgRepo,
	}
}

func (osvc *ObjectService) Refresh(ctx context.Context) error {
	dss, err := osvc.dsRepo.GetDatastores(ctx)
	if err != nil {
		return fmt.Errorf("failed to get datastores: %w", err)
	}
	for _, ds := range dss {
		osvc.chunkRepos[ds.ID] = osvc.crFactory.New(ds)
	}
	return nil
}

func (osvc *ObjectService) CreateObject(ctx context.Context, name, bucket string, r io.Reader) error {
	slog.Debug("ObjectService::CreateObject called.", "name", name, "bucket", bucket)
	// FIXME: Should avoid per request refresh for performance.
	err := osvc.Refresh(ctx)
	if err != nil {
		return fmt.Errorf("failed to refresh: %w", err)
	}
	om, err := osvc.getObjectMetadataByName(ctx, name, bucket)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return fmt.Errorf("failed to get object metadata by name: %w", err)
		}
		_, err := osvc.createObjectMetadata(ctx, name, bucket)
		if err != nil {
			return fmt.Errorf("failed to create object metadata: %w", err)
		}
		om, err = osvc.getObjectMetadataByName(ctx, name, bucket)
		if err != nil {
			return fmt.Errorf("failed to get object metadata: %w", err)
		}
	}
	lg, err := osvc.lgRepo.GetLocationGroup(ctx, om.LocationGroupID)
	if err != nil {
		return fmt.Errorf("failed to get location group: %w", err)
	}
	if !slices.Equal(lg.CurrentDatastores, lg.DesiredDatastores) {
		// FIXME: double write.
		return fmt.Errorf("unsupported behavior")
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}
	// Should remove the assumption of 2D1P.
	m := ec.NewManager(2, 1, minECChunkSizeInByte)
	codes, err := m.Encode(data)
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}
	for i, ds := range lg.CurrentDatastores {
		err = osvc.chunkRepos[ds].CreateObject(ctx, bucket, name, bytes.NewBuffer(codes[i]))
		if err != nil {
			return fmt.Errorf("CreateObject failed: %w", err)
		}
	}
	return nil
}

func (osvc *ObjectService) GetObject(ctx context.Context, name, bucket string) ([]byte, error) {
	slog.Debug("ObjectService::GetObject called.", "name", name, "bucket", bucket)
	om, err := osvc.getObjectMetadataByName(ctx, name, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get object metadata: %w", err)
	}
	lg, err := osvc.lgRepo.GetLocationGroup(ctx, om.LocationGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get location group: %w", err)
	}

	eg := new(errgroup.Group)
	// FIXME: remove the 2D1P assumption.
	// FIXME: should decode even if some datastores are down.
	codes := make([][]byte, 3)
	for i, ds := range lg.CurrentDatastores[:2] {
		i := i
		ds := ds
		eg.Go(func() error {
			rc, err := osvc.chunkRepos[ds].GetObject(ctx, bucket, name)
			if err != nil {
				return fmt.Errorf("GetObject failed: %w", err)
			}
			defer rc.Close()
			data, err := io.ReadAll(rc)
			if err != nil {
				return fmt.Errorf("failed to read data: %w", err)
			}
			codes[i] = data
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get object chunk: %w", err)
	}

	m := ec.NewManager(2, 1, minECChunkSizeInByte)
	data, err := m.Decode(codes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}
	return data, nil
}

func (osvc *ObjectService) createObjectMetadata(
	ctx context.Context,
	name string,
	bucketName string,
) (int64, error) {
	bucket, err := osvc.bucketRepo.GetBucketByName(ctx, bucketName)
	if err != nil {
		return 0, fmt.Errorf("failed to get bucket by name: %w", err)
	}
	lgs, err := osvc.lgRepo.GetLocationGroups(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get location groups: %w", err)
	}
	if len(lgs) == 0 {
		return 0, ErrLocationGroupNotFound
	}
	lg := lgs[randv2.IntN(len(lgs))]

	id, err := osvc.omRepo.CreateObjectMetadata(ctx, &CreateObjectMetadataRequest{
		Name:            name,
		BucketID:        bucket.ID,
		LocationGroupID: lg.ID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get object metadata: %w", err)
	}

	return id, nil
}

func (osvc *ObjectService) getObjectMetadataByName(
	ctx context.Context,
	name, bucketName string,
) (*entity.ObjectMetadata, error) {
	bucket, err := osvc.bucketRepo.GetBucketByName(ctx, bucketName)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrInvalidParameter
		}
		return nil, fmt.Errorf("failed to get bucket by name: %w", err)
	}

	om, err := osvc.omRepo.GetObjectMetadataByName(ctx, name, bucket.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object metadata by name: %w", err)
	}
	return om, nil
}

func (osvc *ObjectService) DeleteObject(ctx context.Context, name, bucket string, r io.Reader) error {
	slog.Debug("ObjectService::DeleteObject called.", "name", name, "bucket", bucket)
	om, err := osvc.getObjectMetadataByName(ctx, name, bucket)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed to get object metadata by name: %w", err)
	}
	lg, err := osvc.lgRepo.GetLocationGroup(ctx, om.LocationGroupID)
	if err != nil {
		return fmt.Errorf("failed to get location group: %w", err)
	}
	if !slices.Equal(lg.CurrentDatastores, lg.DesiredDatastores) {
		// FIXME: double delete.
		return fmt.Errorf("unsupported behavior")
	}
	for _, ds := range lg.CurrentDatastores {
		err = osvc.chunkRepos[ds].DeleteObject(ctx, bucket, name)
		if err != nil {
			return fmt.Errorf("DeleteObject failed: %w", err)
		}
	}
	err = osvc.omRepo.DeleteObjectMetadata(ctx, om.ID)
	if err != nil {
		return fmt.Errorf("failed to delete object metadata: %w", err)
	}
	return nil
}
