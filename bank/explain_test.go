package bank

import (
	"testing"

	"github.com/neumathe/la-dsl/dsl"
)

func TestGenerateBankExplanationAlignedWithQuestion(t *testing.T) {
	key := "Chapter1_1"
	seed, salt := "explain-bank-seed", "explain-salt"
	g, err := GenerateBankQuestion(key, seed, salt)
	if err != nil {
		t.Fatal(err)
	}
	ex, err := GenerateBankExplanation(key, seed, salt)
	if err != nil {
		t.Fatal(err)
	}
	if ex.Title != g.Title {
		t.Fatalf("title: %q vs %q", ex.Title, g.Title)
	}
	if len(ex.AnswerSteps) != len(g.AnswerFields) {
		t.Fatalf("len %d vs %d", len(ex.AnswerSteps), len(g.AnswerFields))
	}
	for i, f := range g.AnswerFields {
		if ex.AnswerSteps[i].Expected != dsl.ValueToCanonicalString(f.Value) {
			t.Fatalf("field %s", f.ID)
		}
	}
}
