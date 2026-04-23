package dsl

import (
	"fmt"
	"math/big"
)

// IntegerKernelVectorOne 求整数矩阵 A（m×n）零空间中的一个非零整数向量；
// 要求消元后核维数为 1（常见于 m=n-1 且行满秩）。
func IntegerKernelVectorOne(A *MatrixInt) (*VectorInt, error) {
	r := A.R
	c := A.C
	if r == 0 || c == 0 {
		return nil, fmt.Errorf("kernel: empty matrix")
	}
	M := make([][]*big.Rat, r)
	for i := 0; i < r; i++ {
		M[i] = make([]*big.Rat, c)
		for j := 0; j < c; j++ {
			M[i][j] = new(big.Rat).SetInt64(A.A[i][j])
		}
	}
	row := 0
	pivotCols := make([]int, 0, c)
	for col := 0; col < c && row < r; col++ {
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
		for j := col; j < c; j++ {
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
			for j := col; j < c; j++ {
				tmp := new(big.Rat).Mul(f, M[row][j])
				M[i][j].Sub(M[i][j], tmp)
			}
		}
		pivotCols = append(pivotCols, col)
		row++
	}
	rank := len(pivotCols)
	if rank >= c {
		return nil, fmt.Errorf("kernel: trivial kernel")
	}
	freeCol := -1
	pc := 0
	for col := 0; col < c; col++ {
		if pc < len(pivotCols) && pivotCols[pc] == col {
			pc++
			continue
		}
		freeCol = col
		break
	}
	if freeCol == -1 {
		return nil, fmt.Errorf("kernel: no free column")
	}
	x := make([]*big.Rat, c)
	for j := 0; j < c; j++ {
		x[j] = big.NewRat(0, 1)
	}
	x[freeCol] = big.NewRat(1, 1)
	for i := rank - 1; i >= 0; i-- {
		col := pivotCols[i]
		s := big.NewRat(0, 1)
		for j := col + 1; j < c; j++ {
			tmp := new(big.Rat).Mul(M[i][j], x[j])
			s.Add(s, tmp)
		}
		x[col] = new(big.Rat).Neg(s)
	}
	lcmVal := big.NewInt(1)
	for j := 0; j < c; j++ {
		if x[j].Sign() == 0 {
			continue
		}
		d := x[j].Denom()
		lcmVal = lcmPositiveBig(lcmVal, d)
	}
	if lcmVal.Sign() == 0 {
		lcmVal = big.NewInt(1)
	}
	out := NewVectorInt(c)
	for j := 0; j < c; j++ {
		t := new(big.Rat).Mul(x[j], new(big.Rat).SetFrac(lcmVal, big.NewInt(1)))
		if !t.IsInt() {
			return nil, fmt.Errorf("kernel: scale failed at col %d", j)
		}
		out.V[j] = t.Num().Int64()
	}
	nonZero := false
	for j := 0; j < c; j++ {
		if out.V[j] != 0 {
			nonZero = true
			break
		}
	}
	if !nonZero {
		return nil, fmt.Errorf("kernel: zero vector")
	}
	return out, nil
}

func lcmPositiveBig(a, b *big.Int) *big.Int {
	if a.Sign() == 0 {
		return new(big.Int).Set(b)
	}
	if b.Sign() == 0 {
		return new(big.Int).Set(a)
	}
	g := new(big.Int).GCD(nil, nil, new(big.Int).Set(a), new(big.Int).Set(b))
	return new(big.Int).Div(new(big.Int).Mul(a, b), g)
}
