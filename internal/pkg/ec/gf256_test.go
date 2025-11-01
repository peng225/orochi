package ec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	x := gf256(0x82)
	require.Equal(t, gf256(0x21), x.Add(0xa3))
}

func TestMul(t *testing.T) {
	testCases := []struct {
		name     string
		a        gf256
		b        gf256
		expected gf256
	}{
		{
			name:     "multiplied by zero",
			a:        0x5,
			b:        0x0,
			expected: 0x0,
		},
		{
			name:     "multiplied by one",
			a:        0xb,
			b:        0x1,
			expected: 0xb,
		},
		{
			name:     "a == b",
			a:        0xa,
			b:        0xa,
			expected: 0x44,
		},
		{
			name:     "a > b",
			a:        0x2a,
			b:        0x5,
			expected: 0x82,
		},
		{
			name:     "a < b",
			a:        0x5,
			b:        0x2a,
			expected: 0x82,
		},
		{
			name:     "large",
			a:        0xc5,
			b:        0xd9,
			expected: 0x8f,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, tc.a.Mul(tc.b))
		})
	}
}

func TestDiv(t *testing.T) {
	testCases := []struct {
		name     string
		a        gf256
		b        gf256
		expected gf256
	}{
		{
			name:     "divide zero by some value",
			a:        0x0,
			b:        0x3,
			expected: 0x0,
		},
		{
			name:     "divide by one",
			a:        0x10,
			b:        0x1,
			expected: 0x10,
		},
		{
			name:     "a == b",
			a:        0xa,
			b:        0xa,
			expected: 0x1,
		},
		{
			name:     "other",
			a:        0xe5,
			b:        0x26,
			expected: 0x1c,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, tc.a.Div(tc.b))
		})
	}
}

func TestPow(t *testing.T) {
	testCases := []struct {
		name     string
		a        gf256
		i        int
		expected gf256
	}{
		{
			name:     "power of zero",
			a:        0x0,
			i:        3,
			expected: 0x0,
		},
		{
			name:     "power of one",
			a:        0x1,
			i:        5,
			expected: 0x1,
		},
		{
			name:     "powered by zero",
			a:        0x5,
			i:        0,
			expected: 0x1,
		},
		{
			name:     "other",
			a:        0xa5,
			i:        5,
			expected: 0x91,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, tc.a.Pow(tc.i))
		})
	}
}

func TestInv(t *testing.T) {
	testCases := []struct {
		name     string
		a        gf256
		expected gf256
	}{
		{
			name:     "inverse of one",
			a:        0x1,
			expected: 0x1,
		},
		{
			name:     "other",
			a:        0xa5,
			expected: 0xb8,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, tc.a.Inv())
		})
	}
}
