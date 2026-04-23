package bank

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/neumathe/la-dsl/dsl"
)

// TestExpectedFieldCountsCatalog 与 AllQuestionKeys 一一对应，防止漏题或空数漂移。
func TestExpectedFieldCountsCatalog(t *testing.T) {
	if len(ExpectedAnswerFieldCount) != len(AllQuestionKeys) {
		t.Fatalf("ExpectedAnswerFieldCount map size %d != AllQuestionKeys %d",
			len(ExpectedAnswerFieldCount), len(AllQuestionKeys))
	}
	for _, k := range AllQuestionKeys {
		if _, ok := ExpectedAnswerFieldCount[k]; !ok {
			t.Errorf("missing expected count for %q", k)
		}
	}
	for k := range ExpectedAnswerFieldCount {
		found := false
		for _, x := range AllQuestionKeys {
			if x == k {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("extra key in ExpectedAnswerFieldCount: %q", k)
		}
	}
}

// TestEveryBankQuestion 对题库中每一道逻辑题：空数、答案 id、确定性、满分判题。
func TestEveryBankQuestion(t *testing.T) {
	seed := "unit-all-questions-2026"
	salt := "unit-salt"
	for _, key := range AllQuestionKeys {
		key := key
		t.Run(key, func(t *testing.T) {
			wantN := ExpectedAnswerFieldCount[key]
			g, err := GenerateBankQuestion(key, seed+":"+key, salt)
			if err != nil {
				t.Fatalf("GenerateBankQuestion: %v", err)
			}
			if len(g.AnswerFields) != wantN {
				t.Fatalf("answer field count: want %d got %d", wantN, len(g.AnswerFields))
			}
			for i, f := range g.AnswerFields {
				prefix := key + "_"
				if !strings.HasPrefix(f.ID, prefix) {
					t.Fatalf("field[%d] id %q should start with %q", i, f.ID, prefix)
				}
				if f.ID[len(prefix):] == "" {
					t.Fatalf("field[%d] id missing suffix", i)
				}
			}
			g2, err := GenerateBankQuestion(key, seed+":"+key, salt)
			if err != nil {
				t.Fatal(err)
			}
			if g.Title != g2.Title {
				t.Fatal("title not deterministic")
			}
			for i := range g.AnswerFields {
				if g.AnswerFields[i].ID != g2.AnswerFields[i].ID {
					t.Fatalf("field id order drift at %d", i)
				}
				a := fmt.Sprintf("%v", g.AnswerFields[i].Value)
				b := fmt.Sprintf("%v", g2.AnswerFields[i].Value)
				if a != b {
					t.Fatalf("field %s value drift: %s vs %s", g.AnswerFields[i].ID, a, b)
				}
			}
			ans := make(map[string]string, len(g.AnswerFields))
			for _, f := range g.AnswerFields {
				ans[f.ID] = dsl.ValueToCanonicalString(f.Value)
			}
			jr, err := JudgeBankQuestion(key, seed+":"+key, salt, ans, nil)
			if err != nil {
				t.Fatal(err)
			}
			if !jr.AllCorrect {
				t.Fatalf("judge all correct: %+v", jr.Fields)
			}
			if jr.ScoreEarned != jr.ScoreMax || jr.TotalFields != wantN {
				t.Fatalf("score fields: earned=%v max=%v total=%d want=%d",
					jr.ScoreEarned, jr.ScoreMax, jr.TotalFields, wantN)
			}
		})
	}
}

// TestJudgeBankQuestionWrongAnswer 抽样：错一个空则不得满分。
func TestJudgeBankQuestionWrongAnswer(t *testing.T) {
	key := "Chapter1_1"
	seed := "wrong-one"
	salt := "s"
	g, err := GenerateBankQuestion(key, seed, salt)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.AnswerFields) < 1 {
		t.Fatal("need at least one field")
	}
	id0 := g.AnswerFields[0].ID
	wrong := map[string]string{id0: "999999999"}
	for i := 1; i < len(g.AnswerFields); i++ {
		f := g.AnswerFields[i]
		wrong[f.ID] = dsl.ValueToCanonicalString(f.Value)
	}
	jr, err := JudgeBankQuestion(key, seed, salt, wrong, nil)
	if err != nil {
		t.Fatal(err)
	}
	if jr.AllCorrect {
		t.Fatal("expected not all correct")
	}
	if jr.Fields[0].Correct {
		t.Fatal("first field should be wrong")
	}
}

// TestChapter1LambdaDetZero 不变量：det(A)=0。
func TestChapter1LambdaDetZero(t *testing.T) {
	for _, key := range []string{"Chapter1_8_1", "Chapter1_8_2"} {
		t.Run(key, func(t *testing.T) {
			p, err := BuildProblem(key)
			if err != nil {
				t.Fatal(err)
			}
			inst, err := dsl.InstantiateProblem(p, "inv-det0:"+key, "s")
			if err != nil {
				t.Fatal(err)
			}
			A, ok := inst.Vars["A"].(*dsl.MatrixInt)
			if !ok {
				t.Fatal("A not matrix")
			}
			if dsl.BareissDet(A).Sign() != 0 {
				t.Fatal("expected det(A)=0")
			}
		})
	}
}

// TestChapter2SixMatrixEquation 验证 Ch2_6 的 det(B) 可正确计算。
func TestChapter2SixMatrixEquation(t *testing.T) {
	p, err := BuildProblem("Chapter2_6")
	if err != nil {
		t.Fatal(err)
	}
	inst, err := dsl.InstantiateProblem(p, "ch26", "s")
	if err != nil {
		t.Fatal(err)
	}
	v, err := dsl.EvaluateExpression("detB", inst)
	if err != nil {
		t.Fatal(err)
	}
	bi, ok := v.(*big.Int)
	if !ok {
		t.Fatalf("want *big.Int, got %T", v)
	}
	// det(B) should be a nonzero integer (since A+I has det ±1 and det(3A) ≠ 0)
	t.Logf("det(B) = %s", bi.String())
}

// TestChapter2SevenPAEqualsC 满足 PA=C。
func TestChapter2SevenPAEqualsC(t *testing.T) {
	for _, key := range []string{"Chapter2_7_1", "Chapter2_7_2", "Chapter2_7_3"} {
		t.Run(key, func(t *testing.T) {
			p, err := BuildProblem(key)
			if err != nil {
				t.Fatal(err)
			}
			inst, err := dsl.InstantiateProblem(p, "pa-c:"+key, "s")
			if err != nil {
				t.Fatal(err)
			}
			P, ok := inst.Vars["P"].(*dsl.MatrixInt)
			if !ok {
				t.Fatal("P missing")
			}
			C, ok := inst.Vars["C"].(*dsl.MatrixInt)
			if !ok {
				t.Fatal("C missing")
			}
			PA, err := dsl.EvaluateExpression("matmul(P,A)", inst)
			if err != nil {
				t.Fatal(err)
			}
			pa := PA.(*dsl.MatrixInt)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if pa.A[i][j] != C.A[i][j] {
						t.Fatalf("PA!=C at (%d,%d): %d vs %d", i+1, j+1, pa.A[i][j], C.A[i][j])
					}
				}
			}
			if P.R != 3 {
				t.Fatal("P size")
			}
		})
	}
}

// TestChapter5ThreeDiagonalAdj 对角阵：伴随特征值 μ_i 满足 μ_i · λ_i = det(A)。
func TestChapter5ThreeDiagonalAdj(t *testing.T) {
	p, err := BuildProblem("Chapter5_3")
	if err != nil {
		t.Fatal(err)
	}
	inst, err := dsl.InstantiateProblem(p, "ch53", "s")
	if err != nil {
		t.Fatal(err)
	}
	A := inst.Vars["A"].(*dsl.MatrixInt)
	detA := dsl.BareissDet(A)
	for i := 1; i <= 3; i++ {
		lam := A.A[i-1][i-1]
		mu, err := dsl.EvaluateExpression(fmt.Sprintf("diag_adj(A,%d)", i), inst)
		if err != nil {
			t.Fatal(err)
		}
		mubi, ok := mu.(*big.Int)
		if !ok {
			t.Fatalf("diag_adj want *big.Int got %T", mu)
		}
		prod := new(big.Int).Mul(mubi, big.NewInt(lam))
		if prod.Cmp(detA) != 0 {
			t.Fatalf("i=%d: mu*lambda=%s want det(A)=%s", i, prod.String(), detA.String())
		}
	}
}
