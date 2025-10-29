package service

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"slices"

	"github.com/peng225/orochi/internal/entity"
)

const (
	targetLGNumPerDatastore = 100
)

type DatastoreService struct {
	dsRepo DatastoreRepository
	lgRepo LocationGroupRepository
}

func NewDatastoreService(dsRepo DatastoreRepository, lgRepo LocationGroupRepository) *DatastoreService {
	return &DatastoreService{
		dsRepo: dsRepo,
		lgRepo: lgRepo,
	}
}

func (dss *DatastoreService) GetDatastore(ctx context.Context, id int64) (*entity.Datastore, error) {
	return dss.dsRepo.GetDatastore(ctx, id)
}

func (dss *DatastoreService) CreateDatastore(ctx context.Context, baseURL string) (int64, error) {
	if !isValidURL(baseURL) {
		return 0, ErrInvalidParameter
	}
	// FIXME: Transaction required.
	id, err := dss.dsRepo.CreateDatastore(ctx, &CreateDatastoreRequest{
		BaseURL: baseURL,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create datastore: %w", err)
	}

	err = dss.reconstructLocationGroups(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to reconstruct location groups: %w", err)
	}

	return id, nil
}

func (dss *DatastoreService) reconstructLocationGroups(ctx context.Context) error {
	dsIDs, err := dss.dsRepo.GetDatastoreIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get datastore IDs: %w", err)
	}
	lgs, err := dss.lgRepo.GetLocationGroups(ctx)
	if err != nil {
		return fmt.Errorf("failed to get location groups: %w", err)
	}

	// FIXME: Should remove the assumption of 2D1P erasure coding.
	stripeWidth := 3
	if len(dsIDs) < stripeWidth {
		slog.Info("Too small number of datastores.", "# of datastores", len(dsIDs))
		return nil
	}
	perm := permutation(len(dsIDs), stripeWidth, targetLGNumPerDatastore)
	targetNum := min(perm, targetLGNumPerDatastore)
	currentLGNumPerDatastore := len(lgs) / len(dsIDs)
	if targetNum/2 <= currentLGNumPerDatastore &&
		currentLGNumPerDatastore <= targetNum*3/2 {
		// Don't need to reconstruct location groups.
		return nil
	}
	if targetNum < len(lgs) {
		err := dss.shrinkLocationGroup()
		if err != nil {
			return fmt.Errorf("failed to shrink location group: %w", err)
		}
	} else {
		newDesiredDatastores := generateNewDesiredDSs(dsIDs, stripeWidth, targetNum)
		err := dss.expandLocationGroup(ctx, lgs, newDesiredDatastores)
		if err != nil {
			return fmt.Errorf("failed to expand location group: %w", err)
		}
	}
	return nil
}

func (dss *DatastoreService) shrinkLocationGroup() error {
	// FIXME: implement
	return nil
}

func (dss *DatastoreService) expandLocationGroup(ctx context.Context, lgs []*entity.LocationGroup, newDesiredDSs [][]int64) error {
	i := 0
	for _, lg := range lgs {
		err := dss.lgRepo.UpdateDesiredDatastores(ctx, lg.ID, newDesiredDSs[i])
		if err != nil {
			return err
		}
		i++
	}
	for _, desiredDs := range newDesiredDSs[i:] {
		_, err := dss.lgRepo.CreateLocationGroup(ctx, &CreateLocationGroupRequest{
			Datastores: desiredDs,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func isValidURL(s string) bool {
	re := regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+(:[0-9]+)?$`)
	return re.MatchString(s)
}

func permutation(n, k, upperBound int) int {
	if k < 0 || k > n {
		return 0
	}
	result := 1
	for i := range k {
		result *= n - i
		if result >= upperBound {
			return upperBound
		}
	}
	return result
}

func generateNewDesiredDSs(dsIDs []int64, stripeWidth, targetNum int) [][]int64 {
	current := make([]int64, 0, stripeWidth)
	result := make([][]int64, 0, targetNum)
	generateNewDesiredDSsHelper(dsIDs, stripeWidth, targetNum, current, &result)
	return result
}

func generateNewDesiredDSsHelper(dsIDs []int64, stripeWidth, targetNum int, current []int64, result *[][]int64) {
	if len(*result) == targetNum {
		return
	}
	if len(current) == stripeWidth {
		*result = append(*result, slices.Clone(current))
		return
	}
	for i, dsID := range dsIDs {
		current = append(current, dsID)
		generateNewDesiredDSsHelper(append(slices.Clone(dsIDs[0:i]), dsIDs[i+1:]...), stripeWidth, targetNum, current, result)
		current = current[:len(current)-1]
	}
}
