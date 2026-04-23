package dsl

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

// QuestionExplanation 与出题、判题同一套 seed/salt 下的结构化解析，供前端展示或学习核对。
type QuestionExplanation struct {
	ProblemID   int64                 `json:"problem_id,omitempty"`
	Version     string                `json:"version,omitempty"`
	Title       string                `json:"title"`
	Summary     string                `json:"summary,omitempty"`
	TitlePlugs  map[string]string     `json:"title_plugs,omitempty"`
	Variables   []ExplainNamedValue   `json:"variables,omitempty"`
	Derived     []ExplainDerivedValue `json:"derived,omitempty"`
	AnswerSteps []AnswerExplainStep   `json:"answer_steps"`
	Solution    string                `json:"solution_zh,omitempty"`
}

// ExplainNamedValue 本题随机实例中的基础变量取值。
type ExplainNamedValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ExplainDerivedValue 派生量：DSL 表达式及其在本实例下的值。
type ExplainDerivedValue struct {
	Name  string `json:"name"`
	Expr  string `json:"expr"`
	Value string `json:"value"`
}

// AnswerExplainStep 单个填空的计算式与标准答案（规范字符串，与判分一致）。
type AnswerExplainStep struct {
	FieldID  string `json:"field_id"`
	Expr     string `json:"expr"`
	Expected string `json:"expected"`
	Note     string `json:"note,omitempty"`
}

// ValueToExplainString 将矩阵/向量等展开为易读文本；标量等与 ValueToCanonicalString 一致。
func ValueToExplainString(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case *MatrixInt:
		rows := make([]string, 0, t.R)
		for i := 0; i < t.R; i++ {
			cells := make([]string, t.C)
			for j := 0; j < t.C; j++ {
				cells[j] = strconv.FormatInt(t.A[i][j], 10)
			}
			rows = append(rows, "["+strings.Join(cells, " ")+"]")
		}
		return strings.Join(rows, "; ")
	case *VectorInt:
		cells := make([]string, t.N)
		for i := 0; i < t.N; i++ {
			cells[i] = strconv.FormatInt(t.V[i], 10)
		}
		return "(" + strings.Join(cells, ", ") + ")"
	case []*big.Rat:
		parts := make([]string, 0, len(t))
		for _, r := range t {
			if r == nil {
				parts = append(parts, "")
				continue
			}
			parts = append(parts, r.RatString())
		}
		return "[" + strings.Join(parts, ", ") + "]"
	default:
		return ValueToCanonicalString(v)
	}
}

// GenerateExplanation 与 GenerateQuestion 使用相同实例化路径，额外输出变量、派生量与各空解析。
func GenerateExplanation(p Problem, seedStr, serverSalt string) (*QuestionExplanation, error) {
	inst, err := InstantiateProblem(p, seedStr, serverSalt)
	if err != nil {
		return nil, err
	}
	title, err := RenderTitle(p, inst)
	if err != nil {
		return nil, err
	}
	fields, err := ExtractAnswerWithMeta(p, inst)
	if err != nil {
		return nil, err
	}
	return buildQuestionExplanation(p, inst, title, fields), nil
}

func buildQuestionExplanation(p Problem, inst *Instance, title string, fields []AnswerField) *QuestionExplanation {
	out := &QuestionExplanation{
		ProblemID:   p.ID,
		Version:     p.Version,
		Title:       title,
		AnswerSteps: make([]AnswerExplainStep, 0, len(fields)),
	}
	if len(fields) > 0 {
		out.Summary = fmt.Sprintf("本题共 %d 处填空；下列为本题实例下的已知量、派生定义及各空标准答案（与判分规范字符串一致）。", len(fields))
	}

	renderStrings := map[string]string{}
	if len(p.Render) > 0 {
		r, err := RenderInst(p, inst)
		if err == nil {
			out.TitlePlugs = make(map[string]string, len(r))
			keys := make([]string, 0, len(r))
			for k := range r {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				s := FormatValueForTitle(r[k])
				out.TitlePlugs[k] = s
				renderStrings[k] = s
			}
		}
	}

	varNames := make([]string, 0, len(p.Variables))
	for name := range p.Variables {
		varNames = append(varNames, name)
	}
	sort.Strings(varNames)
	for _, name := range varNames {
		val, ok := inst.Vars[name]
		if !ok {
			continue
		}
		s := ValueToExplainString(val)
		out.Variables = append(out.Variables, ExplainNamedValue{
			Name:  name,
			Value: s,
		})
		// 把变量值也加入占位符池（若 Render 已有同名 key 则 Render 的 LaTeX 格式优先）。
		if _, exists := renderStrings[name]; !exists {
			renderStrings[name] = s
		}
	}

	derivedNames := make([]string, 0, len(p.Derived))
	for name := range p.Derived {
		derivedNames = append(derivedNames, name)
	}
	sort.Strings(derivedNames)
	for _, name := range derivedNames {
		expr := strings.TrimSpace(p.Derived[name])
		val, ok := inst.Derived[name]
		if !ok {
			continue
		}
		s := ValueToExplainString(val)
		out.Derived = append(out.Derived, ExplainDerivedValue{
			Name:  name,
			Expr:  expr,
			Value: s,
		})
		// 也把派生量的文本值加入占位符池，方便 solution_steps 引用。
		renderStrings[name] = s
	}

	for _, f := range fields {
		note := f.Note
		if note == "" {
			note = fmt.Sprintf("填空 %q 的值为 %s。", f.ID, ValueToCanonicalString(f.Value))
		}
		expected := ValueToCanonicalString(f.Value)
		out.AnswerSteps = append(out.AnswerSteps, AnswerExplainStep{
			FieldID:  f.ID,
			Expr:     f.Expr,
			Expected: expected,
			Note:     note,
		})
		// 把答案的期望值加入占位符池，用 field ID 作为 key。
		renderStrings[f.ID] = expected
	}

	// 从 Problem.Meta 中提取 solution_zh，做 {{key}} 占位符展开和 {{expr:...}} 内联求值。
	if p.Meta != nil {
		if sol, ok := p.Meta["solution_zh"]; ok {
			if s, ok2 := sol.(string); ok2 {
				s = expandPlaceholders(s, renderStrings)
				s = expandExprPlaceholders(s, p, inst)
				out.Solution = s
			}
		}
	}

	return out
}

// expandPlaceholders 将字符串中的 {{key}} 替换为 vars[key] 的值。
func expandPlaceholders(s string, vars map[string]string) string {
	for k, v := range vars {
		s = strings.ReplaceAll(s, "{{"+k+"}}", v)
	}
	return s
}

// expandExprPlaceholders 将字符串中的 {{expr:DSL_EXPRESSION}} 模板求值并替换。
// 例如 {{expr:mget(A,1,1)}} 在运行时对实例化的 A 求其 (1,1) 元素，
// {{expr:det(A)}} 求行列式值。这使得 solution_zh 可以引用任意中间计算量，
// 无需为每个中间量单独加 Derived 条目。
func expandExprPlaceholders(s string, p Problem, inst *Instance) string {
	for {
		start := strings.Index(s, "{{expr:")
		if start == -1 {
			break
		}
		end := strings.Index(s[start:], "}}")
		if end == -1 {
			break
		}
		end += start // absolute position
		expr := s[start+7 : end] // extract expression between {{expr: and }}
		val, err := EvaluateExpression(expr, inst)
		var replacement string
		if err != nil {
			replacement = "⟨求值失败:" + expr + "⟩"
		} else {
			replacement = ValueToExplainString(val)
		}
		s = s[:start] + replacement + s[end+2:]
	}
	return s
}
