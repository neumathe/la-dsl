package dsl

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
)

// genMatrixParamInfinitSolution generates a 3×3 coefficient matrix A and 3×1 RHS vector b
// for a system Ax=b where one entry of A depends on λ and one entry of b depends on μ,
// such that the system has infinitely many solutions when λ and μ take specific integer values.
//
// Strategy:
//   1. Generate first 2 rows randomly (full rank 2×3)
//   2. Make row 3 a linear combination of rows 1 and 2, with λ replacing one coefficient
//   3. Choose b such that when λ makes rank(A)=2, b is in the column space (rank(A|b)=2 too)
//   4. Specifically: b[3] = μ which equals the same linear combo of b[1], b[2]
//
// The specific values λ_val and μ_val are stored in inst.Vars.
// The instantiated matrix (with λ_val, μ_val substituted) is returned.
func genMatrixParamInfinitSolution(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	r := v.Rows
	c := v.Cols
	if r != 3 || c != 3 {
		return nil, errors.New("param_infinit_solution: need 3×3")
	}

	lmin := int64(defaultInt(v.Generator, "entry_min", -5))
	lmax := int64(defaultInt(v.Generator, "entry_max", 5))
	if lmax < lmin {
		lmin, lmax = lmax, lmin
	}
	lambdaMin := int64(defaultInt(v.Generator, "lambda_min", -10))
	lambdaMax := int64(defaultInt(v.Generator, "lambda_max", 10))
	if lambdaMax < lambdaMin {
		lambdaMin, lambdaMax = lambdaMax, lambdaMin
	}

	for attempt := 0; attempt < 500; attempt++ {
		attemptSeed := rng.Int63()
		attemptRng := rand.New(rand.NewSource(attemptSeed))

		// Generate first 2 rows of A (3×3), ensuring row 1 and row 2 are linearly independent
		A := NewMatrixInt(3, 3)
		for i := 0; i < 2; i++ {
			for j := 0; j < 3; j++ {
				A.A[i][j] = int64(attemptRng.Intn(int(lmax-lmin+1)) + int(lmin))
			}
		}
		row1 := []int64{A.A[0][0], A.A[0][1], A.A[0][2]}
		row2 := []int64{A.A[1][0], A.A[1][1], A.A[1][2]}
		if isProportionalInt64(row1, row2) {
			continue
		}

		// Generate b[1], b[2]
		b := NewVectorInt(3)
		b.V[0] = int64(attemptRng.Intn(int(lmax-lmin+1)) + int(lmin))
		b.V[1] = int64(attemptRng.Intn(int(lmax-lmin+1)) + int(lmin))

		// Row 3 = k1 * row1 + k2 * row2, with λ at position (3, colLambda)
		k1 := int64(attemptRng.Intn(5) - 2) // -2 to 2
		k2 := int64(attemptRng.Intn(5) - 2)
		if k1 == 0 && k2 == 0 {
			continue
		}

		colLambda := attemptRng.Intn(3) // 0,1,2 → col 1,2,3

		for j := 0; j < 3; j++ {
			if j == colLambda {
				continue
			}
			A.A[2][j] = k1*A.A[0][j] + k2*A.A[1][j]
		}

		lambdaVal := k1*A.A[0][colLambda] + k2*A.A[1][colLambda]
		if lambdaVal < lambdaMin || lambdaVal > lambdaMax {
			continue
		}

		muVal := k1*b.V[0] + k2*b.V[1]
		if muVal < -20 || muVal > 20 {
			continue
		}

		A.A[2][colLambda] = lambdaVal
		b.V[2] = muVal

		if matrixRankRat(A) != 2 {
			continue
		}
		aug := NewMatrixInt(3, 4)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				aug.A[i][j] = A.A[i][j]
			}
			aug.A[i][3] = b.V[i]
		}
		if matrixRankRat(aug) != 2 {
			continue
		}

		// Compute RREF of augmented matrix [A|b]
		rref, err := RrefRatSafe(aug)
		if err != nil {
			continue
		}

		// Verify RREF entries are all integers
		allInt := true
		for i := 0; i < 3; i++ {
			for j := 0; j < 4; j++ {
				d := rref[i][j].Denom()
				if d.Sign() != 0 && d.Cmp(big.NewInt(1)) != 0 {
					allInt = false
					break
				}
			}
			if !allInt {
				break
			}
		}
		if !allInt {
			continue
		}

		rrefMatrix := NewMatrixInt(3, 4)
		for i := 0; i < 3; i++ {
			for j := 0; j < 4; j++ {
				rrefMatrix.A[i][j] = rref[i][j].Num().Int64()
			}
		}

		// Compute nullspace basis of A (should have 1 vector since 3-2=1)
		nb, err := NullspaceBasisRational(A)
		if err != nil || len(nb) != 1 {
			continue
		}

		// Compute a particular solution from RREF
		x0 := make([]*big.Rat, 3)
		for i := 0; i < 3; i++ {
			x0[i] = rref[i][3] // b column in RREF (free vars set to 0)
		}

		// Verify x0 is a valid solution: A*x0 = b
		valid := true
		for i := 0; i < 3; i++ {
			sum := new(big.Rat)
			for j := 0; j < 3; j++ {
				sum.Add(sum, new(big.Rat).Mul(new(big.Rat).SetInt64(A.A[i][j]), x0[j]))
			}
			if sum.Cmp(new(big.Rat).SetInt64(b.V[i])) != 0 {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}

		// Store all data
		if inst.Vars == nil {
			inst.Vars = map[string]interface{}{}
		}
		inst.Vars["_param_lambda_val"] = lambdaVal
		inst.Vars["_param_mu_val"] = muVal
		inst.Vars["_param_b"] = b
		inst.Vars["_param_rref"] = rrefMatrix
		inst.Vars["_param_x0"] = ratVecToIntVec(x0)
		inst.Vars["_param_nb1"] = ratVecToIntVec(nb[0])
		inst.Vars["_param_colLambda"] = int64(colLambda + 1) // 1-based
		inst.Vars["_param_k1"] = k1
		inst.Vars["_param_k2"] = k2
		inst.Vars["_param_A"] = A

		return A, nil
	}
	return nil, fmt.Errorf("param_infinit_solution: failed after attempts")
}

func isProportionalInt64(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	// Find first nonzero pair to determine ratio
	ratio := new(big.Rat)
	found := false
	for i := 0; i < len(a); i++ {
		if a[i] == 0 && b[i] == 0 {
			continue
		}
		if a[i] == 0 || b[i] == 0 {
			return false // one is zero, other isn't → not proportional
		}
		ratio.SetInt64(a[i])
		ratio.Quo(ratio, new(big.Rat).SetInt64(b[i]))
		found = true
		break
	}
	if !found {
		return true // both all zeros
	}
	for i := 0; i < len(a); i++ {
		if a[i] == 0 && b[i] == 0 {
			continue
		}
		if a[i] == 0 || b[i] == 0 {
			return false
		}
		actual := new(big.Rat).SetInt64(a[i])
		actual.Quo(actual, new(big.Rat).SetInt64(b[i]))
		if actual.Cmp(ratio) != 0 {
			return false
		}
	}
	return true
}

func containsInt(arr []int, v int) bool {
	for _, a := range arr {
		if a == v {
			return true
		}
	}
	return false
}

func findPivotCols3x4(M [][]*big.Rat) []int {
	pivots := []int{}
	row := 0
	for col := 0; col < 4 && row < 3; col++ {
		if M[row][col].Sign() != 0 {
			pivots = append(pivots, col)
			row++
		}
	}
	return pivots
}

func ratVecToIntVec(ratVec []*big.Rat) *VectorInt {
	lcm := big.NewInt(1)
	for _, r := range ratVec {
		d := r.Denom()
		if d.Sign() != 0 {
			gcd := new(big.Int).GCD(nil, nil, lcm, d)
			newLcm := new(big.Int).Mul(lcm, d)
			newLcm.Quo(newLcm, gcd)
			lcm = newLcm
		}
	}
	out := NewVectorInt(len(ratVec))
	for i, r := range ratVec {
		scaled := new(big.Rat).Mul(r, new(big.Rat).SetInt(lcm))
		out.V[i] = scaled.Num().Int64()
	}
	return out
}

// RrefRat computes the reduced row echelon form of a rational matrix.
func RrefRat(M *MatrixInt) [][]*big.Rat {
	r, c := M.R, M.C
	mat := make([][]*big.Rat, r)
	for i := 0; i < r; i++ {
		mat[i] = make([]*big.Rat, c)
		for j := 0; j < c; j++ {
			mat[i][j] = new(big.Rat).SetInt64(M.A[i][j])
		}
	}

	row := 0
	for col := 0; col < c && row < r; col++ {
		pivotRow := -1
		for i := row; i < r; i++ {
			if mat[i][col].Sign() != 0 {
				pivotRow = i
				break
			}
		}
		if pivotRow == -1 {
			continue
		}
		if pivotRow != row {
			mat[pivotRow], mat[row] = mat[pivotRow], mat[row]
		}
		pv := new(big.Rat).Set(mat[row][col])
		for j := col; j < c; j++ {
			mat[row][j] = new(big.Rat).Quo(mat[row][j], pv)
		}
		for i := 0; i < r; i++ {
			if i == row {
				continue
			}
			f := new(big.Rat).Set(mat[i][col])
			if f.Sign() == 0 {
				continue
			}
			for j := col; j < c; j++ {
				tmp := new(big.Rat).Mul(f, mat[row][j])
				mat[i][j] = new(big.Rat).Sub(mat[i][j], tmp)
			}
		}
		row++
	}

	return mat
}

// RrefRatSafe computes RREF with explicit new allocations to avoid aliasing issues.
func RrefRatSafe(M *MatrixInt) ([][]*big.Rat, error) {
	r, c := M.R, M.C
	mat := make([][]*big.Rat, r)
	for i := 0; i < r; i++ {
		mat[i] = make([]*big.Rat, c)
		for j := 0; j < c; j++ {
			mat[i][j] = new(big.Rat).SetInt64(M.A[i][j])
		}
	}

	row := 0
	for col := 0; col < c && row < r; col++ {
		pivotRow := -1
		for i := row; i < r; i++ {
			if mat[i][col].Sign() != 0 {
				pivotRow = i
				break
			}
		}
		if pivotRow == -1 {
			continue
		}
		if pivotRow != row {
			mat[pivotRow], mat[row] = mat[pivotRow], mat[row]
		}
		pv := new(big.Rat).Set(mat[row][col])
		if pv.Sign() == 0 {
			return nil, fmt.Errorf("RREF: zero pivot")
		}
		for j := col; j < c; j++ {
			old := new(big.Rat).Set(mat[row][j])
			mat[row][j] = new(big.Rat).Quo(old, pv)
		}
		for i := 0; i < r; i++ {
			if i == row {
				continue
			}
			f := new(big.Rat).Set(mat[i][col])
			if f.Sign() == 0 {
				continue
			}
			for j := col; j < c; j++ {
				tmp := new(big.Rat).Mul(f, mat[row][j])
				old := new(big.Rat).Set(mat[i][j])
				mat[i][j] = new(big.Rat).Sub(old, tmp)
			}
		}
		row++
	}

	return mat, nil
}