package bank

import (
	"fmt"

	"github.com/neumathe/la-dsl/dsl"
)

func buildChapter5_1() dsl.Problem {
	k := "Chapter5_1"
	ids := BlankIDs(k, 12)
	fds := make([]dsl.AnswerFieldDef, 0, 12)
	// 3 eigenvalues — 与 3 组 eigenvector 组成 eigen_pair，允许列顺序任意。
	for i := 1; i <= 3; i++ {
		fds = append(fds, dsl.AnswerFieldDef{
			ID:     ids[i-1],
			Expr:   fmt.Sprintf("eigenval(A,%d)", i),
			Layout: dsl.LayoutVectorComponent("lambda", i, "eigenvalues"),
			Judge: &dsl.AnswerJudgeSpec{
				Kind: "eigen_pair", EigenGroup: "E", EigenRole: "lambda",
				EigenColumn: i, MatrixVar: "A",
			},
		})
	}
	// 3 eigenvectors (3 components each)
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			kk := 3 + (i-1)*3 + j
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[kk-1],
				Expr:   fmt.Sprintf("eigenvec_comp(A,%d,%d)", i, j),
				Layout: dsl.LayoutVectorComponent(fmt.Sprintf("alpha_%d", i), j, "eigenvectors"),
				Judge: &dsl.AnswerJudgeSpec{
					Kind: "eigen_pair", EigenGroup: "E", EigenRole: "vec",
					EigenColumn: i, EigenComponent: j, MatrixVar: "A",
				},
			})
		}
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $A={{A}}$，求 $A$ 的三个特征值 $\lambda_1,\lambda_2,\lambda_3$ 及对应的特征向量 $\alpha_1,\alpha_2,\alpha_3$ 各分量（特征值与特征向量的顺序可任意对应；特征向量允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "eigen_reverse_3x3"}},
		},
		Render: map[string]string{"A": "A"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：求解矩阵特征值与特征向量，先由特征方程求特征值，再对每个特征值求解齐次线性方程组得到特征向量。

**步骤 1**：求特征值。由特征方程 $|\lambda E - A| = 0$，代入 $A = {{A}}$ 得：
$$|\lambda E - A| = \begin{vmatrix} \lambda-{{expr:mget(A,1,1)}} & -{{expr:mget(A,1,2)}} & -{{expr:mget(A,1,3)}} \\ -{{expr:mget(A,2,1)}} & \lambda-{{expr:mget(A,2,2)}} & -{{expr:mget(A,2,3)}} \\ -{{expr:mget(A,3,1)}} & -{{expr:mget(A,3,2)}} & \lambda-{{expr:mget(A,3,3)}} \end{vmatrix} = 0$$
展开为关于 $\lambda$ 的三次方程，解得三个特征值：
$$\lambda_1 = {{expr:eigenval(A,1)}},\quad \lambda_2 = {{expr:eigenval(A,2)}},\quad \lambda_3 = {{expr:eigenval(A,3)}}$$

**步骤 2**：对每个特征值 $\lambda_i$，解齐次线性方程组 $(\lambda_i E - A)x = 0$ 求特征向量。

当 $\lambda_1 = {{expr:eigenval(A,1)}}$ 时，$\lambda_1 E - A = \begin{bmatrix} {{expr:eigenval(A,1)}}-{{expr:mget(A,1,1)}} & -{{expr:mget(A,1,2)}} & -{{expr:mget(A,1,3)}} \\ -{{expr:mget(A,2,1)}} & {{expr:eigenval(A,1)}}-{{expr:mget(A,2,2)}} & -{{expr:mget(A,2,3)}} \\ -{{expr:mget(A,3,1)}} & -{{expr:mget(A,3,2)}} & {{expr:eigenval(A,1)}}-{{expr:mget(A,3,3)}} \end{bmatrix}$，行化简后求得特征向量 $\alpha_1 = ({{expr:eigenvec_comp(A,1,1)}},{{expr:eigenvec_comp(A,1,2)}},{{expr:eigenvec_comp(A,1,3)}})^T$。

当 $\lambda_2 = {{expr:eigenval(A,2)}}$ 时，同理得 $\alpha_2 = ({{expr:eigenvec_comp(A,2,1)}},{{expr:eigenvec_comp(A,2,2)}},{{expr:eigenvec_comp(A,2,3)}})^T$。

当 $\lambda_3 = {{expr:eigenval(A,3)}}$ 时，同理得 $\alpha_3 = ({{expr:eigenvec_comp(A,3,1)}},{{expr:eigenvec_comp(A,3,2)}},{{expr:eigenvec_comp(A,3,3)}})^T$。

**步骤 3**：将三个特征值与对应的特征向量填入各空，特征值与特征向量可任意配对填写，特征向量允许整体乘以非零整数不影响判分。`,
		},
	}
}

func buildChapter5_2() dsl.Problem {
	k := "Chapter5_2"
	ids := BlankIDs(k, 12)
	fds := make([]dsl.AnswerFieldDef, 0, 12)
	for i := 1; i <= 3; i++ {
		fds = append(fds, dsl.AnswerFieldDef{
			ID:     ids[i-1],
			Expr:   fmt.Sprintf("eigenval(A,%d)", i),
			Layout: dsl.LayoutVectorComponent("lambda", i, "eigenvalues"),
			Judge: &dsl.AnswerJudgeSpec{
				Kind: "eigen_pair", EigenGroup: "E", EigenRole: "lambda",
				EigenColumn: i, MatrixVar: "A",
			},
		})
	}
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			kk := 3 + (i-1)*3 + j
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[kk-1],
				Expr:   fmt.Sprintf("eigenvec_comp(A,%d,%d)", i, j),
				Layout: dsl.LayoutVectorComponent(fmt.Sprintf("alpha_%d", i), j, "eigenvectors"),
				Judge: &dsl.AnswerJudgeSpec{
					Kind: "eigen_pair", EigenGroup: "E", EigenRole: "vec",
					EigenColumn: i, EigenComponent: j, MatrixVar: "A",
				},
			})
		}
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $A={{A}}$，求 $A$ 的三个特征值 $\lambda_1,\lambda_2,\lambda_3$ 及对应的特征向量 $\alpha_1,\alpha_2,\alpha_3$ 各分量（特征值与特征向量的顺序可任意对应；特征向量允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "eigen_reverse_3x3"}},
		},
		Render: map[string]string{"A": "A"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：求解矩阵特征值与特征向量，先由特征方程求特征值，再对每个特征值求解齐次线性方程组得到特征向量。

**步骤 1**：求特征值。由特征方程 $|\lambda E - A| = 0$，代入 $A = {{A}}$ 得：
$$|\lambda E - A| = \begin{vmatrix} \lambda-{{expr:mget(A,1,1)}} & -{{expr:mget(A,1,2)}} & -{{expr:mget(A,1,3)}} \\ -{{expr:mget(A,2,1)}} & \lambda-{{expr:mget(A,2,2)}} & -{{expr:mget(A,2,3)}} \\ -{{expr:mget(A,3,1)}} & -{{expr:mget(A,3,2)}} & \lambda-{{expr:mget(A,3,3)}} \end{vmatrix} = 0$$
展开为关于 $\lambda$ 的三次方程，解得三个特征值：
$$\lambda_1 = {{expr:eigenval(A,1)}},\quad \lambda_2 = {{expr:eigenval(A,2)}},\quad \lambda_3 = {{expr:eigenval(A,3)}}$$

**步骤 2**：对每个特征值 $\lambda_i$，解齐次线性方程组 $(\lambda_i E - A)x = 0$ 求特征向量。

当 $\lambda_1 = {{expr:eigenval(A,1)}}$ 时，$\lambda_1 E - A = \begin{bmatrix} {{expr:eigenval(A,1)}}-{{expr:mget(A,1,1)}} & -{{expr:mget(A,1,2)}} & -{{expr:mget(A,1,3)}} \\ -{{expr:mget(A,2,1)}} & {{expr:eigenval(A,1)}}-{{expr:mget(A,2,2)}} & -{{expr:mget(A,2,3)}} \\ -{{expr:mget(A,3,1)}} & -{{expr:mget(A,3,2)}} & {{expr:eigenval(A,1)}}-{{expr:mget(A,3,3)}} \end{bmatrix}$，行化简后求得特征向量 $\alpha_1 = ({{expr:eigenvec_comp(A,1,1)}},{{expr:eigenvec_comp(A,1,2)}},{{expr:eigenvec_comp(A,1,3)}})^T$。

当 $\lambda_2 = {{expr:eigenval(A,2)}}$ 时，同理得 $\alpha_2 = ({{expr:eigenvec_comp(A,2,1)}},{{expr:eigenvec_comp(A,2,2)}},{{expr:eigenvec_comp(A,2,3)}})^T$。

当 $\lambda_3 = {{expr:eigenval(A,3)}}$ 时，同理得 $\alpha_3 = ({{expr:eigenvec_comp(A,3,1)}},{{expr:eigenvec_comp(A,3,2)}},{{expr:eigenvec_comp(A,3,3)}})^T$。

**步骤 3**：将三个特征值与对应的特征向量填入各空，特征值与特征向量可任意配对填写，特征向量允许整体乘以非零整数不影响判分。`,
		},
	}
}

func buildChapter5_3() dsl.Problem {
	k := "Chapter5_3"
	ids := BlankIDs(k, 4)
	// 前 3 个是 A* 的特征值，顺序无关；最后一个是 det(...)，标量判分。
	jmu := &dsl.AnswerJudgeSpec{Kind: "permutation_multiset", PermGroup: "mu"}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`已知三阶矩阵 $A$ 的三个特征值分别为 {{eig_text}}，则 $A^\*$ 的三个特征值及 $\det(A^2+6A-2E)$（前三空顺序任意，共 4 空）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "diagonal_distinct", "min": -6, "max": 6}},
		},
		Derived: map[string]string{
			"eig_text": "eigenval_list_text(A)",
			"mu1": "diag_adj(A,1)",
			"mu2": "diag_adj(A,2)",
			"mu3": "diag_adj(A,3)",
			"d4":  "det_quad_shift_diag(A,6,-2)",
		},
		Render: map[string]string{"eig_text": "eig_text"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "mu1", Judge: jmu, Layout: dsl.LayoutVectorComponent("mu", 1, "eigenvalues of A^*")},
			{ID: ids[1], Expr: "mu2", Judge: jmu, Layout: dsl.LayoutVectorComponent("mu", 2, "eigenvalues of A^*")},
			{ID: ids[2], Expr: "mu3", Judge: jmu, Layout: dsl.LayoutVectorComponent("mu", 3, "eigenvalues of A^*")},
			{ID: ids[3], Expr: "d4", Layout: dsl.LayoutVectorComponent("det", 1, "\\det(A^2+6A-2I)")},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：利用特征值的函数性质——若 $\lambda$ 是 $A$ 的特征值，则 $f(\lambda)$ 是 $f(A)$ 的特征值。

**步骤 1**：矩阵 $A = {{A}}$ 为对角矩阵，三个特征值即为对角元：$\lambda_1 = {{expr:mget(A,1,1)}},\lambda_2 = {{expr:mget(A,2,2)}},\lambda_3 = {{expr:mget(A,3,3)}}$。

**步骤 2**：计算 $\det(A) = \lambda_1 \cdot \lambda_2 \cdot \lambda_3 = {{expr:det(A)}}$。由伴随矩阵公式 $A^* = \det(A) \cdot A^{-1}$，$A^*$ 的特征值为 $\mu_i = \dfrac{\det(A)}{\lambda_i}$：
$$\mu_1 = \frac{{expr:det(A)}}{{expr:mget(A,1,1)}} = {{expr:diag_adj(A,1)}},\quad \mu_2 = \frac{{expr:det(A)}}{{expr:mget(A,2,2)}} = {{expr:diag_adj(A,2)}},\quad \mu_3 = \frac{{expr:det(A)}}{{expr:mget(A,3,3)}} = {{expr:diag_adj(A,3)}}$$

**步骤 3**：对于多项式 $f(A) = A^2+6A-2E$，若 $\lambda_i$ 是 $A$ 的特征值，则 $f(\lambda_i) = \lambda_i^2+6\lambda_i-2$ 是 $f(A)$ 的特征值。分别计算：
$$f(\lambda_1) = {{expr:mget(A,1,1)}}^2 + 6\cdot{{expr:mget(A,1,1)}} - 2,\quad f(\lambda_2) = {{expr:mget(A,2,2)}}^2 + 6\cdot{{expr:mget(A,2,2)}} - 2,\quad f(\lambda_3) = {{expr:mget(A,3,3)}}^2 + 6\cdot{{expr:mget(A,3,3)}} - 2$$
$$\det(A^2+6A-2E) = f(\lambda_1)\cdot f(\lambda_2)\cdot f(\lambda_3) = {{expr:det_quad_shift_diag(A,6,-2)}}$$

**步骤 4**：综上，$A^*$ 的三个特征值为 {{expr:diag_adj(A,1)}}, {{expr:diag_adj(A,2)}}, {{expr:diag_adj(A,3)}}（顺序任意），$\det(A^2+6A-2E) = {{expr:det_quad_shift_diag(A,6,-2)}}$。`,
		},
	}
}

func buildChapter5_5() dsl.Problem {
	k := "Chapter5_5"
	ids := BlankIDs(k, 5)
	je := &dsl.AnswerJudgeSpec{Kind: "permutation_multiset", PermGroup: "eig"}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`已知五阶矩阵 $A$ 满足 $${{conditions}}$$，其特征值为（顺序任意）%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"s": {Kind: "scalar", Generator: map[string]interface{}{"rule": "eigen_rank_inference_5", "lambda_min": -8, "lambda_max": 8}},
		},
		Derived: map[string]string{
			"conditions": "eigen_rank_condition_text(s)",
			"e1":         "eigenval_rank(s,1)",
			"e2":         "eigenval_rank(s,2)",
			"e3":         "eigenval_rank(s,3)",
			"e4":         "eigenval_rank(s,4)",
			"e5":         "eigenval_rank(s,5)",
		},
		Render: map[string]string{"conditions": "conditions"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "e1", Judge: je, Layout: dsl.LayoutVectorComponent("lambda", 1, "eigenvalues")},
			{ID: ids[1], Expr: "e2", Judge: je, Layout: dsl.LayoutVectorComponent("lambda", 2, "eigenvalues")},
			{ID: ids[2], Expr: "e3", Judge: je, Layout: dsl.LayoutVectorComponent("lambda", 3, "eigenvalues")},
			{ID: ids[3], Expr: "e4", Judge: je, Layout: dsl.LayoutVectorComponent("lambda", 4, "eigenvalues")},
			{ID: ids[4], Expr: "e5", Judge: je, Layout: dsl.LayoutVectorComponent("lambda", 5, "eigenvalues")},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：利用秩-零度定理——若 $R(A-\lambda_i E) = r_i$，则 $\lambda_i$ 作为特征值的几何重数（零度）为 $n - r_i = 5 - r_i$。

**步骤 1**：题目给出的条件为 ${{conditions}}$，即三个关于 $A-\lambda_i E$ 的秩条件。

**步骤 2**：由每个秩条件 $R(A-\lambda_i E) = r_i$ 得 $\dim N(A-\lambda_i E) = 5 - r_i$，即特征值 $\lambda_i$ 的几何重数为 $5-r_i$：
- 由第一个秩条件得一个特征值，其重数为 $5 - r_1$
- 由第二个秩条件得一个特征值，其重数为 $5 - r_2$
- 由第三个秩条件得一个特征值，其重数为 $5 - r_3$

**步骤 3**：将各特征值的重数相加验证总和为 $5$，确认没有遗漏的特征值（对于对称矩阵，几何重数等于代数重数）。

**步骤 4**：综合得五个特征值（从小到大排列）为 {{e1}}, {{e2}}, {{e3}}, {{e4}}, {{e5}}（顺序任意）。`,
		},
	}
}

func buildChapter5_4() dsl.Problem {
	k := "Chapter5_4"
	ids := BlankIDs(k, 4)
	je := &dsl.AnswerJudgeSpec{Kind: "permutation_multiset", PermGroup: "eig"}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`已知四阶矩阵 $A$ 的秩等于 {{r}}，$A$ 的各行元素之和都等于 {{s}}，且 $|A+{{k}}E|=0$，则矩阵 $A$ 的特征值为（顺序任意）%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"s": {Kind: "scalar", Generator: map[string]interface{}{"rule": "eigen_row_sum_rank4", "row_sum_min": 2, "row_sum_max": 6, "k_min": 2, "k_max": 10}},
		},
		Derived: map[string]string{
			"r":  "eigen_rowsum_r(s)",
			"s":  "eigen_rowsum_s(s)",
			"k":  "eigen_rowsum_k(s)",
			"e1": "eigenval_rowsum(s,1)",
			"e2": "eigenval_rowsum(s,2)",
			"e3": "eigenval_rowsum(s,3)",
			"e4": "eigenval_rowsum(s,4)",
		},
		Render: map[string]string{"r": "r", "s": "s", "k": "k"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "e1", Judge: je, Layout: dsl.LayoutVectorComponent("lambda", 1, "eigenvalues")},
			{ID: ids[1], Expr: "e2", Judge: je, Layout: dsl.LayoutVectorComponent("lambda", 2, "eigenvalues")},
			{ID: ids[2], Expr: "e3", Judge: je, Layout: dsl.LayoutVectorComponent("lambda", 3, "eigenvalues")},
			{ID: ids[3], Expr: "e4", Judge: je, Layout: dsl.LayoutVectorComponent("lambda", 4, "eigenvalues")},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：利用秩、行和、行列式三个条件分别推出不同的特征值，共得到四个特征值。

**步骤 1**：由 $R(A) = {{r}}$，可知 $\dim N(A) = 4 - {{r}} = 4-{{r}}$，即特征值 $\lambda = 0$ 的代数重数为 $4-{{r}}$。

**步骤 2**：由 $A$ 的各行元素之和都等于 ${{s}}$，可知 $A(1,1,1,1)^T = {{s}}(1,1,1,1)^T$，即 $\lambda = {{s}}$ 是一个特征值，对应特征向量为 $(1,1,1,1)^T$。

**步骤 3**：由 $|A+{{k}}E| = 0$，即 $|A-(-{{k}})E| = 0$，可知 $\lambda = -{{k}}$ 是 $A$ 的一个特征值。

**步骤 4**：综合以上，四个特征值分别为：$0$（$4-{{r}}$ 重）、${{s}}$（1 重）、$-{{k}}$（1 重），共计 $4-{{r}} + 1 + 1 = 4$ 个。从小到大排列为：{{e1}}, {{e2}}, {{e3}}, {{e4}}。`,
		},
	}
}

func buildChapter5_8() dsl.Problem {
	k := "Chapter5_8"
	ids := BlankIDs(k, 24)
	fds := make([]dsl.AnswerFieldDef, 0, 24)
	// 3 eigenvalues（顺序任意）与 α、Q 共同形成 eigen_pair。
	for i := 1; i <= 3; i++ {
		fds = append(fds, dsl.AnswerFieldDef{
			ID:     ids[i-1],
			Expr:   fmt.Sprintf("sym_eigenval(S,%d)", i),
			Layout: dsl.LayoutVectorComponent("lambda", i, "eigenvalues"),
			Judge: &dsl.AnswerJudgeSpec{
				Kind: "eigen_pair", EigenGroup: "EA", EigenRole: "lambda",
				EigenColumn: i, MatrixVar: "S",
			},
		})
	}
	// 3 eigenvectors α_1, α_2, α_3（每个 3 分量）与 λ 配对。
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			idx := 3 + (i-1)*3 + j
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[idx-1],
				Expr:   fmt.Sprintf("sym_eigenvec_comp(S,%d,%d)", i, j),
				Layout: dsl.LayoutVectorComponent(fmt.Sprintf("alpha_%d", i), j, "eigenvectors"),
				Judge: &dsl.AnswerJudgeSpec{
					Kind: "eigen_pair", EigenGroup: "EA", EigenRole: "vec",
					EigenColumn: i, EigenComponent: j, MatrixVar: "S",
				},
			})
		}
	}
	// Q 的 9 个元素：列 i 是第 i 个特征对的向量（借用 EA 组的 λ）。
	// 矩阵 (行 i, 列 j) 的元素 = sym_eigenvec_comp(S, j, i)
	for row := 1; row <= 3; row++ {
		for col := 1; col <= 3; col++ {
			idx := 12 + (row-1)*3 + col
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[idx-1],
				Expr:   fmt.Sprintf("sym_eigenvec_comp(S,%d,%d)", col, row),
				Layout: dsl.LayoutMatrixCell("Q", row, col, 3, 3, "Q"),
				Judge: &dsl.AnswerJudgeSpec{
					Kind: "eigen_pair", EigenGroup: "EQ", EigenRole: "vec",
					EigenColumn: col, EigenComponent: row, MatrixVar: "S",
					RefLambdaGroup: "EA",
				},
			})
		}
	}
	// Λ = Q^{-1}SQ 的对角元（3 空；与 λ multiset 一致即可）。
	jlam := &dsl.AnswerJudgeSpec{Kind: "permutation_multiset", PermGroup: "diag_lambda"}
	for i := 1; i <= 3; i++ {
		idx := 21 + i
		fds = append(fds, dsl.AnswerFieldDef{
			ID:     ids[idx-1],
			Expr:   fmt.Sprintf("sym_eigenval(S,%d)", i),
			Layout: dsl.LayoutMatrixCell("Lambda", i, i, 3, 3, "Lambda"),
			Judge:  jlam,
		})
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设对称阵 $S={{S}}$，求正交矩阵 $Q$ 使 $Q^{-1}SQ$ 为对角阵。填写三个特征值、三个特征向量 $\alpha_1,\alpha_2,\alpha_3$ 各分量、$Q$ 的 9 个元素（按行自上而下、从左到右），以及 $Q^{-1}SQ$ 对角元（特征值与对应特征向量的顺序可任意对应；每个特征向量/Q 每一列允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"S": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "symmetric_eigen_reverse_3x3", "lambda_min": -5, "lambda_max": 5, "max_entry": 12}},
		},
		Render: map[string]string{"S": "S"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：实对称矩阵必可正交相似对角化。先求特征值，再求特征向量并单位化得到正交矩阵 $Q$，最后验证 $Q^{-1}SQ = \Lambda$。

**步骤 1**：求特征值。解特征方程 $|\lambda E - S| = 0$，代入 $S = {{S}}$ 得：
$$|\lambda E - S| = \begin{vmatrix} \lambda-{{expr:mget(S,1,1)}} & -{{expr:mget(S,1,2)}} & -{{expr:mget(S,1,3)}} \\ -{{expr:mget(S,2,1)}} & \lambda-{{expr:mget(S,2,2)}} & -{{expr:mget(S,2,3)}} \\ -{{expr:mget(S,3,1)}} & -{{expr:mget(S,3,2)}} & \lambda-{{expr:mget(S,3,3)}} \end{vmatrix} = 0$$
解得三个特征值：
$$\lambda_1 = {{expr:sym_eigenval(S,1)}},\quad \lambda_2 = {{expr:sym_eigenval(S,2)}},\quad \lambda_3 = {{expr:sym_eigenval(S,3)}}$$

**步骤 2**：对每个特征值 $\lambda_i$，解 $(\lambda_i E - S)x = 0$ 得特征向量 $\alpha_i$。

当 $\lambda_1 = {{expr:sym_eigenval(S,1)}}$ 时，解 $(\lambda_1 E - S)x = 0$ 得 $\alpha_1 = ({{expr:sym_eigenvec_comp(S,1,1)}},{{expr:sym_eigenvec_comp(S,1,2)}},{{expr:sym_eigenvec_comp(S,1,3)}})^T$。

当 $\lambda_2 = {{expr:sym_eigenval(S,2)}}$ 时，解 $(\lambda_2 E - S)x = 0$ 得 $\alpha_2 = ({{expr:sym_eigenvec_comp(S,2,1)}},{{expr:sym_eigenvec_comp(S,2,2)}},{{expr:sym_eigenvec_comp(S,2,3)}})^T$。

当 $\lambda_3 = {{expr:sym_eigenval(S,3)}}$ 时，解 $(\lambda_3 E - S)x = 0$ 得 $\alpha_3 = ({{expr:sym_eigenvec_comp(S,3,1)}},{{expr:sym_eigenvec_comp(S,3,2)}},{{expr:sym_eigenvec_comp(S,3,3)}})^T$。

**步骤 3**：单位化。由于 $S$ 为实对称矩阵，不同特征值对应的特征向量天然正交。将每个 $\alpha_i$ 单位化得 $q_i = \dfrac{\alpha_i}{\|\alpha_i\|}$。

**步骤 4**：构造正交矩阵 $Q = (q_1, q_2, q_3)$，即 $Q$ 的第 $i$ 列为单位化后的第 $i$ 个特征向量。$Q$ 的 9 个元素按行自上而下、从左到右依次填写。

**步骤 5**：对角阵 $\Lambda = Q^{-1}SQ = \operatorname{diag}(\lambda_1,\lambda_2,\lambda_3)$，对角元依次为 {{expr:sym_eigenval(S,1)}}, {{expr:sym_eigenval(S,2)}}, {{expr:sym_eigenval(S,3)}}（与特征向量顺序对应）。`,
		},
	}
}

func buildChapter5_6() dsl.Problem {
	k := "Chapter5_6"
	ids := BlankIDs(k, 9)
	fds := make([]dsl.AnswerFieldDef, 0, 9)
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[(i-1)*3+j-1],
				Expr:   fmt.Sprintf("mget(A,%d,%d)", i, j),
				Layout: dsl.LayoutMatrixCell("A", i, j, 3, 3, "A"),
			})
		}
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $M={{M}}$ 可逆，三列依次为 $\alpha_1,\alpha_2,\alpha_3$。求满足 $A\alpha_k=\alpha_1(k=1,2,3)$ 的矩阵 $A$：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"M": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "upper_unit", "min": -4, "max": 4}},
		},
		Derived: map[string]string{
			"0Minv": "inv(M)",
			"1C":    "triple_first_col(M)",
			"A":     "matmul(1C, 0Minv)",
		},
		Render: map[string]string{"M": "M"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：将条件写成矩阵方程 $AM = C$，其中 $C$ 的三列均为 $\alpha_1$，从而 $A = CM^{-1}$。

**步骤 1**：记 $M = (\alpha_1,\alpha_2,\alpha_3) = {{M}}$，则条件 $A\alpha_k = \alpha_1$（$k=1,2,3$）等价于：
$$AM = (\alpha_1,\alpha_1,\alpha_1) = C$$
其中 $C$ 为三列相同的矩阵，每列均为 $\alpha_1 = ({{expr:mget(M,1,1)}},{{expr:mget(M,2,1)}},{{expr:mget(M,3,1)}})^T$。

**步骤 2**：由 $AM = C$ 且 $M$ 可逆，得 $A = CM^{-1}$。先计算 $M$ 的逆矩阵：
$$M^{-1} = {{0Minv}}$$

**步骤 3**：构造 $C$ 矩阵（三列均为 $\alpha_1$）：
$$C = {{1C}} = \begin{bmatrix} {{expr:mget(M,1,1)}} & {{expr:mget(M,1,1)}} & {{expr:mget(M,1,1)}} \\ {{expr:mget(M,2,1)}} & {{expr:mget(M,2,1)}} & {{expr:mget(M,2,1)}} \\ {{expr:mget(M,3,1)}} & {{expr:mget(M,3,1)}} & {{expr:mget(M,3,1)}} \end{bmatrix}$$

**步骤 4**：计算 $A = CM^{-1}$：
$$A = {{A}} = \begin{bmatrix} {{expr:mget(A,1,1)}} & {{expr:mget(A,1,2)}} & {{expr:mget(A,1,3)}} \\ {{expr:mget(A,2,1)}} & {{expr:mget(A,2,2)}} & {{expr:mget(A,2,3)}} \\ {{expr:mget(A,3,1)}} & {{expr:mget(A,3,2)}} & {{expr:mget(A,3,3)}} \end{bmatrix}$$
九个元素按行依次填入各空。`,
		},
	}
}

func buildChapter5_7() dsl.Problem {
	k := "Chapter5_7"
	ids := BlankIDs(k, 12)
	fds := make([]dsl.AnswerFieldDef, 0, 12)
	// 3 eigenvalues (顺序任意；与 Q 的 9 个元素成 eigen_pair)
	for i := 1; i <= 3; i++ {
		fds = append(fds, dsl.AnswerFieldDef{
			ID:     ids[i-1],
			Expr:   fmt.Sprintf("sym_eigenval(S,%d)", i),
			Layout: dsl.LayoutVectorComponent("lambda", i, "特征值"),
			Judge: &dsl.AnswerJudgeSpec{
				Kind: "eigen_pair", EigenGroup: "E", EigenRole: "lambda",
				EigenColumn: i, MatrixVar: "S",
			},
		})
	}
	// Q 的 9 个元素（按行 i、列 j 展开）：Q_ij = sym_eigenvec_comp(S, j, i)
	for row := 1; row <= 3; row++ {
		for col := 1; col <= 3; col++ {
			idx := 3 + (row-1)*3 + col
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[idx-1],
				Expr:   fmt.Sprintf("sym_eigenvec_comp(S,%d,%d)", col, row),
				Layout: dsl.LayoutMatrixCell("Q", row, col, 3, 3, "Q"),
				Judge: &dsl.AnswerJudgeSpec{
					Kind: "eigen_pair", EigenGroup: "E", EigenRole: "vec",
					EigenColumn: col, EigenComponent: row, MatrixVar: "S",
				},
			})
		}
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设对称阵 $S={{S}}$，求正交矩阵 $Q$ 使 $Q^{-1}SQ=\Lambda$（对角阵），填写三个特征值及 $Q$ 的 9 个元素（按行自上而下、从左到右；特征值与 $Q$ 的列顺序可任意对应，每一列允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"S": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "symmetric_eigen_reverse_3x3", "lambda_min": -5, "lambda_max": 5, "max_entry": 12}},
		},
		Render: map[string]string{"S": "S"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：对称矩阵可正交相似对角化。先求特征值和特征向量，将特征向量单位化后构成正交矩阵 $Q$。

**步骤 1**：求特征值。解特征方程 $|\lambda E - S| = 0$，代入 $S = {{S}}$ 得：
$$|\lambda E - S| = \begin{vmatrix} \lambda-{{expr:mget(S,1,1)}} & -{{expr:mget(S,1,2)}} & -{{expr:mget(S,1,3)}} \\ -{{expr:mget(S,2,1)}} & \lambda-{{expr:mget(S,2,2)}} & -{{expr:mget(S,2,3)}} \\ -{{expr:mget(S,3,1)}} & -{{expr:mget(S,3,2)}} & \lambda-{{expr:mget(S,3,3)}} \end{vmatrix} = 0$$
展开为关于 $\lambda$ 的三次方程，解得三个特征值：
$$\lambda_1 = {{expr:sym_eigenval(S,1)}},\quad \lambda_2 = {{expr:sym_eigenval(S,2)}},\quad \lambda_3 = {{expr:sym_eigenval(S,3)}}$$

**步骤 2**：对每个特征值 $\lambda_i$，解 $(\lambda_i E - S)x = 0$ 得特征向量 $\alpha_i$。

当 $\lambda_1 = {{expr:sym_eigenval(S,1)}}$ 时，$\alpha_1 = ({{expr:sym_eigenvec_comp(S,1,1)}},{{expr:sym_eigenvec_comp(S,1,2)}},{{expr:sym_eigenvec_comp(S,1,3)}})^T$。

当 $\lambda_2 = {{expr:sym_eigenval(S,2)}}$ 时，$\alpha_2 = ({{expr:sym_eigenvec_comp(S,2,1)}},{{expr:sym_eigenvec_comp(S,2,2)}},{{expr:sym_eigenvec_comp(S,2,3)}})^T$。

当 $\lambda_3 = {{expr:sym_eigenval(S,3)}}$ 时，$\alpha_3 = ({{expr:sym_eigenvec_comp(S,3,1)}},{{expr:sym_eigenvec_comp(S,3,2)}},{{expr:sym_eigenvec_comp(S,3,3)}})^T$。

**步骤 3**：由于 $S$ 是实对称矩阵，不同特征值对应的特征向量天然正交。将每个特征向量单位化得 $q_i = \dfrac{\alpha_i}{\|\alpha_i\|}$。

**步骤 4**：构造正交矩阵 $Q = (q_1, q_2, q_3)$，$Q$ 的第 $i$ 列为单位化后的第 $i$ 个特征向量。$Q$ 的 9 个元素按行自上而下、从左到右依次填入各空。`,
		},
	}
}
