package dsl

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"strings"
)

// InstantiateProblem 根据 DSL Problem 与 seed 实例化一道题
func InstantiateProblem(p Problem, seedStr string, serverSalt string) (*Instance, error) {
	inst := &Instance{
		ProblemID: p.ID,
		Seed:      seedStr,
		Vars:      map[string]interface{}{},
		Derived:   map[string]interface{}{},
	}
	seed := deriveSeed(seedStr, fmt.Sprintf("%d", p.ID), p.Version, serverSalt)
	rng := rand.New(rand.NewSource(seed))

	// 对变量名排序，确保确定性的生成顺序（Go map 迭代顺序是随机的）
	varNames := make([]string, 0, len(p.Variables))
	for name := range p.Variables {
		varNames = append(varNames, name)
	}
	sort.Strings(varNames)

	for _, name := range varNames {
		v := p.Variables[name]
		if v.Fixed != nil {
			// 将 fixed 值转换为正确的类型
			fixedVal, err := convertFixed(v)
			if err != nil {
				return nil, fmt.Errorf("convert fixed variable %s error: %w", name, err)
			}
			inst.Vars[name] = fixedVal
			continue
		}
		val, err := generateVariable(rng, name, v, inst)
		if err != nil {
			if v.Generator != nil {
				if rule, ok := v.Generator["rule"].(string); ok && rule == "integer_solution" {
					continue
				}
			}
			return nil, fmt.Errorf("generate variable %s error: %w", name, err)
		}
		inst.Vars[name] = val
	}

	// 对派生变量名排序，确保确定性
	derivedNames := make([]string, 0, len(p.Derived))
	for name := range p.Derived {
		derivedNames = append(derivedNames, name)
	}
	sort.Strings(derivedNames)

	for _, name := range derivedNames {
		expr := strings.TrimSpace(p.Derived[name])
		if strings.HasPrefix(expr, "integer_solution(") {
			ins := insideParens(expr)
			parts := splitArgs(ins)
			if len(parts) != 2 {
				return nil, fmt.Errorf("integer_solution expects 2 args A,x")
			}
			Aname := strings.TrimSpace(parts[0])
			xname := strings.TrimSpace(parts[1])
			if _, ok := inst.Vars[xname]; !ok {
				if xv, ok := p.Variables[xname]; ok {
					val, err := generateVariable(rng, xname, xv, inst)
					if err != nil {
						return nil, err
					}
					inst.Vars[xname] = val
				} else {
					return nil, fmt.Errorf("x var %s not found", xname)
				}
			}
			if _, ok := inst.Vars[Aname]; !ok {
				if Av, ok := p.Variables[Aname]; ok {
					val, err := generateVariable(rng, Aname, Av, inst)
					if err != nil {
						return nil, err
					}
					inst.Vars[Aname] = val
				} else {
					return nil, fmt.Errorf("matrix var %s not found", Aname)
				}
			}
			A, ok1 := inst.Vars[Aname].(*MatrixInt)
			x, ok2 := inst.Vars[xname].(*VectorInt)
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("integer_solution expects matrix and vector")
			}
			b := NewVectorInt(A.R)
			for i := 0; i < A.R; i++ {
				var s int64
				for j := 0; j < A.C; j++ {
					s += A.A[i][j] * x.V[j]
				}
				b.V[i] = s
			}
			inst.Vars[name] = b
			inst.Derived[name] = b
			continue
		}
		val, err := EvaluateExpression(expr, inst)
		if err != nil {
			return nil, fmt.Errorf("evaluate derived %s (%s): %w", name, expr, err)
		}
		inst.Vars[name] = val
		inst.Derived[name] = val
	}

	for name, v := range p.Variables {
		if v.Generator != nil {
			if rule, ok := v.Generator["rule"].(string); ok && rule == "integer_solution" {
				Aref, _ := v.Generator["A"].(string)
				Xref, _ := v.Generator["x"].(string)
				if Aref == "" || Xref == "" {
					return nil, fmt.Errorf("integer_solution var %s missing A/x refs", name)
				}
				if _, ok := inst.Vars[Aref]; !ok {
					if Av, ok := p.Variables[Aref]; ok {
						val, err := generateVariable(rng, Aref, Av, inst)
						if err != nil {
							return nil, err
						}
						inst.Vars[Aref] = val
					} else {
						return nil, fmt.Errorf("aref %s not found", Aref)
					}
				}
				if _, ok := inst.Vars[Xref]; !ok {
					if Xv, ok := p.Variables[Xref]; ok {
						val, err := generateVariable(rng, Xref, Xv, inst)
						if err != nil {
							return nil, err
						}
						inst.Vars[Xref] = val
					} else {
						return nil, fmt.Errorf("xref %s not found", Xref)
					}
				}
				A := inst.Vars[Aref].(*MatrixInt)
				x := inst.Vars[Xref].(*VectorInt)
				b := NewVectorInt(A.R)
				for i := 0; i < A.R; i++ {
					var s int64
					for j := 0; j < A.C; j++ {
						s += A.A[i][j] * x.V[j]
					}
					b.V[i] = s
				}
				inst.Vars[name] = b
			}
		}
	}
	return inst, nil
}

// RenderInst 根据 Problem.Render 生成前端可用的渲染变量
func RenderInst(p Problem, inst *Instance) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	for key, expr := range p.Render {
		val, err := EvaluateExpression(expr, inst)
		if err != nil {
			if v, ok := inst.Vars[expr]; ok {
				out[key] = v
				continue
			}
			return nil, fmt.Errorf("render %s expr %s: %w", key, expr, err)
		}
		out[key] = val
	}
	return out, nil
}

// ExtractAnswer 根据 AnswerSchema 计算标准答案
func ExtractAnswer(p Problem, inst *Instance) (interface{}, error) {
	if p.Answer.Expression != "" {
		return EvaluateExpression(p.Answer.Expression, inst)
	}
	if len(p.Answer.Fields) > 0 {
		res := []interface{}{}
		for _, f := range p.Answer.Fields {
			val, err := EvaluateExpression(f, inst)
			if err != nil {
				if strings.Contains(f, "[") {
					idxOpen := strings.Index(f, "[")
					idxClose := strings.Index(f, "]")
					if idxOpen > 0 && idxClose > idxOpen {
						vname := strings.TrimSpace(f[:idxOpen])
						idx := mustAtoi(f[idxOpen+1 : idxClose])
						if vv, ok := inst.Vars[vname]; ok {
							switch t := vv.(type) {
							case *VectorInt:
								if idx >= 1 && idx <= t.N {
									res = append(res, t.V[idx-1])
									continue
								}
							case []*big.Rat:
								if idx >= 1 && idx <= len(t) {
									res = append(res, t[idx-1])
									continue
								}
							}
						}
					}
				}
				return nil, fmt.Errorf("field %s eval error: %w", f, err)
			}
			res = append(res, val)
		}
		return res, nil
	}
	return nil, errors.New("no answer schema")
}

// ExtractAnswerWithMeta 在 ExtractAnswer 的基础上，提供按空 ID 封装好的答案列表：
// - 对于单一 Expression，返回一个 ID 固定为 "ans" 的字段
// - 对于 FieldDefs（推荐），使用 fieldDef.ID 与 fieldDef.Expr
// - 对于老的 Fields，自动生成 ID：field_1, field_2, ...
func ExtractAnswerWithMeta(p Problem, inst *Instance) ([]AnswerField, error) {
	// 1) Expression 场景：单个答案
	if p.Answer.Expression != "" {
		val, err := EvaluateExpression(p.Answer.Expression, inst)
		if err != nil {
			return nil, err
		}
		id := "ans"
		if len(p.Answer.FieldDefs) > 0 && p.Answer.FieldDefs[0].ID != "" {
			id = p.Answer.FieldDefs[0].ID
		}
		return []AnswerField{
			{
				ID:    id,
				Expr:  p.Answer.Expression,
				Value: val,
			},
		}, nil
	}

	// 2) 新版 FieldDefs：显式指定每个空的 ID 和表达式
	if len(p.Answer.FieldDefs) > 0 {
		fields := make([]AnswerField, 0, len(p.Answer.FieldDefs))
		for i, fd := range p.Answer.FieldDefs {
			if fd.Expr == "" {
				continue
			}
			val, err := EvaluateExpression(fd.Expr, inst)
			if err != nil {
				return nil, err
			}
			id := fd.ID
			if id == "" {
				id = fmt.Sprintf("field_%d", i+1)
			}
			fields = append(fields, AnswerField{
				ID:    id,
				Expr:  fd.Expr,
				Value: val,
			})
		}
		return fields, nil
	}

	// 3) 兼容老的 Fields：复用 ExtractAnswer 的逻辑并生成默认 ID
	if len(p.Answer.Fields) > 0 {
		raw, err := ExtractAnswer(p, inst)
		if err != nil {
			return nil, err
		}
		var vals []interface{}
		switch t := raw.(type) {
		case []interface{}:
			vals = t
		default:
			vals = []interface{}{t}
		}
		n := len(p.Answer.Fields)
		if len(vals) < n {
			n = len(vals)
		}
		fields := make([]AnswerField, 0, n)
		for i := 0; i < n; i++ {
			expr := p.Answer.Fields[i]
			id := fmt.Sprintf("field_%d", i+1)
			fields = append(fields, AnswerField{
				ID:    id,
				Expr:  expr,
				Value: vals[i],
			})
		}
		return fields, nil
	}

	return nil, errors.New("no answer schema")
}

// solveLinearSystemRat 解 Ax=b，返回有理数向量
func solveLinearSystemRat(A *MatrixInt, b *VectorInt) ([]*big.Rat, error) {
	n := A.R
	if A.C != n {
		return nil, errors.New("matrix A must be square for solve")
	}
	if b.N != n {
		return nil, errors.New("b size mismatch")
	}
	M := make([][]*big.Rat, n)
	for i := 0; i < n; i++ {
		M[i] = make([]*big.Rat, n+1)
		for j := 0; j < n; j++ {
			M[i][j] = new(big.Rat).SetInt64(A.A[i][j])
		}
		M[i][n] = new(big.Rat).SetInt64(b.V[i])
	}
	row := 0
	for col := 0; col < n && row < n; col++ {
		pivot := -1
		for r := row; r < n; r++ {
			if M[r][col].Sign() != 0 {
				pivot = r
				break
			}
		}
		if pivot == -1 {
			continue
		}
		if pivot != row {
			M[pivot], M[row] = M[row], M[pivot]
		}
		pv := new(big.Rat).Set(M[row][col])
		for j := col; j < n+1; j++ {
			M[row][j].Quo(M[row][j], pv)
		}
		for r := 0; r < n; r++ {
			if r == row {
				continue
			}
			f := new(big.Rat).Set(M[r][col])
			if f.Sign() == 0 {
				continue
			}
			for j := col; j < n+1; j++ {
				tmp := new(big.Rat).Mul(f, M[row][j])
				M[r][j].Sub(M[r][j], tmp)
			}
		}
		row++
	}
	x := make([]*big.Rat, n)
	for i := 0; i < n; i++ {
		x[i] = new(big.Rat).Set(M[i][n])
	}
	return x, nil
}

// convertFixed 将 Variable.Fixed 从 interface{} 转换为正确的类型
func convertFixed(v Variable) (interface{}, error) {
	if v.Fixed == nil {
		return nil, errors.New("fixed is nil")
	}

	switch v.Kind {
	case "scalar":
		// 可能是 int, int64, float64 等
		switch val := v.Fixed.(type) {
		case int:
			return int64(val), nil
		case int64:
			return val, nil
		case float64:
			return int64(val), nil
		default:
			return nil, fmt.Errorf("unsupported scalar fixed type: %T", v.Fixed)
		}

	case "vector":
		// Fixed 应该是 []interface{} 或 *VectorInt
		if vec, ok := v.Fixed.(*VectorInt); ok {
			return vec, nil
		}
		arr, ok := v.Fixed.([]interface{})
		if !ok {
			return nil, fmt.Errorf("vector fixed expects []interface{}, got %T", v.Fixed)
		}
		size := v.Size
		if size == 0 {
			size = len(arr)
		}
		vec := NewVectorInt(size)
		for i := 0; i < size && i < len(arr); i++ {
			vec.V[i] = toInt64(arr[i])
		}
		return vec, nil

	case "matrix":
		// Fixed 应该是 [][]interface{} 或 *MatrixInt
		if mat, ok := v.Fixed.(*MatrixInt); ok {
			return mat, nil
		}
		arr2d, ok := v.Fixed.([][]interface{})
		if !ok {
			return nil, fmt.Errorf("matrix fixed expects [][]interface{}, got %T", v.Fixed)
		}
		rows := v.Rows
		cols := v.Cols
		if rows == 0 {
			rows = len(arr2d)
		}
		if cols == 0 && rows > 0 {
			cols = len(arr2d[0])
		}
		mat := NewMatrixInt(rows, cols)
		for i := 0; i < rows && i < len(arr2d); i++ {
			for j := 0; j < cols && j < len(arr2d[i]); j++ {
				mat.A[i][j] = toInt64(arr2d[i][j])
			}
		}
		return mat, nil

	default:
		return nil, fmt.Errorf("unknown variable kind: %s", v.Kind)
	}
}

// toInt64 将 interface{} 转换为 int64
func toInt64(v interface{}) int64 {
	switch val := v.(type) {
	case int:
		return int64(val)
	case int64:
		return val
	case float64:
		return int64(val)
	default:
		return 0
	}
}
