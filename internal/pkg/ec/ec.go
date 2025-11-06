package ec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"maps"
	"math"
	"slices"
)

var (
	ErrInvalidParameter    error = errors.New("invalid parameter")
	ErrTooManyMissingCodes error = errors.New("too many missing codes")
)

const (
	headerSizeInByte = 4
)

type Manager struct {
	numData            int
	numParity          int
	minChunkSizeInByte int
	generatorMatrix    *GF256Matrix
}

func NewManager(
	numData, numParity, minChunkSizeInByte int,
) *Manager {
	if numData <= 0 {
		panic("numData should be positive")
	}
	if numParity <= 0 {
		panic("numParity should be positive")
	}
	if minChunkSizeInByte <= 0 {
		panic("minChunkSizeInByte should be positive")
	}
	// Calc parity blocks.
	// Choose 1, a, a^2, ... a^(d+p-1) as evaluation points.
	generatorMatrix := NewGF256Matrix(numData+numParity, numData)
	for i := range numData {
		for j := range numData {
			if i == j {
				generatorMatrix.Set(i, j, 1)
			}
		}
	}
	for i := range numParity {
		val := gf256(1)
		evalPoint := GF256Generator().Pow(numData + i)
		for j := range numData {
			generatorMatrix.Set(numData+i, j, val)
			val = val.Mul(evalPoint)
		}
	}

	return &Manager{
		numData:            numData,
		numParity:          numParity,
		minChunkSizeInByte: minChunkSizeInByte,
		generatorMatrix:    generatorMatrix,
	}
}

func (m *Manager) Encode(data []byte) ([][]byte, error) {
	result := make([][]byte, 0)
	dataSize := uint32(len(data))
	// Set data chunks.
	if len(data)+headerSizeInByte < m.numData*m.minChunkSizeInByte {
		// Write the header at the top of the first data chunk.
		chunkEndOffset := min(m.minChunkSizeInByte-headerSizeInByte, len(data))
		buf := make([]byte, headerSizeInByte, headerSizeInByte+chunkEndOffset)
		binary.LittleEndian.PutUint32(buf, dataSize)
		buf = append(buf, data[0:chunkEndOffset]...)
		result = append(result, buf)
		data = data[chunkEndOffset:]
		for len(data) != 0 {
			chunkEndOffset := min(m.minChunkSizeInByte, len(data))
			buf := make([]byte, 0, chunkEndOffset)
			buf = append(buf, data[0:chunkEndOffset]...)
			result = append(result, buf)
			data = data[chunkEndOffset:]
		}
		for range m.numData - len(result) {
			buf := make([]byte, 0)
			result = append(result, buf)
		}
	} else {
		chunkSize := int(
			math.Ceil(
				(float64(len(data)) + headerSizeInByte) / float64(m.numData),
			),
		)
		// Write the header at the top of the first data chunk.
		chunkEndOffset := min(chunkSize-headerSizeInByte, len(data))
		buf := make([]byte, headerSizeInByte, headerSizeInByte+chunkEndOffset)
		binary.LittleEndian.PutUint32(buf, dataSize)
		buf = append(buf, data[0:chunkEndOffset]...)
		result = append(result, buf)
		data = data[chunkEndOffset:]
		for len(data) != 0 {
			chunkEndOffset := min(chunkSize, len(data))
			buf := make([]byte, 0, chunkEndOffset)
			buf = append(buf, data[0:chunkEndOffset]...)
			result = append(result, buf)
			data = data[chunkEndOffset:]
		}
	}

	// Set parity chunks.
	for range m.numParity {
		result = append(result, make([]byte, 0))
	}
	for offsetInData := range len(result[0]) {
		for parityID := range m.numParity {
			parityCode := gf256(0)
			for dataID := range m.numData {
				if len(result[dataID]) <= offsetInData {
					continue
				}
				parityCode = gf256(parityCode).Add(
					m.generatorMatrix.Get(m.numData+parityID, dataID).Mul(
						gf256(result[dataID][offsetInData]),
					),
				)
			}
			result[m.numData+parityID] = append(result[m.numData+parityID], byte(parityCode))
		}
	}

	return result, nil
}

func (m *Manager) Decode(codes [][]byte) ([]byte, error) {
	if len(codes) != m.numData+m.numParity {
		return nil, ErrInvalidParameter
	}
	remainingChunkIDs := make(map[int]struct{})
	missingChunkIDs := make(map[int]struct{})
	missingDataChunkExists := false
	for i, v := range codes {
		if v == nil {
			missingChunkIDs[i] = struct{}{}
			if i < m.numData {
				missingDataChunkExists = true
			}
		} else {
			remainingChunkIDs[i] = struct{}{}
		}
	}
	if len(remainingChunkIDs) < m.numData {
		return nil, ErrTooManyMissingCodes
	}

	chunkSize := 0
	for _, c := range codes {
		if chunkSize < len(c) {
			chunkSize = len(c)
		}
	}
	data := make([]byte, 0, chunkSize*m.numData-headerSizeInByte)
	if !missingDataChunkExists {
		for i := range m.numData {
			startOffset := 0
			// Eliminate the header from the first data chunk.
			if i == 0 {
				startOffset = headerSizeInByte
			}
			data = append(data, codes[i][startOffset:]...)
		}
		return data, nil
	}

	remainingGenMatrix := NewGF256Matrix(m.numData, m.numData)
	row := 0
	for i := range len(codes) {
		if _, ok := remainingChunkIDs[i]; !ok {
			continue
		}
		for j := range m.numData {
			remainingGenMatrix.Set(row, j, m.generatorMatrix.Get(i, j))
		}
		row += 1
		if row == remainingGenMatrix.numRows {
			break
		}
	}
	inv, err := remainingGenMatrix.Inverse()
	if err != nil {
		return nil, fmt.Errorf("failed to get the inverse of remaining generator matrix: %w", err)
	}

	for id := range missingChunkIDs {
		codes[id] = make([]byte, 0, chunkSize)
	}
	sortedRemChunkIDs := slices.Sorted(maps.Keys(remainingChunkIDs))
	for i := range chunkSize {
		remainingSymbolsVec := NewGF256Matrix(m.numData, 1)
		for j, rcID := range sortedRemChunkIDs[:m.numData] {
			if i < len(codes[rcID]) {
				remainingSymbolsVec.Set(j, 0, gf256(codes[rcID][i]))
			}
		}
		decodedData, err := inv.MulRight(remainingSymbolsVec)
		if err != nil {
			return nil, fmt.Errorf("matrix multiplication failed: %w", err)
		}
		if decodedData.numRows != m.numData {
			return nil, fmt.Errorf("unexpected decoded data vector size")
		}
		for id := range missingChunkIDs {
			if id < m.numData {
				codes[id] = append(codes[id], byte(decodedData.Get(id, 0)))
			}
		}
	}

	for i := range m.numData {
		startOffset := 0
		// Eliminate the header from the first data chunk.
		if i == 0 {
			startOffset = headerSizeInByte
		}
		data = append(data, codes[i][startOffset:]...)
	}
	dataSize := binary.LittleEndian.Uint32(codes[0][:headerSizeInByte])
	return data[:dataSize], nil
}
