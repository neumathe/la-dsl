package dsl

import (
	"strings"
	"testing"
)

func TestGenerateExplanationMatchesGenerateQuestion(t *testing.T) {
	p := Problem{
		ID:      91001,
		Version: "v1",
		Title:   `设矩阵 $M={{M}}$，求其 $(1,1)$ 元`,
		Variables: map[string]Variable{
			"M": {Kind: "matrix", Rows: 2, Cols: 2, Generator: map[string]interface{}{"rule": "range", "min": 1, "max": 3}},
		},
		Render:  map[string]string{"M": "M"},
		Derived: map[string]string{"t": "transpose(M)"},
		Answer: AnswerSchema{
			FieldDefs: []AnswerFieldDef{{ID: "f1", Expr: "mget(M,1,1)"}},
		},
	}
	seed, salt := "exp-seed", "exp-salt"
	g, err := GenerateQuestion(p, seed, salt)
	if err != nil {
		t.Fatal(err)
	}
	ex, err := GenerateExplanation(p, seed, salt)
	if err != nil {
		t.Fatal(err)
	}
	if ex.Title != g.Title {
		t.Fatalf("title mismatch: %q vs %q", ex.Title, g.Title)
	}
	if len(ex.AnswerSteps) != len(g.AnswerFields) {
		t.Fatalf("steps %d vs fields %d", len(ex.AnswerSteps), len(g.AnswerFields))
	}
	for i := range g.AnswerFields {
		if ex.AnswerSteps[i].FieldID != g.AnswerFields[i].ID {
			t.Fatalf("id[%d]", i)
		}
		if ex.AnswerSteps[i].Expr != g.AnswerFields[i].Expr {
			t.Fatalf("expr[%d]", i)
		}
		if ex.AnswerSteps[i].Expected != ValueToCanonicalString(g.AnswerFields[i].Value) {
			t.Fatalf("expected[%d]: %q vs %q", i, ex.AnswerSteps[i].Expected, ValueToCanonicalString(g.AnswerFields[i].Value))
		}
	}
	if len(ex.Derived) != 1 || ex.Derived[0].Name != "t" {
		t.Fatalf("derived: %+v", ex.Derived)
	}
	if ex.TitlePlugs["M"] == "" {
		t.Fatal("missing title plug M")
	}
}

func TestGenerateExplanationDeterministic(t *testing.T) {
	p := Problem{
		ID: 91002, Version: "v1",
		Variables: map[string]Variable{
			"s": {Kind: "scalar", Generator: map[string]interface{}{"rule": "range", "min": -2, "max": 4}},
		},
		Answer: AnswerSchema{FieldDefs: []AnswerFieldDef{{ID: "a1", Expr: "s"}}},
	}
	a, _ := GenerateExplanation(p, "d", "s")
	b, _ := GenerateExplanation(p, "d", "s")
	if a.AnswerSteps[0].Expected != b.AnswerSteps[0].Expected {
		t.Fatal("drift")
	}
}

// TestSolutionZhPlaceholderExpansion 验证 solution_zh 中的 {{key}} 占位符被正确展开。
func TestSolutionZhPlaceholderExpansion(t *testing.T) {
	p := Problem{
		ID:      91003,
		Version: "v1",
		Title:   `三阶行列式 $D={{vA}}$ 等于 {{blank:f1}}`,
		Variables: map[string]Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "range", "min": -2, "max": 2}},
		},
		Derived: map[string]string{"d": "det(A)", "vA": "vmatrix_title(A)"},
		Render:  map[string]string{"vA": "vA"},
		Answer: AnswerSchema{FieldDefs: []AnswerFieldDef{{ID: "f1", Expr: "d"}}},
		Meta: map[string]interface{}{
			"solution_zh": `行列式 $D={{vA}}$ 的值为 {{d}}，填空 f1 的答案为 {{f1}}。`,
		},
	}
	ex, err := GenerateExplanation(p, "test-sol-seed", "test-sol-salt")
	if err != nil {
		t.Fatal(err)
	}
	if ex.Solution == "" {
		t.Fatal("solution_zh is empty")
	}
	// {{vA}} 应被 LaTeX 格式矩阵替换（含 \begin{vmatrix}，行列式专用）
	if !strings.Contains(ex.Solution, "\\begin{vmatrix}") {
		t.Fatalf("solution_zh should contain LaTeX matrix from Render key, got: %q", ex.Solution)
	}
	// {{d}} 应被行列式数值替换（与标准答案一致）
	if !strings.Contains(ex.Solution, ex.AnswerSteps[0].Expected) {
		t.Fatalf("solution_zh should contain det value {{d}}=%s matching expected, got: %q", ex.AnswerSteps[0].Expected, ex.Solution)
	}
	// {{f1}} 应被答案字段期望值替换
	if !strings.Contains(ex.Solution, ex.AnswerSteps[0].Expected) {
		t.Fatalf("solution_zh should contain field ID placeholder expanded to expected value, got: %q", ex.Solution)
	}
	// 不应再包含未展开的 {{...}} 模板
	if strings.Contains(ex.Solution, "{{") {
		t.Fatalf("solution_zh should not contain unexpanded placeholders, got: %q", ex.Solution)
	}
}

// TestSolutionZhExprPlaceholderExpansion 验证 {{expr:DSL_EXPRESSION}} 内联求值展开。
func TestSolutionZhExprPlaceholderExpansion(t *testing.T) {
	p := Problem{
		ID:      91005,
		Version: "v1",
		Title:   `设 $A={{A}}$，求行列式`,
		Variables: map[string]Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "range", "min": -3, "max": 3}},
		},
		Derived: map[string]string{"d": "det(A)"},
		Render:  map[string]string{"A": "A"},
		Answer:  AnswerSchema{FieldDefs: []AnswerFieldDef{{ID: "f1", Expr: "d"}}},
		Meta: map[string]interface{}{
			"solution_zh": `按第1行展开：$a_{11}={{expr:mget(A,1,1)}}$, $a_{12}={{expr:mget(A,1,2)}}$, $a_{13}={{expr:mget(A,1,3)}}$；行列式 $={{expr:det(A)}}$。`,
		},
	}
	ex, err := GenerateExplanation(p, "expr-seed", "expr-salt")
	if err != nil {
		t.Fatal(err)
	}
	if ex.Solution == "" {
		t.Fatal("solution_zh is empty")
	}
	// {{expr:mget(A,1,1)}} 应被实际矩阵元素值替换（整数）
	// 验证不含未展开的 {{expr: 模板
	if strings.Contains(ex.Solution, "{{expr:") {
		t.Fatalf("solution_zh should not contain unexpanded {{expr:...}} placeholders, got: %q", ex.Solution)
	}
	// {{expr:det(A)}} 的值应与 {{d}} 一致（Derived 已有 det(A)）
	if !strings.Contains(ex.Solution, ex.AnswerSteps[0].Expected) {
		t.Fatalf("solution_zh should contain det(A) value matching expected, got: %q", ex.Solution)
	}
	// 不应有任何未展开的 {{ 模板残留
	if strings.Contains(ex.Solution, "{{") {
		t.Fatalf("solution_zh should not contain any unexpanded {{, got: %q", ex.Solution)
	}
}

// TestSolutionZhDerivedKeyExpansion 验证 Derived 键也可作为占位符使用。
func TestSolutionZhDerivedKeyExpansion(t *testing.T) {
	p := Problem{
		ID:      91004,
		Version: "v1",
		Title:   `设 $A={{A}}$，求 $AB$`,
		Variables: map[string]Variable{
			"A": {Kind: "matrix", Rows: 2, Cols: 2, Generator: map[string]interface{}{"rule": "range", "min": 1, "max": 3}},
			"B": {Kind: "matrix", Rows: 2, Cols: 2, Generator: map[string]interface{}{"rule": "range", "min": 1, "max": 3}},
		},
		Derived: map[string]string{
			"AB": "matmul(A,B)",
		},
		Render: map[string]string{"A": "A", "B": "B"},
		Answer: AnswerSchema{FieldDefs: []AnswerFieldDef{
			{ID: "f1", Expr: "mget(AB,1,1)"},
			{ID: "f2", Expr: "mget(AB,1,2)"},
			{ID: "f3", Expr: "mget(AB,2,1)"},
			{ID: "f4", Expr: "mget(AB,2,2)"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `$AB={{AB}}$，答案 f1={{f1}} f2={{f2}} f3={{f3}} f4={{f4}}。`,
		},
	}
	ex, err := GenerateExplanation(p, "sol2", "salt2")
	if err != nil {
		t.Fatal(err)
	}
	if ex.Solution == "" {
		t.Fatal("solution_zh is empty")
	}
	// {{AB}} 应被 Derived 值替换（ValueToExplainString 格式）
	abDerived := ""
	for _, d := range ex.Derived {
		if d.Name == "AB" {
			abDerived = d.Value
		}
	}
	if abDerived == "" {
		t.Fatal("missing derived AB")
	}
	if !strings.Contains(ex.Solution, abDerived) {
		t.Fatalf("solution_zh should contain derived AB value %q, got: %q", abDerived, ex.Solution)
	}
	// 各 field ID 占位符也应展开为期望值
	for _, step := range ex.AnswerSteps {
		if !strings.Contains(ex.Solution, step.Expected) {
			t.Fatalf("solution_zh should contain field %s expected value %q, got: %q", step.FieldID, step.Expected, ex.Solution)
		}
	}
}
