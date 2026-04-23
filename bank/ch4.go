package bank

import (
	"fmt"

	"github.com/neumathe/la-dsl/dsl"
)

func buildChapter4_8() dsl.Problem {
	k := "Chapter4_8"
	ids := BlankIDs(k, 12)
	// 通解 = η₁ + k₁(η₂-η₁) + k₂(η₃-η₁)，其中 η₂-η₁ 和 η₃-η₁ 是基础解系向量
	ja := &dsl.AnswerJudgeSpec{Kind: "affine_rational", AffineGroup: "p0", MatrixVar: "A", BVecVar: "b2"}
	j1 := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: "xi1"}
	j2 := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: "xi2"}
	fds := []dsl.AnswerFieldDef{
		{ID: ids[0], Expr: "eta1[1]", Judge: ja}, {ID: ids[1], Expr: "eta1[2]", Judge: ja},
		{ID: ids[2], Expr: "eta1[3]", Judge: ja}, {ID: ids[3], Expr: "eta1[4]", Judge: ja},
		{ID: ids[4], Expr: "nullbasis_comp(A,1,1)", Judge: j1},
		{ID: ids[5], Expr: "nullbasis_comp(A,1,2)", Judge: j1},
		{ID: ids[6], Expr: "nullbasis_comp(A,1,3)", Judge: j1},
		{ID: ids[7], Expr: "nullbasis_comp(A,1,4)", Judge: j1},
		{ID: ids[8], Expr: "nullbasis_comp(A,2,1)", Judge: j2},
		{ID: ids[9], Expr: "nullbasis_comp(A,2,2)", Judge: j2},
		{ID: ids[10], Expr: "nullbasis_comp(A,2,3)", Judge: j2},
		{ID: ids[11], Expr: "nullbasis_comp(A,2,4)", Judge: j2},
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $Ax=b$ 的三个解为 $\eta_1={{eta1}}$、$\eta_2={{eta2}}$、$\eta_3={{eta3}}$，且 $\mathrm{rank}(A)=2$。写出通解：一特解 + 基础解系中两个向量（各分量；基础解系向量允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A":  {Kind: "matrix", Rows: 2, Cols: 4, Generator: map[string]interface{}{"rule": "rank2_2x4", "min": -4, "max": 4}},
			"x0": {Kind: "vector", Size: 4, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
		},
		Derived: map[string]string{
			"b2":   "A * x0",
			"nb1":  "nullbasis_vec(A,1)",
			"nb2":  "nullbasis_vec(A,2)",
			"eta1": "x0",
			"eta2": "vecadd(x0,nb1)",
			"eta3": "vecadd(x0,nb2)",
		},
		Render: map[string]string{"eta1": "eta1", "eta2": "eta2", "eta3": "eta3"},
		Answer:  dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：已知非齐次方程组 $Ax=b$ 的三个特解 $\eta_1,\eta_2,\eta_3$，利用非齐次解的性质：任意两个非齐次解之差是齐次方程 $Ax=0$ 的解。由于 $\mathrm{rank}(A)=2$（$A$ 为 $2\times 4$ 矩阵），基础解系含 $4-2=2$ 个线性无关向量，恰好 $\eta_2-\eta_1$ 和 $\eta_3-\eta_1$ 线性无关，构成一组基础解系。

**步骤 1**：取已知特解 $\eta_1={{eta1}}$，其 4 个分量依次填写前 4 个空。

**步骤 2**：计算齐次基础解系第一个向量 $\xi_1=\eta_2-\eta_1$，各分量为：
- $\xi_{1,1}={{expr:nullbasis_comp(A,1,1)}}$
- $\xi_{1,2}={{expr:nullbasis_comp(A,1,2)}}$
- $\xi_{1,3}={{expr:nullbasis_comp(A,1,3)}}$
- $\xi_{1,4}={{expr:nullbasis_comp(A,1,4)}}$

**步骤 3**：计算齐次基础解系第二个向量 $\xi_2=\eta_3-\eta_1$，各分量为：
- $\xi_{2,1}={{expr:nullbasis_comp(A,2,1)}}$
- $\xi_{2,2}={{expr:nullbasis_comp(A,2,2)}}$
- $\xi_{2,3}={{expr:nullbasis_comp(A,2,3)}}$
- $\xi_{2,4}={{expr:nullbasis_comp(A,2,4)}}$

**步骤 4**：通解为 $x=\eta_1+k_1\xi_1+k_2\xi_2={{eta1}}+k_1\xi_1+k_2\xi_2$。`,
		},
	}
}

func buildChapter4_3_1() dsl.Problem { return homBasis4x5("Chapter4_3_1") }
func buildChapter4_3_2() dsl.Problem { return homBasis4x5("Chapter4_3_2") }
func buildChapter4_3_3() dsl.Problem { return homBasis4x5("Chapter4_3_3") }

func homBasis4x5(key string) dsl.Problem {
	ids := BlankIDs(key, 20)
	fds := make([]dsl.AnswerFieldDef, 0, 20)
	for g := 1; g <= 4; g++ {
		j := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: fmt.Sprintf("xi%d", g)}
		for i := 1; i <= 5; i++ {
			k := (g-1)*5 + i
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[k-1],
				Expr:   fmt.Sprintf("nullbasis_comp(A,%d,%d)", g, i),
				Judge:  j,
			})
		}
	}
	return dsl.Problem{
		ID:      ProblemID(key),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $A={{A}}$，填写齐次方程 $Ax=0$ 的一组基础解系中四个向量（每个向量 5 个分量；答案允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 5, Generator: map[string]interface{}{"rule": "rank1_outer", "min": -4, "max": 4}},
		},
		Render: map[string]string{"A": "A"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：$A$ 是 $4\times 5$ 秩 1 矩阵，由秩-零化度定理，齐次方程 $Ax=0$ 的解空间维数（零化度）为 $5-\mathrm{rank}(A)=5-1=4$，故基础解系含 4 个线性无关的 5 维向量。

**步骤 1**：计算系数矩阵的秩 $\mathrm{rank}(A)={{expr:rank(A)}}=1$，零化度 $\mathrm{nullity}(A)={{expr:nullity(A)}}=4$，即需要求 4 个基础解系向量。

**步骤 2**：对 $A$ 做初等行变换化为行最简形（RREF），确定主元列和自由变量。4×5 秩 1 矩阵有 1 个主元列、4 个自由变量。

**步骤 3**：依次令每个自由变量为 1、其余为 0，回代求出主元变量，得到四个基础解系向量 $\xi_1,\xi_2,\xi_3,\xi_4$。

**步骤 4**：四个基础解系向量的分量依次为：

向量 $\xi_1$（第 1–5 空）：
- $\xi_{1,1}={{expr:nullbasis_comp(A,1,1)}}$
- $\xi_{1,2}={{expr:nullbasis_comp(A,1,2)}}$
- $\xi_{1,3}={{expr:nullbasis_comp(A,1,3)}}$
- $\xi_{1,4}={{expr:nullbasis_comp(A,1,4)}}$
- $\xi_{1,5}={{expr:nullbasis_comp(A,1,5)}}$

向量 $\xi_2$（第 6–10 空）：
- $\xi_{2,1}={{expr:nullbasis_comp(A,2,1)}}$
- $\xi_{2,2}={{expr:nullbasis_comp(A,2,2)}}$
- $\xi_{2,3}={{expr:nullbasis_comp(A,2,3)}}$
- $\xi_{2,4}={{expr:nullbasis_comp(A,2,4)}}$
- $\xi_{2,5}={{expr:nullbasis_comp(A,2,5)}}$

向量 $\xi_3$（第 11–15 空）：
- $\xi_{3,1}={{expr:nullbasis_comp(A,3,1)}}$
- $\xi_{3,2}={{expr:nullbasis_comp(A,3,2)}}$
- $\xi_{3,3}={{expr:nullbasis_comp(A,3,3)}}$
- $\xi_{3,4}={{expr:nullbasis_comp(A,3,4)}}$
- $\xi_{3,5}={{expr:nullbasis_comp(A,3,5)}}$

向量 $\xi_4$（第 16–20 空）：
- $\xi_{4,1}={{expr:nullbasis_comp(A,4,1)}}$
- $\xi_{4,2}={{expr:nullbasis_comp(A,4,2)}}$
- $\xi_{4,3}={{expr:nullbasis_comp(A,4,3)}}$
- $\xi_{4,4}={{expr:nullbasis_comp(A,4,4)}}$
- $\xi_{4,5}={{expr:nullbasis_comp(A,4,5)}}$`,
		},
	}
}

func buildChapter4_1() dsl.Problem {
	k := "Chapter4_1"
	ids := BlankIDs(k, 5)
	jb := &dsl.AnswerJudgeSpec{Kind: "sorted_basis_columns", BasisGroup: "bc41", MatrixVar: "M", Ncols: 4}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $M={{M}}$ 的四列为 $\alpha_1,\ldots,\alpha_4$，填写 $\dim L(\alpha_1,\ldots,\alpha_4)$ 及一个极大线性无关列组下标（1–4，不足填 0）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"M": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
		},
		Derived: map[string]string{"r": "rank(M)", "piv": "basis_cols(M)"},
		Render:  map[string]string{"M": "M"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "r"},
			{ID: ids[1], Expr: "piv[1]", Judge: jb},
			{ID: ids[2], Expr: "piv[2]", Judge: jb},
			{ID: ids[3], Expr: "piv[3]", Judge: jb},
			{ID: ids[4], Expr: "piv[4]", Judge: jb},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：向量组 $\alpha_1,\ldots,\alpha_4$ 张成的子空间维数 $\dim L(\alpha_1,\ldots,\alpha_4)$ 等于矩阵 $M$ 的列秩，即 $\mathrm{rank}(M)$。对 $M$ 做初等行变换化为行阶梯形，主元列（pivot columns）的下标即为一个极大线性无关列组的下标。

**步骤 1**：矩阵 $M={{M}}$，对 $M$ 做初等行变换化为行阶梯形，非零行个数即为 $\mathrm{rank}(M)={{expr:rank(M)}}$，故 $\dim L(\alpha_1,\ldots,\alpha_4)={{expr:rank(M)}}$。

**步骤 2**：行阶梯形中主元所在的列下标即为极大线性无关列组下标。由 $\mathrm{basis\_cols}(M)={{expr:basis_cols(M)}}$ 得到主元列下标，从小到大排列依次为：
- 第 1 个下标：{{expr:piv[1]}}
- 第 2 个下标：{{expr:piv[2]}}
- 第 3 个下标：{{expr:piv[3]}}
- 第 4 个下标：{{expr:piv[4]}}（若不足则补 0）

**步骤 3**：依次填写 1 个维数值和 4 个下标值即可。`,
		},
	}
}

func buildChapter4_4() dsl.Problem {
	k := "Chapter4_4"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $A={{A}}$ 为三阶方阵，齐次方程组 $Ax=0$ 有非零解，求参数 $\lambda={{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"A": {
				Kind: "matrix", Rows: 3, Cols: 3,
				Generator: map[string]interface{}{
					"rule": "lambda_linear_det_zero", "param_var": "lambda",
					"param_row": 3, "param_col": 3,
					"entry_min": -4, "entry_max": 4, "lambda_min": -10, "lambda_max": 10,
					"max_attempts": 120,
				},
			},
		},
		Render: map[string]string{"A": "A"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "lambda"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：齐次方程组 $Ax=0$ 有非零解的充要条件是系数行列式为零，即 $\det(A)=0$。矩阵 $A$ 中含参数 $\lambda$（位于第 $(3,3)$ 位置），解方程 $\det(A)=0$ 即可求出 $\lambda$。

**步骤 1**：系数矩阵 $A={{A}}$，其中 $\lambda$ 为待求参数。

**步骤 2**：计算行列式 $\det(A)={{expr:det(A)}}$。由于 $\lambda$ 只出现在一个位置，行列式关于 $\lambda$ 是线性的，展开后得到关于 $\lambda$ 的一次方程。

**步骤 3**：令 $\det(A)=0$，解出 $\lambda={{expr:lambda}}$。`,
		},
	}
}

func buildChapter4_2() dsl.Problem {
	k := "Chapter4_2"
	ids := BlankIDs(k, 3)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设基矩阵 $B={{B}}$，$\beta={{beta}}$，求 $\beta$ 在基 $B$ 列下的坐标：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"B":    {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "upper_unit", "min": -4, "max": 4}},
			"beta": {Kind: "vector", Size: 3, Generator: map[string]interface{}{"rule": "range", "min": -5, "max": 5}},
		},
		Derived: map[string]string{"x": "solve(B,beta)"},
		Render:  map[string]string{"B": "B", "beta": "beta"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "x[1]"}, {ID: ids[1], Expr: "x[2]"}, {ID: ids[2], Expr: "x[3]"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：向量 $\beta$ 在基 $B$ 各列下的坐标 $x=(x_1,x_2,x_3)^T$ 满足线性方程组 $Bx=\beta$，即求解 $x=B^{-1}\beta$。由于 $B$ 是上三角单位矩阵（对角元全为 1 的上三角阵），可用回代法高效求解。

**步骤 1**：基矩阵 $B={{B}}$，目标向量 $\beta={{beta}}$。设坐标为 $x=(x_1,x_2,x_3)^T$，则方程为 $Bx=\beta$。

**步骤 2**：由于 $B$ 为上三角矩阵，从最后一行开始回代：
- 第 3 行：$x_3 = \beta_3 = {{expr:beta[3]}}$
- 第 2 行：$b_{22}x_2 + b_{23}x_3 = \beta_2$，代入 $x_3$ 解出 $x_2$
- 第 1 行：$b_{11}x_1 + b_{12}x_2 + b_{13}x_3 = \beta_1$，代入 $x_2,x_3$ 解出 $x_1$

**步骤 3**：解得坐标向量 $x=(x_1,x_2,x_3)^T$，三个分量依次为：
- $x_1 = {{expr:x[1]}}$
- $x_2 = {{expr:x[2]}}$
- $x_3 = {{expr:x[3]}}$`,
		},
	}
}

func buildChapter4_5_1() dsl.Problem { return nonHom12("Chapter4_5_1") }
func buildChapter4_5_2() dsl.Problem { return nonHom12("Chapter4_5_2") }

func nonHom12(key string) dsl.Problem {
	ids := BlankIDs(key, 12)
	ja := &dsl.AnswerJudgeSpec{Kind: "affine_rational", AffineGroup: "eta", MatrixVar: "A", BVecVar: "b2"}
	j1 := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: "z1"}
	j2 := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: "z2"}
	fds := []dsl.AnswerFieldDef{
		{ID: ids[0], Expr: "x0[1]", Judge: ja}, {ID: ids[1], Expr: "x0[2]", Judge: ja},
		{ID: ids[2], Expr: "x0[3]", Judge: ja}, {ID: ids[3], Expr: "x0[4]", Judge: ja},
		{ID: ids[4], Expr: "nullbasis_comp(A,1,1)", Judge: j1},
		{ID: ids[5], Expr: "nullbasis_comp(A,1,2)", Judge: j1},
		{ID: ids[6], Expr: "nullbasis_comp(A,1,3)", Judge: j1},
		{ID: ids[7], Expr: "nullbasis_comp(A,1,4)", Judge: j1},
		{ID: ids[8], Expr: "nullbasis_comp(A,2,1)", Judge: j2},
		{ID: ids[9], Expr: "nullbasis_comp(A,2,2)", Judge: j2},
		{ID: ids[10], Expr: "nullbasis_comp(A,2,3)", Judge: j2},
		{ID: ids[11], Expr: "nullbasis_comp(A,2,4)", Judge: j2},
	}
	return dsl.Problem{
		ID:      ProblemID(key),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $A={{A}}$，$b={{b2}}$，填写 $Ax=b$ 的一特解 $\eta$ 与齐次基础解系 $\xi_1,\xi_2$（各分量；$\xi_i$ 允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A":  {Kind: "matrix", Rows: 2, Cols: 4, Generator: map[string]interface{}{"rule": "rank2_2x4", "min": -4, "max": 4}},
			"x0": {Kind: "vector", Size: 4, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
		},
		Derived: map[string]string{"b2": "A * x0"},
		Render:  map[string]string{"A": "A", "b2": "b2"},
		Answer:  dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：非齐次方程组 $Ax=b$ 的通解 = 一个特解 $\eta$ + 齐次方程 $Ax=0$ 的通解。$A$ 为 $2\times 4$ 秩 2 矩阵，基础解系含 $4-2=2$ 个线性无关的向量。

**步骤 1**：对增广矩阵 $(A|b)$ 做初等行变换化为行最简形（RREF），其中 $A={{A}}$，$b={{b2}}$。

**步骤 2**：令自由变量为 0，从 RREF 读出特解 $\eta$，4 个分量依次为：
- $\eta_1 = {{expr:x0[1]}}$
- $\eta_2 = {{expr:x0[2]}}$
- $\eta_3 = {{expr:x0[3]}}$
- $\eta_4 = {{expr:x0[4]}}$

**步骤 3**：对系数矩阵 $A$ 做行变换化为行最简形，令第一个自由变量为 1、第二个为 0，求得齐次基础解系第一个向量 $\xi_1$，分量依次为：
- $\xi_{1,1} = {{expr:nullbasis_comp(A,1,1)}}$
- $\xi_{1,2} = {{expr:nullbasis_comp(A,1,2)}}$
- $\xi_{1,3} = {{expr:nullbasis_comp(A,1,3)}}$
- $\xi_{1,4} = {{expr:nullbasis_comp(A,1,4)}}$

**步骤 4**：令第一个自由变量为 0、第二个为 1，求得齐次基础解系第二个向量 $\xi_2$，分量依次为：
- $\xi_{2,1} = {{expr:nullbasis_comp(A,2,1)}}$
- $\xi_{2,2} = {{expr:nullbasis_comp(A,2,2)}}$
- $\xi_{2,3} = {{expr:nullbasis_comp(A,2,3)}}$
- $\xi_{2,4} = {{expr:nullbasis_comp(A,2,4)}}$

**步骤 5**：通解为 $x=\eta+k_1\xi_1+k_2\xi_2$。依次填写 $\eta$ 的 4 个分量、$\xi_1$ 的 4 个分量、$\xi_2$ 的 4 个分量，共 12 个空。`,
		},
	}
}

func buildChapter4_7() dsl.Problem {
	k := "Chapter4_7"
	ids := BlankIDs(k, 8)
	ja := &dsl.AnswerJudgeSpec{Kind: "affine_rational", AffineGroup: "xp", MatrixVar: "A4", BVecVar: "b4"}
	jn := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: "ker"}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $A={{A4}}$，$b={{b4}}$，$\mathrm{rank}(A)=3$。填写 $Ax=b$ 的一特解与齐次方程基础解系中一个向量（各分量；基础解系向量允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A4": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "rank334_last_dep", "min": -4, "max": 4}},
			"xp": {Kind: "vector", Size: 4, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
		},
		Derived: map[string]string{"b4": "A4 * xp"},
		Render:  map[string]string{"A4": "A4", "b4": "b4"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "xp[1]", Judge: ja}, {ID: ids[1], Expr: "xp[2]", Judge: ja},
			{ID: ids[2], Expr: "xp[3]", Judge: ja}, {ID: ids[3], Expr: "xp[4]", Judge: ja},
			{ID: ids[4], Expr: "nullbasis_comp(A4,1,1)", Judge: jn},
			{ID: ids[5], Expr: "nullbasis_comp(A4,1,2)", Judge: jn},
			{ID: ids[6], Expr: "nullbasis_comp(A4,1,3)", Judge: jn},
			{ID: ids[7], Expr: "nullbasis_comp(A4,1,4)", Judge: jn},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：$A$ 为 $4\times 4$ 矩阵且 $\mathrm{rank}(A)=3$，由秩-零化度定理，基础解系含 $4-3=1$ 个向量。非齐次通解 = 一个特解 + 齐次通解。

**步骤 1**：对增广矩阵 $(A|b)=({{A4}}|{{b4}})$ 做初等行变换化为行最简形。

**步骤 2**：令自由变量为 0，从 RREF 读出特解 $\eta$，其 4 个分量依次为：
- $\eta_1 = {{expr:xp[1]}}$
- $\eta_2 = {{expr:xp[2]}}$
- $\eta_3 = {{expr:xp[3]}}$
- $\eta_4 = {{expr:xp[4]}}$
即为前 4 个填空的答案。

**步骤 3**：对系数矩阵 $A$ 做行变换化为行最简形，令唯一自由变量为 1，回代求出基础解系向量 $\xi$（4 维），分量依次为：
- $\xi_1 = {{expr:nullbasis_comp(A4,1,1)}}$
- $\xi_2 = {{expr:nullbasis_comp(A4,1,2)}}$
- $\xi_3 = {{expr:nullbasis_comp(A4,1,3)}}$
- $\xi_4 = {{expr:nullbasis_comp(A4,1,4)}}$

**步骤 4**：通解为 $x=\eta+k\xi$。依次填写 $\eta$ 和 $\xi$ 的各分量，共 8 个空。`,
		},
	}
}

func buildChapter4_6() dsl.Problem {
	k := "Chapter4_6"
	ids := BlankIDs(k, 22)
	ja := &dsl.AnswerJudgeSpec{Kind: "affine_rational", AffineGroup: "eta46", MatrixVar: "A", BVecVar: "b3"}
	jn := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: "xi146"}
	fds := []dsl.AnswerFieldDef{
		// λ (1 blank)
		{ID: ids[0], Expr: "param_lambda_val(A)"},
		// μ (1 blank)
		{ID: ids[1], Expr: "param_mu_val(A)"},
		// rank(A) (1 blank)
		{ID: ids[2], Expr: "param_rank_A(A)"},
		// rank(A|b) (1 blank)
		{ID: ids[3], Expr: "param_rank_aug(A)"},
		// RREF 3×4 (12 blanks)
		{ID: ids[4], Expr: "param_rref_comp(A,1,1)"},
		{ID: ids[5], Expr: "param_rref_comp(A,1,2)"},
		{ID: ids[6], Expr: "param_rref_comp(A,1,3)"},
		{ID: ids[7], Expr: "param_rref_comp(A,1,4)"},
		{ID: ids[8], Expr: "param_rref_comp(A,2,1)"},
		{ID: ids[9], Expr: "param_rref_comp(A,2,2)"},
		{ID: ids[10], Expr: "param_rref_comp(A,2,3)"},
		{ID: ids[11], Expr: "param_rref_comp(A,2,4)"},
		{ID: ids[12], Expr: "param_rref_comp(A,3,1)"},
		{ID: ids[13], Expr: "param_rref_comp(A,3,2)"},
		{ID: ids[14], Expr: "param_rref_comp(A,3,3)"},
		{ID: ids[15], Expr: "param_rref_comp(A,3,4)"},
		// 特解 η (3 blanks) + 基础解系 ξ (3 blanks)
		{ID: ids[16], Expr: "param_x0_comp(A,1)", Judge: ja},
		{ID: ids[17], Expr: "param_x0_comp(A,2)", Judge: ja},
		{ID: ids[18], Expr: "param_x0_comp(A,3)", Judge: ja},
		{ID: ids[19], Expr: "param_nb_comp(A,1)", Judge: jn},
		{ID: ids[20], Expr: "param_nb_comp(A,2)", Judge: jn},
		{ID: ids[21], Expr: "param_nb_comp(A,3)", Judge: jn},
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`已知线性方程组 $${{system_title}}$$，则 $\lambda$={{blank:%s}}、$\mu$={{blank:%s}} 时方程组有无穷多解；此时系数矩阵的秩为 {{blank:%s}}，增广矩阵的秩为 {{blank:%s}}；增广阵行最简形（按行）{{blank:%s}} {{blank:%s}} {{blank:%s}} {{blank:%s}} {{blank:%s}} {{blank:%s}} {{blank:%s}} {{blank:%s}} {{blank:%s}} {{blank:%s}} {{blank:%s}} {{blank:%s}}；通解为 $x={{blank:%s}} {{blank:%s}} {{blank:%s}} + k\,{{blank:%s}} {{blank:%s}} {{blank:%s}}$`,
			ids[0], ids[1], ids[2], ids[3],
			ids[4], ids[5], ids[6], ids[7], ids[8], ids[9], ids[10], ids[11], ids[12], ids[13], ids[14], ids[15],
			ids[16], ids[17], ids[18], ids[19], ids[20], ids[21]),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{
				"rule": "param_infinit_solution",
				"entry_min": -5, "entry_max": 5,
				"lambda_min": -10, "lambda_max": 10,
			}},
		},
		Derived: map[string]string{
			"system_title": "param_system_title(A)",
			"b3":           "_param_b",
		},
		Render: map[string]string{"system_title": "system_title"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：线性方程组 $Ax=b$ 有无穷多解的充要条件是 $\mathrm{rank}(A)=\mathrm{rank}(A|b)<n$（未知数个数）。先取参数 $\lambda$ 使 $\det(A)=0$（系数矩阵降秩），再确定 $\mu$ 使增广矩阵秩不变，然后对增广矩阵做行变换求通解。

**步骤 1**：由方程组 $Ax=b$（具体形式见题面），计算 $\det(A)$ 并令其为 0，求得使系数矩阵降秩的参数值 $\lambda={{expr:param_lambda_val(A)}}$。

**步骤 2**：将 $\lambda={{expr:param_lambda_val(A)}}$ 代入方程组，确定常数项中参数 $\mu={{expr:param_mu_val(A)}}$ 使得 $\mathrm{rank}(A|b)=\mathrm{rank}(A)$，此时方程组有无穷多解。

**步骤 3**：此时系数矩阵的秩 $\mathrm{rank}(A)={{expr:param_rank_A(A)}}$，增广矩阵的秩 $\mathrm{rank}(A|b)={{expr:param_rank_aug(A)}}$，二者相等且小于未知数个数 3，确认有无穷多解。

**步骤 4**：对增广矩阵 $(A|b)$（3 行 × 4 列，含常数项列）做初等行变换化为行最简形（RREF），得到 12 个元素，按行依次填写：

第一行（第 1–4 个 RREF 元素）：
- $r_{11} = {{expr:param_rref_comp(A,1,1)}}$
- $r_{12} = {{expr:param_rref_comp(A,1,2)}}$
- $r_{13} = {{expr:param_rref_comp(A,1,3)}}$
- $r_{14} = {{expr:param_rref_comp(A,1,4)}}$

第二行（第 5–8 个 RREF 元素）：
- $r_{21} = {{expr:param_rref_comp(A,2,1)}}$
- $r_{22} = {{expr:param_rref_comp(A,2,2)}}$
- $r_{23} = {{expr:param_rref_comp(A,2,3)}}$
- $r_{24} = {{expr:param_rref_comp(A,2,4)}}$

第三行（第 9–12 个 RREF 元素）：
- $r_{31} = {{expr:param_rref_comp(A,3,1)}}$
- $r_{32} = {{expr:param_rref_comp(A,3,2)}}$
- $r_{33} = {{expr:param_rref_comp(A,3,3)}}$
- $r_{34} = {{expr:param_rref_comp(A,3,4)}}$

**步骤 5**：从行最简形读出特解 $\eta$（3 维）和齐次基础解系向量 $\xi$（3 维）：

特解 $\eta$（第 17–19 空）：
- $\eta_1 = {{expr:param_x0_comp(A,1)}}$
- $\eta_2 = {{expr:param_x0_comp(A,2)}}$
- $\eta_3 = {{expr:param_x0_comp(A,3)}}$

基础解系向量 $\xi$（第 20–22 空）：
- $\xi_1 = {{expr:param_nb_comp(A,1)}}$
- $\xi_2 = {{expr:param_nb_comp(A,2)}}$
- $\xi_3 = {{expr:param_nb_comp(A,3)}}$

**步骤 6**：通解为 $x=\eta+k\xi$，即：
$$x = \begin{pmatrix} {{expr:param_x0_comp(A,1)}} \\ {{expr:param_x0_comp(A,2)}} \\ {{expr:param_x0_comp(A,3)}} \end{pmatrix} + k \begin{pmatrix} {{expr:param_nb_comp(A,1)}} \\ {{expr:param_nb_comp(A,2)}} \\ {{expr:param_nb_comp(A,3)}} \end{pmatrix}$$`,
		},
	}
}
