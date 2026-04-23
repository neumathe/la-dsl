package bank

import (
	"fmt"

	"github.com/neumathe/la-dsl/dsl"
)

// ──────────────────────────────────────────────────────────────
// Chapter 7 — 线性空间与线性变换（12 题）
// ──────────────────────────────────────────────────────────────

func buildChapter7_7() dsl.Problem {
	k := "Chapter7_7"
	ids := BlankIDs(k, 4)
	jw := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: "w"}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $\alpha_1,\alpha_2,\alpha_3$ 为矩阵 $V={{V}}$ 的三列，求与它们都正交的一个非零向量（填 4 个分量）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"V": {Kind: "matrix", Rows: 4, Cols: 3, Generator: map[string]interface{}{"rule": "full_rank", "min": -5, "max": 5}},
		},
		Derived: map[string]string{"MT": "transpose(V)", "w": "nullvec(MT)"},
		Render:  map[string]string{"V": "V"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "w[1]", Judge: jw},
			{ID: ids[1], Expr: "w[2]", Judge: jw},
			{ID: ids[2], Expr: "w[3]", Judge: jw},
			{ID: ids[3], Expr: "w[4]", Judge: jw},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一步：写出矩阵 V 的三列向量。V = {{V}}，即
  α₁ = ({{expr:mget(V,1,1)}}, {{expr:mget(V,2,1)}}, {{expr:mget(V,3,1)}}, {{expr:mget(V,4,1)}})^T
  α₂ = ({{expr:mget(V,1,2)}}, {{expr:mget(V,2,2)}}, {{expr:mget(V,3,2)}}, {{expr:mget(V,4,2)}})^T
  α₃ = ({{expr:mget(V,1,3)}}, {{expr:mget(V,2,3)}}, {{expr:mget(V,3,3)}}, {{expr:mget(V,4,3)}})^T

第二步：求与 α₁, α₂, α₃ 都正交的非零向量 w ∈ R⁴，即满足 αᵢ · w = 0（i=1,2,3）。
将三个正交条件写成齐次线性方程组 V^T w = 0，其中
  V^T = {{expr:transpose(V)}}（3×4 矩阵）

第三步：由于 V 列满秩（rank(V)=3），故 rank(V^T)=3。V^T 为 3×4 矩阵，零空间维数 = 4 − 3 = 1，即解空间为一维，存在一个自由变量。

第四步：对 V^T 作行初等变换化为行最简形（RREF），确定主元列和自由变量位置。令自由变量为 1，回代求得特解。

第五步：求得 w = ({{expr:w}})^T，即
  w₁ = {{expr:w[1]}}，w₂ = {{expr:w[2]}}，w₃ = {{expr:w[3]}}，w₄ = {{expr:w[4]}}

第六步：验证正交性：
  α₁ · w = {{expr:mget(V,1,1)}}·w₁ + {{expr:mget(V,2,1)}}·w₂ + {{expr:mget(V,3,1)}}·w₃ + {{expr:mget(V,4,1)}}·w₄ = 0
  α₂ · w = {{expr:mget(V,1,2)}}·w₁ + {{expr:mget(V,2,2)}}·w₂ + {{expr:mget(V,3,2)}}·w₃ + {{expr:mget(V,4,2)}}·w₄ = 0
  α₃ · w = {{expr:mget(V,1,3)}}·w₁ + {{expr:mget(V,2,3)}}·w₂ + {{expr:mget(V,3,3)}}·w₃ + {{expr:mget(V,4,3)}}·w₄ = 0

答案：w = ({{expr:w[1]}}, {{expr:w[2]}}, {{expr:w[3]}}, {{expr:w[4]}})^T（允许整体乘以任意非零常数）`,
		},
	}
}

func buildChapter7_4() dsl.Problem {
	k := "Chapter7_4"
	ids := BlankIDs(k, 3)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设基矩阵 $B={{B}}$，向量 $\alpha={{a}}$，求 $\alpha$ 在基 $B$ 列下的坐标：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"B": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "upper_unit", "min": -4, "max": 4}},
			"a": {Kind: "vector", Size: 3, Generator: map[string]interface{}{"rule": "range", "min": -6, "max": 6}},
		},
		Derived: map[string]string{"x": "solve(B,a)"},
		Render:  map[string]string{"B": "B", "a": "a"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "x[1]"}, {ID: ids[1], Expr: "x[2]"}, {ID: ids[2], Expr: "x[3]"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一步：理解题意。基矩阵 B 的三个列向量构成 R³ 的一组基。求 α 在这组基下的坐标 x = (x₁, x₂, x₃)^T，即满足 Bx = α。

第二步：写出基矩阵 B 和各列向量：
  B = {{expr:B}}，即
  β₁ = ({{expr:mget(B,1,1)}}, {{expr:mget(B,2,1)}}, {{expr:mget(B,3,1)}})^T
  β₂ = ({{expr:mget(B,1,2)}}, {{expr:mget(B,2,2)}}, {{expr:mget(B,3,2)}})^T
  β₃ = ({{expr:mget(B,1,3)}}, {{expr:mget(B,2,3)}}, {{expr:mget(B,3,3)}})^T
  α = ({{expr:a}})

第三步：建立方程组 Bx = α，展开为：
  {{expr:mget(B,1,1)}}·x₁ + {{expr:mget(B,1,2)}}·x₂ + {{expr:mget(B,1,3)}}·x₃ = {{expr:a[1]}}
  {{expr:mget(B,2,1)}}·x₁ + {{expr:mget(B,2,2)}}·x₂ + {{expr:mget(B,2,3)}}·x₃ = {{expr:a[2]}}
  {{expr:mget(B,3,1)}}·x₁ + {{expr:mget(B,3,2)}}·x₂ + {{expr:mget(B,3,3)}}·x₃ = {{expr:a[3]}}

第四步：由于 B 为上三角单位对角矩阵（对角元全为 1），可用回代法求解。
从最后一行开始：
  第 3 行：x₃ = {{expr:a[3]}}（因为 B₃₃ = 1 且 B₃₁ = B₃₂ = 0）
  第 2 行：x₂ + {{expr:mget(B,2,3)}}·x₃ = {{expr:a[2]}} → x₂ = {{expr:a[2]}} − {{expr:mget(B,2,3)}}·x₃
  第 1 行：x₁ + {{expr:mget(B,1,2)}}·x₂ + {{expr:mget(B,1,3)}}·x₃ = {{expr:a[1]}} → x₁ = {{expr:a[1]}} − {{expr:mget(B,1,2)}}·x₂ − {{expr:mget(B,1,3)}}·x₃

第五步：代入计算得坐标 x = ({{expr:solve(B,a)}})^T，即
  x₁ = {{expr:x[1]}}，x₂ = {{expr:x[2]}}，x₃ = {{expr:x[3]}}

第六步：验证 Bx = α：将 x 各分量代入左边，结果应等于 α。

答案：x = ({{expr:x[1]}}, {{expr:x[2]}}, {{expr:x[3]}})^T`,
		},
	}
}

func buildChapter7_10() dsl.Problem {
	k := "Chapter7_10"
	ids := BlankIDs(k, 12)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`在 $\mathbb{R}[x]_3$ 中，由基 $1,x,x^2$ 到基 $\beta_1={{b1}},\beta_2={{b2}},\beta_3={{b3}}$ 的过渡矩阵 $P$ 及多项式 $p(x)={{px}}$ 在基 $\beta$ 下坐标（过渡矩阵第一行：{{blank:%s}} {{blank:%s}} {{blank:%s}}；第二行：{{blank:%s}} {{blank:%s}} {{blank:%s}}；第三行：{{blank:%s}} {{blank:%s}} {{blank:%s}}；坐标：{{blank:%s}} {{blank:%s}} {{blank:%s}}）`,
			ids[0], ids[1], ids[2], ids[3], ids[4], ids[5], ids[6], ids[7], ids[8], ids[9], ids[10], ids[11]),
		Variables: map[string]dsl.Variable{
			"B":  {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "upper_triangular_nonzero_diag", "min": -4, "max": 4}},
			"p":  {Kind: "vector", Size: 3, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
		},
		Derived: map[string]string{
			"coord": "solve(B,p)",
			"b1":    "poly_from_matcol(B,1)",
			"b2":    "poly_from_matcol(B,2)",
			"b3":    "poly_from_matcol(B,3)",
			"px":    "poly_from_vec(p)",
		},
		Render: map[string]string{"px": "px", "b1": "b1", "b2": "b2", "b3": "b3"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "mget(B,1,1)"}, {ID: ids[1], Expr: "mget(B,1,2)"}, {ID: ids[2], Expr: "mget(B,1,3)"},
			{ID: ids[3], Expr: "mget(B,2,1)"}, {ID: ids[4], Expr: "mget(B,2,2)"}, {ID: ids[5], Expr: "mget(B,2,3)"},
			{ID: ids[6], Expr: "mget(B,3,1)"}, {ID: ids[7], Expr: "mget(B,3,2)"}, {ID: ids[8], Expr: "mget(B,3,3)"},
			{ID: ids[9], Expr: "coord[1]"}, {ID: ids[10], Expr: "coord[2]"}, {ID: ids[11], Expr: "coord[3]"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一部分：求过渡矩阵 P

第一步：理解过渡矩阵的定义。从旧基 (1, x, x²) 到新基 (β₁, β₂, β₃) 的过渡矩阵 P 满足 (β₁, β₂, β₃) = (1, x, x²) · P。即 P 的第 j 列就是 βⱼ 在旧基 (1, x, x²) 下的坐标向量。

第二步：将 β₁, β₂, β₃ 写成旧基下的坐标列：
  β₁ = {{expr:poly_from_matcol(B,1)}}，坐标为 B 的第 1 列 = ({{expr:mget(B,1,1)}}, {{expr:mget(B,2,1)}}, {{expr:mget(B,3,1)}})^T
  β₂ = {{expr:poly_from_matcol(B,2)}}，坐标为 B 的第 2 列 = ({{expr:mget(B,1,2)}}, {{expr:mget(B,2,2)}}, {{expr:mget(B,3,2)}})^T
  β₃ = {{expr:poly_from_matcol(B,3)}}，坐标为 B 的第 3 列 = ({{expr:mget(B,1,3)}}, {{expr:mget(B,2,3)}}, {{expr:mget(B,3,3)}})^T

第三步：因此过渡矩阵 P = B = {{expr:B}}：
  P = | {{expr:mget(B,1,1)}}  {{expr:mget(B,1,2)}}  {{expr:mget(B,1,3)}} |
      | {{expr:mget(B,2,1)}}  {{expr:mget(B,2,2)}}  {{expr:mget(B,2,3)}} |
      | {{expr:mget(B,3,1)}}  {{expr:mget(B,3,2)}}  {{expr:mget(B,3,3)}} |

第二部分：求 p(x) 在基 β 下的坐标

第四步：设 p(x) 在基 β 下的坐标为 y = (y₁, y₂, y₃)^T，则 p(x) = y₁β₁ + y₂β₂ + y₃β₃。在旧基坐标下即 p = B · y，所以 y = B⁻¹p = solve(B, p)。

第五步：已知 p(x) = {{expr:poly_from_vec(p)}}，对应旧基坐标向量 p = ({{expr:p[1]}}, {{expr:p[2]}}, {{expr:p[3]}})^T。解方程组 By = p。

由于 B 为上三角矩阵且对角元非零（{{expr:mget(B,1,1)}}, {{expr:mget(B,2,2)}}, {{expr:mget(B,3,3)}} 均 ≠ 0），可用回代法求解：
  从第 3 行解出 y₃ = ({{expr:p[3]}}) / {{expr:mget(B,3,3)}}
  代入第 2 行解出 y₂
  代入第 1 行解出 y₁

第六步：解得坐标 y = ({{expr:solve(B,p)}})^T，即
  y₁ = {{expr:coord[1]}}，y₂ = {{expr:coord[2]}}，y₃ = {{expr:coord[3]}}

验证：计算 B · y 应等于 p，即 {{expr:mget(B,1,1)}}·y₁ + {{expr:mget(B,1,2)}}·y₂ + {{expr:mget(B,1,3)}}·y₃ = {{expr:p[1]}}，其余两行同理。

答案：过渡矩阵 P = B（9 个元素如上），p(x) 在基 β 下坐标 = ({{expr:coord[1]}}, {{expr:coord[2]}}, {{expr:coord[3]}})^T`,
		},
	}
}

func buildChapter7_9() dsl.Problem {
	k := "Chapter7_9"
	ids := BlankIDs(k, 9)
	fds := make([]dsl.AnswerFieldDef, 0, 9)
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[(i-1)*3+j-1],
				Expr:   fmt.Sprintf("mget(Tnew,%d,%d)", i, j),
				Layout: dsl.LayoutMatrixCell("P^{-1}AP", i, j, 3, 3, "P^{-1}AP"),
			})
		}
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设线性变换 $T$ 在基 $\varepsilon_1,\varepsilon_2,\varepsilon_3$ 下的矩阵为 $A={{A}}$，求 $T$ 在基 ${{basis_text}}$ 下的矩阵：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
			"P": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "upper_unit", "min": -4, "max": 4}},
		},
		Derived: map[string]string{
			"basis_text": "basis_linear_combo_title(P)",
			"0Pi":  "inv(P)",
			"0PiA": "matmul(0Pi, A)",
			"Tnew": "matmul(0PiA, P)",
		},
		Render: map[string]string{"A": "A", "basis_text": "basis_text"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一步：设旧基为 ε₁, ε₂, ε₃，新基为 η₁, η₂, η₃。它们的关系为 (η₁, η₂, η₃) = (ε₁, ε₂, ε₃) · P，其中过渡矩阵
  P = {{expr:P}} = | {{expr:mget(P,1,1)}}  {{expr:mget(P,1,2)}}  {{expr:mget(P,1,3)}} |
                   | {{expr:mget(P,2,1)}}  {{expr:mget(P,2,2)}}  {{expr:mget(P,2,3)}} |
                   | {{expr:mget(P,3,1)}}  {{expr:mget(P,3,2)}}  {{expr:mget(P,3,3)}} |

第二步：线性变换 T 在旧基下的矩阵为 A = {{expr:A}}。T 在新基下的矩阵 T_new 与 A 的关系为：T_new = P⁻¹AP。

第三步：求 P⁻¹。P 为上三角单位对角矩阵（对角元全为 1），其逆矩阵也是上三角单位对角矩阵。
  P⁻¹ = {{expr:inv(P)}} = | {{expr:mget(0Pi,1,1)}}  {{expr:mget(0Pi,1,2)}}  {{expr:mget(0Pi,1,3)}} |
                           | {{expr:mget(0Pi,2,1)}}  {{expr:mget(0Pi,2,2)}}  {{expr:mget(0Pi,2,3)}} |
                           | {{expr:mget(0Pi,3,1)}}  {{expr:mget(0Pi,3,2)}}  {{expr:mget(0Pi,3,3)}} |

第四步：计算 P⁻¹A = {{expr:matmul(0Pi,A)}} =
  | {{expr:mget(0PiA,1,1)}}  {{expr:mget(0PiA,1,2)}}  {{expr:mget(0PiA,1,3)}} |
  | {{expr:mget(0PiA,2,1)}}  {{expr:mget(0PiA,2,2)}}  {{expr:mget(0PiA,2,3)}} |
  | {{expr:mget(0PiA,3,1)}}  {{expr:mget(0PiA,3,2)}}  {{expr:mget(0PiA,3,3)}} |

第五步：计算 T_new = (P⁻¹A) · P = {{expr:matmul(0PiA,P)}} =
  | {{expr:mget(Tnew,1,1)}}  {{expr:mget(Tnew,1,2)}}  {{expr:mget(Tnew,1,3)}} |
  | {{expr:mget(Tnew,2,1)}}  {{expr:mget(Tnew,2,2)}}  {{expr:mget(Tnew,2,3)}} |
  | {{expr:mget(Tnew,3,1)}}  {{expr:mget(Tnew,3,2)}}  {{expr:mget(Tnew,3,3)}} |

第六步：验证相似不变量：
  tr(A) = {{expr:mget(A,1,1)}} + {{expr:mget(A,2,2)}} + {{expr:mget(A,3,3)}}
  tr(T_new) = {{expr:mget(Tnew,1,1)}} + {{expr:mget(Tnew,2,2)}} + {{expr:mget(Tnew,3,3)}}，应与 tr(A) 相等。

答案：T 在新基下的矩阵 T_new = {{expr:matmul(0PiA,P)}}（9 个元素依次填入）`,
		},
	}
}

func buildChapter7_1() dsl.Problem {
	k := "Chapter7_1"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $V=\left\{\begin{bmatrix}0&-x\\-2y&-2x-y+1\end{bmatrix}\mid x,y\in\mathbb{R}\right\}$，判断 $V$ 是否构成 $\mathbb{R}$ 上的线性空间（1=是，0=否）：{{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"s": {Kind: "scalar", Fixed: 0},
		},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "s"}}},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一步：回顾线性空间的定义。一个集合 V 构成 R 上的线性空间，必须满足若干公理，其中最基本的要求是：V 必须包含零向量（零元）。对于 2×2 矩阵空间，零向量即零矩阵 O = [[0,0],[0,0]]。

第二步：检查零矩阵是否在 V 中。V 中矩阵的一般形式为 [[0, -x], [-2y, -2x-y+1]]。若此矩阵等于零矩阵，则各位置元素须满足：
  左上角：0 = 0（恒成立）
  右上角：-x = 0 → x = 0
  左下角：-2y = 0 → y = 0
  右下角：-2x - y + 1 = 0

第三步：将 x = 0, y = 0 代入右下角：-2(0) - 0 + 1 = 1 ≠ 0。矛盾！

第四步：因此，不存在任何 x, y ∈ R 使得 V 中的矩阵等于零矩阵。V 不包含零元。

第五步：V 中的元素当 x = 0, y = 0 时为 [[0, 0], [0, 1]]，这并非零矩阵。零矩阵 [[0,0],[0,0]] 不在 V 中。

结论：V 不包含零元，因此 V 不构成 R 上的线性空间。

答案：{{expr:s}}（即 0 = 否）`,
		},
	}
}

func buildChapter7_2() dsl.Problem {
	k := "Chapter7_2"
	ids := BlankIDs(k, 5)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`线性空间 $V=\{ax^1+bx^3+cx^4 \mid a,b,c\in \mathbb{R}\}$ 的维数是 {{blank:%s}}，一组基可以取为 $x^{{blank:%s}}, x^{{blank:%s}}, x^{{blank:%s}}, x^{{blank:%s}}$（从左至右依次填写，多余的空填 0）`,
			ids[0], ids[1], ids[2], ids[3], ids[4]),
		Variables: map[string]dsl.Variable{
			"dim": {Kind: "scalar", Fixed: 3},
			"e1":  {Kind: "scalar", Fixed: 1},
			"e2":  {Kind: "scalar", Fixed: 3},
			"e3":  {Kind: "scalar", Fixed: 4},
			"e4":  {Kind: "scalar", Fixed: 0},
		},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "dim"},
			{ID: ids[1], Expr: "e1"},
			{ID: ids[2], Expr: "e2"},
			{ID: ids[3], Expr: "e3"},
			{ID: ids[4], Expr: "e4"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一步：分析空间 V 的结构。V = {ax¹ + bx³ + cx⁴ | a, b, c ∈ R}，即 V 中的每个元素都是 x¹, x³, x⁴ 这三个单项式的线性组合，系数 a, b, c 可以取任意实数。

第二步：判断生成元。V 由 {x¹, x³, x⁴} 生成，即 V = L(x¹, x³, x⁴)。这三个多项式张成整个空间 V。

第三步：验证线性无关性。设 k₁·x¹ + k₂·x³ + k₃·x⁴ = 0（零多项式）。零多项式的所有系数均为 0，即：
  x¹ 的系数：k₁ = 0
  x³ 的系数：k₂ = 0
  x⁴ 的系数：k₃ = 0
只有零解 k₁ = k₂ = k₃ = 0，因此 {x¹, x³, x⁴} 线性无关。

第四步：{x¹, x³, x⁴} 既是 V 的生成元组，又线性无关，所以它是 V 的一组基。基中含有 {{expr:dim}} = 3 个向量，故 dim(V) = {{expr:dim}}。

第五步：为什么 x² 不在 V 中？因为 V 的定义中只包含 x¹, x³, x⁴ 的线性组合，x² 的系数始终为 0，无法通过调整 a, b, c 得到 x² 项。即对任意 a, b, c ∈ R，ax¹ + bx³ + cx⁴ 中 x² 的系数恒为 0。

第六步：题目要求填写 4 个基向量指数，基只有 3 个元素 {x¹, x³, x⁴}，多余的最后一个空填 {{expr:e4}} = 0。

答案：维数 = {{expr:dim}}，基为 x^{{expr:e1}}, x^{{expr:e2}}, x^{{expr:e3}}，多余空填 {{expr:e4}}。`,
		},
	}
}

// Ch7_3: R^{2×2} 中矩阵在给定基下的坐标（同构于 R⁴）
func buildChapter7_3() dsl.Problem {
	k := "Chapter7_3"
	ids := BlankIDs(k, 4)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`线性空间 $\mathbb{R}^{2\times2}$ 中，将 $2\times2$ 矩阵的元素按行依次排列可看作 $\mathbb{R}^4$ 中的向量。元素 $\alpha={{a}}$ 在基 $B={{B}}$ 各列下的坐标：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"B": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "upper_unit", "min": -4, "max": 4}},
			"a": {Kind: "vector", Size: 4, Generator: map[string]interface{}{"rule": "range", "min": -6, "max": 6}},
		},
		Derived: map[string]string{"x": "solve(B,a)"},
		Render:  map[string]string{"B": "B", "a": "a"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "x[1]"},
			{ID: ids[1], Expr: "x[2]"},
			{ID: ids[2], Expr: "x[3]"},
			{ID: ids[3], Expr: "x[4]"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一步：将 R^{2×2} 中的矩阵按行展开为 R⁴ 中的向量。即 2×2 矩阵 [[a,b],[c,d]] 对应 R⁴ 中的向量 (a, b, c, d)^T。
  α = ({{expr:a[1]}}, {{expr:a[2]}}, {{expr:a[3]}}, {{expr:a[4]}})^T

第二步：基矩阵 B 的四列 β₁, β₂, β₃, β₄ 构成 R⁴ 的一组基。
  B = {{expr:B}}，即
  β₁ = ({{expr:mget(B,1,1)}}, {{expr:mget(B,2,1)}}, {{expr:mget(B,3,1)}}, {{expr:mget(B,4,1)}})^T
  β₂ = ({{expr:mget(B,1,2)}}, {{expr:mget(B,2,2)}}, {{expr:mget(B,3,2)}}, {{expr:mget(B,4,2)}})^T
  β₃ = ({{expr:mget(B,1,3)}}, {{expr:mget(B,2,3)}}, {{expr:mget(B,3,3)}}, {{expr:mget(B,4,3)}})^T
  β₄ = ({{expr:mget(B,1,4)}}, {{expr:mget(B,2,4)}}, {{expr:mget(B,3,4)}}, {{expr:mget(B,4,4)}})^T

B 为上三角单位对角矩阵，det(B) = 1 ≠ 0，故 B 可逆，四列确实构成 R⁴ 的基。

第三步：求 α 在基 B 下的坐标 x = (x₁, x₂, x₃, x₄)^T，即满足 Bx = α：
  {{expr:mget(B,1,1)}}·x₁ + {{expr:mget(B,1,2)}}·x₂ + {{expr:mget(B,1,3)}}·x₃ + {{expr:mget(B,1,4)}}·x₄ = {{expr:a[1]}}
  {{expr:mget(B,2,1)}}·x₁ + {{expr:mget(B,2,2)}}·x₂ + {{expr:mget(B,2,3)}}·x₃ + {{expr:mget(B,2,4)}}·x₄ = {{expr:a[2]}}
  {{expr:mget(B,3,1)}}·x₁ + {{expr:mget(B,3,2)}}·x₂ + {{expr:mget(B,3,3)}}·x₃ + {{expr:mget(B,3,4)}}·x₄ = {{expr:a[3]}}
  {{expr:mget(B,4,1)}}·x₁ + {{expr:mget(B,4,2)}}·x₂ + {{expr:mget(B,4,3)}}·x₃ + {{expr:mget(B,4,4)}}·x₄ = {{expr:a[4]}}

第四步：由于 B 为上三角单位对角矩阵（对角元全为 1，下三角全为 0），用回代法从第 4 行开始：
  第 4 行：x₄ = {{expr:a[4]}}（B 第 4 行仅有对角元为 1）
  第 3 行：x₃ + {{expr:mget(B,3,4)}}·x₄ = {{expr:a[3]}} → x₃ = {{expr:a[3]}} − {{expr:mget(B,3,4)}}·x₄
  第 2 行：x₂ + {{expr:mget(B,2,3)}}·x₃ + {{expr:mget(B,2,4)}}·x₄ = {{expr:a[2]}} → x₂ = {{expr:a[2]}} − {{expr:mget(B,2,3)}}·x₃ − {{expr:mget(B,2,4)}}·x₄
  第 1 行：x₁ + {{expr:mget(B,1,2)}}·x₂ + {{expr:mget(B,1,3)}}·x₃ + {{expr:mget(B,1,4)}}·x₄ = {{expr:a[1]}} → x₁ = {{expr:a[1]}} − {{expr:mget(B,1,2)}}·x₂ − {{expr:mget(B,1,3)}}·x₃ − {{expr:mget(B,1,4)}}·x₄

第五步：解得坐标 x = ({{expr:solve(B,a)}})^T，即
  x₁ = {{expr:x[1]}}，x₂ = {{expr:x[2]}}，x₃ = {{expr:x[3]}}，x₄ = {{expr:x[4]}}

第六步：验证：Bx = α。将各坐标代入左边，结果应等于 α。

答案：x = ({{expr:x[1]}}, {{expr:x[2]}}, {{expr:x[3]}}, {{expr:x[4]}})^T`,
		},
	}
}

func buildChapter7_5_1() dsl.Problem {
	k := "Chapter7_5_1"
	ids := BlankIDs(k, 5)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $\alpha_1,\alpha_2,\alpha_3,\alpha_4$ 为矩阵 $A={{A}}$ 的四列，向量空间 $L(\alpha_1,\alpha_2,\alpha_3,\alpha_4)$ 的维数是 {{blank:%s}}，一组基可取为 $\alpha_{{blank:%s}}, \alpha_{{blank:%s}}, \alpha_{{blank:%s}}, \alpha_{{blank:%s}}$（多余的空填 0）`,
			ids[0], ids[1], ids[2], ids[3], ids[4]),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "range", "min": -5, "max": 5}},
		},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "space_rank(A)"},
			{ID: ids[1], Expr: "basis_index(A,1)"},
			{ID: ids[2], Expr: "basis_index(A,2)"},
			{ID: ids[3], Expr: "basis_index(A,3)"},
			{ID: ids[4], Expr: "basis_index(A,4)"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一步：矩阵 A 的四个列向量为：
  α₁ = ({{expr:mget(A,1,1)}}, {{expr:mget(A,2,1)}}, {{expr:mget(A,3,1)}}, {{expr:mget(A,4,1)}})^T
  α₂ = ({{expr:mget(A,1,2)}}, {{expr:mget(A,2,2)}}, {{expr:mget(A,3,2)}}, {{expr:mget(A,4,2)}})^T
  α₃ = ({{expr:mget(A,1,3)}}, {{expr:mget(A,2,3)}}, {{expr:mget(A,3,3)}}, {{expr:mget(A,4,3)}})^T
  α₄ = ({{expr:mget(A,1,4)}}, {{expr:mget(A,2,4)}}, {{expr:mget(A,3,4)}}, {{expr:mget(A,4,4)}})^T

向量空间 L(α₁, α₂, α₃, α₄) 的维数等于 A 的列秩，即 rank(A)。

第二步：对 A 作行初等变换，化为行最简形（RREF）。行初等变换不改变列向量之间的线性关系。

第三步：统计 RREF 中主元（pivot，即每行首个非零元为 1 的位置）的个数，即为 A 的秩 r = {{expr:space_rank(A)}}。这就是向量空间的维数。

第四步：RREF 中主元所在的列号，对应原矩阵 A 中的极大线性无关组的列下标。这些列号通过 basis_index(A, k) 给出。
  第 1 个主元列：α_{basis_index(A,1)}
  第 2 个主元列：α_{basis_index(A,2)}
  第 3 个主元列：α_{basis_index(A,3)}
  第 4 个主元列：α_{basis_index(A,4)}（如果秩 < 4，则填 0 表示无此基向量）

第五步：从 A 中选取这些列构成一组基，不足的列用 0 补齐。

答案：维数 = {{expr:space_rank(A)}}，基下标依次为 α_{basis_index(A,1)}, α_{basis_index(A,2)}, ...，不足填 0`,
		},
	}
}

func buildChapter7_5_2() dsl.Problem {
	k := "Chapter7_5_2"
	ids := BlankIDs(k, 9)
	fds := make([]dsl.AnswerFieldDef, 0, 9)
	for col := 1; col <= 3; col++ {
		jb := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: fmt.Sprintf("b%d", col)}
		for row := 1; row <= 3; row++ {
			fds = append(fds, dsl.AnswerFieldDef{
				ID:    ids[len(fds)],
				Expr:  fmt.Sprintf("gs_comp(V,%d,%d)", col, row),
				Judge: jb,
			})
		}
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $V={{V}}$ 的三列为 $\mathbb{R}^3$ 中线性无关向量，Gram–Schmidt 正交化（不单位化），依次填 $\beta_1,\beta_2,\beta_3$ 各分量（答案可能含分数，请填写 a/b 格式；每个向量允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"V": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "full_rank", "min": -5, "max": 5}},
		},
		Render: map[string]string{"V": "V"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一步：设 V 的三个列向量为：
  v₁ = ({{expr:mget(V,1,1)}}, {{expr:mget(V,2,1)}}, {{expr:mget(V,3,1)}})^T
  v₂ = ({{expr:mget(V,1,2)}}, {{expr:mget(V,2,2)}}, {{expr:mget(V,3,2)}})^T
  v₃ = ({{expr:mget(V,1,3)}}, {{expr:mget(V,2,3)}}, {{expr:mget(V,3,3)}})^T
已知它们线性无关（V 满秩，rank(V) = 3）。

第二步：Gram-Schmidt 正交化（不单位化）公式：
  β₁ = v₁
  β₂ = v₂ − (⟨v₂, β₁⟩/⟨β₁, β₁⟩) · β₁
  β₃ = v₃ − (⟨v₃, β₁⟩/⟨β₁, β₁⟩) · β₁ − (⟨v₃, β₂⟩/⟨β₂, β₂⟩) · β₂
其中 ⟨·,·⟩ 表示 R³ 的标准内积（点积）。

第三步：计算 β₁ = v₁：
  β₁ = ({{expr:gs_comp(V,1,1)}}, {{expr:gs_comp(V,1,2)}}, {{expr:gs_comp(V,1,3)}})^T

第四步：计算 β₂。先求内积：
  ⟨v₂, β₁⟩ = {{expr:mget(V,1,2)}}·{{expr:gs_comp(V,1,1)}} + {{expr:mget(V,2,2)}}·{{expr:gs_comp(V,1,2)}} + {{expr:mget(V,3,2)}}·{{expr:gs_comp(V,1,3)}}
  ⟨β₁, β₁⟩ = {{expr:gs_comp(V,1,1)}}² + {{expr:gs_comp(V,1,2)}}² + {{expr:gs_comp(V,1,3)}}²
  投影系数 = ⟨v₂, β₁⟩/⟨β₁, β₁⟩

  β₂ = v₂ − 投影系数 · β₁
     = ({{expr:gs_comp(V,2,1)}}, {{expr:gs_comp(V,2,2)}}, {{expr:gs_comp(V,2,3)}})^T

第五步：计算 β₃。先求内积：
  ⟨v₃, β₁⟩ 和 ⟨v₃, β₂⟩，以及 ⟨β₁, β₁⟩ 和 ⟨β₂, β₂⟩。
  投影系数₁ = ⟨v₃, β₁⟩/⟨β₁, β₁⟩
  投影系数₂ = ⟨v₃, β₂⟩/⟨β₂, β₂⟩

  β₃ = v₃ − 投影系数₁ · β₁ − 投影系数₂ · β₂
     = ({{expr:gs_comp(V,3,1)}}, {{expr:gs_comp(V,3,2)}}, {{expr:gs_comp(V,3,3)}})^T

第六步：验证正交性（⟨βᵢ, βⱼ⟩ = 0 对 i ≠ j）：
  ⟨β₁, β₂⟩ = {{expr:gs_comp(V,1,1)}}·{{expr:gs_comp(V,2,1)}} + {{expr:gs_comp(V,1,2)}}·{{expr:gs_comp(V,2,2)}} + {{expr:gs_comp(V,1,3)}}·{{expr:gs_comp(V,2,3)}} = 0
  ⟨β₁, β₃⟩ = {{expr:gs_comp(V,1,1)}}·{{expr:gs_comp(V,3,1)}} + {{expr:gs_comp(V,1,2)}}·{{expr:gs_comp(V,3,2)}} + {{expr:gs_comp(V,1,3)}}·{{expr:gs_comp(V,3,3)}} = 0
  ⟨β₂, β₃⟩ = {{expr:gs_comp(V,2,1)}}·{{expr:gs_comp(V,3,1)}} + {{expr:gs_comp(V,2,2)}}·{{expr:gs_comp(V,3,2)}} + {{expr:gs_comp(V,2,3)}}·{{expr:gs_comp(V,3,3)}} = 0

答案：β₁ = ({{expr:gs_comp(V,1,1)}}, {{expr:gs_comp(V,1,2)}}, {{expr:gs_comp(V,1,3)}})^T
      β₂ = ({{expr:gs_comp(V,2,1)}}, {{expr:gs_comp(V,2,2)}}, {{expr:gs_comp(V,2,3)}})^T
      β₃ = ({{expr:gs_comp(V,3,1)}}, {{expr:gs_comp(V,3,2)}}, {{expr:gs_comp(V,3,3)}})^T`,
		},
	}
}

func buildChapter7_5_3() dsl.Problem {
	k := "Chapter7_5_3"
	ids := BlankIDs(k, 5)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $\alpha_1,\alpha_2,\alpha_3,\alpha_4$ 为矩阵 $A={{A}}$ 的四列，向量空间 $L(\alpha_1,\alpha_2,\alpha_3,\alpha_4)$ 的维数是 {{blank:%s}}，一组基可取为 $\alpha_{{blank:%s}}, \alpha_{{blank:%s}}, \alpha_{{blank:%s}}, \alpha_{{blank:%s}}$（多余的空填 0）`,
			ids[0], ids[1], ids[2], ids[3], ids[4]),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
		},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "space_rank(A)"},
			{ID: ids[1], Expr: "basis_index(A,1)"},
			{ID: ids[2], Expr: "basis_index(A,2)"},
			{ID: ids[3], Expr: "basis_index(A,3)"},
			{ID: ids[4], Expr: "basis_index(A,4)"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一步：矩阵 A 的四个列向量为：
  α₁ = ({{expr:mget(A,1,1)}}, {{expr:mget(A,2,1)}}, {{expr:mget(A,3,1)}}, {{expr:mget(A,4,1)}})^T
  α₂ = ({{expr:mget(A,1,2)}}, {{expr:mget(A,2,2)}}, {{expr:mget(A,3,2)}}, {{expr:mget(A,4,2)}})^T
  α₃ = ({{expr:mget(A,1,3)}}, {{expr:mget(A,2,3)}}, {{expr:mget(A,3,3)}}, {{expr:mget(A,4,3)}})^T
  α₄ = ({{expr:mget(A,1,4)}}, {{expr:mget(A,2,4)}}, {{expr:mget(A,3,4)}}, {{expr:mget(A,4,4)}})^T

向量空间 L(α₁, α₂, α₃, α₄) 的维数等于 A 的列秩。

第二步：对 A 作行初等变换化为行最简形（RREF）。行初等变换不改变列向量之间的线性关系。

第三步：RREF 中主元的个数即为 A 的秩 r = {{expr:space_rank(A)}}，这也是向量空间的维数。

第四步：RREF 中主元所在列号给出原矩阵 A 中极大线性无关组的列下标：
  第 1 个基向量：α_{basis_index(A,1)}
  第 2 个基向量：α_{basis_index(A,2)}
  第 3 个基向量：α_{basis_index(A,3)}
  第 4 个基向量：α_{basis_index(A,4)}（若秩 < 4 则填 0）

答案：维数 = {{expr:space_rank(A)}}，基下标依次为 basis_index(A,1), basis_index(A,2), ...，不足填 0`,
		},
	}
}

// Ch7_6: 在内积 (f,g)=∫₋₁¹ f·g dx 下对三个多项式做 Schmidt 正交化，
// 填写 g₁,g₂,g₃ 各系数（常数项、x系数、x²系数，共 9 空）。
func buildChapter7_6() dsl.Problem {
	k := "Chapter7_6"
	ids := BlankIDs(k, 9)
	fds := make([]dsl.AnswerFieldDef, 0, 9)
	for j := 1; j <= 3; j++ {
		jg := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: fmt.Sprintf("g%d", j)}
		for c := 1; c <= 3; c++ {
			fds = append(fds, dsl.AnswerFieldDef{
				ID:    ids[(j-1)*3+c-1],
				Expr:  fmt.Sprintf("poly_schmidt_comp(P,%d,%d)", j, c),
				Judge: jg,
			})
		}
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`在内积 $\langle f,g\rangle=\int_{-1}^{1}f(x)g(x)\,dx$ 下，对 {{input_text}} 做 Schmidt 正交化，依次填 $g_1,g_2,g_3$ 各系数（常数项、$x$ 系数、$x^2$ 系数；答案可能含分数，请填写 a/b 格式；每个多项式允许整体乘以非零数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"P": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "poly_schmidt_integral", "coef_min": -5, "coef_max": 5}},
		},
		Derived: map[string]string{
			"input_text": "poly_schmidt_input_text(P)",
		},
		Render: map[string]string{"input_text": "input_text"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一步：输入的三个多项式为 {{input_text}}。即
  f₁(x) 的系数向量 = ({{expr:mget(P,1,1)}}, {{expr:mget(P,2,1)}}, {{expr:mget(P,3,1)}})^T，f₁(x) = {{expr:poly_from_matcol(P,1)}}
  f₂(x) 的系数向量 = ({{expr:mget(P,1,2)}}, {{expr:mget(P,2,2)}}, {{expr:mget(P,3,2)}})^T，f₂(x) = {{expr:poly_from_matcol(P,2)}}
  f₃(x) 的系数向量 = ({{expr:mget(P,1,3)}}, {{expr:mget(P,2,3)}}, {{expr:mget(P,3,3)}})^T，f₃(x) = {{expr:poly_from_matcol(P,3)}}
其中系数向量的三个分量分别对应常数项、x 系数、x² 系数。

第二步：在内积 ⟨f,g⟩ = ∫₋₁¹ f(x)g(x)dx 下，基本内积值（利用奇偶性简化计算）：
  ⟨1, 1⟩ = ∫₋₁¹ 1 dx = 2
  ⟨1, x⟩ = ∫₋₁¹ x dx = 0（奇函数在对称区间上积分为 0）
  ⟨1, x²⟩ = ∫₋₁¹ x² dx = 2/3
  ⟨x, x⟩ = ∫₋₁¹ x² dx = 2/3
  ⟨x, x²⟩ = ∫₋₁¹ x³ dx = 0（奇函数）
  ⟨x², x²⟩ = ∫₋₁¹ x⁴ dx = 2/5

对于任意两个多项式 f(x) = a₀ + a₁x + a₂x² 和 g(x) = b₀ + b₁x + b₂x²，有：
  ⟨f, g⟩ = a₀b₀·2 + (a₀b₂ + a₂b₀)·(2/3) + a₁b₁·(2/3) + a₂b₂·(2/5)

第三步：g₁ = f₁ = {{expr:poly_from_matcol(P,1)}}
  g₁ 的系数：常数项 = {{expr:poly_schmidt_comp(P,1,1)}}，x 系数 = {{expr:poly_schmidt_comp(P,1,2)}}，x² 系数 = {{expr:poly_schmidt_comp(P,1,3)}}

第四步：g₂ = f₂ − (⟨f₂, g₁⟩/⟨g₁, g₁⟩) · g₁
  计算 ⟨f₂, g₁⟩ 和 ⟨g₁, g₁⟩（利用第二步的基本内积公式逐项展开）
  投影系数 c₂₁ = ⟨f₂, g₁⟩/⟨g₁, g₁⟩
  g₂ = f₂ − c₂₁ · g₁

  g₂ 的系数：常数项 = {{expr:poly_schmidt_comp(P,2,1)}}，x 系数 = {{expr:poly_schmidt_comp(P,2,2)}}，x² 系数 = {{expr:poly_schmidt_comp(P,2,3)}}

第五步：g₃ = f₃ − (⟨f₃, g₁⟩/⟨g₁, g₁⟩) · g₁ − (⟨f₃, g₂⟩/⟨g₂, g₂⟩) · g₂
  计算 ⟨f₃, g₁⟩, ⟨f₃, g₂⟩, ⟨g₂, g₂⟩
  投影系数 c₃₁ = ⟨f₃, g₁⟩/⟨g₁, g₁⟩，c₃₂ = ⟨f₃, g₂⟩/⟨g₂, g₂⟩
  g₃ = f₃ − c₃₁ · g₁ − c₃₂ · g₂

  g₃ 的系数：常数项 = {{expr:poly_schmidt_comp(P,3,1)}}，x 系数 = {{expr:poly_schmidt_comp(P,3,2)}}，x² 系数 = {{expr:poly_schmidt_comp(P,3,3)}}

第六步：验证正交性：⟨gᵢ, gⱼ⟩ = 0（i ≠ j）。利用第二步的内积公式，将 gᵢ, gⱼ 的系数代入逐项计算，结果应为 0。

答案：g₁ 系数 = ({{expr:poly_schmidt_comp(P,1,1)}}, {{expr:poly_schmidt_comp(P,1,2)}}, {{expr:poly_schmidt_comp(P,1,3)}})
      g₂ 系数 = ({{expr:poly_schmidt_comp(P,2,1)}}, {{expr:poly_schmidt_comp(P,2,2)}}, {{expr:poly_schmidt_comp(P,2,3)}})
      g₃ 系数 = ({{expr:poly_schmidt_comp(P,3,1)}}, {{expr:poly_schmidt_comp(P,3,2)}}, {{expr:poly_schmidt_comp(P,3,3)}})`,
		},
	}
}

func buildChapter7_8() dsl.Problem {
	k := "Chapter7_8"
	ids := BlankIDs(k, 9)
	fds := make([]dsl.AnswerFieldDef, 0, 9)
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[(i-1)*3+j-1],
				Expr:   fmt.Sprintf("mget(Tnew,%d,%d)", i, j),
				Layout: dsl.LayoutMatrixCell("B^{-1}A_0B", i, j, 3, 3, "B^{-1}A_0B"),
			})
		}
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设线性空间 $\mathbb{R}^3$ 上线性变换 ${{T_formula}}$，在基 $\xi_1={{b1}},\xi_2={{b2}},\xi_3={{b3}}$ 下的矩阵为：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A0": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "range", "min": -5, "max": 5}},
			"B":  {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "upper_unit", "min": -4, "max": 4}},
		},
		Derived: map[string]string{
			"T_formula": "linear_transform_title(A0)",
			"b1":        "col(B,1)",
			"b2":        "col(B,2)",
			"b3":        "col(B,3)",
			"0Bi":    "inv(B)",
			"0BiA0":  "matmul(0Bi, A0)",
			"Tnew":   "matmul(0BiA0, B)",
		},
		Render: map[string]string{"T_formula": "T_formula", "b1": "b1", "b2": "b2", "b3": "b3"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `解题步骤：

第一步：线性变换 T 在标准基下的矩阵为 A₀ = {{expr:A0}}，即 T(α) = A₀α：
  A₀ = | {{expr:mget(A0,1,1)}}  {{expr:mget(A0,1,2)}}  {{expr:mget(A0,1,3)}} |
       | {{expr:mget(A0,2,1)}}  {{expr:mget(A0,2,2)}}  {{expr:mget(A0,2,3)}} |
       | {{expr:mget(A0,3,1)}}  {{expr:mget(A0,3,2)}}  {{expr:mget(A0,3,3)}} |

第二步：新基 ξ₁, ξ₂, ξ₃ 构成的基矩阵为：
  B = {{expr:B}}，即
  ξ₁ = ({{expr:mget(B,1,1)}}, {{expr:mget(B,2,1)}}, {{expr:mget(B,3,1)}})^T
  ξ₂ = ({{expr:mget(B,1,2)}}, {{expr:mget(B,2,2)}}, {{expr:mget(B,3,2)}})^T
  ξ₃ = ({{expr:mget(B,1,3)}}, {{expr:mget(B,2,3)}}, {{expr:mget(B,3,3)}})^T

第三步：T 在新基下的矩阵 T_new = B⁻¹A₀B。推导：
  T(ξⱼ) = A₀ · ξⱼ = A₀ · (B 的第 j 列) = (A₀B) 的第 j 列
  将 T(ξⱼ) 用新基 ξ₁, ξ₂, ξ₃ 表示，其系数矩阵即为 B⁻¹A₀B。

第四步：计算 B⁻¹。B 为上三角单位对角矩阵（对角元全为 1），其逆也是上三角单位对角矩阵。
  B⁻¹ = {{expr:inv(B)}} = | {{expr:mget(0Bi,1,1)}}  {{expr:mget(0Bi,1,2)}}  {{expr:mget(0Bi,1,3)}} |
                           | {{expr:mget(0Bi,2,1)}}  {{expr:mget(0Bi,2,2)}}  {{expr:mget(0Bi,2,3)}} |
                           | {{expr:mget(0Bi,3,1)}}  {{expr:mget(0Bi,3,2)}}  {{expr:mget(0Bi,3,3)}} |

第五步：计算 B⁻¹A₀ = {{expr:matmul(0Bi,A0)}}：
  | {{expr:mget(0BiA0,1,1)}}  {{expr:mget(0BiA0,1,2)}}  {{expr:mget(0BiA0,1,3)}} |
  | {{expr:mget(0BiA0,2,1)}}  {{expr:mget(0BiA0,2,2)}}  {{expr:mget(0BiA0,2,3)}} |
  | {{expr:mget(0BiA0,3,1)}}  {{expr:mget(0BiA0,3,2)}}  {{expr:mget(0BiA0,3,3)}} |

第六步：计算 T_new = (B⁻¹A₀) · B = {{expr:matmul(0BiA0,B)}}：
  | {{expr:mget(Tnew,1,1)}}  {{expr:mget(Tnew,1,2)}}  {{expr:mget(Tnew,1,3)}} |
  | {{expr:mget(Tnew,2,1)}}  {{expr:mget(Tnew,2,2)}}  {{expr:mget(Tnew,2,3)}} |
  | {{expr:mget(Tnew,3,1)}}  {{expr:mget(Tnew,3,2)}}  {{expr:mget(Tnew,3,3)}} |

第七步：验证相似不变量：
  tr(A₀) = {{expr:mget(A0,1,1)}} + {{expr:mget(A0,2,2)}} + {{expr:mget(A0,3,3)}}
  tr(T_new) = {{expr:mget(Tnew,1,1)}} + {{expr:mget(Tnew,2,2)}} + {{expr:mget(Tnew,3,3)}}（应与 tr(A₀) 相等）

答案：T 在新基下的矩阵 T_new 的 9 个元素依次填入`,
		},
	}
}
