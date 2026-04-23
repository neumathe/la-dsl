package dsl

import (
	"fmt"
	"math/big"
	"sort"
)

// eigenCollectLambdasByColumn 从组内字段下标中挑出 role="lambda" 的字段，按 EigenColumn 升序返回下标数组。
// 要求 EigenColumn 从 1 连续（若不连续返回 nil）。
func eigenCollectLambdasByColumn(g *GeneratedQuestion, idxs []int) []int {
	type entry struct {
		col int
		idx int
	}
	var es []entry
	for _, i := range idxs {
		j := g.AnswerFields[i].Judge
		if j == nil || j.EigenRole != "lambda" {
			continue
		}
		es = append(es, entry{col: j.EigenColumn, idx: i})
	}
	if len(es) == 0 {
		return nil
	}
	sort.Slice(es, func(a, b int) bool { return es[a].col < es[b].col })
	out := make([]int, 0, len(es))
	for k, e := range es {
		if e.col != k+1 {
			return nil
		}
		out = append(out, e.idx)
	}
	return out
}

// eigenCollectVecColumns 将 role="vec" 的字段按 EigenColumn+EigenComponent 整理成 columns[j-1][component-1]=fieldIdx。
// 返回 columns（二维下标数组）与 dim（向量分量数）。若分组不完整，返回 ok=false。
func eigenCollectVecColumns(g *GeneratedQuestion, idxs []int) (cols [][]int, dim int, ok bool) {
	type entry struct {
		col  int
		comp int
		idx  int
	}
	var es []entry
	maxCol := 0
	maxComp := 0
	for _, i := range idxs {
		j := g.AnswerFields[i].Judge
		if j == nil || j.EigenRole != "vec" {
			continue
		}
		es = append(es, entry{col: j.EigenColumn, comp: j.EigenComponent, idx: i})
		if j.EigenColumn > maxCol {
			maxCol = j.EigenColumn
		}
		if j.EigenComponent > maxComp {
			maxComp = j.EigenComponent
		}
	}
	if len(es) == 0 {
		return nil, 0, false
	}
	cols = make([][]int, maxCol)
	for c := range cols {
		cols[c] = make([]int, maxComp)
		for k := range cols[c] {
			cols[c][k] = -1
		}
	}
	for _, e := range es {
		if e.col < 1 || e.comp < 1 {
			return nil, 0, false
		}
		cols[e.col-1][e.comp-1] = e.idx
	}
	for c := 0; c < maxCol; c++ {
		for k := 0; k < maxComp; k++ {
			if cols[c][k] < 0 {
				return nil, 0, false
			}
		}
	}
	return cols, maxComp, true
}

// judgeEigenPairGroup 判分一组 eigen_pair。
// 返回 done=false 表示参数不全，应退化为标量判分。
func judgeEigenPairGroup(
	g *GeneratedQuestion, inst *Instance, user map[string]string,
	gname string, idxs []int,
	eigenLambdasByGroup map[string][]int,
) (ok bool, note string, done bool) {
	vecCols, _, vecOK := eigenCollectVecColumns(g, idxs)
	if !vecOK {
		return false, "eigen_pair: incomplete vec fields", true
	}

	// 定位 λ 字段下标（优先本组；若本组无，则借用 RefLambdaGroup）。
	lambdaIdxs := eigenLambdasByGroup[gname]
	if len(lambdaIdxs) == 0 {
		// 找任一 vec 字段的 RefLambdaGroup
		ref := ""
		for _, i := range idxs {
			j := g.AnswerFields[i].Judge
			if j != nil && j.RefLambdaGroup != "" {
				ref = j.RefLambdaGroup
				break
			}
		}
		if ref == "" {
			return false, "eigen_pair: no lambda fields and no ref", true
		}
		lambdaIdxs = eigenLambdasByGroup[ref]
	}
	if len(lambdaIdxs) == 0 {
		return false, "eigen_pair: missing lambda fields", true
	}
	if len(lambdaIdxs) != len(vecCols) {
		return false, fmt.Sprintf("eigen_pair: λ count %d != column count %d", len(lambdaIdxs), len(vecCols)), true
	}

	// MatrixVar 由本组任一字段提供。
	matrixVar := ""
	for _, i := range idxs {
		j := g.AnswerFields[i].Judge
		if j != nil && j.MatrixVar != "" {
			matrixVar = j.MatrixVar
			break
		}
	}
	if matrixVar == "" {
		return false, "eigen_pair: missing matrix_var", true
	}
	if inst == nil {
		return false, "eigen_pair: no instance", true
	}
	vA, exists := inst.Vars[matrixVar]
	if !exists {
		return false, fmt.Sprintf("eigen_pair: var %q not found", matrixVar), true
	}
	A, isMat := vA.(*MatrixInt)
	if !isMat {
		return false, "eigen_pair: matrix_var not a matrix", true
	}

	// 解析用户 λ 值。
	userLambdas := make([]*big.Rat, len(lambdaIdxs))
	for k, li := range lambdaIdxs {
		f := g.AnswerFields[li]
		sub := ""
		if user != nil {
			sub = user[f.ID]
		}
		if NormalizeUserAnswer(sub) == "" {
			return false, fmt.Sprintf("eigen_pair: empty λ at column %d", k+1), true
		}
		r, err := ParseUserRational(sub)
		if err != nil {
			return false, fmt.Sprintf("eigen_pair: bad λ at column %d: %v", k+1, err), true
		}
		userLambdas[k] = r
	}

	// 解析用户向量并校验 A·v == λ·v。
	for col := 0; col < len(vecCols); col++ {
		subs := make([]string, len(vecCols[col]))
		for comp, fi := range vecCols[col] {
			if user != nil {
				subs[comp] = user[g.AnswerFields[fi].ID]
			}
		}
		v, err := parseUserRatVector(subs)
		if err != nil {
			return false, fmt.Sprintf("eigen_pair: bad vec at column %d: %v", col+1, err), true
		}
		if !ratVectorNonZero(v) {
			return false, fmt.Sprintf("eigen_pair: zero vector at column %d", col+1), true
		}
		// A 可能是方阵（A·v 维数匹配）
		if A.C != len(v) {
			return false, fmt.Sprintf("eigen_pair: dim mismatch at column %d", col+1), true
		}
		Av, err := matrixIntTimesVectorRat(A, v)
		if err != nil {
			return false, fmt.Sprintf("eigen_pair: matvec error: %v", err), true
		}
		// λ·v
		lv := make([]*big.Rat, len(v))
		for i := range v {
			lv[i] = new(big.Rat).Mul(userLambdas[col], v[i])
		}
		for i := range v {
			if Av[i].Cmp(lv[i]) != 0 {
				return false, fmt.Sprintf("eigen_pair: A·v ≠ λ·v at column %d", col+1), true
			}
		}
	}

	// 校验用户 λ multiset 与标准答案 multiset 相等。
	expVals := make([]*big.Rat, len(lambdaIdxs))
	for k, li := range lambdaIdxs {
		r, err := toRat(g.AnswerFields[li].Value)
		if err != nil {
			return false, fmt.Sprintf("eigen_pair: expected λ parse: %v", err), true
		}
		expVals[k] = r
	}
	if !ratMultisetEqual(userLambdas, expVals) {
		return false, "eigen_pair: λ multiset mismatch", true
	}
	return true, "", true
}

func ratVectorNonZero(v []*big.Rat) bool {
	for _, x := range v {
		if x != nil && x.Sign() != 0 {
			return true
		}
	}
	return false
}

func ratMultisetEqual(a, b []*big.Rat) bool {
	if len(a) != len(b) {
		return false
	}
	aa := append([]*big.Rat(nil), a...)
	bb := append([]*big.Rat(nil), b...)
	sort.Slice(aa, func(i, j int) bool { return aa[i].Cmp(aa[j]) < 0 })
	sort.Slice(bb, func(i, j int) bool { return bb[i].Cmp(bb[j]) < 0 })
	for i := range aa {
		if aa[i].Cmp(bb[i]) != 0 {
			return false
		}
	}
	return true
}

// judgePermutationMultisetGroup 允许同组内答案任意顺序，只要 multiset 相等即判对。
func judgePermutationMultisetGroup(g *GeneratedQuestion, user map[string]string, idxs []int) (bool, string) {
	exp := make([]*big.Rat, 0, len(idxs))
	got := make([]*big.Rat, 0, len(idxs))
	for _, i := range idxs {
		f := g.AnswerFields[i]
		r, err := toRat(f.Value)
		if err != nil {
			return false, fmt.Sprintf("perm: expected parse: %v", err)
		}
		exp = append(exp, r)
		sub := ""
		if user != nil {
			sub = user[f.ID]
		}
		if NormalizeUserAnswer(sub) == "" {
			return false, "perm: empty submission"
		}
		u, err := ParseUserRational(sub)
		if err != nil {
			return false, fmt.Sprintf("perm: user parse: %v", err)
		}
		got = append(got, u)
	}
	if !ratMultisetEqual(exp, got) {
		return false, "perm: multiset mismatch"
	}
	return true, ""
}
