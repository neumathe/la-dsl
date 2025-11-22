package qa

import (
	"fmt"
	"testing"

	"github.com/neumathe/la-dsl/dsl"
)

// 这个包用于测试高层封装：仅给 DSL + seed，返回题目文本和带空 ID 的答案。
// TODO
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
