package dsl

import (
	"fmt"
	"math/big"
)

// NullspaceBasisRational 求整数矩阵 A（m×n）在 Q^n 上的零空间一组基（每个基向量为 n 维有理向量）。
// 约定：Ax=0，x 为列向量。
func NullspaceBasisRational(A *MatrixInt) ([][]*big.Rat, error) {
	if A == nil || A.R == 0 || A.C == 0 {
		return nil, fmt.Errorf("nullspace: empty matrix")
	}
	r, n := A.R, A.C
	M := make([][]*big.Rat, r)
	for i := 0; i < r; i++ {
		M[i] = make([]*big.Rat, n)
		for j := 0; j < n; j++ {
			M[i][j] = new(big.Rat).SetInt64(A.A[i][j])
		}
	}
	pivotCols := make([]int, 0, n)
	row := 0
	for col := 0; col < n && row < r; col++ {
		pivot := -1
		for i := row; i < r; i++ {
			if M[i][col].Sign() != 0 {
				pivot = i
				break
			}
		}
		if pivot == -1 {
			continue
		}
		if pivot != row {
			M[pivot], M[row] = M[row], M[pivot]
		}
		pv := new(big.Rat).Set(M[row][col])
		for j := col; j < n; j++ {
			M[row][j].Quo(M[row][j], pv)
		}
		for i := 0; i < r; i++ {
			if i == row {
				continue
			}
			f := new(big.Rat).Set(M[i][col])
			if f.Sign() == 0 {
				continue
			}
			for j := col; j < n; j++ {
				tmp := new(big.Rat).Mul(f, M[row][j])
				M[i][j].Sub(M[i][j], tmp)
			}
		}
		pivotCols = append(pivotCols, col)
		row++
	}
	rank := len(pivotCols)
	dim := n - rank
	if dim <= 0 {
		return nil, fmt.Errorf("nullspace: trivial kernel")
	}
	freeCols := make([]int, 0, dim)
	pc := 0
	for col := 0; col < n; col++ {
		if pc < len(pivotCols) && pivotCols[pc] == col {
			pc++
			continue
		}
		freeCols = append(freeCols, col)
	}
	if len(freeCols) != dim {
		return nil, fmt.Errorf("nullspace: free col count mismatch")
	}

	basis := make([][]*big.Rat, 0, dim)
	for _, fc := range freeCols {
		x := make([]*big.Rat, n)
		for j := 0; j < n; j++ {
			x[j] = big.NewRat(0, 1)
		}
		x[fc] = big.NewRat(1, 1)
		for i := rank - 1; i >= 0; i-- {
			col := pivotCols[i]
			s := big.NewRat(0, 1)
			for j := col + 1; j < n; j++ {
				tmp := new(big.Rat).Mul(M[i][j], x[j])
				s.Add(s, tmp)
			}
			x[col] = new(big.Rat).Neg(s)
		}
		basis = append(basis, x)
	}
	return basis, nil
}
