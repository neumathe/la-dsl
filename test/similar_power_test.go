package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/neumathe/la-dsl/dsl"
)

// TestSimilarPowerRandom
// 题型：设 A、P 为 2x2 矩阵，则 (P^{-1}AP)^8 和 A^7 各是多少。
// 这里：
// - A 随机生成为 λI（2x2），λ 为小整数，保证计算简单；
// - P 随机生成为 2x2 满秩整数矩阵；
// - 利用 DSL 的 pow/mget 等通用能力计算 A^7 和 (P^{-1}AP)^8；
// - 数学上对任意可逆 P，有 A=λI 时 (P^{-1}AP)^8 = A^8 = λ^8 I。
func TestSimilarPowerRandom(t *testing.T) {
	prob := dsl.Problem{
		ID:      90001,
		Version: "v1",
		Title:   "设 A={{A}}, P={{P}}，则 (P^{-1}AP)^8 = \\begin{bmatrix}{{blank:b11}}&{{blank:b12}}\\\\{{blank:b21}}&{{blank:b22}}\\end{bmatrix}，A^7 = \\begin{bmatrix}{{blank:a11}}&{{blank:a12}}\\\\{{blank:a21}}&{{blank:a22}}\\end{bmatrix}",
		Variables: map[string]dsl.Variable{
			"A": {
				Kind: "matrix",
				Rows: 2,
				Cols: 2,
				Generator: map[string]interface{}{
					"rule":       "scalar_identity",
					"lambda_min": -3,
					"lambda_max": 3,
					"lambda_var": "lambda",
				},
			},
			"P": {
				Kind: "matrix",
				Rows: 2,
				Cols: 2,
				Generator: map[string]interface{}{
					"rule": "full_rank",
					"min":  -5,
					"max":  5,
				},
			},
		},
		Derived: map[string]string{
			"A7": "pow(A,7)",
			"A8": "pow(A,8)", // 对 A=λI，有 (P^{-1}AP)^8 = A^8 = λ^8 I
		},
		Render: map[string]string{
			"A": "A",
			"P": "P",
		},
		Answer: dsl.AnswerSchema{
			FieldDefs: []dsl.AnswerFieldDef{
				// (P^{-1}AP)^8 的 4 个元素，等于 A8 的 4 个元素
				{ID: "b11", Expr: "mget(A8,1,1)"},
				{ID: "b12", Expr: "mget(A8,1,2)"},
				{ID: "b21", Expr: "mget(A8,2,1)"},
				{ID: "b22", Expr: "mget(A8,2,2)"},
				// A^7 的 4 个元素
				{ID: "a11", Expr: "mget(A7,1,1)"},
				{ID: "a12", Expr: "mget(A7,1,2)"},
				{ID: "a21", Expr: "mget(A7,2,1)"},
				{ID: "a22", Expr: "mget(A7,2,2)"},
			},
		},
	}

	// 使用时间作为 seed，保证每次运行随机，但同一 seed 行为确定
	seedStr := fmt.Sprintf("similar-power-%d", time.Now().UnixNano())

	// 先实例化，检查 A 是否为 λI，并计算理论上的 λ^7 和 λ^8
	inst, err := dsl.InstantiateProblem(prob, seedStr, "test-salt")
	if err != nil {
		t.Fatalf("InstantiateProblem error: %v", err)
	}

	rawA, ok := inst.Vars["A"]
	if !ok {
		t.Fatalf("A not generated")
	}
	A, ok := rawA.(*dsl.MatrixInt)
	if !ok || A.R != 2 || A.C != 2 {
		t.Fatalf("A is not 2x2 MatrixInt, got %T", rawA)
	}

	rawLambda, ok := inst.Vars["lambda"]
	if !ok {
		t.Fatalf("lambda not generated")
	}
	lambda, ok := rawLambda.(int64)
	if !ok {
		t.Fatalf("lambda is not int64, got %T", rawLambda)
	}

	// 校验 A = λI
	if A.A[0][0] != lambda || A.A[1][1] != lambda || A.A[0][1] != 0 || A.A[1][0] != 0 {
		t.Fatalf("A is not lambda I, A=%+v, lambda=%d", A.A, lambda)
	}

	// 计算理论上的 λ^7 和 λ^8
	l7 := int64(1)
	l8 := int64(1)
	for i := 0; i < 7; i++ {
		l7 *= lambda
	}
	for i := 0; i < 8; i++ {
		l8 *= lambda
	}

	// 使用高层 API 生成题目与答案
	g, err := dsl.GenerateQuestion(prob, seedStr, "test-salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}
	if len(g.AnswerFields) != 8 {
		t.Fatalf("expected 8 answer fields, got %d", len(g.AnswerFields))
	}

	// 检查答案是否符合 λ^7 I 和 λ^8 I
	expected := map[string]int64{
		"b11": l8, "b12": 0, "b21": 0, "b22": l8,
		"a11": l7, "a12": 0, "a21": 0, "a22": l7,
	}
	for _, f := range g.AnswerFields {
		v, ok := f.Value.(int64)
		if !ok {
			t.Fatalf("answer field %s value type %T, expected int64", f.ID, f.Value)
		}
		if exp, ok := expected[f.ID]; ok {
			if v != exp {
				t.Fatalf("field %s value=%d, expected %d (lambda=%d)", f.ID, v, exp, lambda)
			}
		}
	}

	// 输出示例题目与答案，便于人工查看
	t.Logf("示例题目（LaTeX）：%s", g.Title)
	for _, f := range g.AnswerFields {
		t.Logf("答案字段 ID=%s, Value=%v", f.ID, f.Value)
	}
}
