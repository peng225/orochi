package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseECConfig_Success(t *testing.T) {
	testCases := []struct {
		param             string
		expectedNumData   int
		expectedNumParity int
	}{
		{
			param:             "2D1P",
			expectedNumData:   2,
			expectedNumParity: 1,
		},
		{
			param:             "20D10P",
			expectedNumData:   20,
			expectedNumParity: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.param, func(t *testing.T) {
			t.Parallel()
			numData, numParity, err := ParseECConfig(tc.param)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedNumData, numData)
			assert.Equal(t, tc.expectedNumParity, numParity)
		})
	}
}
