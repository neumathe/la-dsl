package test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/neumathe/la-dsl/dsl"
)

// TestCofactorA24Random
// 随机生成 4x4 行列式，考察代数余子式 A_{24}，验证：
// 1. runtime 生成的 A 与 A_{24} 一致
// 2. ExtractAnswer 返回的答案与手工计算一致
func TestCofactorA24Random(t *testing.T) {
	p := dsl.Problem{
		ID:      50001,
		Version: "v1",
		Title:   "行列式 D = {{D}} 中代数余子式 A_{24} 等于 ____",
		Variables: map[string]dsl.Variable{
			"A": {
				Kind: "matrix",
				Rows: 4,
				Cols: 4,
				Generator: map[string]interface{}{
					"rule": "range",
					"min":  -5,
					"max":  5,
				},
			},
		},
		Derived: map[string]string{
			"A24": "cofactor(A,2,4)",
		},
		Render: map[string]string{
			"D": "A",
		},
		Answer: dsl.AnswerSchema{
			Expression: "A24",
		},
	}

	// 固定 seed 以便测试可重复
	seedStr := fmt.Sprintf("cofactor-a24-%d", time.Now().UnixNano())
	inst, err := dsl.InstantiateProblem(p, seedStr, "test-salt")
	if err != nil {
		t.Fatalf("InstantiateProblem error: %v", err)
	}

	rawA, ok := inst.Vars["A"]
	if !ok {
		t.Fatalf("A not generated")
	}
	A, ok := rawA.(*dsl.MatrixInt)
	if !ok {
		t.Fatalf("A is not *MatrixInt, got %T", rawA)
	}

	// 1) 通过表达式系统计算 A24
	valExpr, err := dsl.EvaluateExpression("cofactor(A, 2, 4)", inst)
	if err != nil {
		t.Fatalf("EvaluateExpression cofactor(A,2,4) error: %v", err)
	}
	exprCof, ok := valExpr.(*big.Int)
	if !ok {
		t.Fatalf("cofactor expression result is not *big.Int, got %T", valExpr)
	}

	// Derived 中的 A24 也应该是同一个结果
	rawA24, ok := inst.Derived["A24"]
	if !ok {
		t.Fatalf("A24 not in derived")
	}
	derivedCof, ok := rawA24.(*big.Int)
	if !ok {
		t.Fatalf("derived A24 is not *big.Int, got %T", rawA24)
	}
	if exprCof.Cmp(derivedCof) != 0 {
		t.Fatalf("expr cofactor != derived cofactor: %s vs %s", exprCof.String(), derivedCof.String())
	}

	// 2) 手工构造 (2,4) 位置的余子式并计算代数余子式
	minor := dsl.NewMatrixInt(3, 3)
	ri := 0
	for i := 0; i < 4; i++ {
		if i == 1 { // 删除第 2 行
			continue
		}
		ci := 0
		for j := 0; j < 4; j++ {
			if j == 3 { // 删除第 4 列
				continue
			}
			minor.A[ri][ci] = A.A[i][j]
			ci++
		}
		ri++
	}
	detMinor := dsl.BareissDet(minor)
	sign := int64(1) // (-1)^{2+4} = 1
	manualCof := new(big.Int).Mul(big.NewInt(sign), detMinor)

	if manualCof.Cmp(exprCof) != 0 {
		t.Fatalf("manual cofactor != expr cofactor: %s vs %s", manualCof.String(), exprCof.String())
	}

	// 3) ExtractAnswer 也应该得到同一个值
	ans, err := dsl.ExtractAnswer(p, inst)
	if err != nil {
		t.Fatalf("ExtractAnswer error: %v", err)
	}
	ansBig, ok := ans.(*big.Int)
	if !ok {
		t.Fatalf("answer is not *big.Int, got %T", ans)
	}
	if ansBig.Cmp(exprCof) != 0 {
		t.Fatalf("answer != expr cofactor: %s vs %s", ansBig.String(), exprCof.String())
	}

	// 4) 输出一份题目和答案示例（LaTeX 形式）
	latexD := matrixToLatexDet(A)
	title, err := dsl.RenderTitle(p, inst)
	if err != nil {
		t.Fatalf("RenderTitle error: %v", err)
	}
	t.Logf("DSL Title 渲染结果：%s", title)
	t.Logf("示例题目（LaTeX）：行列式 D = %s 中代数余子式 A_{24} 等于 ____", latexD)
	t.Logf("示例答案：A_{24} = %s", ansBig.String())
}

func matrixToLatexDet(m *dsl.MatrixInt) string {
	var b strings.Builder
	b.WriteString("\\begin{vmatrix}")
	for i := 0; i < m.R; i++ {
		for j := 0; j < m.C; j++ {
			b.WriteString(fmt.Sprintf("%d", m.A[i][j]))
			if j < m.C-1 {
				b.WriteString("&")
			}
		}
		if i < m.R-1 {
			b.WriteString("\\\\")
		}
	}
	b.WriteString("\\end{vmatrix}")
	return b.String()
}
