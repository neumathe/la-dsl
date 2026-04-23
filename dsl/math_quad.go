package dsl

import (
	"fmt"
	"math/big"
)

// SymmetricDefClass 对称矩阵的 Sylvester 定号分类（3 阶）。
type SymmetricDefClass int

const (
	SymDefPD SymmetricDefClass = iota // 正定
	SymDefND                          // 负定
	SymDefInd                         // 不定（或半正定/半负定等非上述两类）
)

func principalMinorDet(m *MatrixInt, k int) *big.Int {
	if k < 1 || k > m.R || k > m.C {
		return big.NewInt(0)
	}
	sub := NewMatrixInt(k, k)
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			sub.A[i][j] = m.A[i][j]
		}
	}
	return BareissDet(sub)
}

// ClassifySymmetric3 判定 3×3 整数对称矩阵的定号类。
func ClassifySymmetric3(S *MatrixInt) (SymmetricDefClass, error) {
	if S.R != 3 || S.C != 3 {
		return SymDefInd, fmt.Errorf("classify: need 3×3")
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if S.A[i][j] != S.A[j][i] {
				return SymDefInd, fmt.Errorf("classify: matrix not symmetric")
			}
		}
	}
	d1 := principalMinorDet(S, 1)
	d2 := principalMinorDet(S, 2)
	d3 := principalMinorDet(S, 3)
	if d1.Sign() > 0 && d2.Sign() > 0 && d3.Sign() > 0 {
		return SymDefPD, nil
	}
	if d1.Sign() < 0 && d2.Sign() > 0 && d3.Sign() < 0 {
		return SymDefND, nil
	}
	return SymDefInd, nil
}
