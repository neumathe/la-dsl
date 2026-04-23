package dsl

import (
	"fmt"
	"math/big"
	"strings"
)

// RenderTitle 根据 Problem.Title 和 RenderInst 的结果输出完整题面
// Title 中使用 {{key}} 占位符，key 来自 Problem.Render 的键。
func RenderTitle(p Problem, inst *Instance) (string, error) {
	if p.Title == "" {
		return "", nil
	}
	rendered, err := RenderInst(p, inst)
	if err != nil {
		return "", err
	}
	title := p.Title
	for k, v := range rendered {
		ph := "{{" + k + "}}"
		title = strings.ReplaceAll(title, ph, FormatValueForTitle(v))
	}
	return title, nil
}

// FormatValueForTitle 将内部值转成适合题干中的字符串（如 LaTeX bmatrix）。
func FormatValueForTitle(v interface{}) string {
	switch t := v.(type) {
	case *MatrixInt:
		// LaTeX 矩阵（列向量题常用），方便直接嵌入 $...$ 中
		var rows []string
		for i := 0; i < t.R; i++ {
			var cols []string
			for j := 0; j < t.C; j++ {
				cols = append(cols, fmt.Sprintf("%d", t.A[i][j]))
			}
			rows = append(rows, strings.Join(cols, "&"))
		}
		return `\begin{bmatrix}` + strings.Join(rows, `\\`) + `\end{bmatrix}`
	case *VectorInt:
		// 列向量形式
		var rows []string
		for i := 0; i < t.N; i++ {
			rows = append(rows, fmt.Sprintf("%d", t.V[i]))
		}
		return `\begin{bmatrix}` + strings.Join(rows, `\\`) + `\end{bmatrix}`
	case *big.Int:
		return t.String()
	case string:
		return t // LaTeX string directly
	case int64, int, float64:
		return fmt.Sprintf("%v", t)
	default:
		return fmt.Sprintf("%v", v)
	}
}




func formatLambdaCasesTitle(A *MatrixInt, paramRow int64, paramCol int64, constC int64) string {
	n := A.R
	m := A.C
	var rows []string
	for i := 0; i < n; i++ {
		var parts []string
		for j := 0; j < m; j++ {
			isFirst := len(parts) == 0
			xv := fmt.Sprintf("x_{%d}", j+1)
			if int64(i+1) == paramRow && int64(j+1) == paramCol {
				// This is the λ position: render as λ + constC
				if constC == 0 {
					if isFirst {
						parts = append(parts, "\\lambda "+xv)
					} else {
						parts = append(parts, "+\\lambda "+xv)
					}
				} else if constC > 0 {
					if isFirst {
						parts = append(parts, fmt.Sprintf("\\lambda+%d%s", constC, xv))
					} else {
						parts = append(parts, fmt.Sprintf("+\\lambda+%d%s", constC, xv))
					}
				} else {
					// constC < 0, render as λ-k
					if isFirst {
						parts = append(parts, fmt.Sprintf("\\lambda%d%s", constC, xv))
					} else {
						parts = append(parts, fmt.Sprintf("+\\lambda%d%s", constC, xv))
					}
				}
				continue
			}
			coef := A.A[i][j]
			if coef == 0 {
				continue
			}
			if isFirst {
				if coef == 1 {
					parts = append(parts, xv)
				} else if coef == -1 {
					parts = append(parts, "-" + xv)
				} else {
					parts = append(parts, fmt.Sprintf("%d%s", coef, xv))
				}
			} else {
				if coef == 1 {
					parts = append(parts, "+" + xv)
				} else if coef == -1 {
					parts = append(parts, "-" + xv)
				} else if coef > 0 {
					parts = append(parts, fmt.Sprintf("+%d%s", coef, xv))
				} else {
					parts = append(parts, fmt.Sprintf("%d%s", coef, xv))
				}
			}
		}
		if len(parts) == 0 {
			parts = append(parts, "0")
		}
		parts = append(parts, "=0")
		rows = append(rows, strings.Join(parts, ""))
	}
	return `\begin{cases}` + strings.Join(rows, `\\`) + `\end{cases}`
}

func formatVmatrixTitle(m *MatrixInt) string {
	var rows []string
	for i := 0; i < m.R; i++ {
		var cols []string
		for j := 0; j < m.C; j++ {
			cols = append(cols, fmt.Sprintf("%d", m.A[i][j]))
		}
		rows = append(rows, strings.Join(cols, "&"))
	}
	return `\begin{vmatrix}` + strings.Join(rows, `\\`) + `\end{vmatrix}`
}

func formatEquidiagonalTitle(m *MatrixInt) string {
	n := m.R
	if n <= 4 {
		// Small enough to render fully
		return formatVmatrixTitle(m)
	}
	// For large equidiagonal matrices, use \cdots format
	// Show first 2 and last row/column, with \cdots in between
	a := m.A[0][0] // diagonal value
	b := m.A[0][1] // off-diagonal value
	var rows []string
	// Row 1: a, b, ..., b
	row1 := []string{fmt.Sprintf("%d", a), fmt.Sprintf("%d", b)}
	if n > 3 {
		row1 = append(row1, `\cdots`)
	}
	row1 = append(row1, fmt.Sprintf("%d", b))
	rows = append(rows, strings.Join(row1, "&"))
	// Row 2: b, a, ..., b
	row2 := []string{fmt.Sprintf("%d", b), fmt.Sprintf("%d", a)}
	if n > 3 {
		row2 = append(row2, `\cdots`)
	}
	row2 = append(row2, fmt.Sprintf("%d", b))
	rows = append(rows, strings.Join(row2, "&"))
	// \cdots row
	if n > 3 {
		cdots_row := []string{`\cdots`, `\cdots`, `\cdots`, `\cdots`}
		rows = append(rows, strings.Join(cdots_row, "&"))
	}
	// Last row: b, b, ..., a
	rowLast := []string{fmt.Sprintf("%d", b), fmt.Sprintf("%d", b)}
	if n > 3 {
		rowLast = append(rowLast, `\cdots`)
	}
	rowLast = append(rowLast, fmt.Sprintf("%d", a))
	rows = append(rows, strings.Join(rowLast, "&"))
	return `\begin{vmatrix}` + strings.Join(rows, `\\`) + `\end{vmatrix}`
}

func formatCasesTitle(A *MatrixInt, b *VectorInt) string {
	n := A.R // number of equations
	m := A.C // number of variables
	var rows []string
	for i := 0; i < n; i++ {
		var parts []string
		for j := 0; j < m; j++ {
			coef := A.A[i][j]
			xv := fmt.Sprintf("x_{%d}", j+1)
			if coef == 0 {
				continue
			}
			isFirst := len(parts) == 0
			if isFirst {
				if coef == 1 {
					parts = append(parts, xv)
				} else if coef == -1 {
					parts = append(parts, "-" + xv)
				} else {
					parts = append(parts, fmt.Sprintf("%d%s", coef, xv))
				}
			} else {
				if coef == 1 {
					parts = append(parts, "+" + xv)
				} else if coef == -1 {
					parts = append(parts, "-" + xv)
				} else if coef > 0 {
					parts = append(parts, fmt.Sprintf("+%d%s", coef, xv))
				} else {
					parts = append(parts, fmt.Sprintf("%d%s", coef, xv))
				}
			}
		}
		if len(parts) == 0 {
			parts = append(parts, "0")
		}
		parts = append(parts, fmt.Sprintf("=%d", b.V[i]))
		rows = append(rows, strings.Join(parts, ""))
	}
	return `\begin{cases}` + strings.Join(rows, `\\`) + `\end{cases}`
}

func formatQuadraticExpr(S *MatrixInt) string {
	n := S.R
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if S.A[i][j] != S.A[j][i] {
				return "non-symmetric"
			}
		}
	}

	xi := func(k int) string {
		return fmt.Sprintf("x_{%d}", k)
	}

	var terms []string
	for i := 0; i < n; i++ {
		v := S.A[i][i]
		if v == 0 {
			continue
		}
		if v == 1 {
			terms = append(terms, xi(i+1)+"^2")
		} else if v == -1 {
			terms = append(terms, "-"+xi(i+1)+"^2")
		} else {
			terms = append(terms, fmt.Sprintf("%d%s^2", v, xi(i+1)))
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			coef := 2 * S.A[i][j]
			if coef == 0 {
				continue
			}
			if coef == 1 {
				terms = append(terms, xi(i+1)+xi(j+1))
			} else if coef == -1 {
				terms = append(terms, "-"+xi(i+1)+xi(j+1))
			} else {
				terms = append(terms, fmt.Sprintf("%d%s%s", coef, xi(i+1), xi(j+1)))
			}
		}
	}

	if len(terms) == 0 {
		return "0"
	}
	result := terms[0]
	for i := 1; i < len(terms); i++ {
		if strings.HasPrefix(terms[i], "-") {
			result += terms[i]
		} else {
			result += "+" + terms[i]
		}
	}
	return result
}

// evalParamSystemTitle 生成含参方程组的 LaTeX 题面字符串。
// 格式：\begin{cases} a₁x₁ + a₂x₂ + a₃x₃ = b₁ \\ ... \\ -4x₁+15x₂+\\lambda x₃=\mu \end{cases}
func evalParamSystemTitle(inst *Instance) (interface{}, error) {
	Av, ok := inst.Vars["_param_A"]
	if !ok {
		return nil, fmt.Errorf("param_system_title: A not found")
	}
	A, ok := Av.(*MatrixInt)
	if !ok {
		return nil, fmt.Errorf("param_system_title: A not matrix")
	}
	bv, ok := inst.Vars["_param_b"]
	if !ok {
		return nil, fmt.Errorf("param_system_title: b not found")
	}
	b, ok := bv.(*VectorInt)
	if !ok {
		return nil, fmt.Errorf("param_system_title: b not vector")
	}
	colLambdaV, ok := inst.Vars["_param_colLambda"]
	if !ok {
		return nil, fmt.Errorf("param_system_title: colLambda not found")
	}
	colLambda := int64(0)
	switch t := colLambdaV.(type) {
	case int64:
		colLambda = t
	case int:
		colLambda = int64(t)
	}

	formatCoeff := func(val int64, varIdx int, isFirst bool) string {
		xv := fmt.Sprintf("x_{%d}", varIdx)
		if val == 0 {
			return ""
		}
		if isFirst {
			if val == 1 {
				return xv
			}
			if val == -1 {
				return "-" + xv
			}
			return fmt.Sprintf("%d%s", val, xv)
		}
		if val == 1 {
			return "+" + xv
		}
		if val == -1 {
			return "-" + xv
		}
		return fmt.Sprintf("+%d%s", val, xv)
	}

	var rows []string
	for i := 0; i < 3; i++ {
		var parts []string
		for j := 0; j < 3; j++ {
			if i == 2 && j+1 == int(colLambda) {
				// This position uses λ instead of a number
				if len(parts) == 0 {
					parts = append(parts, fmt.Sprintf("\\lambda %s", fmt.Sprintf("x_{%d}", j+1)))
				} else {
					parts = append(parts, fmt.Sprintf("+\\lambda %s", fmt.Sprintf("x_{%d}", j+1)))
				}
			} else {
				s := formatCoeff(A.A[i][j], j+1, len(parts) == 0)
				if s != "" {
					parts = append(parts, s)
				}
			}
		}
		if i == 2 {
			// b[3] uses μ
			parts = append(parts, fmt.Sprintf("=%s", "\\mu"))
		} else {
			parts = append(parts, fmt.Sprintf("=%d", b.V[i]))
		}
		rows = append(rows, strings.Join(parts, ""))
	}

	return fmt.Sprintf("\\begin{cases}%s\\end{cases}", strings.Join(rows, `\\`)), nil
}

// formatPolynomialFromVec 将向量 [a,b,c] 渲染为多项式 a+bx+cx² 的 LaTeX 字符串
func formatPolynomialFromVec(v *VectorInt) string {
	var terms []string
	// k=0: constant, k=1: x coefficient, k=2: x² coefficient
	for k := 0; k < v.N && k < 3; k++ {
		val := v.V[k]
		if val == 0 {
			continue
		}
		isFirst := len(terms) == 0
		switch k {
		case 0: // constant
			if isFirst {
				terms = append(terms, fmt.Sprintf("%d", val))
			} else {
				terms = append(terms, fmt.Sprintf("+%d", val))
			}
		case 1: // x
			if isFirst {
				if val == 1 {
					terms = append(terms, "x")
				} else if val == -1 {
					terms = append(terms, "-x")
				} else {
					terms = append(terms, fmt.Sprintf("%dx", val))
				}
			} else {
				if val == 1 {
					terms = append(terms, "+x")
				} else if val == -1 {
					terms = append(terms, "-x")
				} else {
					terms = append(terms, fmt.Sprintf("+%dx", val))
				}
			}
		case 2: // x²
			if isFirst {
				if val == 1 {
					terms = append(terms, "x^2")
				} else if val == -1 {
					terms = append(terms, "-x^2")
				} else {
					terms = append(terms, fmt.Sprintf("%dx^2", val))
				}
			} else {
				if val == 1 {
					terms = append(terms, "+x^2")
				} else if val == -1 {
					terms = append(terms, "-x^2")
				} else {
					terms = append(terms, fmt.Sprintf("+%dx^2", val))
				}
			}
		}
	}
	if len(terms) == 0 {
		return "0"
	}
	return strings.Join(terms, "")
}

// formatLinearTransformTitle renders matrix A₀ as T(x₁,x₂,x₃)ᵀ=(...)ᵀ
// Example: T\begin{pmatrix}x_1\\x_2\\x_3\end{pmatrix}=\begin{pmatrix}2x_1+4x_3\\-3x_1+4x_2+4x_3\\...\end{pmatrix}
func formatLinearTransformTitle(m *MatrixInt) string {
	n := m.R
	xi := func(k int) string { return fmt.Sprintf("x_{%d}", k) }
	var rowExprs []string
	for i := 0; i < n; i++ {
		var parts []string
		for j := 0; j < m.C; j++ {
		 coef := m.A[i][j]
			if coef == 0 {
				continue
			}
			isFirst := len(parts) == 0
			if isFirst {
				if coef == 1 {
					parts = append(parts, xi(j+1))
				} else if coef == -1 {
					parts = append(parts, "-"+xi(j+1))
				} else {
					parts = append(parts, fmt.Sprintf("%d%s", coef, xi(j+1)))
				}
			} else {
				if coef == 1 {
					parts = append(parts, "+"+xi(j+1))
				} else if coef == -1 {
					parts = append(parts, "-"+xi(j+1))
				} else if coef > 0 {
					parts = append(parts, fmt.Sprintf("+%d%s", coef, xi(j+1)))
				} else {
					parts = append(parts, fmt.Sprintf("%d%s", coef, xi(j+1)))
				}
			}
		}
		if len(parts) == 0 {
			parts = append(parts, "0")
		}
		rowExprs = append(rowExprs, strings.Join(parts, ""))
	}
	left := `T\begin{pmatrix}x_1\\x_2\\x_3\end{pmatrix}`
	right := `\begin{pmatrix}` + strings.Join(rowExprs, `\\`) + `\end{pmatrix}`
	return left + "=" + right
}

// formatBasisLinearComboTitle renders the columns of matrix B as new basis vectors
// expressed as linear combinations of ε₁, ε₂, ε₃.
// Column j is rendered as: c_{1j}ε₁ + c_{2j}ε₂ + c_{3j}ε₃
// First column (which should be ε₁ for upper_unit) is just ε₁.
func formatBasisLinearComboTitle(m *MatrixInt) string {
	n := m.R
	eps := func(k int) string { return fmt.Sprintf("\\varepsilon_{%d}", k) }
	var vecExprs []string
	for j := 0; j < m.C; j++ {
		var parts []string
		for i := 0; i < n; i++ {
		 coef := m.A[i][j]
			if coef == 0 {
				continue
			}
			isFirst := len(parts) == 0
			if isFirst {
				if coef == 1 {
					parts = append(parts, eps(i+1))
				} else if coef == -1 {
					parts = append(parts, "-"+eps(i+1))
				} else {
					parts = append(parts, fmt.Sprintf("%d%s", coef, eps(i+1)))
				}
			} else {
				if coef == 1 {
					parts = append(parts, "+"+eps(i+1))
				} else if coef == -1 {
					parts = append(parts, "-"+eps(i+1))
				} else if coef > 0 {
					parts = append(parts, fmt.Sprintf("+%d%s", coef, eps(i+1)))
				} else {
					parts = append(parts, fmt.Sprintf("%d%s", coef, eps(i+1)))
				}
			}
		}
		if len(parts) == 0 {
			parts = append(parts, "0")
		}
		vecExprs = append(vecExprs, strings.Join(parts, ""))
	}
	return strings.Join(vecExprs, ", ")
}
