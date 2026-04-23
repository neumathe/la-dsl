package dsl

import (
	"math/big"
	"testing"
)

func TestJudgeGeneratedQuestionRat(t *testing.T) {
	g := &GeneratedQuestion{
		AnswerFields: []AnswerField{
			{ID: "a", Value: big.NewRat(1, 2)},
			{ID: "b", Value: int64(3)},
		},
	}
	r := JudgeGeneratedQuestion(g, map[string]string{"a": "0.5", "b": "3"}, nil)
	if !r.AllCorrect || r.CorrectCount != 2 {
		t.Fatalf("got %+v", r)
	}
}

func TestParseUserRationalUnicodeMinusAndFractionSlash(t *testing.T) {
	r, err := ParseUserRational("\u22123\u20442") // −3⁄2
	if err != nil {
		t.Fatal(err)
	}
	if r.Cmp(big.NewRat(-3, 2)) != 0 {
		t.Fatalf("got %s", r.RatString())
	}
}

func TestJudgeUnicodeInput(t *testing.T) {
	g := &GeneratedQuestion{
		AnswerFields: []AnswerField{{ID: "x", Value: big.NewRat(-3, 7)}},
	}
	r := JudgeGeneratedQuestion(g, map[string]string{"x": "\u22123／7"}, nil)
	if !r.AllCorrect {
		t.Fatalf("%+v", r.Fields[0])
	}
}

func TestJudgeWrong(t *testing.T) {
	g := &GeneratedQuestion{
		AnswerFields: []AnswerField{{ID: "x", Value: int64(7)}},
	}
	r := JudgeGeneratedQuestion(g, map[string]string{"x": "8"}, nil)
	if r.AllCorrect || r.Fields[0].Correct {
		t.Fatal("expected wrong")
	}
}
