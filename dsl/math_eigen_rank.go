package dsl

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

// genScalarEigenRankInference5 generates a *scalar* placeholder variable
// whose eigenvalues (5 total) are determined by rank conditions.
// The student is not shown any matrix — they receive rank conditions
// like "R(A+8E)=3, R(A-3E)=3, R(A-7E)=4" and deduce all 5 eigenvalues.
//
// Strategy: pick 3 distinct nonzero eigenvalues λ₁,λ₂,λ₃ with multiplicities
// m₁,m₂,m₃ (summing to 5). The rank conditions are:
//   R(A-λᵢI) = 5-mᵢ  (since nullity(A-λᵢI) = mᵢ)
//
// Stores _eigen_rank_lambdas, _eigen_rank_mults, _eigen_rank_ranks in inst.Vars.
func genScalarEigenRankInference5(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	lmin := int64(defaultInt(v.Generator, "lambda_min", -8))
	lmax := int64(defaultInt(v.Generator, "lambda_max", 8))
	if lmax < lmin {
		lmin, lmax = lmax, lmin
	}

	for attempt := 0; attempt < 200; attempt++ {
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

		// Pick multiplicities that sum to 5, each ≥ 1
		// Distribute 2 extra among the 3 eigenvalues (each starts at 1)
		mults := []int64{1, 1, 1}
		remaining := int64(2)
		for remaining > 0 {
			idx := attemptRng.Intn(3)
			mults[idx]++
			remaining--
		}

		// Sort eigenvalues and their multiplicities together (ascending)
		pairs := make([][2]int64, 3)
		for i := 0; i < 3; i++ {
			pairs[i] = [2]int64{lambdas[i], mults[i]}
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i][0] < pairs[j][0] })
		for i := 0; i < 3; i++ {
			lambdas[i] = pairs[i][0]
			mults[i] = pairs[i][1]
		}

		// Verify total multiplicities = 5
		total := int64(0)
		for _, m := range mults {
			total += m
		}
		if total != 5 {
			continue
		}

		// Verify all eigenvalues nonzero and distinct
		allNonzero := true
		for _, l := range lambdas {
			if l == 0 {
				allNonzero = false
			}
		}
		if !allNonzero || lambdas[0] == lambdas[1] || lambdas[1] == lambdas[2] {
			continue
		}

		// Compute rank conditions: R(A - λᵢI) = 5 - mults[i]
		ranks := make([]int64, 3)
		for i := 0; i < 3; i++ {
			ranks[i] = 5 - mults[i]
		}

		// Verify ranks are between 1 and 4 (not 0 or 5)
		allValid := true
		for _, r := range ranks {
			if r < 1 || r > 4 {
				allValid = false
			}
		}
		if !allValid {
			continue
		}

		// Store data for answer extraction
		if inst.Vars == nil {
			inst.Vars = map[string]interface{}{}
		}
		inst.Vars["_eigen_rank_lambdas"] = lambdas
		inst.Vars["_eigen_rank_mults"] = mults
		inst.Vars["_eigen_rank_ranks"] = ranks

		// Return a placeholder scalar (the student never sees a matrix)
		return int64(0), nil
	}
	return nil, fmt.Errorf("eigen_rank_inference_5: failed to generate")
}

// genScalarEigenRowSumRank4 generates conditions for a 4×4 matrix A
// such that: rank(A)=r, all row sums = s, det(A+kE)=0,
// from which the student infers all 4 eigenvalues.
//
// Strategy (all 4 eigenvalues are determined):
//   - Row sum s ⟹ λ=s is an eigenvalue (1,1,1,1 eigenvector)
//   - rank(A)=2 ⟹ λ=0 has multiplicity 2
//   - det(A+kE)=0 ⟹ λ=-k is an eigenvalue
//   - Total: 0, 0, s, -k (all 4 determined)
//
// For rank=3: 0, s, -k, λ₄ (need an extra eigenvalue)
// We only generate rank=2 cases (matching the HTML template structure).
//
// Stores _eigen_rowsum_lambdas (sorted), _eigen_rowsum_s, _eigen_rowsum_r,
// _eigen_rowsum_k in inst.Vars.
func genScalarEigenRowSumRank4(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	smin := int64(defaultInt(v.Generator, "row_sum_min", 2))
	smax := int64(defaultInt(v.Generator, "row_sum_max", 6))
	kmin := int64(defaultInt(v.Generator, "k_min", 2))
	kmax := int64(defaultInt(v.Generator, "k_max", 10))

	for attempt := 0; attempt < 200; attempt++ {
		attemptSeed := rng.Int63()
		attemptRng := rand.New(rand.NewSource(attemptSeed))

		// Pick row sum s (eigenvalue from row-sum property, must be nonzero)
		s := int64(attemptRng.Intn(int(smax-smin+1)) + int(smin))
		if s == 0 {
			continue
		}

		// rank = 2 (matching HTML template), so λ=0 has multiplicity 2
		r := 2

		// Pick k such that λ=-k is an eigenvalue, -k ≠ 0 and -k ≠ s
		k := int64(0)
		for t := 0; t < 80; t++ {
			kk := int64(attemptRng.Intn(int(kmax-kmin+1)) + int(kmin))
			if kk != s && kk != 0 {
				k = kk
				break
			}
		}
		if k == 0 {
			continue
		}

		// Eigenvalues: 0, 0, s, -k (sorted ascending)
		lambdas := []int64{0, 0, s, -k}
		sort.Slice(lambdas, func(i, j int) bool { return lambdas[i] < lambdas[j] })

		// Store
		if inst.Vars == nil {
			inst.Vars = map[string]interface{}{}
		}
		inst.Vars["_eigen_rowsum_lambdas"] = lambdas
		inst.Vars["_eigen_rowsum_s"] = s
		inst.Vars["_eigen_rowsum_r"] = int64(r)
		inst.Vars["_eigen_rowsum_k"] = k

		return int64(0), nil
	}
	return nil, fmt.Errorf("eigen_row_sum_rank4: failed to generate")
}

// formatRankConditionList generates the LaTeX string for rank conditions.
// Input: pairs of (lambda, rank) like [(-8,3), (3,3), (7,4)]
func formatRankConditionList(lambdas []int64, ranks []int64) string {
	var parts []string
	for i := 0; i < len(lambdas); i++ {
		l := lambdas[i]
		r := ranks[i]
		if l < 0 {
			parts = append(parts, fmt.Sprintf("R(A%dE)=%d", l, r))
		} else {
			parts = append(parts, fmt.Sprintf("R(A+%dE)=%d", l, r))
		}
	}
	return strings.Join(parts, "，")
}

// eigenval_rank(S,i): returns the i-th eigenvalue (1-based, sorted ascending)
// from _eigen_rank_lambdas/mults stored by genScalarEigenRankInference5
func evalEigenvalRank(inst *Instance, i int) (int64, error) {
	lv, ok := inst.Vars["_eigen_rank_lambdas"]
	if !ok {
		return 0, fmt.Errorf("eigenval_rank: data not found")
	}
	lambdas, ok := lv.([]int64)
	if !ok {
		return 0, fmt.Errorf("eigenval_rank: unexpected type %T", lv)
	}
	mv, ok2 := inst.Vars["_eigen_rank_mults"]
	if !ok2 {
		return 0, fmt.Errorf("eigenval_rank: multiplicities not found")
	}
	mults, ok2 := mv.([]int64)
	if !ok2 {
		return 0, fmt.Errorf("eigenval_rank: multiplicities type %T", mv)
	}
	// Build full list of eigenvalues (with multiplicities), sorted ascending
	var full []int64
	for j := 0; j < len(lambdas); j++ {
		for t := 0; t < int(mults[j]); t++ {
			full = append(full, lambdas[j])
		}
	}
	sort.Slice(full, func(a, b int) bool { return full[a] < full[b] })
	if i < 1 || i > len(full) {
		return 0, fmt.Errorf("eigenval_rank: index %d out of range (len=%d)", i, len(full))
	}
	return full[i-1], nil
}

// eigenval_rowsum(S,i): returns the i-th eigenvalue (1-based, sorted ascending)
// from _eigen_rowsum_lambdas stored by genScalarEigenRowSumRank4
func evalEigenvalRowsum(inst *Instance, i int) (int64, error) {
	lv, ok := inst.Vars["_eigen_rowsum_lambdas"]
	if !ok {
		return 0, fmt.Errorf("eigenval_rowsum: data not found")
	}
	lambdas, ok := lv.([]int64)
	if !ok {
		return 0, fmt.Errorf("eigenval_rowsum: unexpected type %T", lv)
	}
	if i < 1 || i > len(lambdas) {
		return 0, fmt.Errorf("eigenval_rowsum: index %d out of range", i)
	}
	return lambdas[i-1], nil
}

// evalEigenRankConditionText returns LaTeX string of rank conditions
func evalEigenRankConditionText(inst *Instance) (string, error) {
	lv, ok := inst.Vars["_eigen_rank_lambdas"]
	if !ok {
		return "", fmt.Errorf("eigen_rank_condition_text: data not found")
	}
	lambdas, ok := lv.([]int64)
	if !ok {
		return "", fmt.Errorf("eigen_rank_condition_text: type %T", lv)
	}
	rv, ok2 := inst.Vars["_eigen_rank_ranks"]
	if !ok2 {
		return "", fmt.Errorf("eigen_rank_condition_text: ranks not found")
	}
	ranks, ok2 := rv.([]int64)
	if !ok2 {
		return "", fmt.Errorf("eigen_rank_condition_text: ranks type %T", rv)
	}
	return formatRankConditionList(lambdas, ranks), nil
}

// evalEigenRowsumConditionText generates the condition text for rowsum-rank problems
func evalEigenRowsumConditionText(inst *Instance) (string, error) {
	sv, ok := inst.Vars["_eigen_rowsum_s"]
	if !ok {
		return "", fmt.Errorf("eigen_rowsum_condition: s not found")
	}
	s, ok := sv.(int64)
	if !ok {
		return "", fmt.Errorf("eigen_rowsum_condition: s type %T", sv)
	}
	rv, ok2 := inst.Vars["_eigen_rowsum_r"]
	if !ok2 {
		return "", fmt.Errorf("eigen_rowsum_condition: r not found")
	}
	r, ok2 := rv.(int64)
	if !ok2 {
		return "", fmt.Errorf("eigen_rowsum_condition: r type %T", rv)
	}
	kv, ok3 := inst.Vars["_eigen_rowsum_k"]
	if !ok3 {
		return "", fmt.Errorf("eigen_rowsum_condition: k not found")
	}
	k, ok3 := kv.(int64)
	if !ok3 {
		return "", fmt.Errorf("eigen_rowsum_condition: k type %T", kv)
	}
	return fmt.Sprintf("R(A)=%d，A 的各行元素之和都等于 %d，且 |A+%dE|=0", r, s, k), nil
}