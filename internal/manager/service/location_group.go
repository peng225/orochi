package service

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/peng225/orochi/internal/entity"
)

const (
	targetLGNumPerDatastore = 100
)

type LocationGroupService struct {
	tx      Transaction
	dsRepo  DatastoreRepository
	lgRepo  LocationGroupRepository
	eccRepo ECConfigRepository
}

func NewLocationGroupService(
	tx Transaction, dsRepo DatastoreRepository, lgRepo LocationGroupRepository, eccRepo ECConfigRepository,
) *LocationGroupService {
	return &LocationGroupService{
		tx:      tx,
		dsRepo:  dsRepo,
		lgRepo:  lgRepo,
		eccRepo: eccRepo,
	}
}

func (lgs *LocationGroupService) ReconstructLocationGroupsForAllECConfigs(ctx context.Context) error {
	ecConfigs, err := lgs.eccRepo.GetECConfigs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get EC configs: %w", err)
	}
	for _, ecConfig := range ecConfigs {
		err := lgs.ReconstructLocationGroups(ctx, ecConfig)
		if err != nil {
			return fmt.Errorf("failed to reconstruct location group for EC ID %d: %w", ecConfig.ID, err)
		}
	}
	return nil
}

func (lgs *LocationGroupService) ReconstructLocationGroups(ctx context.Context, ecConfig *entity.ECConfig) error {
	dsIDs, err := lgs.dsRepo.GetDatastoreIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get datastore IDs: %w", err)
	}
	locationGroups, err := lgs.lgRepo.GetLocationGroupsByECConfigID(ctx, ecConfig.ID)
	if err != nil {
		return fmt.Errorf("failed to get location groups: %w", err)
	}

	stripeWidth := ecConfig.NumData + ecConfig.NumParity
	if len(dsIDs) < stripeWidth {
		slog.Info("Too small number of datastores.", "# of datastores", len(dsIDs))
		return nil
	}
	perm := permutation(len(dsIDs), stripeWidth, targetLGNumPerDatastore)
	targetNumPerDatastore := min(perm/len(dsIDs), targetLGNumPerDatastore)
	currentLGNumPerDatastore := len(locationGroups) / len(dsIDs)
	if targetNumPerDatastore/2 <= currentLGNumPerDatastore &&
		currentLGNumPerDatastore <= targetNumPerDatastore*3/2 {
		// Don't need to reconstruct location groups.
		return nil
	}
	slog.Info("Location group update required.",
		"len(dsIDs)", len(dsIDs),
		"stripeWidth", stripeWidth,
		"targetNumPerDatastore", targetNumPerDatastore,
		"currentLGNumPerDatastore", currentLGNumPerDatastore)
	if targetNumPerDatastore < currentLGNumPerDatastore {
		err := lgs.shrinkLocationGroup()
		if err != nil {
			return fmt.Errorf("failed to shrink location group: %w", err)
		}
	} else {
		newDesiredDatastores := generateNewDesiredDSs(dsIDs, stripeWidth, targetNumPerDatastore*len(dsIDs))
		// FIXME: should I create a location group for EC configs that have the same stripe width,
		//        not per EC config?
		err := lgs.expandLocationGroup(ctx, locationGroups, newDesiredDatastores, ecConfig.ID)
		if err != nil {
			return fmt.Errorf("failed to expand location group: %w", err)
		}
	}
	return nil
}

func (lgs *LocationGroupService) shrinkLocationGroup() error {
	// FIXME: implement
	return nil
}

func (lgs *LocationGroupService) expandLocationGroup(
	ctx context.Context,
	locationGroups []*entity.LocationGroup,
	newDesiredDSs [][]int64,
	ecConfigID int64,
) error {
	i := 0
	// Avoid to update location groups that already have
	// one of the new desired datastores.
	for j := range len(locationGroups) {
		for k := range newDesiredDSs[i:] {
			if slices.Equal(locationGroups[j].DesiredDatastores, newDesiredDSs[k]) {
				locationGroups[i], locationGroups[j] = locationGroups[j], locationGroups[i]
				newDesiredDSs[i], newDesiredDSs[k] = newDesiredDSs[k], newDesiredDSs[i]
				i++
				break
			}
		}
	}

	for _, lg := range locationGroups[i:] {
		// FIXME: should move the actual object asynchronously
		// according to the new desired datastores.
		err := lgs.lgRepo.UpdateDesiredDatastores(ctx, lg.ID, newDesiredDSs[i])
		if err != nil {
			return err
		}
		i++
	}
	for _, desiredDs := range newDesiredDSs[i:] {
		_, err := lgs.lgRepo.CreateLocationGroup(ctx, &CreateLocationGroupRequest{
			Datastores: desiredDs,
			ECConfigID: ecConfigID,
		})
		if err != nil {
			return err
		}
	}
	return nil
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

// FIXME: should resolve the unbalance allocations among datastores.
// Shuffling dsIDs before iteration may resolve the issue.
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
