package dsl

import (
	"math/big"
	"testing"
)

func TestVectorsRationalCollinear_parallelAndZero(t *testing.T) {
	u := []*big.Rat{big.NewRat(1, 1), big.NewRat(2, 1), big.NewRat(3, 1)}
	v := []*big.Rat{big.NewRat(2, 1), big.NewRat(4, 1), big.NewRat(6, 1)}
	ok, note := vectorsRationalCollinear(u, v)
	if !ok || note != "" {
		t.Fatalf("parallel: ok=%v note=%q", ok, note)
	}
	z0 := []*big.Rat{big.NewRat(0, 1), big.NewRat(0, 1), big.NewRat(0, 1)}
	ok, _ = vectorsRationalCollinear(z0, z0)
	if !ok {
		t.Fatal("both zero vectors should match")
	}
	bad := []*big.Rat{big.NewRat(1, 1), big.NewRat(0, 1), big.NewRat(0, 1)}
	ok, _ = vectorsRationalCollinear(u, bad)
	if ok {
		t.Fatal("not collinear")
	}
}

func TestAffineSolutionOK_2x4(t *testing.T) {
	A := NewMatrixInt(2, 4)
	A.A[0] = []int64{1, 0, 1, 0}
	A.A[1] = []int64{0, 1, 0, 1}
	b := NewVectorInt(2)
	b.V[0], b.V[1] = 5, 7
	x := []*big.Rat{big.NewRat(1, 1), big.NewRat(2, 1), big.NewRat(4, 1), big.NewRat(5, 1)}
	ok, note := affineSolutionOK(A, b, x)
	if !ok || note != "" {
		t.Fatalf("expected solution ok=%v note=%q", ok, note)
	}
	bad := []*big.Rat{big.NewRat(1, 1), big.NewRat(0, 1), big.NewRat(0, 1), big.NewRat(0, 1)}
	ok, _ = affineSolutionOK(A, b, bad)
	if ok {
		t.Fatal("wrong x should fail")
	}
}

func TestSortedBasisColumnsOK_identity3(t *testing.T) {
	M := NewMatrixInt(3, 3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == j {
				M.A[i][j] = 1
			}
		}
	}
	ok, note := sortedBasisColumnsOK(M, 3, []string{"3", "1", "2"})
	if !ok || note != "" {
		t.Fatalf("perm basis: ok=%v note=%q", ok, note)
	}
	ok, _ = sortedBasisColumnsOK(M, 3, []string{"1", "2"})
	if ok {
		t.Fatal("too few cols should fail")
	}
	ok, _ = sortedBasisColumnsOK(M, 3, []string{"0", "1", "2"})
	if ok {
		t.Fatal("only two non-zero cols should not span rank 3")
	}
}

func TestColumnSubsetRank_subset(t *testing.T) {
	M := NewMatrixInt(2, 3)
	M.A[0] = []int64{1, 0, 2}
	M.A[1] = []int64{0, 1, 0}
	if r := columnSubsetRank(M, []int{1, 2}); r != 2 {
		t.Fatalf("col1 col2 indep want rank 2 got %d", r)
	}
	if r := columnSubsetRank(M, []int{1, 3}); r != 1 {
		t.Fatalf("col3 dep on col1 want rank 1 got %d", r)
	}
}

func TestMatrixIntTimesVectorRat(t *testing.T) {
	A := NewMatrixInt(2, 2)
	A.A[0] = []int64{1, 2}
	A.A[1] = []int64{3, 4}
	x := []*big.Rat{big.NewRat(1, 2), big.NewRat(-1, 3)}
	y, err := matrixIntTimesVectorRat(A, x)
	if err != nil {
		t.Fatal(err)
	}
	if y[0].Cmp(big.NewRat(-1, 6)) != 0 || y[1].Cmp(big.NewRat(1, 6)) != 0 {
		t.Fatalf("bad y %v %v", y[0].RatString(), y[1].RatString())
	}
}
