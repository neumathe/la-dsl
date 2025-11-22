package test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/neumathe/la-dsl/dsl"
)

// TestLambdaLinearDetZeroGenerator
// 验证 lambda_linear_det_zero 生成的矩阵：
// 1. 行列式为 0
// 2. 求得的 λ 为整数，且在给定区间内
// 3. 用另一种方式（通过 det(λ) 的线性性质）重新算出的 λ 与生成值一致
func TestLambdaLinearDetZeroGenerator(t *testing.T) {
	p := dsl.Problem{
		ID:      10001,
		Version: "v1",
		Title:   "",
		Variables: map[string]dsl.Variable{
			"A": {
				Kind: "matrix",
				Rows: 3,
				Cols: 3,
				Generator: map[string]interface{}{
					"rule":         "lambda_linear_det_zero",
					"param_var":    "lambda",
					"param_row":    3,
					"param_col":    3,
					"entry_min":    -8,
					"entry_max":    8,
					"lambda_min":   -10,
					"lambda_max":   10,
					"max_attempts": 100,
				},
			},
		},
		Render: map[string]string{
			"A": "A",
		},
		Answer: dsl.AnswerSchema{
			Expression: "lambda",
		},
	}

	// 每次运行使用不同的 seedStr，从而得到不同的 λ 和矩阵
	seedStr := fmt.Sprintf("lambda-det-zero-%d", time.Now().UnixNano())
	inst, err := dsl.InstantiateProblem(p, seedStr, "test-salt")
	if err != nil {
		t.Fatalf("InstantiateProblem error: %v", err)
	}

	// 检查 A 和 lambda 是否生成
	rawA, ok := inst.Vars["A"]
	if !ok {
		t.Fatalf("A not generated in instance vars")
	}
	A, ok := rawA.(*dsl.MatrixInt)
	if !ok {
		t.Fatalf("A is not *MatrixInt, got %T", rawA)
	}

	rawLambda, ok := inst.Vars["lambda"]
	if !ok {
		t.Fatalf("lambda not generated in instance vars")
	}
	lambdaVal, ok := rawLambda.(int64)
	if !ok {
		t.Fatalf("lambda is not int64, got %T", rawLambda)
	}

	// 1) det(A) 必须为 0
	d := dsl.BareissDet(A)
	if d.Sign() != 0 {
		t.Fatalf("det(A) != 0, got %s", d.String())
	}

	// 2) 使用另一种方式从矩阵恢复 λ：
	//    已知只有 (3,3) 元素依赖 λ，且形式为 λ + c
	//    已知当前矩阵是在 λ=lambdaVal 下的数值 A(λ)
	//    于是 c = A(3,3) - λ
	//    再构造 A(0) 和 A(1)，利用 det 的线性性重新求 λ = -b/a
	const row, col = 3, 3
	constC := A.A[row-1][col-1] - lambdaVal

	A0 := dsl.NewMatrixInt(A.R, A.C)
	A1 := dsl.NewMatrixInt(A.R, A.C)
	for i := 0; i < A.R; i++ {
		for j := 0; j < A.C; j++ {
			A0.A[i][j] = A.A[i][j]
			A1.A[i][j] = A.A[i][j]
		}
	}
	A0.A[row-1][col-1] = constC     // λ = 0
	A1.A[row-1][col-1] = constC + 1 // λ = 1

	d0 := dsl.BareissDet(A0)
	d1 := dsl.BareissDet(A1)

	aBig := new(big.Int).Sub(d1, d0)
	if aBig.Sign() == 0 {
		t.Fatalf("det is independent of lambda in recovered check")
	}
	bBig := new(big.Int).Set(d0)

	q := new(big.Int)
	rBig := new(big.Int)
	q.DivMod(bBig, aBig, rBig)
	if rBig.Sign() != 0 {
		t.Fatalf("recovered lambda is not integer, remainder=%s", rBig.String())
	}
	lambdaRecovered := new(big.Int).Neg(q)
	if !lambdaRecovered.IsInt64() {
		t.Fatalf("recovered lambda not int64")
	}
	if lambdaRecovered.Int64() != lambdaVal {
		t.Fatalf("lambda mismatch: generated=%d recovered=%d", lambdaVal, lambdaRecovered.Int64())
	}

	// 3) 用 runtime 的 ExtractAnswer 再取一次答案，应该等于 lambdaVal
	ans, err := dsl.ExtractAnswer(p, inst)
	if err != nil {
		t.Fatalf("ExtractAnswer error: %v", err)
	}
	ansInt, ok := ans.(int64)
	if !ok {
		t.Fatalf("answer is not int64, got %T", ans)
	}
	if ansInt != lambdaVal {
		t.Fatalf("answer mismatch: lambda=%d answer=%d", lambdaVal, ansInt)
	}

	// 4) 输出一份题目和答案示例（严格按照原始题型样式，只随机数值）
	question := formatHomogeneousSystemWithLambda(A, lambdaVal)
	t.Logf("示例题目（LaTeX）：%s", question)
	t.Logf("示例答案：λ = %d", lambdaVal)
}

// 将生成的 3x3 矩阵 A 和 λ 格式化为原题型：
// 已知齐次线性方程组
// \( \begin{cases} a11 x_1 + a12 x_2 + a13 x_3 = 0 \\ ... \\ 0 x_1 + c2 x_2 + (\lambda + d) x_3 = 0 \end{cases} \)
// 有非零解，则 \lambda =
func formatHomogeneousSystemWithLambda(m *dsl.MatrixInt, lambdaVal int64) string {
	if m.R != 3 || m.C != 3 {
		return ""
	}
	a11, a12, a13 := m.A[0][0], m.A[0][1], m.A[0][2]
	b11, b12, b13 := m.A[1][0], m.A[1][1], m.A[1][2]
	c11, c12, c13 := m.A[2][0], m.A[2][1], m.A[2][2]

	// 假设第三行第三列系数来源于 (lambda + d)
	delta := c13 - lambdaVal // d

	row1 := fmt.Sprintf("%d x_1%+d x_2%+d x_3=0", a11, a12, a13)
	row2 := fmt.Sprintf("%d x_1%+d x_2%+d x_3=0", b11, b12, b13)

	// 第三行：c11 x1 + c12 x2 + (λ ± k) x3 = 0
	// 这里保持一般形式，不强制 c11=0
	var lambdaTerm string
	if delta == 0 {
		lambdaTerm = "\\lambda"
	} else if delta > 0 {
		lambdaTerm = fmt.Sprintf("\\lambda+%d", delta)
	} else {
		lambdaTerm = fmt.Sprintf("\\lambda%d", delta) // delta 为负，会变成 \lambda-...
	}

	row3Prefix := fmt.Sprintf("%d x_1%+d x_2+", c11, c12)
	row3 := fmt.Sprintf("%s(%s) x_3=0", row3Prefix, lambdaTerm)

	var b strings.Builder
	b.WriteString("已知齐次线性方程组\\(\\begin{cases}")
	b.WriteString(row1)
	b.WriteString("\\\\")
	b.WriteString(row2)
	b.WriteString("\\\\")
	b.WriteString(row3)
	b.WriteString("\\end{cases}\\)有非零解，则\\(\\lambda\\)=")
	return b.String()
}
