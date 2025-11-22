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
		title = strings.ReplaceAll(title, ph, formatValueForTitle(v))
	}
	return title, nil
}

// formatValueForTitle 将内部值转成适合题干中的字符串
func formatValueForTitle(v interface{}) string {
	switch t := v.(type) {
	case *MatrixInt:
		// LaTeX 矩阵（列向量题常用），方便直接嵌入 \( {{A}} \) 中
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
	case int64, int, float64:
		return fmt.Sprintf("%v", t)
	default:
		return fmt.Sprintf("%v", v)
	}
}
