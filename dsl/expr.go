package dsl

import (
	"fmt"
	"math/big"
	"strings"
)

// evaluateExpression 支持 det(A)、solve(A,b)、cofactor(A,i,j)、col(A,1)、A * x 等基础表达式
// EvaluateExpression 对外导出，方便上层做表达式求值（例如 demo / 测试）
func EvaluateExpression(expr string, inst *Instance) (interface{}, error) {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "det(") && strings.HasSuffix(expr, ")") {
		arg := insideParens(expr)
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("det expects matrix, got %T", v)
		}
		return BareissDet(m), nil
	}

	if strings.HasPrefix(expr, "rank(") && strings.HasSuffix(expr, ")") {
		arg := insideParens(expr)
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("rank expects matrix, got %T", v)
		}
		return int64(matrixRankRat(m)), nil
	}

	if strings.HasPrefix(expr, "solve(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("solve expects 2 args")
		}
		vA, ok := inst.Vars[strings.TrimSpace(parts[0])]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", parts[0])
		}
		vB, ok := inst.Vars[strings.TrimSpace(parts[1])]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", parts[1])
		}
		A, ok := vA.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("solve expects matrix A")
		}
		bvec, ok := vB.(*VectorInt)
		if !ok {
			return nil, fmt.Errorf("solve expects vector b")
		}
		xrat, err := solveLinearSystemRat(A, bvec)
		if err != nil {
			return nil, err
		}
		return xrat, nil
	}

	if strings.HasPrefix(expr, "cofactor(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("cofactor expects 3 args")
		}
		Aname := strings.TrimSpace(parts[0])
		i := mustAtoi(strings.TrimSpace(parts[1]))
		j := mustAtoi(strings.TrimSpace(parts[2]))
		v, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", Aname)
		}
		A, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("cofactor expects matrix")
		}
		if i < 1 || j < 1 || i > A.R || j > A.C {
			return nil, fmt.Errorf("index out of range")
		}
		minor := NewMatrixInt(A.R-1, A.C-1)
		ri := 0
		for r := 0; r < A.R; r++ {
			if r == i-1 {
				continue
			}
			ci := 0
			for c := 0; c < A.C; c++ {
				if c == j-1 {
					continue
				}
				minor.A[ri][ci] = A.A[r][c]
				ci++
			}
			ri++
		}
		d := BareissDet(minor)
		sign := int64(1)
		if ((i + j) % 2) == 1 {
			sign = -1
		}
		return new(big.Int).Mul(big.NewInt(sign), d), nil
	}

	if strings.HasPrefix(expr, "basis_cols(") && strings.HasSuffix(expr, ")") {
		arg := insideParens(expr)
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("basis_cols expects matrix, got %T", v)
		}
		pivots := columnPivotsRat(m)
		vec := NewVectorInt(m.C)
		for i := 0; i < len(pivots) && i < m.C; i++ {
			vec.V[i] = int64(pivots[i])
		}
		// 多余位置保持为 0，表示“无”
		return vec, nil
	}

	if strings.HasPrefix(expr, "col(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("col expects 2 args: col(A, 1)")
		}
		Aname := strings.TrimSpace(parts[0])
		k := mustAtoi(strings.TrimSpace(parts[1]))
		v, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", Aname)
		}
		A, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("col expects matrix")
		}
		if k < 1 || k > A.C {
			return nil, fmt.Errorf("col index out of range")
		}
		vec := NewVectorInt(A.R)
		for i := 0; i < A.R; i++ {
			vec.V[i] = A.A[i][k-1]
		}
		return vec, nil
	}

	if strings.Contains(expr, "*") {
		parts := strings.Split(expr, "*")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			vL, lok := inst.Vars[left]
			vR, rok := inst.Vars[right]
			if lok && rok {
				if m, ok := vL.(*MatrixInt); ok {
					if x, ok2 := vR.(*VectorInt); ok2 {
						if m.C != x.N {
							return nil, fmt.Errorf("dimension mismatch in %s * %s", left, right)
						}
						out := NewVectorInt(m.R)
						for i := 0; i < m.R; i++ {
							var sum int64
							for j := 0; j < m.C; j++ {
								sum += m.A[i][j] * x.V[j]
							}
							out.V[i] = sum
						}
						return out, nil
					}
				}
			}
		}
	}

	// 支持简单的下标访问：var[index]，如 x[1]、coord[3]
	if strings.Contains(expr, "[") && strings.HasSuffix(expr, "]") {
		idxOpen := strings.Index(expr, "[")
		idxClose := strings.LastIndex(expr, "]")
		if idxOpen > 0 && idxClose > idxOpen {
			vname := strings.TrimSpace(expr[:idxOpen])
			idx := mustAtoi(expr[idxOpen+1 : idxClose])
			if vv, ok := inst.Vars[vname]; ok {
				switch t := vv.(type) {
				case *VectorInt:
					if idx >= 1 && idx <= t.N {
						return t.V[idx-1], nil
					}
				case []*big.Rat:
					if idx >= 1 && idx <= len(t) {
						return t[idx-1], nil
					}
				}
			}
		}
	}

	if v, ok := inst.Vars[expr]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("unsupported expression: %s", expr)
}

func insideParens(s string) string {
	i := strings.Index(s, "(")
	j := strings.LastIndex(s, ")")
	if i == -1 || j == -1 || j <= i {
		return ""
	}
	return s[i+1 : j]
}

func splitArgs(s string) []string {
	parts := []string{}
	cur := ""
	depth := 0
	for _, ch := range s {
		switch ch {
		case '(':
			depth++
			cur += string(ch)
		case ')':
			depth--
			cur += string(ch)
		case ',':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(cur))
				cur = ""
			} else {
				cur += string(ch)
			}
		default:
			cur += string(ch)
		}
	}
	if strings.TrimSpace(cur) != "" {
		parts = append(parts, strings.TrimSpace(cur))
	}
	return parts
}

func mustAtoi(s string) int {
	s = strings.TrimSpace(s)
	var v int
	fmt.Sscanf(s, "%d", &v)
	return v
}
