package ec

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	mat := NewGF256MatrixWithData(2, []gf256{1, 2, 3, 4})
	require.Equal(t, gf256(1), mat.Get(0, 0))
	require.Equal(t, gf256(2), mat.Get(0, 1))
	require.Equal(t, gf256(3), mat.Get(1, 0))
	require.Equal(t, gf256(4), mat.Get(1, 1))
}

func TestSet(t *testing.T) {
	mat := NewGF256Matrix(2, 3)
	mat.Set(1, 1, 0x5)
	require.Equal(t, gf256(0x5), mat.Get(1, 1))
	mat.Set(1, 2, 0x6)
	require.Equal(t, gf256(0x6), mat.Get(1, 2))
}

func TestInverse(t *testing.T) {
	testCases := []struct {
		name     string
		m        *GF256Matrix
		expected *GF256Matrix
	}{
		{
			name:     "inverse of the unit matrix",
			m:        NewGF256MatrixWithData(2, []gf256{1, 0, 0, 1}),
			expected: NewGF256MatrixWithData(2, []gf256{1, 0, 0, 1}),
		},
		{
			name:     "other1",
			m:        NewGF256MatrixWithData(2, []gf256{0x28, 0x3, 0xc8, 0x4}),
			expected: NewGF256MatrixWithData(2, []gf256{0x81, 0x26, 0x5a, 0x7d}),
		},
		{
			name:     "other2",
			m:        NewGF256MatrixWithData(2, []gf256{0x0, 0x1, 0x1, 0x4}),
			expected: NewGF256MatrixWithData(2, []gf256{0x4, 0x1, 0x1, 0}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			inv, err := tc.m.Inverse()
			require.NoError(t, err)
			require.True(t, slices.Equal(tc.expected.data, inv.data))
		})
	}
}

func TestMulRight(t *testing.T) {
	testCases := []struct {
		name     string
		m1       *GF256Matrix
		m2       *GF256Matrix
		expected *GF256Matrix
	}{
		{
			name:     "I * I",
			m1:       NewGF256MatrixWithData(2, []gf256{1, 0, 0, 1}),
			m2:       NewGF256MatrixWithData(2, []gf256{1, 0, 0, 1}),
			expected: NewGF256MatrixWithData(2, []gf256{1, 0, 0, 1}),
		},
		{
			name:     "not square",
			m1:       NewGF256MatrixWithData(2, []gf256{0x81, 0x2, 0x1, 0x14, 0x0, 0xc1}),
			m2:       NewGF256MatrixWithData(3, []gf256{0x1, 0x20, 0x20, 0x3, 0xc2, 0x5}),
			expected: NewGF256MatrixWithData(2, []gf256{0x3, 0x88, 0x7c, 0x5e}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mul, err := tc.m1.MulRight(tc.m2)
			require.NoError(t, err)
			require.True(t, slices.Equal(tc.expected.data, mul.data))
		})
	}
}
