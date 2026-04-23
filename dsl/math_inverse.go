package dsl

import (
	"fmt"
)

// MatrixInverseInt 求整数矩阵的逆（若逆矩阵非整数则报错）。
func MatrixInverseInt(A *MatrixInt) (*MatrixInt, error) {
	n := A.R
	if n != A.C {
		return nil, fmt.Errorf("inverse: need square matrix")
	}
	inv := NewMatrixInt(n, n)
	for j := 0; j < n; j++ {
		ej := NewVectorInt(n)
		ej.V[j] = 1
		xrat, err := solveLinearSystemRat(A, ej)
		if err != nil {
			return nil, fmt.Errorf("inverse col %d: %w", j, err)
		}
		for i := 0; i < n; i++ {
			if !xrat[i].IsInt() {
				return nil, fmt.Errorf("inverse: non-integer entry at (%d,%d)", i+1, j+1)
			}
			inv.A[i][j] = xrat[i].Num().Int64()
		}
	}
	return inv, nil
}
