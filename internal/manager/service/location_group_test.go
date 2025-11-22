package service_test

import (
	"testing"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/manager/infra/fake"
	"github.com/peng225/orochi/internal/manager/service"
	"github.com/stretchr/testify/require"
)

func TestReconstructLocationGroups(t *testing.T) {
	testCases := []struct {
		name        string
		numDSs      int
		ecConfigStr string
		existingLGs []*service.CreateLocationGroupRequest
		expectedDSs [][]int64
	}{
		{
			name:        "no datastores",
			numDSs:      0,
			ecConfigStr: "2D1P",
			expectedDSs: [][]int64{},
		},
		{
			name:        "small number of datastores without existing LGs",
			numDSs:      3,
			ecConfigStr: "2D1P",
			expectedDSs: [][]int64{
				{1, 2, 3}, {1, 3, 2}, {2, 1, 3},
				{2, 3, 1}, {3, 1, 2}, {3, 2, 1},
			},
		},
		{
			name:        "small number of datastores with existing LGs",
			numDSs:      3,
			ecConfigStr: "2D1P",
			existingLGs: []*service.CreateLocationGroupRequest{
				{
					Datastores: []int64{2, 1, 3},
					ECConfigID: 1,
				},
			},
			expectedDSs: [][]int64{
				{1, 2, 3}, {1, 3, 2}, {2, 1, 3},
				{2, 3, 1}, {3, 1, 2}, {3, 2, 1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dsRepo := fake.NewFakeDatastoreRepository()
			lgRepo := fake.NewFakeLocationGroupRepository()
			for range tc.numDSs {
				_, err := dsRepo.CreateDatastore(t.Context(), &service.CreateDatastoreRequest{
					BaseURL: "http://example.com",
				})
				require.NoError(t, err)
			}
			for _, lg := range tc.existingLGs {
				_, err := lgRepo.CreateLocationGroup(t.Context(), lg)
				require.NoError(t, err)
			}
			lgService := service.NewLocationGroupService(nil, dsRepo, lgRepo, nil)

			d, p, err := entity.ParseECConfig(tc.ecConfigStr)
			require.NoError(t, err)
			ecConfig := &entity.ECConfig{
				ID:        1,
				NumData:   d,
				NumParity: p,
			}
			err = lgService.ReconstructLocationGroups(t.Context(), ecConfig)
			require.NoError(t, err)
			lgs, err := lgRepo.GetLocationGroupsByECConfigID(t.Context(), ecConfig.ID)
			require.NoError(t, err)
			for _, lg := range lgs {
				require.Contains(t, tc.expectedDSs, lg.Datastores)
			}
		})
	}
}
