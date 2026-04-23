package dsl

import (
	"math/big"
	"testing"
)

// 固定 2×4 秩 2 与特解，覆盖 affine_rational + 两组 rational_line 判题。
func TestJudgeGeneratedQuestionContext_affineAndLineGroups(t *testing.T) {
	ja := &AnswerJudgeSpec{Kind: "affine_rational", AffineGroup: "eta", MatrixVar: "A", BVecVar: "b2"}
	j1 := &AnswerJudgeSpec{Kind: "rational_line", LineGroup: "xi1"}
	j2 := &AnswerJudgeSpec{Kind: "rational_line", LineGroup: "xi2"}
	p := Problem{
		ID:      91001,
		Version: "test-v1",
		Title:   "t",
		Variables: map[string]Variable{
			"A": {Kind: "matrix", Rows: 2, Cols: 4, Fixed: [][]interface{}{
				{1, 0, 1, 0},
				{0, 1, 0, 1},
			}},
			"x0": {Kind: "vector", Size: 4, Fixed: []interface{}{1, 1, 1, 1}},
		},
		Derived: map[string]string{"b2": "A * x0"},
		Answer: AnswerSchema{FieldDefs: []AnswerFieldDef{
			{ID: "p1", Expr: "x0[1]", Judge: ja},
			{ID: "p2", Expr: "x0[2]", Judge: ja},
			{ID: "p3", Expr: "x0[3]", Judge: ja},
			{ID: "p4", Expr: "x0[4]", Judge: ja},
			{ID: "u1", Expr: "nullbasis_comp(A,1,1)", Judge: j1},
			{ID: "u2", Expr: "nullbasis_comp(A,1,2)", Judge: j1},
			{ID: "u3", Expr: "nullbasis_comp(A,1,3)", Judge: j1},
			{ID: "u4", Expr: "nullbasis_comp(A,1,4)", Judge: j1},
			{ID: "v1", Expr: "nullbasis_comp(A,2,1)", Judge: j2},
			{ID: "v2", Expr: "nullbasis_comp(A,2,2)", Judge: j2},
			{ID: "v3", Expr: "nullbasis_comp(A,2,3)", Judge: j2},
			{ID: "v4", Expr: "nullbasis_comp(A,2,4)", Judge: j2},
		}},
	}
	inst, err := InstantiateProblem(p, "seed-jctx-1", "salt")
	if err != nil {
		t.Fatal(err)
	}
	g, err := GenerateQuestionFromInstance(p, inst)
	if err != nil {
		t.Fatal(err)
	}
	user := map[string]string{}
	for _, f := range g.AnswerFields {
		user[f.ID] = ValueToCanonicalString(f.Value)
	}
	res := JudgeGeneratedQuestionContext(g, &p, inst, user, nil)
	if !res.AllCorrect {
		t.Fatalf("exact match should pass: %+v", res.Fields)
	}
	// 齐次方向整体 2 倍仍共线
	user2 := map[string]string{}
	for _, f := range g.AnswerFields {
		user2[f.ID] = ValueToCanonicalString(f.Value)
	}
	for _, id := range []string{"u1", "u2", "u3", "u4", "v1", "v2", "v3", "v4"} {
		for _, f := range g.AnswerFields {
			if f.ID != id {
				continue
			}
			if rv, ok := f.Value.(*big.Rat); ok {
				d := new(big.Rat).Mul(rv, big.NewRat(2, 1))
				user2[id] = d.String()
			}
		}
	}
	res2 := JudgeGeneratedQuestionContext(g, &p, inst, user2, nil)
	if !res2.AllCorrect {
		t.Fatalf("scaled null should pass: %+v", res2.Fields)
	}
}

func TestJudgeGeneratedQuestionContext_sortedBasis(t *testing.T) {
	jb := &AnswerJudgeSpec{Kind: "sorted_basis_columns", BasisGroup: "bc", MatrixVar: "M", Ncols: 3}
	p := Problem{
		ID:      91002,
		Version: "test-v1",
		Title:   "t",
		Variables: map[string]Variable{
			"M": {Kind: "matrix", Rows: 3, Cols: 3, Fixed: [][]interface{}{
				{1, 0, 0},
				{0, 2, 0},
				{0, 0, 3},
			}},
		},
		Answer: AnswerSchema{FieldDefs: []AnswerFieldDef{
			{ID: "c1", Expr: "mget(M,1,1)", Judge: jb},
			{ID: "c2", Expr: "mget(M,2,2)", Judge: jb},
			{ID: "c3", Expr: "mget(M,3,3)", Judge: jb},
		}},
	}
	inst, err := InstantiateProblem(p, "seed-jctx-2", "salt")
	if err != nil {
		t.Fatal(err)
	}
	g, err := GenerateQuestionFromInstance(p, inst)
	if err != nil {
		t.Fatal(err)
	}
	user := map[string]string{"c1": "3", "c2": "1", "c3": "2"}
	res := JudgeGeneratedQuestionContext(g, &p, inst, user, nil)
	if !res.AllCorrect {
		t.Fatalf("basis perm: %+v", res.Fields)
	}
}

func TestJudgeGeneratedQuestionContext_nilInstStructuredFails(t *testing.T) {
	g := &GeneratedQuestion{
		AnswerFields: []AnswerField{
			{ID: "a", Value: int64(1), Judge: &AnswerJudgeSpec{Kind: "rational_line", LineGroup: "w"}},
		},
	}
	res := JudgeGeneratedQuestionContext(g, nil, nil, map[string]string{"a": "1"}, nil)
	if res.AllCorrect || res.Fields[0].Correct {
		t.Fatal("without inst rational_line should not pass as scalar")
	}
}

func TestGenerateQuestionFromInstance_nilInst(t *testing.T) {
	_, err := GenerateQuestionFromInstance(Problem{}, nil)
	if err == nil {
		t.Fatal("expect error")
	}
}
