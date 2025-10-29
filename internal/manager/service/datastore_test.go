package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermutation(t *testing.T) {
	testCases := []struct {
		name       string
		n          int
		k          int
		upperBound int
		expected   int
	}{
		{
			name:       "does not reach upper bound",
			n:          5,
			k:          3,
			upperBound: 1000,
			expected:   60,
		},
		{
			name:       "reach upper bound",
			n:          6,
			k:          3,
			upperBound: 20,
			expected:   20,
		},
	}
	for _, tc := range testCases {
		ret := permutation(tc.n, tc.k, tc.upperBound)
		assert.Equal(t, tc.expected, ret)
	}
}

func TestGenerateNewDesiredDSs(t *testing.T) {
	testCases := []struct {
		name  string
		dsIDs []int64
		stripeWidth,
		targetNum int
	}{
		{
			name:        "stripe width is 3",
			dsIDs:       []int64{1, 2, 3, 4},
			stripeWidth: 3,
			targetNum:   10,
		},
		{
			name:        "stripe width is 6",
			dsIDs:       []int64{1, 2, 3, 4, 5, 6, 7},
			stripeWidth: 6,
			targetNum:   10,
		},
	}
	for _, tc := range testCases {
		result := generateNewDesiredDSs(tc.dsIDs, tc.stripeWidth, tc.targetNum)
		assert.Len(t, result, tc.targetNum)
		for _, r := range result {
			appeared := make(map[int64]struct{})
			for _, dsID := range r {
				_, ok := appeared[dsID]
				require.False(t, ok, r)
				appeared[dsID] = struct{}{}
			}
		}
	}
}
