package bank

import (
	"testing"

	"github.com/neumathe/la-dsl/dsl"
)

func TestJudgeBankQuestion_Chapter4_8_roundTrip(t *testing.T) {
	p, err := BuildProblem("Chapter4_8")
	if err != nil {
		t.Fatal(err)
	}
	seed, salt := "bank-judge-rt-1", "salt"
	inst, err := dsl.InstantiateProblem(p, seed, salt)
	if err != nil {
		t.Fatal(err)
	}
	g, err := dsl.GenerateQuestionFromInstance(p, inst)
	if err != nil {
		t.Fatal(err)
	}
	user := map[string]string{}
	for _, f := range g.AnswerFields {
		user[f.ID] = dsl.ValueToCanonicalString(f.Value)
	}
	res, err := JudgeBankQuestion("Chapter4_8", seed, salt, user, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !res.AllCorrect {
		t.Fatalf("round-trip judge: %+v", res.Fields)
	}
}
