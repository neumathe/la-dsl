package dsl

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
)

// vectorsRationalCollinear 判断 u 与 v 是否在 Q^d 上共线（允许零向量全零匹配）。
func vectorsRationalCollinear(exp, got []*big.Rat) (bool, string) {
	if len(exp) != len(got) {
		return false, "dim mismatch"
	}
	d := len(exp)
	allE0, allG0 := true, true
	for i := 0; i < d; i++ {
		if exp[i] == nil || got[i] == nil {
			return false, "nil component"
		}
		if exp[i].Sign() != 0 {
			allE0 = false
		}
		if got[i].Sign() != 0 {
			allG0 = false
		}
	}
	if allE0 && allG0 {
		return true, ""
	}
	if allE0 != allG0 {
		return false, "zero pattern mismatch"
	}
	pivot := -1
	for i := 0; i < d; i++ {
		if exp[i].Sign() != 0 {
			pivot = i
			break
		}
	}
	if pivot < 0 {
		return false, "no pivot"
	}
	lam := new(big.Rat).Quo(got[pivot], exp[pivot])
	for i := 0; i < d; i++ {
		t := new(big.Rat).Mul(lam, exp[i])
		if t.Cmp(got[i]) != 0 {
			return false, "not collinear"
		}
	}
	return true, ""
}

func interfaceSliceToRatVec(vals []interface{}) ([]*big.Rat, error) {
	out := make([]*big.Rat, len(vals))
	for i, v := range vals {
		r, err := toRat(v)
		if err != nil {
			return nil, err
		}
		out[i] = r
	}
	return out, nil
}

func parseUserRatVector(subs []string) ([]*big.Rat, error) {
	out := make([]*big.Rat, len(subs))
	for i, s := range subs {
		r, err := ParseUserRational(s)
		if err != nil {
			return nil, fmt.Errorf("comp %d: %w", i, err)
		}
		out[i] = r
	}
	return out, nil
}

// matrixIntTimesVectorRat 计算 A*x（列向量），A 整数，x 有理。
func matrixIntTimesVectorRat(A *MatrixInt, x []*big.Rat) ([]*big.Rat, error) {
	if A.C != len(x) {
		return nil, fmt.Errorf("matvec dim")
	}
	out := make([]*big.Rat, A.R)
	for i := 0; i < A.R; i++ {
		s := big.NewRat(0, 1)
		for j := 0; j < A.C; j++ {
			t := new(big.Rat).SetInt64(A.A[i][j])
			t.Mul(t, x[j])
			s.Add(s, t)
		}
		out[i] = s
	}
	return out, nil
}

func vectorIntToRat(b *VectorInt) []*big.Rat {
	out := make([]*big.Rat, b.N)
	for i := 0; i < b.N; i++ {
		out[i] = big.NewRat(b.V[i], 1)
	}
	return out
}

// affineSolutionOK 验证 A*x == b（有理向量 x）。
func affineSolutionOK(A *MatrixInt, b *VectorInt, x []*big.Rat) (bool, string) {
	if A.R != b.N {
		return false, "A rows != b len"
	}
	y, err := matrixIntTimesVectorRat(A, x)
	if err != nil {
		return false, err.Error()
	}
	bb := vectorIntToRat(b)
	for i := 0; i < A.R; i++ {
		if y[i].Cmp(bb[i]) != 0 {
			return false, "Ax!=b"
		}
	}
	return true, ""
}

// columnSubsetRank 取 A 的列子集（1-based下标），返回子矩阵秩。
func columnSubsetRank(A *MatrixInt, cols []int) int {
	if len(cols) == 0 {
		return 0
	}
	sub := NewMatrixInt(A.R, len(cols))
	for j, c := range cols {
		if c < 1 || c > A.C {
			return -1
		}
		for i := 0; i < A.R; i++ {
			sub.A[i][j] = A.A[i][c-1]
		}
	}
	return matrixRankRat(sub)
}

// sortedBasisColumnsOK 用户填入若干列下标（0 表示空），去重排序后应张成与 A 同秩的列空间。
func sortedBasisColumnsOK(A *MatrixInt, ncols int, user []string) (bool, string) {
	rA := matrixRankRat(A)
	var idx []int
	for _, s := range user {
		s = NormalizeUserAnswer(s)
		if s == "" {
			continue
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return false, "not int"
		}
		if int(v) > ncols {
			return false, "col out of range"
		}
		if v <= 0 {
			// 0 或负表示「不填」，与 HTML 多余空一致
			continue
		}
		idx = append(idx, int(v))
	}
	sort.Ints(idx)
	uniq := make([]int, 0, len(idx))
	for i := 0; i < len(idx); i++ {
		if i > 0 && idx[i] == idx[i-1] {
			continue
		}
		uniq = append(uniq, idx[i])
	}
	if len(uniq) != rA {
		return false, "wrong count"
	}
	if columnSubsetRank(A, uniq) != rA {
		return false, "not basis"
	}
	return true, ""
}
