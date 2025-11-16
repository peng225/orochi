package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	randv2 "math/rand/v2"
	"path/filepath"
	"regexp"
	"slices"
	"sync/atomic"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/pkg/ec"
	"golang.org/x/sync/errgroup"
)

const (
	minECChunkSizeInByte = 4 * 1024
)

var (
	validObjectName = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
)

type ObjectService struct {
	tx         Transaction
	dsClients  map[int64]DatastoreClient
	dscFactory DatastoreClientFactory
	dsRepo     DatastoreRepository
	omRepo     ObjectMetadataRepository
	bucketRepo BucketRepository
	lgRepo     LocationGroupRepository
	eccRepo    ECConfigRepository
}

func NewObjectStore(
	tx Transaction,
	dsClients map[int64]DatastoreClient,
	dscFactory DatastoreClientFactory,
	dsRepo DatastoreRepository,
	omRepo ObjectMetadataRepository,
	bucketRepo BucketRepository,
	lgRepo LocationGroupRepository,
	eccRepo ECConfigRepository,
) *ObjectService {
	if dsClients == nil {
		dsClients = make(map[int64]DatastoreClient)
	}
	return &ObjectService{
		tx:         tx,
		dsClients:  dsClients,
		dscFactory: dscFactory,
		dsRepo:     dsRepo,
		omRepo:     omRepo,
		bucketRepo: bucketRepo,
		lgRepo:     lgRepo,
		eccRepo:    eccRepo,
	}
}

func (osvc *ObjectService) Refresh(ctx context.Context) error {
	dss, err := osvc.dsRepo.GetDatastores(ctx)
	if err != nil {
		return fmt.Errorf("failed to get datastores: %w", err)
	}
	for _, ds := range dss {
		osvc.dsClients[ds.ID] = osvc.dscFactory.New(ds)
	}
	return nil
}

func (osvc *ObjectService) CreateObject(ctx context.Context, name, bucketName string, r io.Reader) error {
	slog.Debug("ObjectService::CreateObject called.", "name", name, "bucketName", bucketName)
	if !validObjectName.MatchString(name) {
		return errors.Join(fmt.Errorf("invalid object name: %s", name), ErrInvalidParameter)
	}

	err := osvc.tx.Do(ctx, func(ctx context.Context) error {
		bucket, err := osvc.bucketRepo.GetBucketByName(ctx, bucketName)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return errors.Join(fmt.Errorf("bucket not found"), ErrInvalidParameter)
			}
			return fmt.Errorf("failed to get bucket by name: %w", err)
		}
		if bucket.Status != entity.BucketStatusActive {
			return fmt.Errorf("unexpected bucket status: %s", bucket.Status)
		}
		om, err := osvc.omRepo.GetObjectMetadataByName(ctx, name, bucket.ID)
		if err != nil {
			if !errors.Is(err, ErrNotFound) {
				return fmt.Errorf("failed to get object metadata by name: %w", err)
			}
			_, err := osvc.createObjectMetadata(ctx, name, bucket)
			if err != nil {
				return fmt.Errorf("failed to create object metadata: %w", err)
			}
			om, err = osvc.omRepo.GetObjectMetadataByName(ctx, name, bucket.ID)
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

		ecConfig, err := osvc.eccRepo.GetECConfig(ctx, lg.ECConfigID)
		if err != nil {
			return fmt.Errorf("failed to get EC config: %w", err)
		}
		m := ec.NewManager(ecConfig.NumData, ecConfig.NumParity, minECChunkSizeInByte)
		codes, err := m.Encode(data)
		if err != nil {
			return fmt.Errorf("failed to encode: %w", err)
		}
		if len(codes) != len(lg.CurrentDatastores) {
			return fmt.Errorf("unexpected code length: expected=%d, actual=%d",
				len(lg.CurrentDatastores), len(codes))
		}
		for _, ds := range lg.CurrentDatastores {
			if _, ok := osvc.dsClients[ds]; !ok {
				err := osvc.Refresh(ctx)
				if err != nil {
					return fmt.Errorf("failed to refresh: %w", err)
				}
				break
			}
		}
		eg := new(errgroup.Group)
		var errorCount atomic.Int32
		for i, ds := range lg.CurrentDatastores {
			eg.Go(func() error {
				err = osvc.dsClients[ds].CreateObject(ctx, filepath.Join(bucketName, name), bytes.NewBuffer(codes[i]))
				if err != nil {
					errorCount.Add(1)
					return fmt.Errorf("CreateObject failed: %w", err)
				}
				return nil
			})
		}
		// FIXME: It is dangerous to accept PUT when numParity datastores are down.
		if err := eg.Wait(); err != nil && errorCount.Load() > int32(ecConfig.NumParity) {
			return fmt.Errorf("failed to create object chunk: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

func (osvc *ObjectService) GetObject(ctx context.Context, name, bucketName string) ([]byte, error) {
	slog.Debug("ObjectService::GetObject called.", "name", name, "bucketName", bucketName)
	if !validObjectName.MatchString(name) {
		return nil, errors.Join(fmt.Errorf("invalid object name: %s", name), ErrInvalidParameter)
	}

	var lg *entity.LocationGroup
	var ecConfig *entity.ECConfig
	err := osvc.tx.Do(ctx, func(ctx context.Context) error {
		bucket, err := osvc.bucketRepo.GetBucketByName(ctx, bucketName)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return errors.Join(fmt.Errorf("bucket not found"), ErrInvalidParameter)
			}
			return fmt.Errorf("failed to get bucket by name: %w", err)
		}
		om, err := osvc.omRepo.GetObjectMetadataByName(ctx, name, bucket.ID)
		if err != nil {
			return fmt.Errorf("failed to get object metadata: %w", err)
		}
		lg, err = osvc.lgRepo.GetLocationGroup(ctx, om.LocationGroupID)
		if err != nil {
			return fmt.Errorf("failed to get location group: %w", err)
		}
		ecConfig, err = osvc.eccRepo.GetECConfig(ctx, lg.ECConfigID)
		if err != nil {
			return fmt.Errorf("failed to get EC config: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	eg := new(errgroup.Group)
	var errorCount atomic.Int32
	codes := make([][]byte, ecConfig.NumData+ecConfig.NumParity)
	for i, ds := range lg.CurrentDatastores {
		eg.Go(func() error {
			rc, err := osvc.dsClients[ds].GetObject(ctx, filepath.Join(bucketName, name))
			if err != nil {
				errorCount.Add(1)
				return fmt.Errorf("GetObject failed: %w", err)
			}
			defer rc.Close()
			data, err := io.ReadAll(rc)
			if err != nil {
				errorCount.Add(1)
				return fmt.Errorf("failed to read data: %w", err)
			}
			codes[i] = data
			return nil
		})
	}
	if err := eg.Wait(); err != nil && errorCount.Load() > int32(ecConfig.NumParity) {
		return nil, fmt.Errorf("failed to get object chunk: %w", err)
	}

	m := ec.NewManager(ecConfig.NumData, ecConfig.NumParity, minECChunkSizeInByte)
	data, err := m.Decode(codes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}
	return data, nil
}

func (osvc *ObjectService) createObjectMetadata(
	ctx context.Context,
	name string,
	bucket *entity.Bucket,
) (int64, error) {
	lgs, err := osvc.lgRepo.GetLocationGroupsByECConfigID(ctx, bucket.ECConfigID)
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

func (osvc *ObjectService) DeleteObject(ctx context.Context, name, bucketName string, r io.Reader) error {
	slog.Debug("ObjectService::DeleteObject called.", "name", name, "bucketName", bucketName)
	if !validObjectName.MatchString(name) {
		return errors.Join(fmt.Errorf("invalid object name: %s", name), ErrInvalidParameter)
	}

	err := osvc.tx.Do(ctx, func(ctx context.Context) error {
		bucket, err := osvc.bucketRepo.GetBucketByName(ctx, bucketName)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return errors.Join(fmt.Errorf("bucket not found"), ErrInvalidParameter)
			}
			return fmt.Errorf("failed to get bucket by name: %w", err)
		}
		om, err := osvc.omRepo.GetObjectMetadataByName(ctx, name, bucket.ID)
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

		eg := new(errgroup.Group)
		var errorCount atomic.Int32
		for _, ds := range lg.CurrentDatastores {
			eg.Go(func() error {
				err = osvc.dsClients[ds].DeleteObject(ctx, filepath.Join(bucketName, name))
				if err != nil {
					errorCount.Add(1)
					return fmt.Errorf("DeleteObject failed: %w", err)
				}
				return nil
			})
		}
		// Even if deletes fail for all datastores, logging this fact
		// allows for recovery later. However, to avoid failing
		// to notice if a bug causes deletes to never succeed at all,
		// we require that deletes succeed for at least one datastore.
		if err := eg.Wait(); err != nil && errorCount.Load() == int32(len(lg.CurrentDatastores)) {
			return fmt.Errorf("failed to delete object chunk: %w", err)
		}

		err = osvc.omRepo.DeleteObjectMetadata(ctx, om.ID)
		if err != nil {
			return fmt.Errorf("failed to delete object metadata: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

func (osvc *ObjectService) ListObjects(
	ctx context.Context, bucketName string,
	startFrom int64, limit int,
) ([]string, int64, error) {
	slog.Debug("ObjectService::ListObjects called.",
		"bucketName", bucketName, "startFrom", startFrom, "limit", limit)
	if limit > 1000 {
		return nil, 0, fmt.Errorf("limit must not larger than 1000: %w", ErrInvalidParameter)
	}
	b, err := osvc.bucketRepo.GetBucketByName(ctx, bucketName)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, 0, ErrInvalidParameter
		}
		return nil, 0, fmt.Errorf("failed to get bucket by name: %w", err)
	}
	oms, err := osvc.omRepo.GetObjectMetadatas(ctx, &GetObjectMetadatasRequest{
		BucketID:  b.ID,
		StartFrom: startFrom,
		Limit:     limit + 1,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get object metadatas: %w", err)
	}
	var nextObjectID int64 = math.MaxInt64
	if len(oms) == limit+1 {
		nextObjectID = oms[limit].ID
		oms = oms[:limit]
	}
	objNames := make([]string, 0, len(oms))
	for _, om := range oms {
		objNames = append(objNames, om.Name)
	}

	return objNames, nextObjectID, err
}
