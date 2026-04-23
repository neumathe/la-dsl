package dsl

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

// JudgeOptions 判分选项；全 nil 时每个空权重为 1。
type JudgeOptions struct {
	WeightByID map[string]float64
}

// FieldJudgement 单个填空的判分结果。
type FieldJudgement struct {
	ID         string  `json:"id"`
	Correct    bool    `json:"correct"`
	Expected   string  `json:"expected"`
	Submitted  string  `json:"submitted"`
	Weight     float64 `json:"weight"`
	Score      float64 `json:"score"`
	DetailNote string  `json:"detail_note,omitempty"`
}

// JudgeResult 整题判分结果。
type JudgeResult struct {
	Fields       []FieldJudgement `json:"fields"`
	CorrectCount int              `json:"correct_count"`
	TotalFields  int              `json:"total_fields"`
	ScoreEarned  float64          `json:"score_earned"`
	ScoreMax     float64          `json:"score_max"`
	AllCorrect   bool             `json:"all_correct"`
}

// ValueToCanonicalString 将标准答案格式化为可展示/日志的规范字符串。
func ValueToCanonicalString(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case int:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case *big.Int:
		return t.String()
	case *big.Rat:
		if t.IsInt() {
			return t.Num().String()
		}
		return t.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toRat(v interface{}) (*big.Rat, error) {
	switch t := v.(type) {
	case *big.Rat:
		return new(big.Rat).Set(t), nil
	case *big.Int:
		return new(big.Rat).SetInt(t), nil
	case int64:
		return big.NewRat(t, 1), nil
	case int:
		return big.NewRat(int64(t), 1), nil
	case float64:
		return new(big.Rat).SetFloat64(t), nil
	default:
		return nil, fmt.Errorf("unsupported answer type %T", v)
	}
}

// NormalizeUserAnswer 去掉空白、全角符号等常见输入噪声，并将 Unicode 减号、分数斜杠等归一为 ASCII。
func NormalizeUserAnswer(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\u2212", "-") // minus sign
	s = strings.ReplaceAll(s, "－", "-")     // fullwidth hyphen-minus
	s = strings.ReplaceAll(s, "⁻", "-")     // superscript minus
	s = strings.ReplaceAll(s, "，", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\u00a0", "")
	s = strings.ReplaceAll(s, "（", "(")
	s = strings.ReplaceAll(s, "）", ")")
	s = strings.ReplaceAll(s, "\u2044", "/") // fraction slash ⁄
	s = strings.ReplaceAll(s, "／", "/")       // fullwidth solidus
	return strings.TrimSpace(s)
}

// ParseUserRational 解析用户输入的整数或分数（如 -3/2）。
func ParseUserRational(s string) (*big.Rat, error) {
	s = NormalizeUserAnswer(s)
	for {
		if s == "" {
			return nil, fmt.Errorf("empty")
		}
		if len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')' {
			s = strings.TrimSpace(s[1 : len(s)-1])
			continue
		}
		break
	}
	if s == "" {
		return nil, fmt.Errorf("empty")
	}
	r := new(big.Rat)
	if _, ok := r.SetString(s); ok {
		return r, nil
	}
	// 小数
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return new(big.Rat).SetFloat64(f), nil
	}
	return nil, fmt.Errorf("cannot parse %q", s)
}

// ScalarAnswersEqual 比较标准值与用户输入是否在算术意义下相等。
func ScalarAnswersEqual(expected interface{}, submitted string) (bool, string) {
	sub := NormalizeUserAnswer(submitted)
	if sub == "" {
		return false, "empty"
	}
	ex, err := toRat(expected)
	if err != nil {
		return false, err.Error()
	}
	us, err := ParseUserRational(submitted)
	if err != nil {
		return false, err.Error()
	}
	if ex.Cmp(us) == 0 {
		return true, ""
	}
	return false, "value mismatch"
}

// JudgeGeneratedQuestion 根据 GenerateQuestion 的结果与用户提交的 id->字符串 判分。
func JudgeGeneratedQuestion(g *GeneratedQuestion, user map[string]string, opts *JudgeOptions) JudgeResult {
	return JudgeGeneratedQuestionContext(g, nil, nil, user, opts)
}

func fieldWeight(id string, opts *JudgeOptions) float64 {
	w := 1.0
	if opts != nil && opts.WeightByID != nil {
		if ww, ok := opts.WeightByID[id]; ok && ww > 0 {
			w = ww
		}
	}
	return w
}

// JudgeGeneratedQuestionContext 判分；若题目含结构化 Judge（共线、仿射、列基等价），需传入同一实例的 inst。
func JudgeGeneratedQuestionContext(g *GeneratedQuestion, p *Problem, inst *Instance, user map[string]string, opts *JudgeOptions) JudgeResult {
	if g == nil {
		return JudgeResult{}
	}
	n := len(g.AnswerFields)
	out := JudgeResult{TotalFields: n}
	handled := make([]bool, n)

	collectLineGroups := func() map[string][]int {
		m := map[string][]int{}
		for i := 0; i < n; i++ {
			j := g.AnswerFields[i].Judge
			if j == nil || j.Kind != "rational_line" || j.LineGroup == "" {
				continue
			}
			m[j.LineGroup] = append(m[j.LineGroup], i)
		}
		return m
	}
	collectAffineGroups := func() map[string][]int {
		m := map[string][]int{}
		for i := 0; i < n; i++ {
			j := g.AnswerFields[i].Judge
			if j == nil || j.Kind != "affine_rational" || j.AffineGroup == "" {
				continue
			}
			m[j.AffineGroup] = append(m[j.AffineGroup], i)
		}
		return m
	}
	collectBasisGroups := func() map[string][]int {
		m := map[string][]int{}
		for i := 0; i < n; i++ {
			j := g.AnswerFields[i].Judge
			if j == nil || j.Kind != "sorted_basis_columns" || j.BasisGroup == "" {
				continue
			}
			m[j.BasisGroup] = append(m[j.BasisGroup], i)
		}
		return m
	}
	collectPermGroups := func() map[string][]int {
		m := map[string][]int{}
		for i := 0; i < n; i++ {
			j := g.AnswerFields[i].Judge
			if j == nil || j.Kind != "permutation_multiset" || j.PermGroup == "" {
				continue
			}
			m[j.PermGroup] = append(m[j.PermGroup], i)
		}
		return m
	}
	// eigen_pair：按 EigenGroup 汇总所有 λ / vec 字段。
	collectEigenGroups := func() map[string][]int {
		m := map[string][]int{}
		for i := 0; i < n; i++ {
			j := g.AnswerFields[i].Judge
			if j == nil || j.Kind != "eigen_pair" || j.EigenGroup == "" {
				continue
			}
			m[j.EigenGroup] = append(m[j.EigenGroup], i)
		}
		return m
	}

	appendGroup := func(idxs []int, ok bool, note string) {
		for _, i := range idxs {
			handled[i] = true
			f := g.AnswerFields[i]
			w := fieldWeight(f.ID, opts)
			sub := ""
			if user != nil {
				sub = user[f.ID]
			}
			sc := 0.0
			if ok {
				sc = w
			}
			out.Fields = append(out.Fields, FieldJudgement{
				ID:         f.ID,
				Correct:    ok,
				Expected:   ValueToCanonicalString(f.Value),
				Submitted:  NormalizeUserAnswer(sub),
				Weight:     w,
				Score:      sc,
				DetailNote: note,
			})
		}
	}

	if inst != nil {
		for _, idxs := range collectLineGroups() {
			ok, note := judgeRationalLineGroup(g, user, idxs)
			appendGroup(idxs, ok, note)
		}
		for gname, idxs := range collectAffineGroups() {
			j0 := g.AnswerFields[idxs[0]].Judge
			ok, note := judgeAffineRationalGroup(g, inst, user, idxs, j0.MatrixVar, j0.BVecVar)
			_ = gname
			appendGroup(idxs, ok, note)
		}
		for _, idxs := range collectBasisGroups() {
			j0 := g.AnswerFields[idxs[0]].Judge
			ok, note := judgeSortedBasisGroup(g, inst, user, idxs, j0.MatrixVar, j0.Ncols)
			appendGroup(idxs, ok, note)
		}
		for _, idxs := range collectPermGroups() {
			ok, note := judgePermutationMultisetGroup(g, user, idxs)
			appendGroup(idxs, ok, note)
		}
		// eigen_pair：需要先收集所有 EigenGroup，再处理，因为 vec-only 组可能引用其他组的 λ。
		eigenAll := collectEigenGroups()
		eigenLambdasByGroup := map[string][]int{} // group -> 按 EigenColumn 排序的 lambda field 下标
		for gname, idxs := range eigenAll {
			lambdas := eigenCollectLambdasByColumn(g, idxs)
			if len(lambdas) > 0 {
				eigenLambdasByGroup[gname] = lambdas
			}
		}
		for gname, idxs := range eigenAll {
			ok, note, done := judgeEigenPairGroup(g, inst, user, gname, idxs, eigenLambdasByGroup)
			if !done {
				// 缺参数时，退化为标量判分
				continue
			}
			appendGroup(idxs, ok, note)
		}
	} else {
		// 无实例时，结构化判题无法执行，退化为逐空标量（可能判错含 Judge 的题；请使用 JudgeBankQuestion）。
		for _, idxs := range collectLineGroups() {
			for _, i := range idxs {
				handled[i] = true
				f := g.AnswerFields[i]
				w := fieldWeight(f.ID, opts)
				sub := ""
				if user != nil {
					sub = user[f.ID]
				}
				out.Fields = append(out.Fields, FieldJudgement{
					ID: f.ID, Correct: false, Expected: ValueToCanonicalString(f.Value),
					Submitted: NormalizeUserAnswer(sub), Weight: w, Score: 0,
					DetailNote: "need instance for rational_line judge",
				})
			}
		}
		for _, idxs := range collectAffineGroups() {
			for _, i := range idxs {
				handled[i] = true
				f := g.AnswerFields[i]
				w := fieldWeight(f.ID, opts)
				sub := ""
				if user != nil {
					sub = user[f.ID]
				}
				out.Fields = append(out.Fields, FieldJudgement{
					ID: f.ID, Correct: false, Expected: ValueToCanonicalString(f.Value),
					Submitted: NormalizeUserAnswer(sub), Weight: w, Score: 0,
					DetailNote: "need instance for affine_rational judge",
				})
			}
		}
		for _, idxs := range collectBasisGroups() {
			for _, i := range idxs {
				handled[i] = true
				f := g.AnswerFields[i]
				w := fieldWeight(f.ID, opts)
				sub := ""
				if user != nil {
					sub = user[f.ID]
				}
				out.Fields = append(out.Fields, FieldJudgement{
					ID: f.ID, Correct: false, Expected: ValueToCanonicalString(f.Value),
					Submitted: NormalizeUserAnswer(sub), Weight: w, Score: 0,
					DetailNote: "need instance for sorted_basis_columns judge",
				})
			}
		}
		// permutation_multiset 不依赖 inst，这里直接判分。
		for _, idxs := range collectPermGroups() {
			ok, note := judgePermutationMultisetGroup(g, user, idxs)
			appendGroup(idxs, ok, note)
		}
		for _, idxs := range collectEigenGroups() {
			for _, i := range idxs {
				handled[i] = true
				f := g.AnswerFields[i]
				w := fieldWeight(f.ID, opts)
				sub := ""
				if user != nil {
					sub = user[f.ID]
				}
				out.Fields = append(out.Fields, FieldJudgement{
					ID: f.ID, Correct: false, Expected: ValueToCanonicalString(f.Value),
					Submitted: NormalizeUserAnswer(sub), Weight: w, Score: 0,
					DetailNote: "need instance for eigen_pair judge",
				})
			}
		}
	}

	for i := 0; i < n; i++ {
		if handled[i] {
			continue
		}
		f := g.AnswerFields[i]
		w := fieldWeight(f.ID, opts)
		sub := ""
		if user != nil {
			sub = user[f.ID]
		}
		ok, note := ScalarAnswersEqual(f.Value, sub)
		sc := 0.0
		if ok {
			sc = w
		}
		out.Fields = append(out.Fields, FieldJudgement{
			ID:         f.ID,
			Correct:    ok,
			Expected:   ValueToCanonicalString(f.Value),
			Submitted:  NormalizeUserAnswer(sub),
			Weight:     w,
			Score:      sc,
			DetailNote: note,
		})
	}

	// 恢复与 AnswerFields 相同的顺序
	sort.Slice(out.Fields, func(a, b int) bool {
		ia, ib := out.Fields[a].ID, out.Fields[b].ID
		return fieldOrderIndex(g, ia) < fieldOrderIndex(g, ib)
	})

	correct := 0
	var maxSum float64
	for _, fj := range out.Fields {
		maxSum += fj.Weight
		if fj.Correct {
			correct++
		}
		out.ScoreEarned += fj.Score
	}
	out.CorrectCount = correct
	out.ScoreMax = maxSum
	out.AllCorrect = out.TotalFields > 0 && out.CorrectCount == out.TotalFields
	_ = p
	return out
}

func fieldOrderIndex(g *GeneratedQuestion, id string) int {
	for i, f := range g.AnswerFields {
		if f.ID == id {
			return i
		}
	}
	return 1 << 30
}

func judgeRationalLineGroup(g *GeneratedQuestion, user map[string]string, idxs []int) (bool, string) {
	expVals := make([]interface{}, len(idxs))
	subs := make([]string, len(idxs))
	for k, i := range idxs {
		expVals[k] = g.AnswerFields[i].Value
		if user != nil {
			subs[k] = user[g.AnswerFields[i].ID]
		}
	}
	er, err := interfaceSliceToRatVec(expVals)
	if err != nil {
		return false, err.Error()
	}
	ur, err := parseUserRatVector(subs)
	if err != nil {
		return false, err.Error()
	}
	return vectorsRationalCollinear(er, ur)
}

func judgeAffineRationalGroup(g *GeneratedQuestion, inst *Instance, user map[string]string, idxs []int, matrixVar, bVar string) (bool, string) {
	if matrixVar == "" || bVar == "" {
		return false, "missing matrix/b"
	}
	vA, ok := inst.Vars[matrixVar]
	if !ok {
		return false, "no matrix var"
	}
	vb, ok := inst.Vars[bVar]
	if !ok {
		return false, "no b var"
	}
	A, ok := vA.(*MatrixInt)
	if !ok {
		return false, "matrix type"
	}
	b, ok := vb.(*VectorInt)
	if !ok {
		return false, "b type"
	}
	subs := make([]string, len(idxs))
	for k, i := range idxs {
		if user != nil {
			subs[k] = user[g.AnswerFields[i].ID]
		}
	}
	x, err := parseUserRatVector(subs)
	if err != nil {
		return false, err.Error()
	}
	return affineSolutionOK(A, b, x)
}

func judgeSortedBasisGroup(g *GeneratedQuestion, inst *Instance, user map[string]string, idxs []int, matrixVar string, ncols int) (bool, string) {
	if inst == nil || matrixVar == "" || ncols <= 0 {
		return false, "bad basis params"
	}
	vA, ok := inst.Vars[matrixVar]
	if !ok {
		return false, "no matrix"
	}
	A, ok := vA.(*MatrixInt)
	if !ok {
		return false, "matrix type"
	}
	subs := make([]string, len(idxs))
	for k, i := range idxs {
		if user != nil {
			subs[k] = user[g.AnswerFields[i].ID]
		}
	}
	return sortedBasisColumnsOK(A, ncols, subs)
}
