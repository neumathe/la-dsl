package qa

import (
	"fmt"
	"testing"

	"github.com/neumathe/la-dsl/dsl"
)

// 这个包用于测试高层封装：仅给 DSL + seed，返回题目文本和带空 ID 的答案。

func TestGenerateQuestion_LambdaDetZero(t *testing.T) {
	p := dsl.Problem{
		ID:      1,
		Version: "v1",
		Title:   "已知实数 x= {{x}}，求 x 的值。",
		Variables: map[string]dsl.Variable{
			"x": {
				Kind:  "scalar",
				Fixed: int64(5),
			},
		},
		Render: map[string]string{
			"x": "x",
		},
		Answer: dsl.AnswerSchema{
			Expression: "x",
			FieldDefs: []dsl.AnswerFieldDef{
				{ID: "ans_x", Expr: "x"},
			},
		},
	}

	seed := "simple-demo"
	g1, err := dsl.GenerateQuestion(p, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}
	g2, err := dsl.GenerateQuestion(p, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}
	if g1.Title != g2.Title {
		t.Fatalf("Title not deterministic: %q vs %q", g1.Title, g2.Title)
	}
	if len(g1.AnswerFields) != len(g2.AnswerFields) {
		t.Fatalf("Answer field count mismatch")
	}
	for i := range g1.AnswerFields {
		if g1.AnswerFields[i].ID != g2.AnswerFields[i].ID ||
			g1.AnswerFields[i].Expr != g2.AnswerFields[i].Expr ||
			fmt.Sprintf("%v", g1.AnswerFields[i].Value) != fmt.Sprintf("%v", g2.AnswerFields[i].Value) {
			t.Fatalf("Answer field[%d] mismatch between runs", i)
		}
	}
}

// 1) 齐次线性方程组有非零解，求 λ（使用 lambda_linear_det_zero 生成矩阵与 λ）
func TestGenerateQuestion_LambdaDetZeroDSL(t *testing.T) {
	p := dsl.Problem{
		ID:      10001,
		Version: "v1",
		Title:   "", // 题干在 test 包中已有严格格式化，这里只验证封装与答案
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
		Answer: dsl.AnswerSchema{
			Expression: "lambda",
			FieldDefs: []dsl.AnswerFieldDef{
				{ID: "lambda", Expr: "lambda"},
			},
		},
	}

	seed := "lambda-det-zero-qa"
	g1, err := dsl.GenerateQuestion(p, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}
	g2, err := dsl.GenerateQuestion(p, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}
	if g1.Title != g2.Title {
		t.Fatalf("Title not deterministic: %q vs %q", g1.Title, g2.Title)
	}
	if len(g1.AnswerFields) != 1 || len(g2.AnswerFields) != 1 {
		t.Fatalf("lambda problem expected 1 answer field, got %d and %d", len(g1.AnswerFields), len(g2.AnswerFields))
	}
	a1, a2 := g1.AnswerFields[0], g2.AnswerFields[0]
	if a1.ID != a2.ID || a1.Expr != a2.Expr || fmt.Sprintf("%v", a1.Value) != fmt.Sprintf("%v", a2.Value) {
		t.Fatalf("lambda answer mismatch between runs: %#v vs %#v", a1, a2)
	}
}

// 2) 行列式 D 中代数余子式 A_{24}（4x4 随机矩阵）
func TestGenerateQuestion_CofactorA24DSL(t *testing.T) {
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
			FieldDefs: []dsl.AnswerFieldDef{
				{ID: "A24", Expr: "A24"},
			},
		},
	}

	seed := "cofactor-a24-qa"
	g1, err := dsl.GenerateQuestion(p, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}
	g2, err := dsl.GenerateQuestion(p, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}
	if g1.Title != g2.Title {
		t.Fatalf("Title not deterministic: %q vs %q", g1.Title, g2.Title)
	}
	if len(g1.AnswerFields) != 1 || len(g2.AnswerFields) != 1 {
		t.Fatalf("cofactor problem expected 1 answer field, got %d and %d", len(g1.AnswerFields), len(g2.AnswerFields))
	}
	a1, a2 := g1.AnswerFields[0], g2.AnswerFields[0]
	if a1.ID != a2.ID || a1.Expr != a2.Expr || fmt.Sprintf("%v", a1.Value) != fmt.Sprintf("%v", a2.Value) {
		t.Fatalf("cofactor answer mismatch between runs: %#v vs %#v", a1, a2)
	}
}

// 3) R[x]_3 中基变换 + p(x) 的坐标（过渡矩阵 3x3 + 坐标 3 维，共 12 个空）
func TestGenerateQuestion_BasisTransitionDSL(t *testing.T) {
	// 固定新基的坐标矩阵 B（列是新基在旧基 {1,x,x^2} 下的坐标）
	B := dsl.NewMatrixInt(3, 3)
	B.A[0][0], B.A[1][0], B.A[2][0] = -1, 0, 0
	B.A[0][1], B.A[1][1], B.A[2][1] = -2, 1, 0
	B.A[0][2], B.A[1][2], B.A[2][2] = -1, -2, 2

	e1 := dsl.NewVectorInt(3)
	e1.V[0], e1.V[1], e1.V[2] = 1, 0, 0
	e2 := dsl.NewVectorInt(3)
	e2.V[0], e2.V[1], e2.V[2] = 0, 1, 0
	e3 := dsl.NewVectorInt(3)
	e3.V[0], e3.V[1], e3.V[2] = 0, 0, 1

	p := dsl.Problem{
		ID:      70001,
		Version: "v1",
		Title:   "", // 题干在 test 包中已有严格格式化，这里只验证封装与答案
		Variables: map[string]dsl.Variable{
			"B":  {Kind: "matrix", Rows: 3, Cols: 3, Fixed: B},
			"e1": {Kind: "vector", Size: 3, Fixed: e1},
			"e2": {Kind: "vector", Size: 3, Fixed: e2},
			"e3": {Kind: "vector", Size: 3, Fixed: e3},
			// p(x) = a + bx + cx^2 使用范围内的随机整数系数，由 seed 控制确定性
			"p": {
				Kind: "vector",
				Size: 3,
				Generator: map[string]interface{}{
					"rule": "range",
					"min":  -2,
					"max":  2,
				},
			},
		},
		Derived: map[string]string{
			"Pcol1": "solve(B,e1)",
			"Pcol2": "solve(B,e2)",
			"Pcol3": "solve(B,e3)",
			"coord": "solve(B,p)",
		},
		Answer: dsl.AnswerSchema{
			FieldDefs: []dsl.AnswerFieldDef{
				{ID: "Chapter7_10_1", Expr: "Pcol1[1]"},
				{ID: "Chapter7_10_2", Expr: "Pcol1[2]"},
				{ID: "Chapter7_10_3", Expr: "Pcol1[3]"},
				{ID: "Chapter7_10_4", Expr: "Pcol2[1]"},
				{ID: "Chapter7_10_5", Expr: "Pcol2[2]"},
				{ID: "Chapter7_10_6", Expr: "Pcol2[3]"},
				{ID: "Chapter7_10_7", Expr: "Pcol3[1]"},
				{ID: "Chapter7_10_8", Expr: "Pcol3[2]"},
				{ID: "Chapter7_10_9", Expr: "Pcol3[3]"},
				{ID: "Chapter7_10_10", Expr: "coord[1]"},
				{ID: "Chapter7_10_11", Expr: "coord[2]"},
				{ID: "Chapter7_10_12", Expr: "coord[3]"},
			},
		},
	}

	seed := "basis-transition-qa"
	g1, err := dsl.GenerateQuestion(p, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}
	g2, err := dsl.GenerateQuestion(p, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}
	if g1.Title != g2.Title {
		t.Fatalf("Title not deterministic: %q vs %q", g1.Title, g2.Title)
	}
	if len(g1.AnswerFields) != 12 || len(g2.AnswerFields) != 12 {
		t.Fatalf("basis transition problem expected 12 answer fields, got %d and %d", len(g1.AnswerFields), len(g2.AnswerFields))
	}
	for i := range g1.AnswerFields {
		a1, a2 := g1.AnswerFields[i], g2.AnswerFields[i]
		if a1.ID != a2.ID || a1.Expr != a2.Expr || fmt.Sprintf("%v", a1.Value) != fmt.Sprintf("%v", a2.Value) {
			t.Fatalf("basis transition answer[%d] mismatch between runs: %#v vs %#v", i, a1, a2)
		}
	}
}
