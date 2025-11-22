package dsl

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
)

// generateVariable 按变量类型和生成规则生成值
func generateVariable(rng *rand.Rand, name string, v Variable, inst *Instance) (interface{}, error) {
	switch v.Kind {
	case "scalar":
		return genScalar(rng, v)
	case "vector":
		return genVector(rng, v)
	case "matrix":
		return genMatrix(rng, v, inst)
	default:
		return nil, fmt.Errorf("unknown kind: %s", v.Kind)
	}
}

func genScalar(rng *rand.Rand, v Variable) (interface{}, error) {
	g, ok := v.Generator["rule"].(string)
	if !ok {
		return nil, errors.New("scalar generator rule missing")
	}
	switch g {
	case "range":
		min := int64(defaultInt(v.Generator, "min", -5))
		max := int64(defaultInt(v.Generator, "max", 5))
		return int64(rng.Intn(int(max-min+1)) + int(min)), nil
	case "from_set":
		set := interfaceToIntSlice(v.Generator["set"])
		if len(set) == 0 {
			return nil, errors.New("from_set empty")
		}
		return set[rng.Intn(len(set))], nil
	default:
		return nil, fmt.Errorf("unsupported scalar generator: %s", g)
	}
}

func genVector(rng *rand.Rand, v Variable) (interface{}, error) {
	n := v.Size
	if n == 0 {
		n = v.Rows
	}
	g, _ := v.Generator["rule"].(string)
	switch g {
	case "", "range":
		min := int64(defaultInt(v.Generator, "min", -5))
		max := int64(defaultInt(v.Generator, "max", 5))
		vec := NewVectorInt(n)
		for i := 0; i < n; i++ {
			vec.V[i] = int64(rng.Intn(int(max-min+1)) + int(min))
		}
		return vec, nil
	case "from_set":
		set := interfaceToIntSlice(v.Generator["set"])
		if len(set) == 0 {
			return nil, errors.New("from_set empty")
		}
		vec := NewVectorInt(n)
		for i := 0; i < n; i++ {
			vec.V[i] = set[rng.Intn(len(set))]
		}
		return vec, nil
	default:
		return nil, fmt.Errorf("unsupported vector generator rule: %s", g)
	}
}

func genMatrix(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	r := v.Rows
	c := v.Cols
	if r == 0 || c == 0 {
		return nil, errors.New("matrix rows/cols required")
	}
	g, _ := v.Generator["rule"].(string)
	switch g {
	case "", "range":
		min := int64(defaultInt(v.Generator, "min", -5))
		max := int64(defaultInt(v.Generator, "max", 5))
		m := NewMatrixInt(r, c)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				m.A[i][j] = int64(rng.Intn(int(max-min+1)) + int(min))
			}
		}
		return m, nil
	case "from_set":
		set := interfaceToIntSlice(v.Generator["set"])
		if len(set) == 0 {
			return nil, errors.New("from_set empty")
		}
		m := NewMatrixInt(r, c)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				m.A[i][j] = set[rng.Intn(len(set))]
			}
		}
		return m, nil
	case "sparse":
		values := interfaceToIntSlice(v.Generator["values"])
		density := defaultFloat(v.Generator, "density", 0.3)
		if len(values) == 0 {
			return nil, errors.New("sparse values empty")
		}
		m := NewMatrixInt(r, c)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				if rng.Float64() <= density {
					m.A[i][j] = values[rng.Intn(len(values))]
				} else {
					m.A[i][j] = 0
				}
			}
		}
		return m, nil
	case "full_rank":
		min := int64(defaultInt(v.Generator, "min", -5))
		max := int64(defaultInt(v.Generator, "max", 5))
		// 使用确定性方法：每次尝试用独立的 RNG，避免失败尝试影响随机数序列
		for attempt := 0; attempt < 200; attempt++ {
			// 为每次尝试创建独立的 RNG，基于主 RNG 的下一个 Int63
			attemptSeed := rng.Int63()
			attemptRng := rand.New(rand.NewSource(attemptSeed))

			m := NewMatrixInt(r, c)
			for i := 0; i < r; i++ {
				for j := 0; j < c; j++ {
					m.A[i][j] = int64(attemptRng.Intn(int(max-min+1)) + int(min))
				}
			}
			if r == c {
				d := BareissDet(m)
				if d.Sign() != 0 {
					return m, nil
				}
			} else {
				if matrixRankRat(m) == minInt(r, c) {
					return m, nil
				}
			}
		}
		return nil, errors.New("failed to generate full_rank matrix after attempts")
	case "integer_solution":
		// 该规则在更高层被特殊处理，不在这里直接生成
		return nil, errors.New("integer_solution should be handled as derived/builder (use derived or orchestrator)")
	case "scalar_identity":
		// 生成 λI 型矩阵，所有对角线为同一随机整数 λ，其余为 0
		// 配置：
		//   lambda_min, lambda_max: λ 取值范围（含端点，默认 [-5,5]）
		//   lambda_var: 写入到实例中的标量变量名（可选），例如 "lambda"
		cfg := v.Generator
		lmin := int64(defaultInt(cfg, "lambda_min", -5))
		lmax := int64(defaultInt(cfg, "lambda_max", 5))
		if lmax < lmin {
			lmin, lmax = lmax, lmin
		}
		var lambda int64
		for {
			lambda = int64(rng.Intn(int(lmax-lmin+1)) + int(lmin))
			if lambda != 0 {
				break
			}
		}
		m := NewMatrixInt(r, c)
		for i := 0; i < r && i < c; i++ {
			m.A[i][i] = lambda
		}
		if name, ok := cfg["lambda_var"].(string); ok && name != "" && inst != nil {
			// 记录 λ，供后续表达式或答案使用
			if inst.Vars == nil {
				inst.Vars = map[string]interface{}{}
			}
			inst.Vars[name] = lambda
		}
		return m, nil
	case "lambda_linear_det_zero":
		// 为 det(A)=0 类型题目生成含单一参数 λ 的 3x3 矩阵，并保证 λ 为整数解
		return genMatrixLambdaLinearDetZero(rng, v, inst)
	default:
		return nil, fmt.Errorf("unsupported matrix generator rule: %s", g)
	}
}

// genMatrixLambdaLinearDetZero
// 生成一个矩阵 A，其中只有一个元素线性依赖于某个参数 λ：
//
//	A[r0,c0] = λ + c
//
// 其它元素为整数随机数。通过 det(A(λ))=0 求解 λ，保证 λ 为给定区间内的整数，
// 并将解写入 inst.Vars[paramVar]，返回的是在该 λ 下的数值矩阵 A(λ)。
//
// 配置项（写在 Variable.Generator 里）：
//
//	rule: "lambda_linear_det_zero"
//	param_var: string   // 参数名，默认 "lambda"
//	param_row: int      // 1-based，默认 rows
//	param_col: int      // 1-based，默认 cols
//	entry_min: int      // 其它元素与常数项 c 的取值下界，默认 -5
//	entry_max: int      // 上界，默认 5
//	lambda_min: int     // λ 取值下界，默认 -10
//	lambda_max: int     // 上界，默认 10
//	max_attempts: int   // 最大重试次数，默认 200
func genMatrixLambdaLinearDetZero(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	r := v.Rows
	c := v.Cols
	if r == 0 || c == 0 {
		return nil, errors.New("lambda_linear_det_zero: matrix rows/cols required")
	}

	cfg := v.Generator
	paramVar, _ := cfg["param_var"].(string)
	if paramVar == "" {
		paramVar = "lambda"
	}
	paramRow := defaultInt(cfg, "param_row", r)
	paramCol := defaultInt(cfg, "param_col", c)
	if paramRow < 1 || paramRow > r || paramCol < 1 || paramCol > c {
		return nil, errors.New("lambda_linear_det_zero: param_row/param_col out of range")
	}

	entryMin := defaultInt(cfg, "entry_min", -5)
	entryMax := defaultInt(cfg, "entry_max", 5)
	if entryMax < entryMin {
		entryMin, entryMax = entryMax, entryMin
	}

	lambdaMin := defaultInt(cfg, "lambda_min", -10)
	lambdaMax := defaultInt(cfg, "lambda_max", 10)
	if lambdaMax < lambdaMin {
		lambdaMin, lambdaMax = lambdaMax, lambdaMin
	}

	maxAttempts := defaultInt(cfg, "max_attempts", 200)
	if maxAttempts <= 0 {
		maxAttempts = 200
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		base := NewMatrixInt(r, c)
		// 随机填充除参数位置外的元素
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				if i == paramRow-1 && j == paramCol-1 {
					continue
				}
				base.A[i][j] = int64(rng.Intn(entryMax-entryMin+1) + entryMin)
			}
		}
		// 参数位置常数项 c（表达式为 λ + c，前端可格式化成 λ - k）
		constC := int64(rng.Intn(entryMax-entryMin+1) + entryMin)
		base.A[paramRow-1][paramCol-1] = constC // 这是 λ=0 时的值

		// d(λ) 只在该元素上线性依赖 λ
		// 令 d0 = det(A(λ=0))，d1 = det(A(λ=1))，则
		//   d(λ) = a*λ + b，其中 b=d0, a=d1-d0
		d0 := BareissDet(base)

		// 构造 λ=1 时的矩阵（在 base 的基础上 +1）
		m1 := NewMatrixInt(r, c)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				m1.A[i][j] = base.A[i][j]
			}
		}
		m1.A[paramRow-1][paramCol-1] = constC + 1
		d1 := BareissDet(m1)

		aBig := new(big.Int).Sub(d1, d0)
		if aBig.Sign() == 0 {
			// 行列式与 λ 无关，跳过
			continue
		}
		bBig := new(big.Int).Set(d0)

		// 求 λ = -b/a，要求是整数
		q := new(big.Int)
		rBig := new(big.Int)
		q.DivMod(bBig, aBig, rBig)
		if rBig.Sign() != 0 {
			// 不是整数解
			continue
		}
		lambdaBig := new(big.Int).Neg(q) // λ = -b/a
		if !lambdaBig.IsInt64() {
			continue
		}
		lambdaVal := lambdaBig.Int64()
		if lambdaVal < int64(lambdaMin) || lambdaVal > int64(lambdaMax) {
			continue
		}

		// 构造在该 λ 下的数值矩阵 A(λ)
		A := NewMatrixInt(r, c)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				A.A[i][j] = base.A[i][j]
			}
		}
		A.A[paramRow-1][paramCol-1] = constC + lambdaVal

		// 校验 det(A(λ)) == 0
		dRoot := BareissDet(A)
		if dRoot.Sign() != 0 {
			continue
		}

		// 将 λ 写入实例，供答案与渲染使用
		inst.Vars[paramVar] = lambdaVal
		return A, nil
	}

	return nil, errors.New("lambda_linear_det_zero: failed to generate suitable matrix after attempts")
}

// interfaceToIntSlice 帮助将 interface{} 转为 []int64
func interfaceToIntSlice(i interface{}) []int64 {
	out := []int64{}
	if i == nil {
		return out
	}
	switch v := i.(type) {
	case []interface{}:
		for _, e := range v {
			switch ee := e.(type) {
			case float64:
				out = append(out, int64(ee))
			case int:
				out = append(out, int64(ee))
			case int64:
				out = append(out, ee)
			default:
				// ignore
			}
		}
	case []int:
		for _, e := range v {
			out = append(out, int64(e))
		}
	case []int64:
		return v
	}
	return out
}

func defaultInt(m map[string]interface{}, key string, def int) int {
	if m == nil {
		return def
	}
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case int64:
			return int(v)
		}
	}
	return def
}

func defaultFloat(m map[string]interface{}, key string, def float64) float64 {
	if m == nil {
		return def
	}
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case int64:
			return float64(v)
		}
	}
	return def
}
