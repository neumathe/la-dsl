package dsl

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
)

// SylvesterPositiveRangeSingleT computes the single-sided open interval t > lower
// for a 3×3 symmetric matrix S(t) to be positive definite, where t appears only
// at one diagonal position S[tRow][tRow].
//
// Sylvester's criterion:
//   Δ₁ > 0 (constant — must hold for the fixed diagonal entries)
//   Δ₂ > 0 (constant when t is at row 2+; t-independent when at row 0 or 1)
//   Δ₃ = cofactor · t + det(S₀) > 0
//
// Since the cofactor of the t-position is positive (by construction Δ₂ > 0),
// Δ₃ > 0 gives a single-sided bound: t > lower.
//
// Returns (lower, nil) — upper is always +∞.
func SylvesterPositiveRangeSingleT(Sbase *MatrixInt, tRow int) (lower *big.Rat, upper *big.Rat, err error) {
	if Sbase.R != 3 || Sbase.C != 3 {
		return nil, nil, fmt.Errorf("need 3×3")
	}

	// Δ₁ = S[1][1] (must be > 0)
	d1 := Sbase.A[0][0]
	if d1 <= 0 {
		return nil, nil, fmt.Errorf("Δ₁ = %d ≤ 0", d1)
	}

	// Δ₂ = S[1][1]·S[2][2] - S[1][2]²
	d2 := int64(Sbase.A[0][0]) * int64(Sbase.A[1][1]) - int64(Sbase.A[0][1])*int64(Sbase.A[0][1])
	if d2 <= 0 {
		return nil, nil, fmt.Errorf("Δ₂ = %d ≤ 0", d2)
	}

	// Δ₃ = cofactor[tRow] · t + det(S₀)
	// Compute cofactor by evaluating det at two t values and finding slope.
	S0 := copyMatrix(Sbase)
	S0.A[tRow][tRow] = 0

	S1 := copyMatrix(Sbase)
	S1.A[tRow][tRow] = 1

	det0 := BareissDet(S0)
	det1 := BareissDet(S1)

	// slope = cofactor = det1 - det0
	slope := new(big.Int).Sub(det1, det0)
	b0 := new(big.Int).Set(det0)

	if slope.Sign() == 0 {
		// det is independent of t
		if b0.Sign() > 0 {
			// always PD (no lower bound from Δ₃, but may have from Δ₁/Δ₂)
			return big.NewRat(0, 1), nil, nil
		}
		return nil, nil, fmt.Errorf("Δ₃ independent of t and ≤ 0")
	}

	// Δ₃ > 0: slope·t + b0 > 0
	// For PD we need slope > 0 (which is Δ₂ > 0 by construction when tRow >= 1)
	if slope.Sign() < 0 {
		return nil, nil, fmt.Errorf("cofactor at t-position is negative — no single-sided PD range")
	}

	// slope > 0: Δ₃ > 0 ⟹ t > -b0/slope
	lowerBound := new(big.Rat).SetFrac(new(big.Int).Neg(b0), slope)

	// When tRow = 2 (S[3][3]=t), Δ₁ and Δ₂ are t-independent constants.
	// When tRow = 1 (S[2][2]=t), Δ₂ depends on t: Δ₂ = S₁₁·t - S₁₂² > 0 → t > S₁₂²/S₁₁
	delta2Bound := big.NewRat(0, 1)
	if tRow == 1 {
		s12sq := int64(Sbase.A[0][1]) * int64(Sbase.A[0][1])
		delta2Bound = new(big.Rat).SetFrac(big.NewInt(s12sq), big.NewInt(int64(d1)))
	}

	// The final lower bound is the max of Δ₂ bound and Δ₃ bound.
	if lowerBound.Cmp(delta2Bound) >= 0 {
		return lowerBound, nil, nil
	}
	return delta2Bound, nil, nil
}

// genMatrixSylvesterRange generates a 3×3 symmetric matrix for the
// "find t range for positive definiteness" problem.
//
// Strategy (single-sided interval):
//   - t appears at one diagonal position only: S[tRow][tRow] = t (stored as 0 placeholder)
//   - All other diagonal entries are positive integers (necessary condition for PD)
//   - Off-diagonal entries are small integers
//   - The result is always a single-sided bound: t > lower
//
// The answer has one blank: the lower bound.
func genMatrixSylvesterRange(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	if v.Rows != 3 || v.Cols != 3 {
		return nil, fmt.Errorf("sylvester_range: need 3×3")
	}
	emin := int64(defaultInt(v.Generator, "entry_min", -5))
	emax := int64(defaultInt(v.Generator, "entry_max", 5))
	if emax < emin {
		emin, emax = emax, emin
	}

	// t position: prefer S[3][3] (matching HTML original pattern where t is x₃² coefficient)
	tRow := 2

	for attempt := 0; attempt < 500; attempt++ {
		attemptSeed := rng.Int63()
		attemptRng := rand.New(rand.NewSource(attemptSeed))

		// S[1][1] > 0: small positive integer
		s11 := int64(0)
		for t := 0; t < 40; t++ {
			x := int64(attemptRng.Intn(int(emax-emin+1)) + int(emin))
			if x > 0 && x <= emax {
				s11 = x
				break
			}
		}
		if s11 == 0 {
			continue
		}

		// S[2][2] > 0: small positive integer
		s22 := int64(0)
		for t := 0; t < 40; t++ {
			x := int64(attemptRng.Intn(int(emax-emin+1)) + int(emin))
			if x > 0 && x <= emax {
				s22 = x
				break
			}
		}
		if s22 == 0 {
			continue
		}

		// S[1][2] = S[2][1]: nonzero small integer
		s12 := int64(0)
		for t := 0; t < 40; t++ {
			x := int64(attemptRng.Intn(int(emax-emin+1)) + int(emin))
			if x != 0 {
				s12 = x
				break
			}
		}
		if s12 == 0 {
			continue
		}

		// S[1][3] = S[3][1]: nonzero small integer
		s13 := int64(0)
		for t := 0; t < 40; t++ {
			x := int64(attemptRng.Intn(int(emax-emin+1)) + int(emin))
			if x != 0 {
				s13 = x
				break
			}
		}
		if s13 == 0 {
			continue
		}

		// S[2][3] = S[3][2]: small integer (can be 0)
		s23 := int64(attemptRng.Intn(int(emax-emin+1)) + int(emin))

		// Build Sbase with S[tRow][tRow] = 0 (placeholder for t)
		Sbase := NewMatrixInt(3, 3)
		Sbase.A[0][0] = s11
		Sbase.A[0][1] = s12
		Sbase.A[0][2] = s13
		Sbase.A[1][0] = s12
		Sbase.A[1][1] = s22
		Sbase.A[1][2] = s23
		Sbase.A[2][0] = s13
		Sbase.A[2][1] = s23
		Sbase.A[2][2] = 0 // placeholder for t at S[3][3]

		// Verify Δ₂ > 0 (constant, must hold regardless of t)
		delta2 := s11*s22 - s12*s12
		if delta2 <= 0 {
			continue
		}

		// Compute Sylvester range (single-sided)
		lower, _, err := SylvesterPositiveRangeSingleT(Sbase, tRow)
		if err != nil {
			continue
		}

		// Check that the lower bound is a reasonable number (integer or simple fraction)
		if lower.Cmp(new(big.Rat).SetInt64(-20)) < 0 ||
			lower.Cmp(new(big.Rat).SetInt64(20)) > 0 {
			continue
		}

		// Prefer nice answers: integer or fraction with small denominator
		if lower.Denom().Int64() > 20 {
			continue
		}

		// Store for answer extraction
		if inst.Vars == nil {
			inst.Vars = map[string]interface{}{}
		}
		inst.Vars["_sylvester_Sbase"] = Sbase
		inst.Vars["_sylvester_lower"] = lower
		inst.Vars["_sylvester_s11"] = s11
		inst.Vars["_sylvester_s12"] = s12
		inst.Vars["_sylvester_s13"] = s13
		inst.Vars["_sylvester_s22"] = s22
		inst.Vars["_sylvester_s23"] = s23
		inst.Vars["_sylvester_tRow"] = tRow
		// For formatQuadraticExprWithParam: t at position (2,2) i.e. S[3][3] = t
		inst.Vars["_param_t_pos"] = []interface{}{float64(tRow), float64(tRow)}
		// No second t position — single t only

		return Sbase, nil
	}
	return nil, fmt.Errorf("sylvester_range: failed to generate suitable matrix")
}

// copyMatrix creates a deep copy of a MatrixInt
func copyMatrix(m *MatrixInt) *MatrixInt {
	c := NewMatrixInt(m.R, m.C)
	for i := 0; i < m.R; i++ {
		for j := 0; j < m.C; j++ {
			c.A[i][j] = m.A[i][j]
		}
	}
	return c
}

// formatQuadraticExprWithParam renders a quadratic form f(x₁,x₂,x₃)=... as LaTeX,
// with parameter t appearing at one diagonal position in the matrix.
//
// inst.Vars["_param_t_pos"]: 0-based (row,col) pair where S has t (e.g., (2,2) for S[3][3]=t)
// In the Sbase matrix, the placeholder position has value 0.
func formatQuadraticExprWithParam(Sbase *MatrixInt, inst *Instance) string {
	n := 3

	var tRow int = 2 // default: S[3][3] = t (0-based index 2)

	if tp, ok := inst.Vars["_param_t_pos"]; ok {
		switch t := tp.(type) {
		case []interface{}:
			tRow = int(t[0].(float64))
		}
	}

	xi := func(k int) string {
		return fmt.Sprintf("x_{%d}", k)
	}

	var terms []string

	// Diagonal terms
	for i := 0; i < n; i++ {
		if i == tRow {
			// Parameter position: S[i][i] = t
			terms = append(terms, "t"+xi(i+1)+"^2")
			continue
		}
		v := Sbase.A[i][i]
		if v == 0 {
			continue
		}
		if v == 1 {
			terms = append(terms, xi(i+1)+"^2")
		} else if v == -1 {
			terms = append(terms, "-"+xi(i+1)+"^2")
		} else {
			terms = append(terms, fmt.Sprintf("%d%s^2", v, xi(i+1)))
		}
	}

	// Cross terms: 2·S[i][j]·x_i·x_j
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			coef := 2 * Sbase.A[i][j]
			if coef == 0 {
				continue
			}
			if coef == 1 {
				terms = append(terms, xi(i+1)+xi(j+1))
			} else if coef == -1 {
				terms = append(terms, "-"+xi(i+1)+xi(j+1))
			} else {
				terms = append(terms, fmt.Sprintf("%d%s%s", coef, xi(i+1), xi(j+1)))
			}
		}
	}

	if len(terms) == 0 {
		return "0"
	}
	result := terms[0]
	for i := 1; i < len(terms); i++ {
		if strings.HasPrefix(terms[i], "-") {
			result += terms[i]
		} else {
			result += "+" + terms[i]
		}
	}
	return result
}

// genMatrixParamOrthogonalDiag generates a 3×3 symmetric matrix S = QΛQᵀ
// where S[3][3] = t depends on parameter t, and all other entries are fixed integers.
// The eigenvalues λ₁,λ₂,λ₃ (standard-form coefficients) are known.
// Parameter t can be deduced from tr(S) = λ₁+λ₂+λ₃:
//   t = (λ₁+λ₂+λ₃) - S[1][1] - S[2][2]
//
// Stores _param_Q, _param_Lambda, _param_lambdas, _param_t_val in inst.Vars.
func genMatrixParamOrthogonalDiag(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	if v.Rows != 3 || v.Cols != 3 {
		return nil, fmt.Errorf("param_orthogonal_diag: need 3×3")
	}
	lmin := int64(defaultInt(v.Generator, "lambda_min", -5))
	lmax := int64(defaultInt(v.Generator, "lambda_max", 5))
	if lmax < lmin {
		lmin, lmax = lmax, lmin
	}
	maxEntry := int64(defaultInt(v.Generator, "max_entry", 15))

	for attempt := 0; attempt < 500; attempt++ {
		attemptSeed := rng.Int63()
		attemptRng := rand.New(rand.NewSource(attemptSeed))

		// Pick 3 distinct nonzero eigenvalues
		lambdas := make([]int64, 3)
		used := map[int64]bool{}
		ok := true
		for i := 0; i < 3; i++ {
			found := false
			for t := 0; t < 80; t++ {
				l := int64(attemptRng.Intn(int(lmax-lmin+1)) + int(lmin))
				if l == 0 || used[l] {
					continue
				}
				used[l] = true
				lambdas[i] = l
				found = true
				break
			}
			if !found {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}

		// Generate Q (signed permutation orthogonal matrix)
		perm := []int{0, 1, 2}
		attemptRng.Shuffle(3, func(i, j int) { perm[i], perm[j] = perm[j], perm[i] })
		Q := NewMatrixInt(3, 3)
		for j := 0; j < 3; j++ {
			sign := int64(1)
			if attemptRng.Intn(2) == 0 {
				sign = -1
			}
			Q.A[perm[j]][j] = sign
		}

		// Λ = diag(λ₁, λ₂, λ₃)
		Lambda := NewMatrixInt(3, 3)
		for i := 0; i < 3; i++ {
			Lambda.A[i][i] = lambdas[i]
		}

		// S = Q Λ Qᵀ
		QLambda, err := matrixMulInt(Q, Lambda)
		if err != nil {
			continue
		}
		Qt := NewMatrixInt(3, 3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				Qt.A[i][j] = Q.A[j][i]
			}
		}
		Sfull, err := matrixMulInt(QLambda, Qt)
		if err != nil {
			continue
		}

		// Verify S is symmetric
		for i := 0; i < 3; i++ {
			for j := i + 1; j < 3; j++ {
				if Sfull.A[i][j] != Sfull.A[j][i] {
					continue
				}
			}
		}

		// Check entries are reasonable
		me := int64(0)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if Sfull.A[i][j] > me {
					me = Sfull.A[i][j]
				}
				if -Sfull.A[i][j] > me {
					me = -Sfull.A[i][j]
				}
			}
		}
		if me > maxEntry {
			continue
		}

		// S[3][3] = t (parameter). t value = Sfull[2][2]
		tVal := Sfull.A[2][2]

		// Sbase: S with S[3][3] replaced by 0 (placeholder for t)
		Sbase := copyMatrix(Sfull)
		Sbase.A[2][2] = 0

		// Store for answer extraction
		if inst.Vars == nil {
			inst.Vars = map[string]interface{}{}
		}
		inst.Vars["_param_Q"] = Q
		inst.Vars["_param_Lambda"] = Lambda
		inst.Vars["_param_lambda1"] = lambdas[0]
		inst.Vars["_param_lambda2"] = lambdas[1]
		inst.Vars["_param_lambda3"] = lambdas[2]
		inst.Vars["_param_t_val"] = tVal
		inst.Vars["_param_Sbase"] = Sbase
		inst.Vars["_param_t_pos"] = []interface{}{float64(2), float64(2)}

		return Sbase, nil
	}
	return nil, fmt.Errorf("param_orthogonal_diag: failed")
}