package bank

import (
	"crypto/sha256"
	"fmt"

	"github.com/neumathe/la-dsl/dsl"
)

// ProblemID 为题库键生成稳定的 dsl.Problem.ID（与 Instantiate 的种子派生一致）。
func ProblemID(questionKey string) int64 {
	h := sha256.Sum256([]byte("la-dsl|bank-problem-id|" + questionKey))
	var id int64
	for i := 0; i < 8; i++ {
		id = (id << 8) | int64(h[i])
	}
	if id == 0 {
		return 1
	}
	if id < 0 {
		return -id
	}
	return id
}

// BlankIDs 生成与 HTML 一致的填空 id：{logicalKey}_{1..n}
func BlankIDs(logicalKey string, n int) []string {
	out := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		out = append(out, fmt.Sprintf("%s_%d", logicalKey, i))
	}
	return out
}

// GramSchmidtFieldDefs 列优先：第 j 个正交化向量在 R^{ambientDim} 中的第 i 分量对应 gs_comp(V,j,i)。
func GramSchmidtFieldDefs(logicalKey string, ambientDim, numVec int, matrixVar string) ([]dsl.AnswerFieldDef, []string) {
	n := ambientDim * numVec
	ids := BlankIDs(logicalKey, n)
	fds := make([]dsl.AnswerFieldDef, 0, n)
	for col := 1; col <= numVec; col++ {
		for row := 1; row <= ambientDim; row++ {
			fds = append(fds, dsl.AnswerFieldDef{
				ID:   ids[len(fds)],
				Expr: fmt.Sprintf("gs_comp(%s,%d,%d)", matrixVar, col, row),
			})
		}
	}
	return fds, ids
}
