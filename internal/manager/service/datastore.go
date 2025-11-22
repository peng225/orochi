package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/peng225/orochi/internal/entity"
)

var (
	validURL = regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+(:[0-9]+)?$`)
)

type DatastoreService struct {
	tx        Transaction
	lgService *LocationGroupService
	dsRepo    DatastoreRepository
}

func NewDatastoreService(
	tx Transaction, lgService *LocationGroupService, dsRepo DatastoreRepository,
) *DatastoreService {
	return &DatastoreService{
		tx:        tx,
		lgService: lgService,
		dsRepo:    dsRepo,
	}
}

func (dss *DatastoreService) CreateDatastore(ctx context.Context, baseURL string) (int64, error) {
	if !validURL.MatchString(baseURL) {
		return 0, errors.Join(fmt.Errorf("invalid baseURL: %s", baseURL), ErrInvalidParameter)
	}
	var id int64
	err := dss.tx.Do(ctx, func(ctx context.Context) error {
		ds, err := dss.dsRepo.GetDatastoreByBaseURL(ctx, baseURL)
		if err == nil {
			id = ds.ID
			return nil
		} else if !errors.Is(err, ErrNotFound) {
			return fmt.Errorf("failed to get datastore by base URL: %w", err)
		}
		id, err = dss.dsRepo.CreateDatastore(ctx, &CreateDatastoreRequest{
			BaseURL: baseURL,
		})
		if err != nil {
			return fmt.Errorf("failed to create datastore: %w", err)
		}

		err = dss.lgService.ReconstructLocationGroupsForAllECConfigs(ctx)
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

func (dss *DatastoreService) GetDatastore(ctx context.Context, id int64) (*entity.Datastore, error) {
	return dss.dsRepo.GetDatastore(ctx, id)
}
