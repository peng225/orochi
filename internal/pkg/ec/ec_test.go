package ec

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeDecode_NoMissingChunks(t *testing.T) {
	testCases := []struct {
		name               string
		numData            int
		numParity          int
		minChunkSizeInByte int
		data               []byte
	}{
		{
			name:               "(2D1P) data size is smaller than minChunkSizeInByte-header",
			numData:            2,
			numParity:          1,
			minChunkSizeInByte: 10,
			data:               []byte("abc"),
		},
		{
			name:               "(2D1P) data size is smaller than numData*minChunkSizeInByte-header",
			numData:            2,
			numParity:          1,
			minChunkSizeInByte: 10,
			data:               []byte("abcdefghij"),
		},
		{
			name:               "(2D1P) data size is larger than numData*minChunkSizeInByte-header",
			numData:            2,
			numParity:          1,
			minChunkSizeInByte: 10,
			data:               []byte("abcdefghijklmnopqrstuv"),
		},
		{
			name:               "(3D2P) data size is smaller than numData*minChunkSizeInByte-header",
			numData:            3,
			numParity:          2,
			minChunkSizeInByte: 10,
			data:               []byte("abcdefghij"),
		},
		{
			name:               "(3D2P) data size is larger than numData*minChunkSizeInByte-header",
			numData:            3,
			numParity:          2,
			minChunkSizeInByte: 8,
			data:               []byte("abcdefghijklmnopqrstuvwxyz"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := NewManager(tc.numData, tc.numParity,
				tc.minChunkSizeInByte, 1000)
			codes, err := m.Encode(tc.data)
			require.NoError(t, err)
			require.Len(t, codes, tc.numData+tc.numParity)

			data, err := m.Decode(codes)
			require.NoError(t, err)
			require.True(t, slices.Equal(data, tc.data))
		})
	}
}

func TestEncodeDecode_RecoverFromMissingChunks(t *testing.T) {
	testCases := []struct {
		name               string
		numData            int
		numParity          int
		minChunkSizeInByte int
		data               []byte
		missingChunkIDs    []int
	}{
		{
			name:               "(2D1P) data size is smaller than minChunkSizeInByte-header with first data chunk missing",
			numData:            2,
			numParity:          1,
			minChunkSizeInByte: 10,
			data:               []byte("abc"),
			missingChunkIDs:    []int{0},
		},
		{
			name:               "(2D1P) data size is smaller than numData*minChunkSizeInByte-header with first data chunk missing",
			numData:            2,
			numParity:          1,
			minChunkSizeInByte: 10,
			data:               []byte("abcdefghij"),
			missingChunkIDs:    []int{0},
		},
		{
			name:               "(2D1P) data size is smaller than numData*minChunkSizeInByte-header with second data chunk missing",
			numData:            2,
			numParity:          1,
			minChunkSizeInByte: 10,
			data:               []byte("abcdefghij"),
			missingChunkIDs:    []int{1},
		},
		{
			name:               "(2D1P) data size is smaller than numData*minChunkSizeInByte-header with parity chunk missing",
			numData:            2,
			numParity:          1,
			minChunkSizeInByte: 10,
			data:               []byte("abcdefghij"),
			missingChunkIDs:    []int{1},
		},
		{
			name:               "(2D1P) data size is larger than numData*minChunkSizeInByte-header with first data chunk missing",
			numData:            2,
			numParity:          1,
			minChunkSizeInByte: 10,
			data:               []byte("abcdefghijklmnopqrstuv"),
			missingChunkIDs:    []int{0},
		},
		{
			name:               "(3D2P) data size is smaller than numData*minChunkSizeInByte-header with second data chunk missing",
			numData:            3,
			numParity:          2,
			minChunkSizeInByte: 10,
			data:               []byte("abcdefghij"),
			missingChunkIDs:    []int{1},
		},
		{
			name:               "(3D2P) data size is larger than numData*minChunkSizeInByte-header with first and third data chunk missing",
			numData:            3,
			numParity:          2,
			minChunkSizeInByte: 8,
			data:               []byte("abcdefghijklmnopqrstuvwxyz"),
			missingChunkIDs:    []int{0, 2},
		},
		{
			name:               "(3D2P) data size is larger than numData*minChunkSizeInByte-header with data and parity chunk missing",
			numData:            3,
			numParity:          2,
			minChunkSizeInByte: 8,
			data:               []byte("abcdefghijklmnopqrstuvwxyz"),
			missingChunkIDs:    []int{0, 4},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := NewManager(tc.numData, tc.numParity,
				tc.minChunkSizeInByte, 1000)
			codes, err := m.Encode(tc.data)
			require.NoError(t, err)
			require.Len(t, codes, tc.numData+tc.numParity)

			for _, v := range tc.missingChunkIDs {
				codes[v] = nil
			}
			data, err := m.Decode(codes)
			require.NoError(t, err)
			require.True(t, slices.Equal(data, tc.data))
		})
	}
}
