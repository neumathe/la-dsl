package bank

import (
	"testing"

	"github.com/neumathe/la-dsl/dsl"
)

func TestAllCatalogKeysHaveBuilders(t *testing.T) {
	for _, k := range AllQuestionKeys {
		if _, ok := builders[k]; !ok {
			t.Errorf("missing builder for %q", k)
		}
	}
	for k := range builders {
		found := false
		for _, x := range AllQuestionKeys {
			if x == k {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("extra builder without catalog entry: %q", k)
		}
	}
}

func TestDSLInvLowerUnit(t *testing.T) {
	p := dsl.Problem{
		ID:      99001,
		Version: "v1",
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "lower_unit", "min": -2, "max": 2}},
		},
		Derived: map[string]string{"Inv": "inv(A)"},
		Answer: dsl.AnswerSchema{
			FieldDefs: []dsl.AnswerFieldDef{{ID: "x", Expr: "mget(Inv,1,1)"}},
		},
	}
	g, err := dsl.GenerateQuestion(p, "inv-seed", "s")
	if err != nil {
		t.Fatal(err)
	}
	if len(g.AnswerFields) != 1 {
		t.Fatal("expected 1 field")
	}
}

func TestNullvec3x4(t *testing.T) {
	p := dsl.Problem{
		ID:      99002,
		Version: "v1",
		Variables: map[string]dsl.Variable{
			"V": {Kind: "matrix", Rows: 4, Cols: 3, Generator: map[string]interface{}{"rule": "full_rank", "min": -3, "max": 3}},
		},
		Derived: map[string]string{"MT": "transpose(V)", "w": "nullvec(MT)"},
		Answer: dsl.AnswerSchema{
			FieldDefs: []dsl.AnswerFieldDef{
				{ID: "w1", Expr: "w[1]"}, {ID: "w2", Expr: "w[2]"}, {ID: "w3", Expr: "w[3]"}, {ID: "w4", Expr: "w[4]"},
			},
		},
	}
	_, err := dsl.GenerateQuestion(p, "null-seed", "s")
	if err != nil {
		t.Fatal(err)
	}
}
