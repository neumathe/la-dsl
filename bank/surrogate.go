package bank

import (
	"fmt"
	"strings"

	"github.com/neumathe/la-dsl/dsl"
)

// MatrixRowMajorSurrogate 用 r×c 随机整数矩阵按行优先展开为 n=r*c 个空，
// 用于与 HTML 空位数对齐的结构占位（题干简短说明为数值矩阵展开）。
func MatrixRowMajorSurrogate(key string, rows, cols int, titlePrefix string) dsl.Problem {
	n := rows * cols
	ids := BlankIDs(key, n)
	sb := strings.Builder{}
	sb.WriteString(titlePrefix)
	sb.WriteString("（逐行填写矩阵各元素）")
	for _, id := range ids {
		sb.WriteString(fmt.Sprintf(" ${{blank:%s}}$", id))
	}
	fds := MatrixFieldDefsRectIDs(ids, 0, rows, cols, "A", "A")
	return dsl.Problem{
		ID:      ProblemID(key),
		Version: "bank-v1",
		Title:   sb.String(),
		Variables: map[string]dsl.Variable{
			"A": {
				Kind: "matrix", Rows: rows, Cols: cols,
				Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4},
			},
		},
		Render: map[string]string{"A": "A"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"bank_surrogate": true,
			"matrix_shape":   fmt.Sprintf("%dx%d", rows, cols),
		},
	}
}
