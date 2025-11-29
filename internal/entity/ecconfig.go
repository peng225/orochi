package entity

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	maxNumData   = 20
	maxNumParity = 10
)

var (
	validECParam = regexp.MustCompile(`^([0-9]+)D([0-9]+)P$`)
)

type ECConfig struct {
	ID        int64 `json:"id,omitempty"`
	NumData   int   `json:"numData,omitempty"`
	NumParity int   `json:"numParity,omitempty"`
}

func ParseECConfig(param string) (int, int, error) {
	matches := validECParam.FindStringSubmatch(param)
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("invalid EC param: %s", param)
	}

	numData, err := strconv.ParseInt(matches[1], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse the number of data: %w", err)
	}
	if numData < 1 || maxNumData < numData {
		return 0, 0, fmt.Errorf("invalid number of data: %d", numData)
	}

	numParity, err := strconv.ParseInt(matches[2], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse the number of data: %w", err)
	}
	if numParity < 1 || maxNumParity < numParity {
		return 0, 0, fmt.Errorf("too large number of parity: %d", numParity)
	}
	return int(numData), int(numParity), nil
}

func MustParseECConfig(param string) (int, int) {
	d, p, err := ParseECConfig(param)
	if err != nil {
		panic(err)
	}
	return d, p
}
