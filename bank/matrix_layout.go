package bank

import (
	"fmt"

	"github.com/neumathe/la-dsl/dsl"
)

// MatrixFieldDefsSquare 生成 n×n 矩阵按行优先的填空定义，并附带 matrix_cell 布局（与 mget 下标一致）。
func MatrixFieldDefsSquare(key string, n int, matrixVar, groupLabel string) []dsl.AnswerFieldDef {
	ids := BlankIDs(key, n*n)
	return MatrixFieldDefsRectIDs(ids, 0, n, n, matrixVar, groupLabel)
}

// MatrixFieldDefsRectIDs 在 ids[offset : offset+rows*cols] 上生成 r×c 矩阵元填空及布局。
func MatrixFieldDefsRectIDs(ids []string, offset, rows, cols int, matrixVar, groupLabel string) []dsl.AnswerFieldDef {
	fds := make([]dsl.AnswerFieldDef, 0, rows*cols)
	for i := 1; i <= rows; i++ {
		for j := 1; j <= cols; j++ {
			k := offset + (i-1)*cols + (j-1)
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[k],
				Expr:   fmt.Sprintf("mget(%s,%d,%d)", matrixVar, i, j),
				Layout: dsl.LayoutMatrixCell(matrixVar, i, j, rows, cols, groupLabel),
			})
		}
	}
	return fds
}
