package dsl

import (
	"math/big"
	"testing"
)

func TestEvaluateExpression_rank_nullity_rank_hstack(t *testing.T) {
	inst := &Instance{Vars: map[string]interface{}{}}
	A := NewMatrixInt(3, 4)
	A.A[0] = []int64{1, 0, 0, 0}
	A.A[1] = []int64{0, 1, 0, 0}
	A.A[2] = []int64{1, 1, 0, 0}
	b := NewVectorInt(3)
	b.V[0], b.V[1], b.V[2] = 1, 2, 3
	inst.Vars["A"] = A
	inst.Vars["b3"] = b
	v, err := EvaluateExpression("rank(A)", inst)
	if err != nil {
		t.Fatal(err)
	}
	if v != int64(2) {
		t.Fatalf("rank want 2 got %v", v)
	}
	v, err = EvaluateExpression("nullity(A)", inst)
	if err != nil || v != int64(2) {
		t.Fatalf("nullity want 2 got %v err=%v", v, err)
	}
	v, err = EvaluateExpression("rank_hstack(A,b3)", inst)
	if err != nil || v != int64(2) {
		t.Fatalf("rank_hstack want 2 got %v err=%v", v, err)
	}
}

func TestInertiaCounts22(t *testing.T) {
	pos, neg, err := inertiaCounts22(nil)
	if err == nil {
		t.Fatal("nil expect err")
	}
	_ = pos
	_ = neg
	S := NewMatrixInt(2, 2)
	S.A[0] = []int64{1, 0}
	S.A[1] = []int64{0, 1}
	p, n, err := inertiaCounts22(S)
	if err != nil || p != 2 || n != 0 {
		t.Fatalf("I pos=%d neg=%d err=%v", p, n, err)
	}
	S.A[0][0], S.A[0][1], S.A[1][0], S.A[1][1] = -1, 0, 0, -1
	p, n, err = inertiaCounts22(S)
	if err != nil || p != 0 || n != 2 {
		t.Fatalf("neg def pos=%d neg=%d", p, n)
	}
	S.A[0][0], S.A[0][1], S.A[1][0], S.A[1][1] = 1, 2, 2, 1
	p, n, err = inertiaCounts22(S)
	if err != nil || p != 1 || n != 1 {
		t.Fatalf("indef pos=%d neg=%d", p, n)
	}
}

func TestEvaluateExpression_inertia_expr(t *testing.T) {
	inst := &Instance{Vars: map[string]interface{}{}}
	S := NewMatrixInt(2, 2)
	S.A[0] = []int64{2, 0}
	S.A[1] = []int64{0, 3}
	inst.Vars["S"] = S
	v, err := EvaluateExpression("inertia_pos_22(S)", inst)
	if err != nil || v != int64(2) {
		t.Fatalf("inertia_pos got %v %v", v, err)
	}
	v, err = EvaluateExpression("inertia_neg_22(S)", inst)
	if err != nil || v != int64(0) {
		t.Fatalf("inertia_neg got %v %v", v, err)
	}
}

func TestEvaluateExpression_orthdiag_block4x6(t *testing.T) {
	inst := &Instance{Vars: map[string]interface{}{}}
	P := NewMatrixInt(3, 3)
	P.A[0][0], P.A[1][1], P.A[2][2] = 1, 1, 1
	D := NewMatrixInt(3, 3)
	D.A[0][0], D.A[1][1], D.A[2][2] = 2, 3, 4
	S := NewMatrixInt(3, 3)
	for i := 0; i < 3; i++ {
		S.A[i][i] = D.A[i][i]
	}
	inst.Vars["P"] = P
	inst.Vars["D"] = D
	inst.Vars["S"] = S
	v, err := EvaluateExpression("orthdiag_block4x6(P,D,S)", inst)
	if err != nil {
		t.Fatal(err)
	}
	Z, ok := v.(*MatrixInt)
	if !ok || Z.R != 4 || Z.C != 6 {
		t.Fatalf("Z shape %+v", Z)
	}
	if Z.A[3][0] != 9 || Z.A[3][1] != 24 {
		t.Fatalf("tr/det row want 9,24 got %d,%d", Z.A[3][0], Z.A[3][1])
	}
}

func TestGramSchmidtColsOrthogRatOrthogonal(t *testing.T) {
	V := NewMatrixInt(3, 3)
	V.A[0] = []int64{1, 1, 0}
	V.A[1] = []int64{0, 1, 1}
	V.A[2] = []int64{1, 0, 1}
	u, err := GramSchmidtColsOrthogRat(V)
	if err != nil {
		t.Fatal(err)
	}
	dot := func(a, b []*big.Rat) *big.Rat {
		s := big.NewRat(0, 1)
		for i := 0; i < 3; i++ {
			s.Add(s, new(big.Rat).Mul(a[i], b[i]))
		}
		return s
	}
	if dot(u[0], u[1]).Sign() != 0 || dot(u[0], u[2]).Sign() != 0 || dot(u[1], u[2]).Sign() != 0 {
		t.Fatal("not pairwise orthogonal")
	}
	V4 := NewMatrixInt(4, 2)
	V4.A[0][0], V4.A[1][0], V4.A[2][0], V4.A[3][0] = 1, 0, 0, 0
	V4.A[0][1], V4.A[1][1], V4.A[2][1], V4.A[3][1] = 0, 1, 1, 0
	u4, err := GramSchmidtColsOrthogRatGeneral(V4)
	if err != nil {
		t.Fatal(err)
	}
	if len(u4) != 2 || len(u4[0]) != 4 {
		t.Fatalf("bad dim %d %d", len(u4), len(u4[0]))
	}
	dot4 := func(a, b []*big.Rat) *big.Rat {
		s := big.NewRat(0, 1)
		for i := 0; i < 4; i++ {
			s.Add(s, new(big.Rat).Mul(a[i], b[i]))
		}
		return s
	}
	if dot4(u4[0], u4[1]).Sign() != 0 {
		t.Fatal("4x2 GS not orthogonal")
	}
}
