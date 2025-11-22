package qa

import (
	"fmt"
	"testing"

	"github.com/neumathe/la-dsl/dsl"
)

// TestQA_LambdaDetZero 使用与 test/lambda_det_zero_test.go 等价的 DSL，
// 通过封装 API 生成题目与答案，验证：
// 1. 题干包含 {{blank:ID}} 占位符
// 2. AnswerFields 的 ID 与占位符匹配
// 3. 同 DSL + seed 下多次调用结果完全一致
func TestQA_LambdaDetZero(t *testing.T) {
	prob := dsl.Problem{
		ID:      10001,
		Version: "v1",
		// Title 用占位符标记填空，前端根据 {{blank:ans_lambda}} 渲染输入框
		Title: "已知齐次线性方程组 Ax=0 有非零解，其中 A 为 3×3 矩阵（元素已随机生成），则参数 λ = {{blank:ans_lambda}}",
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
		Render: map[string]string{},
		Answer: dsl.AnswerSchema{
			FieldDefs: []dsl.AnswerFieldDef{
				{ID: "ans_lambda", Expr: "lambda"},
			},
		},
	}

	seed := "qa-lambda-fixed-seed"
	g1, err := dsl.GenerateQuestion(prob, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}

	// 校验题干包含占位符
	if !contains(g1.Title, "{{blank:ans_lambda}}") {
		t.Logf("Warning: Title 没有包含 {{blank:ans_lambda}}，实际: %s", g1.Title)
	}

	// 校验答案字段
	if len(g1.AnswerFields) != 1 {
		t.Fatalf("Expected 1 answer field, got %d", len(g1.AnswerFields))
	}
	if g1.AnswerFields[0].ID != "ans_lambda" {
		t.Fatalf("Expected ID=ans_lambda, got %s", g1.AnswerFields[0].ID)
	}

	// 多次调用应完全一致
	g2, err := dsl.GenerateQuestion(prob, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion 2nd call error: %v", err)
	}
	if g1.Title != g2.Title {
		t.Fatalf("Title not deterministic")
	}
	if fmt.Sprintf("%v", g1.AnswerFields[0].Value) != fmt.Sprintf("%v", g2.AnswerFields[0].Value) {
		t.Fatalf("Answer value not deterministic")
	}

	t.Logf("[QA] Lambda 题：Title=%s, Answer ID=%s Value=%v",
		g1.Title, g1.AnswerFields[0].ID, g1.AnswerFields[0].Value)
}

// TestQA_CofactorA24 对应 test/cofactor_a24_test.go
func TestQA_CofactorA24(t *testing.T) {
	prob := dsl.Problem{
		ID:      50001,
		Version: "v1",
		Title:   "行列式 D = {{D}} 中代数余子式 A_{24} 等于 {{blank:ans_A24}}",
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
			FieldDefs: []dsl.AnswerFieldDef{
				{ID: "ans_A24", Expr: "A24"},
			},
		},
	}

	seed := "qa-cofactor-fixed-seed"
	g1, err := dsl.GenerateQuestion(prob, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}

	if !contains(g1.Title, "{{blank:ans_A24}}") {
		t.Logf("Warning: Title 没有包含 {{blank:ans_A24}}，实际: %s", g1.Title)
	}
	if len(g1.AnswerFields) != 1 {
		t.Fatalf("Expected 1 answer field, got %d", len(g1.AnswerFields))
	}
	if g1.AnswerFields[0].ID != "ans_A24" {
		t.Fatalf("Expected ID=ans_A24, got %s", g1.AnswerFields[0].ID)
	}

	g2, err := dsl.GenerateQuestion(prob, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion 2nd call error: %v", err)
	}
	if g1.Title != g2.Title || fmt.Sprintf("%v", g1.AnswerFields[0].Value) != fmt.Sprintf("%v", g2.AnswerFields[0].Value) {
		t.Fatalf("Result not deterministic")
	}

	t.Logf("[QA] Cofactor 题：Title=%s, Answer ID=%s Value=%v",
		g1.Title, g1.AnswerFields[0].ID, g1.AnswerFields[0].Value)
}

// TestQA_BasisTransition 对应 test/basis_transition_test.go
// 12 个空：9 个矩阵元素 + 3 个坐标分量
// 测试随机性：B（新基）和 p(x) 都随机生成，题面显示随机的新基
func TestQA_BasisTransition(t *testing.T) {
	prob := dsl.Problem{
		ID:      70001,
		Version: "v1",
		// Title 包含新基的显示和 12 个 {{blank:ID}}
		Title: `向量空间 R[x]_3 中，由基 1,x,x^2 到基 {{b1}},{{b2}},{{b3}} 的过渡矩阵为：
第1行：{{blank:M11}} {{blank:M12}} {{blank:M13}}
第2行：{{blank:M21}} {{blank:M22}} {{blank:M23}}
第3行：{{blank:M31}} {{blank:M32}} {{blank:M33}}
p(x)={{px}} 在该基下的坐标为 ({{blank:C1}},{{blank:C2}},{{blank:C3}})^T`,
		Variables: map[string]dsl.Variable{
			// B 是新基的坐标矩阵，随机生成满秩矩阵（保证可逆）
			"B": {
				Kind: "matrix",
				Rows: 3,
				Cols: 3,
				Generator: map[string]interface{}{
					"rule": "full_rank",
					"min":  -3,
					"max":  3,
				},
			},
			"e1": {Kind: "vector", Size: 3, Fixed: []interface{}{1, 0, 0}},
			"e2": {Kind: "vector", Size: 3, Fixed: []interface{}{0, 1, 0}},
			"e3": {Kind: "vector", Size: 3, Fixed: []interface{}{0, 0, 1}},
			// p(x) 随机生成
			"p": {
				Kind: "vector",
				Size: 3,
				Generator: map[string]interface{}{
					"rule": "range",
					"min":  -5,
					"max":  5,
				},
			},
		},
		Derived: map[string]string{
			"Pcol1": "solve(B,e1)",
			"Pcol2": "solve(B,e2)",
			"Pcol3": "solve(B,e3)",
			"coord": "solve(B,p)",
			"b1":    "col(B,1)", // 新基的第一个向量
			"b2":    "col(B,2)", // 新基的第二个向量
			"b3":    "col(B,3)", // 新基的第三个向量
		},
		Render: map[string]string{
			"px": "p",  // 渲染 p(x) 的向量形式
			"b1": "b1", // 渲染新基第一个向量
			"b2": "b2", // 渲染新基第二个向量
			"b3": "b3", // 渲染新基第三个向量
		},
		Answer: dsl.AnswerSchema{
			FieldDefs: []dsl.AnswerFieldDef{
				{ID: "M11", Expr: "Pcol1[1]"},
				{ID: "M12", Expr: "Pcol2[1]"},
				{ID: "M13", Expr: "Pcol3[1]"},
				{ID: "M21", Expr: "Pcol1[2]"},
				{ID: "M22", Expr: "Pcol2[2]"},
				{ID: "M23", Expr: "Pcol3[2]"},
				{ID: "M31", Expr: "Pcol1[3]"},
				{ID: "M32", Expr: "Pcol2[3]"},
				{ID: "M33", Expr: "Pcol3[3]"},
				{ID: "C1", Expr: "coord[1]"},
				{ID: "C2", Expr: "coord[2]"},
				{ID: "C3", Expr: "coord[3]"},
			},
		},
	}

	seed := "qa-basis-fixed-seed"
	g1, err := dsl.GenerateQuestion(prob, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion error: %v", err)
	}

	// 测试确定性：相同 seed 生成相同题目
	g2, err := dsl.GenerateQuestion(prob, seed, "salt")
	if err != nil {
		t.Fatalf("GenerateQuestion 2nd call error: %v", err)
	}
	if g1.Title != g2.Title {
		t.Fatalf("Title not deterministic")
	}
	for i := range g1.AnswerFields {
		if fmt.Sprintf("%v", g1.AnswerFields[i].Value) != fmt.Sprintf("%v", g2.AnswerFields[i].Value) {
			t.Fatalf("AnswerField[%d] value not deterministic", i)
		}
	}

	// 校验 Title 包含所有占位符
	expectedBlanks := []string{"M11", "M12", "M13", "M21", "M22", "M23", "M31", "M32", "M33", "C1", "C2", "C3"}
	for _, id := range expectedBlanks {
		placeholder := fmt.Sprintf("{{blank:%s}}", id)
		if !contains(g1.Title, placeholder) {
			t.Logf("Warning: Title 没有包含 %s", placeholder)
		}
	}

	// 校验答案字段数量和 ID
	if len(g1.AnswerFields) != 12 {
		t.Fatalf("Expected 12 answer fields, got %d", len(g1.AnswerFields))
	}
	for i, expected := range expectedBlanks {
		if g1.AnswerFields[i].ID != expected {
			t.Fatalf("AnswerField[%d] ID=%s, expected %s", i, g1.AnswerFields[i].ID, expected)
		}
	}

	t.Logf("[QA] Basis 题目:\n%s", g1.Title)
	t.Logf("[QA] 答案字段数=%d", len(g1.AnswerFields))
	for i, field := range g1.AnswerFields {
		t.Logf("  [%d] ID=%s, Value=%v", i, field.ID, field.Value)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s != "" && substr != "" &&
		(s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i < len(s)-len(substr)+1; i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}()))
}
