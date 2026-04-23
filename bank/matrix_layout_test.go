package bank

import (
	"testing"

	"github.com/neumathe/la-dsl/dsl"
)

func TestChapter2_4_1MatrixLayouts(t *testing.T) {
	p, err := BuildProblem("Chapter2_4_1")
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Answer.FieldDefs) != 9 {
		t.Fatalf("field count %d", len(p.Answer.FieldDefs))
	}
	for i, fd := range p.Answer.FieldDefs {
		if fd.Layout == nil || fd.Layout.Kind != dsl.LayoutKindMatrixCell {
			t.Fatalf("field %d missing matrix layout", i)
		}
		if fd.Layout.Row < 1 || fd.Layout.Row > 3 || fd.Layout.Col < 1 || fd.Layout.Col > 3 {
			t.Fatalf("bad rc: %+v", fd.Layout)
		}
	}
}
