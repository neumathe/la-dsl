package dsl

import (
	"fmt"
	"math/big"
)

func forEachCombination(n, k int, fn func([]int) bool) {
	if k < 0 || k > n {
		return
	}
	if k == 0 {
		fn([]int{})
		return
	}
	idx := make([]int, k)
	for i := 0; i < k; i++ {
		idx[i] = i
	}
	for {
		cp := append([]int(nil), idx...)
		if fn(cp) {
			return
		}
		t := k - 1
		for t >= 0 && idx[t] == n-k+t {
			t--
		}
		if t < 0 {
			break
		}
		idx[t]++
		for j := t + 1; j < k; j++ {
			idx[j] = idx[j-1] + 1
		}
	}
}

// SolveCoeffCols 设 V 为 n×m 整数矩阵，cols 为互不相同的列下标（1-based），rhs 为另一列下标。
// 若 v_rhs 落在 cols 张成的子空间中，返回系数 x（len(cols)），使 sum_j x[j]*V[:,cols[j]] = V[:,rhs]。
func SolveCoeffCols(V *MatrixInt, cols []int, rhs int) ([]*big.Rat, error) {
	if V == nil {
		return nil, fmt.Errorf("nil V")
	}
	r, m := V.R, V.C
	k := len(cols)
	if k == 0 || rhs < 1 || rhs > m {
		return nil, fmt.Errorf("bad dims")
	}
	seen := map[int]struct{}{}
	for _, c := range cols {
		if c < 1 || c > m {
			return nil, fmt.Errorf("col out of range")
		}
		if _, ok := seen[c]; ok {
			return nil, fmt.Errorf("duplicate col")
		}
		seen[c] = struct{}{}
	}
	bvec := NewVectorInt(r)
	for i := 0; i < r; i++ {
		bvec.V[i] = V.A[i][rhs-1]
	}
	if k > r {
		return nil, fmt.Errorf("k>r")
	}
	var sol []*big.Rat
	var lastErr error
	forEachCombination(r, k, func(rows []int) bool {
		M := NewMatrixInt(k, k)
		b := NewVectorInt(k)
		for i := 0; i < k; i++ {
			ri := rows[i]
			for j := 0; j < k; j++ {
				M.A[i][j] = V.A[ri][cols[j]-1]
			}
			b.V[i] = bvec.V[ri]
		}
		if BareissDet(M).Sign() == 0 {
			return false
		}
		x, err := solveLinearSystemRat(M, b)
		if err != nil {
			lastErr = err
			return false
		}
		if verifyColComb(V, cols, rhs, x) {
			sol = x
			return true
		}
		return false
	})
	if sol == nil {
		if lastErr != nil {
			return nil, lastErr
		}
		return nil, fmt.Errorf("no sol for cols/re rhs")
	}
	return sol, nil
}

func verifyColComb(V *MatrixInt, cols []int, rhs int, x []*big.Rat) bool {
	r := V.R
	for i := 0; i < r; i++ {
		s := big.NewRat(0, 1)
		for j := 0; j < len(cols); j++ {
			t := new(big.Rat).SetInt64(V.A[i][cols[j]-1])
			t.Mul(t, x[j])
			s.Add(s, t)
		}
		rhsv := big.NewRat(V.A[i][rhs-1], 1)
		if s.Cmp(rhsv) != 0 {
			return false
		}
	}
	return true
}
