package service

import (
	"context"
	"errors"
	"fmt"
	"io"
)

type ObjectService struct {
	chunkRepos map[int64]ChunkRepository
	crFactory  ChunkRepositoryFactory
	dsRepo     DatastoreRepository
}

func NewObjectStore(
	chunkRepos map[int64]ChunkRepository,
	crFactory ChunkRepositoryFactory,
	dsRepo DatastoreRepository,
) *ObjectService {
	if chunkRepos == nil {
		chunkRepos = make(map[int64]ChunkRepository)
	}
	return &ObjectService{
		chunkRepos: chunkRepos,
		crFactory:  crFactory,
		dsRepo:     dsRepo,
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

type multiReadCloser struct {
	readClosers []io.ReadCloser
	multiReader io.Reader
}

func newMultiReadCloser(readClosers []io.ReadCloser) *multiReadCloser {
	readers := make([]io.Reader, len(readClosers))
	for i, rc := range readClosers {
		readers[i] = rc
	}
	return &multiReadCloser{
		readClosers: readClosers,
		multiReader: io.MultiReader(readers...),
	}
}

func (mrc *multiReadCloser) Read(p []byte) (n int, err error) {
	return mrc.multiReader.Read(p)
}

func (mrc *multiReadCloser) Close() error {
	var joinedError error
	for _, rc := range mrc.readClosers {
		err := rc.Close()
		if err != nil {
			joinedError = errors.Join(joinedError, err)
		}
	}
	return joinedError
}

func (osvc *ObjectService) GetObject(ctx context.Context, bucket, object string) (io.ReadCloser, error) {
	dsToData := make(map[int64]io.ReadCloser)
	// FIXME: parallelize.
	for id, chunkRepo := range osvc.chunkRepos {
		data, err := chunkRepo.GetObject(ctx, bucket, object)
		if err != nil {
			return nil, fmt.Errorf("GetObject failed: %w", err)
		}
		dsToData[id] = data
	}
	// FIXME: need to order the dataList correctly.
	dataList := make([]io.ReadCloser, 0, len(dsToData))
	for _, d := range dsToData {
		dataList = append(dataList, d)
	}
	return newMultiReadCloser(dataList), nil
}

func (osvc *ObjectService) CreateObject(ctx context.Context, bucket, object string, data io.Reader) error {
	err := osvc.Refresh(ctx)
	if err != nil {
		return fmt.Errorf("failed to refresh: %w", err)
	}
	// FIXME: select the correct client instance.
	err = osvc.chunkRepos[1].CreateObject(ctx, bucket, object, data)
	if err != nil {
		return fmt.Errorf("CreateObject failed: %w", err)
	}
	return nil
}
