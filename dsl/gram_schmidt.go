package dsl

import (
	"fmt"
	"math/big"
)

// GramSchmidtColsOrthogRat 对矩阵 V 的各列（视为 R^{V.R} 中向量）做 Gram–Schmidt 正交化（不单位化），
// 返回 C 个长度均为 V.R 的有理列向量（列下标 1..C）。
func GramSchmidtColsOrthogRat(V *MatrixInt) ([][]*big.Rat, error) {
	return GramSchmidtColsOrthogRatGeneral(V)
}

// GramSchmidtColsOrthogRatGeneral 与 GramSchmidtColsOrthogRat 相同，显式命名供内部使用。
func GramSchmidtColsOrthogRatGeneral(V *MatrixInt) ([][]*big.Rat, error) {
	if V == nil || V.R < 1 || V.C < 1 {
		return nil, fmt.Errorf("gs: bad matrix dimensions")
	}
	r, c := V.R, V.C
	cols := make([][]*big.Rat, c)
	for j := 0; j < c; j++ {
		v := make([]*big.Rat, r)
		for i := 0; i < r; i++ {
			v[i] = new(big.Rat).SetInt64(V.A[i][j])
		}
		cols[j] = v
	}
	dot := func(a, b []*big.Rat) *big.Rat {
		s := big.NewRat(0, 1)
		for i := 0; i < r; i++ {
			t := new(big.Rat).Mul(a[i], b[i])
			s.Add(s, t)
		}
		return s
	}
	scale := func(k *big.Rat, v []*big.Rat) []*big.Rat {
		out := make([]*big.Rat, r)
		for i := 0; i < r; i++ {
			out[i] = new(big.Rat).Mul(k, v[i])
		}
		return out
	}
	sub := func(a, b []*big.Rat) []*big.Rat {
		out := make([]*big.Rat, r)
		for i := 0; i < r; i++ {
			out[i] = new(big.Rat).Sub(a[i], b[i])
		}
		return out
	}

	u := make([][]*big.Rat, c)
	u[0] = cols[0]
	for j := 1; j < c; j++ {
		w := cols[j]
		for k := 0; k < j; k++ {
			uk := u[k]
			num := dot(w, uk)
			den := dot(uk, uk)
			if den.Sign() == 0 {
				return nil, fmt.Errorf("gs: degenerate uk at step %d", k)
			}
			coef := new(big.Rat).Quo(num, den)
			w = sub(w, scale(coef, uk))
		}
		u[j] = w
	}
	return u, nil
}
