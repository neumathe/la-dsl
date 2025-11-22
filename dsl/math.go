package dsl

import (
	"errors"
	"math/big"
)

// BareissDet 使用 Bareiss 算法计算整型矩阵的行列式（返回 big.Int）
func BareissDet(m *MatrixInt) *big.Int {
	n := m.R
	if n != m.C {
		return big.NewInt(0)
	}
	A := make([][]*big.Int, n)
	for i := 0; i < n; i++ {
		A[i] = make([]*big.Int, n)
		for j := 0; j < n; j++ {
			A[i][j] = big.NewInt(m.A[i][j])
		}
	}
	prev := big.NewInt(1)
	for k := 0; k < n-1; k++ {
		if A[k][k].Sign() == 0 {
			swapIdx := -1
			for r := k + 1; r < n; r++ {
				if A[r][k].Sign() != 0 {
					swapIdx = r
					break
				}
			}
			if swapIdx == -1 {
				return big.NewInt(0)
			}
			A[k], A[swapIdx] = A[swapIdx], A[k]
		}
		for i := k + 1; i < n; i++ {
			for j := k + 1; j < n; j++ {
				t1 := new(big.Int).Mul(A[i][j], A[k][k])
				t2 := new(big.Int).Mul(A[i][k], A[k][j])
				num := t1.Sub(t1, t2)
				if prev.Sign() != 0 {
					num.Quo(num, prev)
				}
				A[i][j] = num
			}
		}
		prev = new(big.Int).Set(A[k][k])
	}
	return new(big.Int).Set(A[n-1][n-1])
}

// matrixRankRat 使用 big.Rat 做高斯消元计算矩阵秩
func matrixRankRat(m *MatrixInt) int {
	r := m.R
	c := m.C
	M := make([][]*big.Rat, r)
	for i := 0; i < r; i++ {
		M[i] = make([]*big.Rat, c)
		for j := 0; j < c; j++ {
			M[i][j] = new(big.Rat).SetInt64(m.A[i][j])
		}
	}
	rank := 0
	row := 0
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
		row++
		rank++
	}
	return rank
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// columnPivotsRat 返回矩阵按列做高斯消元时的主元列下标（1-based）
// 用于构造一组由列向量组成的基。
func columnPivotsRat(m *MatrixInt) []int {
	r := m.R
	c := m.C
	M := make([][]*big.Rat, r)
	for i := 0; i < r; i++ {
		M[i] = make([]*big.Rat, c)
		for j := 0; j < c; j++ {
			M[i][j] = new(big.Rat).SetInt64(m.A[i][j])
		}
	}
	pivots := []int{}
	row := 0
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
		pivots = append(pivots, col+1) // 1-based 列下标
		row++
	}
	return pivots
}

// matrixMulInt 计算两个整型矩阵的乘积（a*b）
func matrixMulInt(a, b *MatrixInt) (*MatrixInt, error) {
	if a.C != b.R {
		return nil, errors.New("matrix dimension mismatch in mul")
	}
	res := NewMatrixInt(a.R, b.C)
	for i := 0; i < a.R; i++ {
		for j := 0; j < b.C; j++ {
			var s int64
			for k := 0; k < a.C; k++ {
				s += a.A[i][k] * b.A[k][j]
			}
			res.A[i][j] = s
		}
	}
	return res, nil
}

// matrixPowInt 计算整型矩阵的幂 A^n，n>=0
func matrixPowInt(a *MatrixInt, n int64) (*MatrixInt, error) {
	if a.R != a.C {
		return nil, errors.New("matrix power expects square matrix")
	}
	if n < 0 {
		return nil, errors.New("matrix power with negative exponent is not supported")
	}
	// 单位矩阵
	res := NewMatrixInt(a.R, a.C)
	for i := 0; i < a.R; i++ {
		res.A[i][i] = 1
	}
	if n == 0 {
		return res, nil
	}
	base := NewMatrixInt(a.R, a.C)
	for i := 0; i < a.R; i++ {
		for j := 0; j < a.C; j++ {
			base.A[i][j] = a.A[i][j]
		}
	}
	exp := n
	for exp > 0 {
		if exp&1 == 1 {
			tmp, err := matrixMulInt(res, base)
			if err != nil {
				return nil, err
			}
			res = tmp
		}
		if exp > 1 {
			tmp, err := matrixMulInt(base, base)
			if err != nil {
				return nil, err
			}
			base = tmp
		}
		exp >>= 1
	}
	return res, nil
}
