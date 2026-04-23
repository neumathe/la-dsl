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
		return genScalar(rng, v, inst)
	case "vector":
		return genVector(rng, v)
	case "matrix":
		return genMatrix(rng, v, inst)
	default:
		return nil, fmt.Errorf("unknown kind: %s", v.Kind)
	}
}

func genScalar(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
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
	case "eigen_rank_inference_5":
		return genScalarEigenRankInference5(rng, v, inst)
	case "eigen_row_sum_rank4":
		return genScalarEigenRowSumRank4(rng, v, inst)
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
	case "orthogonal_signed_perm":
		// 3×3：随机带符号置换矩阵（正交、det=±1，元素仅为 -1,0,1）
		if r != 3 || c != 3 {
			return nil, errors.New("orthogonal_signed_perm: need 3×3")
		}
		perm := []int{0, 1, 2}
		rng.Shuffle(3, func(i, j int) { perm[i], perm[j] = perm[j], perm[i] })
		m := NewMatrixInt(3, 3)
		for j := 0; j < 3; j++ {
			sign := int64(1)
			if rng.Intn(2) == 0 {
				sign = -1
			}
			m.A[perm[j]][j] = sign
		}
		return m, nil
	case "rank_minus_one_square":
		// n×n 秩恰为 n-1：前 n-1 行随机，末行复制首行（必相关），再筛秩为 n-1。
		if r != c || r < 2 {
			return nil, errors.New("rank_minus_one_square: need n×n with n≥2")
		}
		n := r
		min := int64(defaultInt(v.Generator, "min", -5))
		max := int64(defaultInt(v.Generator, "max", 5))
		if max < min {
			min, max = max, min
		}
		want := n - 1
		for attempt := 0; attempt < 500; attempt++ {
			attemptSeed := rng.Int63()
			attemptRng := rand.New(rand.NewSource(attemptSeed))
			m := NewMatrixInt(n, n)
			for i := 0; i < n-1; i++ {
				for j := 0; j < n; j++ {
					m.A[i][j] = int64(attemptRng.Intn(int(max-min+1)) + int(min))
				}
			}
			for j := 0; j < n; j++ {
				m.A[n-1][j] = m.A[0][j]
			}
			if matrixRankRat(m) == want {
				return m, nil
			}
		}
		return nil, errors.New("rank_minus_one_square: failed")
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
	case "upper_unit":
		// 上三角、主对角全 1，严格上三角随机整数（det=1，便于 Cramer / 整数解）
		return genMatrixUpperUnit(rng, v)
	case "lower_unit":
		// 下三角、主对角全 1，严格下三角随机整数（det=1，逆仍为整数矩阵）
		return genMatrixLowerUnit(rng, v)
	case "equidiagonal":
		// n×n：主对角同一常数 a，非对角同一常数 b（如「对角为 -1、其余为 -3」类题）
		return genMatrixEquidiagonal(rng, v)
	case "upper_triangular":
		// 上三角矩阵，主对角非零随机，det = 对角乘积
		return genMatrixUpperTriangular(rng, v)
	case "upper_triangular_nonzero_diag":
		// 上三角矩阵，主对角元为非零整数（不含 ±1 以避免平凡化），严格上三角随机整数
		return genMatrixUpperTriangularNonzeroDiag(rng, v)
	case "diagonal_distinct":
		// 对角矩阵，主对角元两两不同且非零
		return genMatrixDiagonalDistinct(rng, v)
	case "rank3_coin":
		// 3×3：随机一半概率秩为 2（列相关），一半满秩
		return genMatrixRank3Coin(rng, v)
	case "symmetric":
		// 整数对称矩阵（上三角随机，下三角镜像）
		return genMatrixSymmetric(rng, v)
	case "rank2_2x4":
		if r != 2 || c != 4 {
			return nil, errors.New("rank2_2x4: need 2×4")
		}
		min := int64(defaultInt(v.Generator, "min", -4))
		max := int64(defaultInt(v.Generator, "max", 4))
		if max < min {
			min, max = max, min
		}
		for attempt := 0; attempt < 200; attempt++ {
			m := NewMatrixInt(2, 4)
			for i := 0; i < 2; i++ {
				for j := 0; j < 4; j++ {
					m.A[i][j] = int64(rng.Intn(int(max-min+1)) + int(min))
				}
			}
			if matrixRankRat(m) == 2 {
				return m, nil
			}
		}
		return nil, errors.New("rank2_2x4: failed")
	case "rank2_3x4":
		// 3×4 秩 2：先取 2×4 秩 2 的两行，第三行为前两行之和（必落在行张成内），故整体秩为 2。
		if r != 3 || c != 4 {
			return nil, errors.New("rank2_3x4: need 3×4")
		}
		min := int64(defaultInt(v.Generator, "min", -4))
		max := int64(defaultInt(v.Generator, "max", 4))
		if max < min {
			min, max = max, min
		}
		for attempt := 0; attempt < 400; attempt++ {
			top := NewMatrixInt(2, 4)
			for i := 0; i < 2; i++ {
				for j := 0; j < 4; j++ {
					top.A[i][j] = int64(rng.Intn(int(max-min+1)) + int(min))
				}
			}
			if matrixRankRat(top) != 2 {
				continue
			}
			m := NewMatrixInt(3, 4)
			for j := 0; j < 4; j++ {
				m.A[0][j] = top.A[0][j]
				m.A[1][j] = top.A[1][j]
				m.A[2][j] = top.A[0][j] + top.A[1][j]
			}
			if matrixRankRat(m) == 2 {
				return m, nil
			}
		}
		return nil, errors.New("rank2_3x4: failed")
	case "rank1_outer":
		// m×n 秩 1：外积 u v^T（u 为 m 维列，v 为 n 维行向量随机整数）。
		if r <= 0 || c <= 0 {
			return nil, errors.New("rank1_outer: bad size")
		}
		min := int64(defaultInt(v.Generator, "min", -4))
		max := int64(defaultInt(v.Generator, "max", 4))
		if max < min {
			min, max = max, min
		}
		for attempt := 0; attempt < 80; attempt++ {
			u := NewVectorInt(r)
			w := NewVectorInt(c)
			for i := 0; i < r; i++ {
				u.V[i] = int64(rng.Intn(int(max-min+1)) + int(min))
			}
			for j := 0; j < c; j++ {
				w.V[j] = int64(rng.Intn(int(max-min+1)) + int(min))
			}
			allu0, allw0 := true, true
			for i := 0; i < r; i++ {
				if u.V[i] != 0 {
					allu0 = false
				}
			}
			for j := 0; j < c; j++ {
				if w.V[j] != 0 {
					allw0 = false
				}
			}
			if allu0 || allw0 {
				continue
			}
			M := NewMatrixInt(r, c)
			for i := 0; i < r; i++ {
				for j := 0; j < c; j++ {
					M.A[i][j] = u.V[i] * w.V[j]
				}
			}
			if matrixRankRat(M) == 1 {
				return M, nil
			}
		}
		return nil, errors.New("rank1_outer: failed")
	case "rank334_last_dep":
		// 4×4 列秩 3：最后一列为前三列的整系数线性组合（用于极大无关组/表示式类题）。
		if r != 4 || c != 4 {
			return nil, errors.New("rank334_last_dep: need 4×4")
		}
		min := int64(defaultInt(v.Generator, "min", -4))
		max := int64(defaultInt(v.Generator, "max", 4))
		if max < min {
			min, max = max, min
		}
		for attempt := 0; attempt < 200; attempt++ {
			U := NewMatrixInt(4, 3)
			for i := 0; i < 4; i++ {
				for j := 0; j < 3; j++ {
					U.A[i][j] = int64(rng.Intn(int(max-min+1)) + int(min))
				}
			}
			if matrixRankRat(U) < 3 {
				continue
			}
			a := int64(rng.Intn(int(max-min+1)) + int(min))
			b := int64(rng.Intn(int(max-min+1)) + int(min))
			cc := int64(rng.Intn(int(max-min+1)) + int(min))
			if a == 0 && b == 0 && cc == 0 {
				continue
			}
			V := NewMatrixInt(4, 4)
			for i := 0; i < 4; i++ {
				for j := 0; j < 3; j++ {
					V.A[i][j] = U.A[i][j]
				}
				V.A[i][3] = a*U.A[i][0] + b*U.A[i][1] + cc*U.A[i][2]
			}
			if matrixRankRat(V) != 3 {
				continue
			}
			return V, nil
		}
		return nil, errors.New("rank334_last_dep: failed")
	case "diagonalizable_2x2":
		// 反向生成：随机选择 2 个互异整数特征值 λ₁,λ₂ 和 2×2 可逆矩阵 V（det=±1），
		// 计算 A = V Λ V⁻¹，要求 A 为整数矩阵且不是对角阵/标量阵。
		// 返回 A 并将 eigenvalues/eigenvectors 写入 inst.Vars。
		return genMatrixDiagonalizable2x2(rng, v, inst)
	case "eigen_reverse_3x3":
		// 反向生成：随机选择 3 个互异整数特征值 λ₁,λ₂,λ₃ 和 3 个整数特征向量 v₁,v₂,v₃，
		// 计算 A = V Λ V⁻¹，要求 A 为整数矩阵（V⁻¹ 各元为整数）。
		// 返回 A 并将 eigenvalues/eigenvectors 写入 inst.Vars。
		return genMatrixEigenReverse3x3(rng, v, inst)
	case "symmetric_eigen_reverse_3x3":
		// 对称矩阵反向生成：随机选择 3 个互异整数特征值 λ₁,λ₂,λ₃，
		// 用带符号置换正交矩阵 Q（det=±1，Qᵀ=Q⁻¹）构造 S = QΛQᵀ（整数对称矩阵）。
		// 返回 S 并将特征值和正交矩阵 Q 写入 inst.Vars。
		return genMatrixSymmetricEigenReverse3x3(rng, v, inst)
	case "similarity_congruence_pair":
		// 生成 3×3 对称阵 A 和对角阵 B，并判断它们是否相似/合同。
		// 策略：约一半概率"相似且合同"，约一半概率"不相似但合同"或"不相似不合同"。
		// 将 is_similar、is_congruent 写入 inst.Vars。
		return genMatrixSimilarityCongruencePair(rng, v, inst)
	case "sylvester_range":
		// 生成含参数 t 的 3×3 对称阵，使正定二次型 t 的取值范围为有限开区间 (a,b)。
		// 将 lower, upper 写入 inst.Vars。
		return genMatrixSylvesterRange(rng, v, inst)
	case "param_orthogonal_diag":
		// 生成含参数 t 的 3×3 对称阵 S=QΛQᵀ（S[3][3]=t），给出标准型系数，
		// 学生由 trace 条件推断 t，并求正交变换矩阵 Q。
		return genMatrixParamOrthogonalDiag(rng, v, inst)
	case "poly_schmidt_integral":
		// 生成 3 个多项式（常/线/二次），在积分内积 ∫₋₁¹ f·g dx 下做 Schmidt 正交化。
		// 存储 _poly_schmidt_input/output 在 inst.Vars。
		return genPolySchmidtIntegral(rng, v, inst)
	case "param_infinit_solution":
		return genMatrixParamInfinitSolution(rng, v, inst)
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
		inst.Vars["_lambda_param_row"] = int64(paramRow)
		inst.Vars["_lambda_param_col"] = int64(paramCol)
		inst.Vars["_lambda_param_constC"] = constC
		return A, nil
	}

	return nil, errors.New("lambda_linear_det_zero: failed to generate suitable matrix after attempts")
}

func genMatrixUpperUnit(rng *rand.Rand, v Variable) (interface{}, error) {
	r := v.Rows
	c := v.Cols
	if r == 0 || c == 0 || r != c {
		return nil, errors.New("upper_unit: need square n×n")
	}
	min := int64(defaultInt(v.Generator, "min", -5))
	max := int64(defaultInt(v.Generator, "max", 5))
	if max < min {
		min, max = max, min
	}
	m := NewMatrixInt(r, c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if i == j {
				m.A[i][j] = 1
			} else if j > i {
				m.A[i][j] = int64(rng.Intn(int(max-min+1)) + int(min))
			} else {
				m.A[i][j] = 0
			}
		}
	}
	return m, nil
}

func genMatrixLowerUnit(rng *rand.Rand, v Variable) (interface{}, error) {
	r := v.Rows
	c := v.Cols
	if r == 0 || c == 0 || r != c {
		return nil, errors.New("lower_unit: need square n×n")
	}
	min := int64(defaultInt(v.Generator, "min", -5))
	max := int64(defaultInt(v.Generator, "max", 5))
	if max < min {
		min, max = max, min
	}
	m := NewMatrixInt(r, c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if i == j {
				m.A[i][j] = 1
			} else if j < i {
				m.A[i][j] = int64(rng.Intn(int(max-min+1)) + int(min))
			} else {
				m.A[i][j] = 0
			}
		}
	}
	return m, nil
}

func genMatrixEquidiagonal(rng *rand.Rand, v Variable) (interface{}, error) {
	cfg := v.Generator
	n := v.Rows
	if n == 0 {
		n = defaultInt(cfg, "n", 3)
	}
	if v.Cols != 0 && v.Cols != n {
		return nil, errors.New("equidiagonal: rows must equal cols")
	}
	if n <= 0 {
		return nil, errors.New("equidiagonal: n must be positive")
	}
	dmin := int64(defaultInt(cfg, "diag_min", -3))
	dmax := int64(defaultInt(cfg, "diag_max", 3))
	omin := int64(defaultInt(cfg, "off_min", -5))
	omax := int64(defaultInt(cfg, "off_max", 5))
	if dmax < dmin {
		dmin, dmax = dmax, dmin
	}
	if omax < omin {
		omin, omax = omax, omin
	}
	for attempt := 0; attempt < 200; attempt++ {
		diag := int64(rng.Intn(int(dmax-dmin+1)) + int(dmin))
		off := int64(rng.Intn(int(omax-omin+1)) + int(omin))
		if diag == off {
			continue
		}
		m := NewMatrixInt(n, n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i == j {
					m.A[i][j] = diag
				} else {
					m.A[i][j] = off
				}
			}
		}
		return m, nil
	}
	return nil, errors.New("equidiagonal: failed to pick diag!=off")
}

func genMatrixUpperTriangular(rng *rand.Rand, v Variable) (interface{}, error) {
	r := v.Rows
	c := v.Cols
	if r == 0 || c == 0 || r != c {
		return nil, errors.New("upper_triangular: need square n×n")
	}
	min := int64(defaultInt(v.Generator, "min", -6))
	max := int64(defaultInt(v.Generator, "max", 6))
	if max < min {
		min, max = max, min
	}
	for attempt := 0; attempt < 100; attempt++ {
		m := NewMatrixInt(r, c)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				if i > j {
					m.A[i][j] = 0
				} else if i == j {
					m.A[i][j] = int64(rng.Intn(int(max-min+1)) + int(min))
					if m.A[i][j] == 0 {
						m.A[i][j] = 1
					}
				} else {
					m.A[i][j] = int64(rng.Intn(int(max-min+1)) + int(min))
				}
			}
		}
		if BareissDet(m).Sign() != 0 {
			return m, nil
		}
	}
	return nil, errors.New("upper_triangular: failed non-singular")
}


func genMatrixUpperTriangularNonzeroDiag(rng *rand.Rand, v Variable) (interface{}, error) {
	r := v.Rows
	c := v.Cols
	if r == 0 || c == 0 || r != c {
		return nil, errors.New("upper_triangular_nonzero_diag: need square n×n")
	}
	min := int64(defaultInt(v.Generator, "min", -4))
	max := int64(defaultInt(v.Generator, "max", 4))
	if max < min {
		min, max = max, min
	}
	// Diagonal values: nonzero integers excluding ±1 to avoid trivialization
	// (e.g., β₁ = c·1 becomes just "c" which is too simple)
	diagChoices := []int64{}
	for x := min; x <= max; x++ {
		if x != 0 && x != 1 && x != -1 {
			diagChoices = append(diagChoices, x)
		}
	}
	// Fallback: if no valid choices in range, use ±2, ±3
	if len(diagChoices) == 0 {
		diagChoices = []int64{-3, -2, 2, 3}
	}
	for attempt := 0; attempt < 100; attempt++ {
		m := NewMatrixInt(r, c)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				if i > j {
					m.A[i][j] = 0
				} else if i == j {
					m.A[i][j] = diagChoices[rng.Intn(len(diagChoices))]
				} else {
					m.A[i][j] = int64(rng.Intn(int(max-min+1)) + int(min))
				}
			}
		}
		if BareissDet(m).Sign() != 0 {
			return m, nil
		}
	}
	return nil, errors.New("upper_triangular_nonzero_diag: failed non-singular")
}
func genMatrixDiagonalDistinct(rng *rand.Rand, v Variable) (interface{}, error) {
	r := v.Rows
	c := v.Cols
	if r == 0 || c == 0 || r != c {
		return nil, errors.New("diagonal_distinct: need square n×n")
	}
	min := int64(defaultInt(v.Generator, "min", -6))
	max := int64(defaultInt(v.Generator, "max", 6))
	if max <= min {
		return nil, errors.New("diagonal_distinct: need min < max")
	}
	for attempt := 0; attempt < 300; attempt++ {
		m := NewMatrixInt(r, c)
		used := map[int64]struct{}{}
		ok := true
		for i := 0; i < r; i++ {
			found := false
			for t := 0; t < 80; t++ {
				d := int64(rng.Intn(int(max-min+1)) + int(min))
				if d == 0 {
					continue
				}
				if _, dup := used[d]; dup {
					continue
				}
				used[d] = struct{}{}
				m.A[i][i] = d
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
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				if i != j {
					m.A[i][j] = 0
				}
			}
		}
		return m, nil
	}
	return nil, errors.New("diagonal_distinct: failed")
}

func genMatrixRank3Coin(rng *rand.Rand, v Variable) (interface{}, error) {
	r := v.Rows
	c := v.Cols
	if r != 3 || c != 3 {
		return nil, errors.New("rank3_coin: need 3×3")
	}
	min := int64(defaultInt(v.Generator, "min", -5))
	max := int64(defaultInt(v.Generator, "max", 5))
	if max < min {
		min, max = max, min
	}
	if rng.Intn(2) == 0 {
		m := NewMatrixInt(3, 3)
		for i := 0; i < 3; i++ {
			m.A[i][0] = int64(rng.Intn(int(max-min+1)) + int(min))
			m.A[i][1] = int64(rng.Intn(int(max-min+1)) + int(min))
			m.A[i][2] = m.A[i][0] + m.A[i][1]
		}
		return m, nil
	}
	for attempt := 0; attempt < 120; attempt++ {
		m := NewMatrixInt(3, 3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				m.A[i][j] = int64(rng.Intn(int(max-min+1)) + int(min))
			}
		}
		if BareissDet(m).Sign() != 0 {
			return m, nil
		}
	}
	return nil, errors.New("rank3_coin: failed full rank")
}

func genMatrixSymmetric(rng *rand.Rand, v Variable) (interface{}, error) {
	n := v.Rows
	if n == 0 || n != v.Cols {
		return nil, errors.New("symmetric: need n×n square")
	}
	min := int64(defaultInt(v.Generator, "min", -6))
	max := int64(defaultInt(v.Generator, "max", 6))
	if max < min {
		min, max = max, min
	}
	m := NewMatrixInt(n, n)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			x := int64(rng.Intn(int(max-min+1)) + int(min))
			m.A[i][j] = x
			m.A[j][i] = x
		}
	}
	return m, nil
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

func genMatrixDiagonalizable2x2(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	// Pick 2 distinct nonzero eigenvalues
	lmin := int64(defaultInt(v.Generator, "lambda_min", -5))
	lmax := int64(defaultInt(v.Generator, "lambda_max", 5))
	if lmax < lmin {
		lmin, lmax = lmax, lmin
	}
	emin := int64(defaultInt(v.Generator, "entry_min", -3))
	emax := int64(defaultInt(v.Generator, "entry_max", 3))
	if emax < emin {
		emin, emax = emax, emin
	}

	for attempt := 0; attempt < 500; attempt++ {
		// Use independent RNG for each attempt to avoid exhaustion
		attemptSeed := rng.Int63()
		attemptRng := rand.New(rand.NewSource(attemptSeed))

		// Pick 2 distinct nonzero eigenvalues
		lambdas := make([]int64, 2)
		used := map[int64]bool{}
		ok := true
		for i := 0; i < 2; i++ {
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

		// Generate V with det=±1: compose upper_unit × lower_unit
		// This guarantees det(V) = 1 without needing to check
		U := NewMatrixInt(2, 2)
		U.A[0][0] = 1
		U.A[1][1] = 1
		U.A[1][0] = 0
		U.A[0][1] = int64(attemptRng.Intn(int(emax-emin+1)) + int(emin))

		L := NewMatrixInt(2, 2)
		L.A[0][0] = 1
		L.A[1][1] = 1
		L.A[0][1] = 0
		L.A[1][0] = int64(attemptRng.Intn(int(emax-emin+1)) + int(emin))

		// V = U × L, det(V) = det(U)×det(L) = 1×1 = 1
		V, err := matrixMulInt(U, L)
		if err != nil {
			continue
		}

		// Compute V⁻¹
		Vinv, err := MatrixInverseInt(V)
		if err != nil {
			continue
		}

		// Construct Λ = diag(λ₁, λ₂)
		Lambda := NewMatrixInt(2, 2)
		for i := 0; i < 2; i++ {
			Lambda.A[i][i] = lambdas[i]
		}

		// A = V Λ V⁻¹
		VL, err := matrixMulInt(V, Lambda)
		if err != nil {
			continue
		}
		A, err := matrixMulInt(VL, Vinv)
		if err != nil {
			continue
		}

		// Verify A has integer entries within reasonable bounds
		maxEntry := int64(0)
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				abs := A.A[i][j]
				if abs < 0 {
					abs = -abs
				}
				if abs > maxEntry {
					maxEntry = abs
				}
			}
		}
		if maxEntry > 15 {
			continue
		}

		// Reject if A is diagonal or scalar matrix (would be trivial)
		if A.A[0][1] == 0 && A.A[1][0] == 0 {
			continue
		}

		// Store eigenvalues and eigenvectors for answer extraction
		if inst.Vars == nil {
			inst.Vars = map[string]interface{}{}
		}
		inst.Vars["_eigen_V"] = V
		inst.Vars["_eigen_lambda1"] = lambdas[0]
		inst.Vars["_eigen_lambda2"] = lambdas[1]

		return A, nil
	}
	return nil, errors.New("diagonalizable_2x2: failed")
}

// randomUnimodular3x3FromElementary 从单位阵经若干次「行_i += k·行_j」得到 det=1 的整数幺模阵。
// 纯随机填元几乎不可能满足 |det|=1，故用初等阵乘积构造。
func randomUnimodular3x3FromElementary(r *rand.Rand, vMax int64) *MatrixInt {
	for t := 0; t < 80; t++ {
		V := NewMatrixInt(3, 3)
		for i := 0; i < 3; i++ {
			V.A[i][i] = 1
		}
		steps := 4 + r.Intn(6)
		for s := 0; s < steps; s++ {
			i := r.Intn(3)
			j := r.Intn(3)
			for j == i {
				j = r.Intn(3)
			}
			k := int64(1)
			if r.Intn(2) == 0 {
				k = -1
			}
			if r.Intn(4) == 0 {
				k *= 2
			}
			for c := 0; c < 3; c++ {
				V.A[i][c] += k * V.A[j][c]
			}
		}
		maxV := int64(0)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				x := V.A[i][j]
				if x < 0 {
					x = -x
				}
				if x > maxV {
					maxV = x
				}
			}
		}
		if maxV > vMax {
			continue
		}
		detV := BareissDet(V)
		if detV.Sign() == 0 {
			continue
		}
		if new(big.Int).Abs(detV).Cmp(big.NewInt(1)) != 0 {
			continue
		}
		return V
	}
	return nil
}

func genMatrixEigenReverse3x3(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	// Choose 3 distinct nonzero eigenvalues
	lmin := int64(defaultInt(v.Generator, "lambda_min", -5))
	lmax := int64(defaultInt(v.Generator, "lambda_max", 5))
	if lmax < lmin {
		lmin, lmax = lmax, lmin
	}

	for attempt := 0; attempt < 500; attempt++ {
		// Use independent RNG for each attempt to avoid exhaustion
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

		// 列向量为特征向量：需 det(V)=±1 才有整数 V⁻¹。随机小整数几乎不满足，用初等变换从 I 生成幺模 V。
		vEntryMax := int64(defaultInt(v.Generator, "v_entry_max", 12))
		V := randomUnimodular3x3FromElementary(attemptRng, vEntryMax)
		if V == nil {
			continue
		}

		// Compute V⁻¹
		Vinv, err := MatrixInverseInt(V)
		if err != nil {
			continue
		}

		// Construct Λ = diag(λ₁, λ₂, λ₃)
		Lambda := NewMatrixInt(3, 3)
		for i := 0; i < 3; i++ {
			Lambda.A[i][i] = lambdas[i]
		}

		// A = V Λ V⁻¹
		VL, err := matrixMulInt(V, Lambda)
		if err != nil {
			continue
		}
		A, err := matrixMulInt(VL, Vinv)
		if err != nil {
			continue
		}

		// Verify A entries are reasonable size
		maxEntry := int64(0)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if A.A[i][j] > maxEntry {
					maxEntry = A.A[i][j]
				}
				if -A.A[i][j] > maxEntry {
					maxEntry = -A.A[i][j]
				}
			}
		}
		if maxEntry > 15 {
			continue
		}

		// Store eigenvalues and eigenvectors for answer extraction
		if inst.Vars == nil {
			inst.Vars = map[string]interface{}{}
		}
		inst.Vars["_eigen_V"] = V
		inst.Vars["_eigen_Vinv"] = Vinv
		inst.Vars["_eigen_Lambda"] = Lambda
		inst.Vars["_eigen_lambda1"] = lambdas[0]
		inst.Vars["_eigen_lambda2"] = lambdas[1]
		inst.Vars["_eigen_lambda3"] = lambdas[2]

		return A, nil
	}
	return nil, errors.New("eigen_reverse_3x3: failed")
}

// genMatrixSymmetricEigenReverse3x3 生成 3×3 整数对称矩阵 S = Q Λ Qᵀ，
// 其中 Λ = diag(λ₁,λ₂,λ₃) 是随机互异非零整数特征值，
// Q 是行列式 ±1 的带符号置换正交矩阵（保证 S 为整数且对称）。
// 将 λ₁,λ₂,λ₃、Q、Λ 写入 inst.Vars 供答案提取。
func genMatrixSymmetricEigenReverse3x3(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	r := v.Rows
	c := v.Cols
	if r != 3 || c != 3 {
		return nil, errors.New("symmetric_eigen_reverse_3x3: need 3×3")
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

		// Generate a signed permutation orthogonal matrix Q (det = ±1)
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

		// S = Q Λ Qᵀ（Qᵀ = Q⁻¹ for signed-perm matrices）
		// Compute QΛ first, then (QΛ)·Qᵀ
		QLambda, err := matrixMulInt(Q, Lambda)
		if err != nil {
			continue
		}
		// Qᵀ: transpose of Q (since Q is signed-perm, Qᵀ = Q⁻¹)
		Qt := NewMatrixInt(3, 3)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				Qt.A[i][j] = Q.A[j][i]
			}
		}
		S, err := matrixMulInt(QLambda, Qt)
		if err != nil {
			continue
		}

		// Verify S is symmetric
		for i := 0; i < 3; i++ {
			for j := i + 1; j < 3; j++ {
				if S.A[i][j] != S.A[j][i] {
					continue // shouldn't happen but just in case
				}
			}
		}

		// Check entries are reasonable
		me := int64(0)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if S.A[i][j] > me {
					me = S.A[i][j]
				}
				if -S.A[i][j] > me {
					me = -S.A[i][j]
				}
			}
		}
		if me > maxEntry {
			continue
		}

		// Store eigenvalues and Q for answer extraction
		if inst.Vars == nil {
			inst.Vars = map[string]interface{}{}
		}
		inst.Vars["_sym_eigen_Q"] = Q
		inst.Vars["_sym_eigen_Lambda"] = Lambda
		inst.Vars["_sym_eigen_lambda1"] = lambdas[0]
		inst.Vars["_sym_eigen_lambda2"] = lambdas[1]
		inst.Vars["_sym_eigen_lambda3"] = lambdas[2]

		return S, nil
	}
	return nil, errors.New("symmetric_eigen_reverse_3x3: failed")
}

// genMatrixSimilarityCongruencePair 生成 3×3 对称阵 A 和对角阵 B，
// 判断是否相似（特征值相同）和合同（惯性指数相同）。
// 策略：
//   - 约 1/3 概率"相似且合同"（B 的对角元恰好是 A 的特征值）
//   - 约 1/3 概率"不相似但合同"（B 与 A 惯性指数相同但特征值不同）
//   - 约 1/3 概率"既不相似也不合同"（B 与 A 惯性指数不同）
//
// 将 is_similar、is_congruent (int64, 0 或 1) 写入 inst.Vars。
// 返回 A（对称矩阵），B 通过 derived 生成。
func genMatrixSimilarityCongruencePair(rng *rand.Rand, v Variable, inst *Instance) (interface{}, error) {
	r := v.Rows
	c := v.Cols
	if r != 3 || c != 3 {
		return nil, errors.New("similarity_congruence_pair: need 3×3")
	}
	lmin := int64(defaultInt(v.Generator, "lambda_min", -6))
	lmax := int64(defaultInt(v.Generator, "lambda_max", 6))
	if lmax < lmin {
		lmin, lmax = lmax, lmin
	}
	maxEntry := int64(defaultInt(v.Generator, "max_entry", 15))

	for attempt := 0; attempt < 300; attempt++ {
		// Pick 3 distinct nonzero eigenvalues for A
		lambdas := make([]int64, 3)
		used := map[int64]bool{}
		ok := true
		for i := 0; i < 3; i++ {
			found := false
			for t := 0; t < 80; t++ {
				l := int64(rng.Intn(int(lmax-lmin+1)) + int(lmin))
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

		// Count A's inertia: positive and negative eigenvalues
		nPosA, nNegA := 0, 0
		for _, l := range lambdas {
			if l > 0 {
				nPosA++
			} else {
				nNegA++
			}
		}

		// Generate signed-perm Q for S = QΛQᵀ
		perm := []int{0, 1, 2}
		rng.Shuffle(3, func(i, j int) { perm[i], perm[j] = perm[j], perm[i] })
		Q := NewMatrixInt(3, 3)
		for j := 0; j < 3; j++ {
			sign := int64(1)
			if rng.Intn(2) == 0 {
				sign = -1
			}
			Q.A[perm[j]][j] = sign
		}

		Lambda := NewMatrixInt(3, 3)
		for i := 0; i < 3; i++ {
			Lambda.A[i][i] = lambdas[i]
		}

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
		A, err := matrixMulInt(QLambda, Qt)
		if err != nil {
			continue
		}

		// Check A entries are reasonable
		me := int64(0)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if A.A[i][j] > me {
					me = A.A[i][j]
				}
				if -A.A[i][j] > me {
					me = -A.A[i][j]
				}
			}
		}
		if me > maxEntry {
			continue
		}

		// Decide scenario
		scenario := rng.Intn(3) // 0=similar+congruent, 1=not_similar+congruent, 2=not_similar+not_congruent
		var BDiag [3]int64
		var isSimilar, isCongruent int64

		switch scenario {
		case 0:
			// Similar and congruent: B's diagonal = A's eigenvalues (possibly reordered)
			bperm := []int{0, 1, 2}
			rng.Shuffle(3, func(i, j int) { bperm[i], bperm[j] = bperm[j], bperm[i] })
			for i := 0; i < 3; i++ {
				BDiag[i] = lambdas[bperm[i]]
			}
			isSimilar = 1
			isCongruent = 1
		case 1:
			// Not similar but congruent: B has different eigenvalues but same inertia
			// Generate 3 distinct nonzero eigenvalues with same n_pos/n_neg but not equal to A's eigenvalues
			BDiag[0] = 0
			BDiag[1] = 0
			BDiag[2] = 0
			for t := 0; t < 200; t++ {
				b0 := int64(rng.Intn(int(lmax-lmin+1)) + int(lmin))
				b1 := int64(rng.Intn(int(lmax-lmin+1)) + int(lmin))
				b2 := int64(rng.Intn(int(lmax-lmin+1)) + int(lmin))
				if b0 == 0 || b1 == 0 || b2 == 0 {
					continue
				}
				if b0 == b1 || b1 == b2 || b0 == b2 {
					continue
				}
				// Same inertia as A
				bnPos, bnNeg := 0, 0
				if b0 > 0 {
					bnPos++
				} else {
					bnNeg++
				}
				if b1 > 0 {
					bnPos++
				} else {
					bnNeg++
				}
				if b2 > 0 {
					bnPos++
				} else {
					bnNeg++
				}
				if bnPos != nPosA || bnNeg != nNegA {
					continue
				}
				// Not similar: at least one eigenvalue different
				bSet := map[int64]bool{b0: true, b1: true, b2: true}
				aSet := map[int64]bool{lambdas[0]: true, lambdas[1]: true, lambdas[2]: true}
				// Same inertia but different sets of eigenvalues
				sameSet := true
				for k := range bSet {
					if !aSet[k] {
						sameSet = false
						break
					}
				}
				if sameSet {
					continue // eigenvalues are the same set (just reordered)
				}
				BDiag[0] = b0
				BDiag[1] = b1
				BDiag[2] = b2
				break
			}
			if BDiag[0] == 0 {
				continue // failed to find suitable B diagonal
			}
			isSimilar = 0
			isCongruent = 1
		case 2:
			// Not similar and not congruent: B has different inertia
			for t := 0; t < 200; t++ {
				b0 := int64(rng.Intn(int(lmax-lmin+1)) + int(lmin))
				b1 := int64(rng.Intn(int(lmax-lmin+1)) + int(lmin))
				b2 := int64(rng.Intn(int(lmax-lmin+1)) + int(lmin))
				if b0 == 0 || b1 == 0 || b2 == 0 {
					continue
				}
				if b0 == b1 || b1 == b2 || b0 == b2 {
					continue
				}
				bnPos, bnNeg := 0, 0
				if b0 > 0 {
					bnPos++
				} else {
					bnNeg++
				}
				if b1 > 0 {
					bnPos++
				} else {
					bnNeg++
				}
				if b2 > 0 {
					bnPos++
				} else {
					bnNeg++
				}
				// Different inertia from A
				if bnPos == nPosA && bnNeg == nNegA {
					continue
				}
				BDiag[0] = b0
				BDiag[1] = b1
				BDiag[2] = b2
				break
			}
			if BDiag[0] == 0 {
				continue
			}
			isSimilar = 0
			isCongruent = 0
		}

		// Construct B as diagonal matrix
		B := NewMatrixInt(3, 3)
		for i := 0; i < 3; i++ {
			B.A[i][i] = BDiag[i]
		}

		// Store for answer extraction
		if inst.Vars == nil {
			inst.Vars = map[string]interface{}{}
		}
		inst.Vars["_sc_A"] = A
		inst.Vars["_sc_B"] = B
		inst.Vars["_sc_is_similar"] = isSimilar
		inst.Vars["_sc_is_congruent"] = isCongruent

		return A, nil
	}
	return nil, errors.New("similarity_congruence_pair: failed")
}
