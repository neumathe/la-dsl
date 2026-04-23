package ladsl

import (
	"testing"

	"github.com/neumathe/la-dsl/bank"
)

func TestDescribeQuestionMatchesRoll(t *testing.T) {
	s := NewService("test-salt")
	key := "Chapter1_7"
	desc, err := DescribeQuestion(key)
	if err != nil {
		t.Fatal(err)
	}
	if desc.BlankCount != 3 || len(desc.Blanks) != 3 {
		t.Fatalf("blank count: %+v", desc)
	}
	pub, err := s.RollQuestion(key, "seed-ch1-7")
	if err != nil {
		t.Fatal(err)
	}
	if len(pub.Blanks) != len(desc.Blanks) {
		t.Fatalf("roll blanks %d vs desc %d", len(pub.Blanks), len(desc.Blanks))
	}
	for i := range desc.Blanks {
		if pub.Blanks[i].ID != desc.Blanks[i].ID {
			t.Fatalf("id[%d]: %q vs %q", i, pub.Blanks[i].ID, desc.Blanks[i].ID)
		}
	}
}

func TestDescribeAllKeysCatalog(t *testing.T) {
	for _, k := range bank.AllQuestionKeys {
		desc, err := DescribeQuestion(k)
		if err != nil {
			t.Fatalf("%s: %v", k, err)
		}
		want, ok := bank.ExpectedAnswerFieldCount[k]
		if !ok {
			t.Fatalf("missing expected count for %s", k)
		}
		if desc.BlankCount != want || len(desc.Blanks) != want {
			t.Fatalf("%s: want %d blanks, got descriptor %+v", k, want, desc)
		}
	}
}

func TestDescribeQuestionMatrixLayout(t *testing.T) {
	desc, err := DescribeQuestion("Chapter2_4_1")
	if err != nil {
		t.Fatal(err)
	}
	if len(desc.Blanks) != 9 {
		t.Fatal(len(desc.Blanks))
	}
	for _, b := range desc.Blanks {
		if b.Layout == nil || b.Layout.Kind != "matrix_cell" {
			t.Fatalf("blank %s layout %+v", b.ID, b.Layout)
		}
	}
}

func TestRollQuestionServer(t *testing.T) {
	s := NewService("srv-salt")
	b, err := s.RollQuestionServer("Chapter1_1", "one-seed")
	if err != nil {
		t.Fatal(err)
	}
	if b.Public.Title == "" {
		t.Fatal("empty title")
	}
	if len(b.Private.AnswerFields) != len(b.Public.Blanks) {
		t.Fatal("bundle mismatch")
	}
}

func TestValidQuestionKey(t *testing.T) {
	if !ValidQuestionKey("Chapter1_1") || ValidQuestionKey("nope") {
		t.Fatal()
	}
}
