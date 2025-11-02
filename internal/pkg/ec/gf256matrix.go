package ec

import (
	"errors"
	"fmt"
)

type GF256Matrix struct {
	numRows int
	data    []gf256
}

func NewGF256Matrix(numRows, numCols int) *GF256Matrix {
	if numRows <= 0 {
		panic("numRows must be positive")
	}
	if numCols <= 0 {
		panic("numCols must be positive")
	}
	return &GF256Matrix{
		numRows: numRows,
		data:    make([]gf256, numRows*numCols),
	}
}

func NewGF256MatrixWithData(numRows int, data []gf256) *GF256Matrix {
	if numRows <= 0 {
		panic("numRows must be positive")
	}
	if len(data)%numRows != 0 {
		panic("len(data) must be divisible by numRows")
	}
	return &GF256Matrix{
		numRows: numRows,
		data:    data,
	}
}

func (mat *GF256Matrix) NumCols() int {
	return len(mat.data) / mat.numRows
}

func (mat *GF256Matrix) Get(i, j int) gf256 {
	return mat.data[i*mat.NumCols()+j]
}

func (mat *GF256Matrix) Set(i, j int, v gf256) {
	mat.data[i*mat.NumCols()+j] = v
}

func (mat *GF256Matrix) swapRows(i, j int) error {
	if i < 0 || j < 0 {
		return fmt.Errorf("both i and j must be positive: i=%d, j=%d", i, j)
	}
	if i >= mat.numRows || j >= mat.numRows {
		return fmt.Errorf("both i and j must be smaller than numRows: numROws=%d i=%d, j=%d",
			mat.numRows, i, j)
	}
	for k := range mat.NumCols() {
		mat.data[i*mat.NumCols()+k], mat.data[j*mat.NumCols()+k] = mat.Get(j, k), mat.Get(i, k)
	}
	return nil
}

func (mat *GF256Matrix) MulRight(m *GF256Matrix) (*GF256Matrix, error) {
	if mat.NumCols() != m.numRows {
		return nil, fmt.Errorf("matrix size mismatch")
	}
	result := NewGF256Matrix(mat.numRows, m.NumCols())
	for i := range mat.numRows {
		for j := range m.NumCols() {
			var sum gf256
			for k := range mat.NumCols() {
				sum = sum.Add(mat.Get(i, k).Mul(m.Get(k, j)))
			}
			result.Set(i, j, sum)
		}
	}
	return result, nil
}

func (mat *GF256Matrix) Inverse() (*GF256Matrix, error) {
	if mat.numRows != len(mat.data)/mat.numRows {
		return nil, errors.New("matrix must be square")
	}

	// Expansion matrix [A | I]
	aug := NewGF256Matrix(mat.numRows, 2*mat.numRows)
	for i := range mat.numRows {
		for j := range mat.NumCols() {
			aug.Set(i, j, mat.Get(i, j))
		}
		aug.Set(i, mat.NumCols()+i, 1)
	}

	// Gauss–Jordan elimination
	for i := range aug.NumCols() / 2 {
		// If the pivot is 0, exchange with another row.
		if aug.Get(i, i) == 0 {
			found := false
			for k := i + 1; k < aug.numRows; k++ {
				if aug.Get(k, i) != 0 {
					err := aug.swapRows(i, k)
					if err != nil {
						return nil, err
					}
					found = true
					break
				}
			}
			if !found {
				return nil, errors.New("matrix is singular")
			}
		}

		// Normalize elements by the inverse of the pivot.
		pivotInv := aug.Get(i, i).Inv()
		for j := range aug.NumCols() {
			aug.Set(i, j, aug.Get(i, j).Mul(pivotInv))
		}

		// For column i, delete values in other rows.
		for k := range aug.numRows {
			if k == i {
				continue
			}
			factor := aug.Get(k, i)
			if factor != 0 {
				for j := range aug.NumCols() {
					aug.Set(k, j, aug.Get(k, j).Sub(factor.Mul(aug.Get(i, j))))
				}
			}
		}
	}

	inv := NewGF256Matrix(mat.numRows, mat.numRows)
	for i := range inv.numRows {
		for j := range inv.NumCols() {
			inv.Set(i, j, aug.Get(i, aug.NumCols()/2+j))
		}
	}
	return inv, nil
}
