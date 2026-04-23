package dsl

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
)

// IntegralInnerProduct computes <f,g> = integral from -1 to 1 of f(x)*g(x) dx
// for polynomials represented as coefficient vectors [a0, a1, a2, ...] (constant first).
//
// <x^n, x^m> = integral(-1,1) x^(n+m) dx
// = 2/(n+m+1) if n+m is even, 0 if n+m is odd
func integralInnerProduct(f, g []int64) *big.Rat {
	result := new(big.Rat)
	for i, ai := range f {
		for j, bj := range g {
			if ai == 0 || bj == 0 {
				continue
			}
			n := i + j // degree of x^i * x^j = x^(i+j)
			if n%2 == 1 {
				continue // odd powers integrate to 0 over [-1,1]
			}
			// integral(-1,1) x^n dx = 2/(n+1)
			coeff := new(big.Rat).SetFrac(big.NewInt(ai*bj*2), big.NewInt(int64(n+1)))
			result.Add(result, coeff)
		}
	}
	return result
}

// polynomialSchmidt performs Gram-Schmidt orthogonalization under the
// integral inner product <f,g> = integral(-1,1) f*g dx.
//
// Input: list of polynomials (coefficient vectors, constant first)
// Output: orthogonal polynomials (rational coefficient vectors)
// Returns coefficients as []*big.Rat vectors.
func polynomialSchmidt(polys [][]int64) [][]*big.Rat {
	n := len(polys)
	result := make([][]*big.Rat, n)

	for k := 0; k < n; k++ {
		// Start with f_k
		fk := make([]*big.Rat, len(polys[k]))
		for i, c := range polys[k] {
			fk[i] = new(big.Rat).SetInt64(c)
		}

		// Subtract projections onto previous orthogonal polynomials
		for j := 0; j < k; j++ {
			gj := result[j]
			// <f_k, g_j>
			fkInt := ratPolyToIntPoly(fk)
			gjInt := ratPolyToIntPoly(gj)
			ip := integralInnerProduct(fkInt, gjInt)
			// <g_j, g_j>
			norm := integralInnerProduct(gjInt, gjInt)
			// projection = (ip/norm) * g_j
			if norm.Sign() == 0 {
				continue
			}
			coef := new(big.Rat).Quo(ip, norm)
			// fk = fk - coef * gj
			for i := 0; i < len(gj); i++ {
				if i >= len(fk) {
					// extend fk
					for len(fk) <= i {
						fk = append(fk, new(big.Rat))
					}
				}
				sub := new(big.Rat).Mul(coef, gj[i])
				fk[i].Sub(fk[i], sub)
			}
		}

		result[k] = fk
	}
	return result
}

// ratPolyToIntPoly converts a []*big.Rat polynomial to []int64
// by finding the LCM of denominators and multiplying through.
func ratPolyToIntPoly(ratCoeffs []*big.Rat) []int64 {
	// Find LCM of denominators
	lcm := big.NewInt(1)
	for _, r := range ratCoeffs {
		if r == nil {
			continue
		}
		d := r.Denom()
		if d.Sign() != 0 {
			newLcm := new(big.Int)
			gcd := new(big.Int).GCD(nil, nil, lcm, d)
			newLcm.Mul(lcm, d)
			newLcm.Quo(newLcm, gcd)
			lcm = newLcm
		}
	}
	result := make([]int64, len(ratCoeffs))
	for i, r := range ratCoeffs {
		if r == nil {
			result[i] = 0
			continue
		}
		scaled := new(big.Rat).Mul(r, new(big.Rat).SetInt(lcm))
		result[i] = scaled.Num().Int64()
	}
	return result
}

// genPolySchmidtIntegral generates 3 polynomials for Gram-Schmidt
// under the integral inner product <f,g> = integral(-1,1) f*g dx.
//
// Strategy: pick 3 polynomials of degree 0, 1, 2 with small integer coefficients,
// verify they are linearly independent under the integral inner product.
// Store: _poly_schmidt_input_polys, _poly_schmidt_output_polys (rational coeffs)
func genPolySchmidtIntegral(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	cmin := int64(defaultInt(v.Generator, "coef_min", -5))
	cmax := int64(defaultInt(v.Generator, "coef_max", 5))
	if cmax < cmin {
		cmin, cmax = cmax, cmin
	}

	for attempt := 0; attempt < 500; attempt++ {
		attemptSeed := rng.Int63()
		attemptRng := rand.New(rand.NewSource(attemptSeed))

		// f1(x) = a0 (constant polynomial, nonzero)
		a0 := int64(0)
		for t := 0; t < 40; t++ {
			x := int64(attemptRng.Intn(int(cmax-cmin+1)) + int(cmin))
			if x != 0 {
				a0 = x
				break
			}
		}
		if a0 == 0 {
			continue
		}
		f1 := []int64{a0}

		// f2(x) = b0 + b1*x (linear polynomial, b1 nonzero)
		b0 := int64(attemptRng.Intn(int(cmax-cmin+1)) + int(cmin))
		b1 := int64(0)
		for t := 0; t < 40; t++ {
			x := int64(attemptRng.Intn(int(cmax-cmin+1)) + int(cmin))
			if x != 0 {
				b1 = x
				break
			}
		}
		if b1 == 0 {
			continue
		}
		f2 := []int64{b0, b1}

		// f3(x) = c0 + c1*x + c2*x^2 (quadratic, c2 nonzero)
		c0_ := int64(attemptRng.Intn(int(cmax-cmin+1)) + int(cmin))
		c1_ := int64(attemptRng.Intn(int(cmax-cmin+1)) + int(cmin))
		c2_ := int64(0)
		for t := 0; t < 40; t++ {
			x := int64(attemptRng.Intn(int(cmax-cmin+1)) + int(cmin))
			if x != 0 {
				c2_ = x
				break
			}
		}
		if c2_ == 0 {
			continue
		}
		f3 := []int64{c0_, c1_, c2_}

		// Verify linear independence under integral inner product
		// The determinant of the Gram matrix must be nonzero
		polys := [][]int64{f1, f2, f3}
		gm := gramMatrixIntegral(polys)
		if gm == nil || gm.Sign() == 0 {
			continue
		}

		// Perform Schmidt orthogonalization
		result := polynomialSchmidt(polys)

		// Convert rational coefficients to integer (find LCM and scale)
		intResults := make([][]int64, 3)
		for i, rp := range result {
			// Find LCM of denominators
			lcm := big.NewInt(1)
			for _, r := range rp {
				if r == nil || r.Sign() == 0 {
					continue
				}
				d := r.Denom()
				if d.Sign() != 0 {
					newLcm := new(big.Int)
					gcd := new(big.Int).GCD(nil, nil, lcm, d)
					newLcm.Mul(lcm, d)
					newLcm.Quo(newLcm, gcd)
					lcm = newLcm
				}
			}
			intPoly := make([]int64, 3) // always 3 components (constant, x, x^2)
			for j := 0; j < 3; j++ {
				if j < len(rp) && rp[j] != nil {
					scaled := new(big.Rat).Mul(rp[j], new(big.Rat).SetInt(lcm))
					intPoly[j] = scaled.Num().Int64()
				} else {
					intPoly[j] = 0
				}
			}
			intResults[i] = intPoly
		}

		// Store
		if inst.Vars == nil {
			inst.Vars = map[string]interface{}{}
		}
		inst.Vars["_poly_schmidt_input"] = polys
		inst.Vars["_poly_schmidt_output"] = intResults

		// Return a placeholder matrix (not shown)
		return NewMatrixInt(3, 3), nil
	}
	return nil, fmt.Errorf("poly_schmidt_integral: failed to generate")
}

// gramMatrixIntegral computes the determinant of the Gram matrix
// for polynomials under the integral inner product.
func gramMatrixIntegral(polys [][]int64) *big.Rat {
	n := len(polys)
	entries := make([][]*big.Rat, n)
	for i := 0; i < n; i++ {
		entries[i] = make([]*big.Rat, n)
		for j := 0; j < n; j++ {
			entries[i][j] = integralInnerProduct(polys[i], polys[j])
		}
	}
	// Compute determinant of n×n rational matrix
	return ratDet(entries, n)
}

// ratDet computes determinant of n×n rational matrix using cofactor expansion
func ratDet(m [][]*big.Rat, n int) *big.Rat {
	if n == 1 {
		return new(big.Rat).Set(m[0][0])
	}
	if n == 2 {
		ad := new(big.Rat).Mul(m[0][0], m[1][1])
		bc := new(big.Rat).Mul(m[0][1], m[1][0])
		return new(big.Rat).Sub(ad, bc)
	}
	det := new(big.Rat)
	for j := 0; j < n; j++ {
		minor := make([][]*big.Rat, n-1)
		for i := 1; i < n; i++ {
			minor[i-1] = make([]*big.Rat, n-1)
			k := 0
			for jj := 0; jj < n; jj++ {
				if jj == j {
					continue
				}
				minor[i-1][k] = m[i][jj]
				k++
			}
		}
		cofactor := ratDet(minor, n-1)
		sign := int64(1)
		if j%2 == 1 {
			sign = -1
		}
		contribution := new(big.Rat).Mul(m[0][j], cofactor)
		contribution.Mul(contribution, new(big.Rat).SetInt64(sign))
		det.Add(det, contribution)
	}
	return det
}

// formatPolyLatex renders a polynomial as LaTeX string.
// coeffs = [a0, a1, a2, ...] where a_k is the coefficient of x^k
func formatPolyLatex(coeffs []int64) string {
	var terms []string
	for k, c := range coeffs {
		if c == 0 {
			continue
		}
		var term string
		switch k {
		case 0:
			term = fmt.Sprintf("%d", c)
		case 1:
			if c == 1 {
				term = "x"
			} else if c == -1 {
				term = "-x"
			} else {
				term = fmt.Sprintf("%dx", c)
			}
		default:
			if c == 1 {
				term = fmt.Sprintf("x^{%d}", k)
			} else if c == -1 {
				term = fmt.Sprintf("-x^{%d}", k)
			} else {
				term = fmt.Sprintf("%dx^{%d}", c, k)
			}
		}
		terms = append(terms, term)
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

// evalPolySchmidtComp returns the k-th coefficient of the j-th Schmidt polynomial
// k=1 is constant, k=2 is x-coeff, k=3 is x^2-coeff; j=1,2,3 for g1,g2,g3
func evalPolySchmidtComp(inst *Instance, j, k int) (int64, error) {
	ov, ok := inst.Vars["_poly_schmidt_output"]
	if !ok {
		return 0, fmt.Errorf("poly_schmidt_comp: data not found")
	}
	output, ok := ov.([][]int64)
	if !ok {
		return 0, fmt.Errorf("poly_schmidt_comp: unexpected type %T", ov)
	}
	if j < 1 || j > len(output) {
		return 0, fmt.Errorf("poly_schmidt_comp: j out of range")
	}
	if k < 1 || k > len(output[j-1]) {
		return 0, fmt.Errorf("poly_schmidt_comp: k out of range")
	}
	return output[j-1][k-1], nil
}

// evalPolySchmidtInputComp returns the k-th coefficient of the j-th input polynomial
func evalPolySchmidtInputComp(inst *Instance, j, k int) (int64, error) {
	iv, ok := inst.Vars["_poly_schmidt_input"]
	if !ok {
		return 0, fmt.Errorf("poly_schmidt_input_comp: data not found")
	}
	input, ok := iv.([][]int64)
	if !ok {
		return 0, fmt.Errorf("poly_schmidt_input_comp: unexpected type %T", iv)
	}
	if j < 1 || j > len(input) {
		return 0, fmt.Errorf("poly_schmidt_input_comp: j out of range")
	}
	poly := input[j-1]
	if k < 1 || k > len(poly) {
		if k > len(poly) {
			return 0, nil // higher degree terms are 0
		}
		return 0, fmt.Errorf("poly_schmidt_input_comp: k out of range")
	}
	return poly[k-1], nil
}

// evalPolySchmidtInputText renders the input polynomials as LaTeX
func evalPolySchmidtInputText(inst *Instance) (string, error) {
	iv, ok := inst.Vars["_poly_schmidt_input"]
	if !ok {
		return "", fmt.Errorf("poly_schmidt_input_text: data not found")
	}
	input, ok := iv.([][]int64)
	if !ok {
		return "", fmt.Errorf("poly_schmidt_input_text: unexpected type %T", iv)
	}
	parts := make([]string, len(input))
	for i, poly := range input {
		parts[i] = fmt.Sprintf("f_{%d}(x)=%s", i+1, formatPolyLatex(poly))
	}
	return strings.Join(parts, "，"), nil
}