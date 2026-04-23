package bank

import (
	"fmt"

	"github.com/neumathe/la-dsl/dsl"
)

func chapter6QuadraticMatrix10(key string, symExpr string, codeHint string) dsl.Problem {
	ids := BlankIDs(key, 10)
	fds := make([]dsl.AnswerFieldDef, 0, 10)
	k := 1
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[k-1],
				Expr:   fmt.Sprintf("mget(S,%d,%d)", i, j),
				Layout: dsl.LayoutMatrixCell("S", i, j, 3, 3, "S"),
			})
			k++
		}
	}
	fds = append(fds, dsl.AnswerFieldDef{ID: ids[9], Expr: symExpr})
	title := fmt.Sprintf(
		`设二次型 $f(x_1,x_2,x_3)={{expr}}$，写出对应的对称矩阵 $S$ 的 9 个元素（按行自上而下、从左到右），最后一空填写正定性代码（%s）：%s`,
		codeHint, joinBlankPlaceholders(ids))
	return dsl.Problem{
		ID:      ProblemID(key),
		Version: "bank-v1",
		Title:   title,
		Variables: map[string]dsl.Variable{
			"S": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "symmetric", "min": -8, "max": 8}},
		},
		Derived: map[string]string{"expr": "quad_expr(S)"},
		Render:  map[string]string{"expr": "expr"},
		Answer:  dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：二次型 $f(x_1,x_2,x_3)={{expr}}$ 的对称矩阵 $S$ 的构造规则为：主对角元 $S_{ii}$ 等于 $x_i^2$ 的系数，非对角元 $S_{ij}$（$i\neq j$）等于交叉项 $x_ix_j$ 系数的一半。

**步骤 1**：从 $f$ 的表达式中读出平方项系数，它们即为主对角元。$x_1^2$ 的系数为 $S_{11}={{expr:mget(S,1,1)}}$，$x_2^2$ 的系数为 $S_{22}={{expr:mget(S,2,2)}}$，$x_3^2$ 的系数为 $S_{33}={{expr:mget(S,3,3)}}$。

**步骤 2**：从 $f$ 的表达式中读出交叉项系数，除以 2 即为非对角元。$x_1x_2$ 系数的一半为 $S_{12}=S_{21}={{expr:mget(S,1,2)}}$，$x_1x_3$ 系数的一半为 $S_{13}=S_{31}={{expr:mget(S,1,3)}}$，$x_2x_3$ 系数的一半为 $S_{23}=S_{32}={{expr:mget(S,2,3)}}$。

**步骤 3**：写出对称矩阵 $S=\begin{pmatrix}{{expr:mget(S,1,1)}}&{{expr:mget(S,1,2)}}&{{expr:mget(S,1,3)}}\\{{expr:mget(S,2,1)}}&{{expr:mget(S,2,2)}}&{{expr:mget(S,2,3)}}\\{{expr:mget(S,3,1)}}&{{expr:mget(S,3,2)}}&{{expr:mget(S,3,3)}}\end{pmatrix}$，其 9 个元素自上而下、从左到右依次为 $S_{11},S_{12},S_{13},S_{21},S_{22},S_{23},S_{31},S_{32},S_{33}$。

**步骤 4**：判断正定性。计算 $S$ 的各阶顺序主子式：$\Delta_1=S_{11}={{expr:mget(S,1,1)}}$，$\Delta_2=\det\begin{pmatrix}S_{11}&S_{12}\\S_{21}&S_{22}\end{pmatrix}={{expr:mget(S,1,1)}}\times{{expr:mget(S,2,2)}}-{{expr:mget(S,1,2)}}^2$，$\Delta_3=\det(S)={{expr:det(S)}}$。根据各阶主子式的符号对照题面编码规则确定正定性代码。`,
		},
	}
}

func buildChapter6_1_1() dsl.Problem {
	return chapter6QuadraticMatrix10("Chapter6_1_1", "symcode_611(S)", "正定填1，负定填0，不定填2")
}

func buildChapter6_1_2() dsl.Problem {
	return chapter6QuadraticMatrix10("Chapter6_1_2", "symcode_612(S)", "正定填1，负定填2，不定填0")
}

func buildChapter6_1_3() dsl.Problem {
	return chapter6QuadraticMatrix10("Chapter6_1_3", "symcode_612(S)", "正定填1，负定填2，不定填0")
}

// Ch6_2: 给定 3×3 对称阵 A 和对角阵 B，判断是否相似（0/1）和是否合同（0/1）。
// 相似 <=> 特征值相同；合同 <=> 惯性指数相同。
func buildChapter6_2() dsl.Problem {
	k := "Chapter6_2"
	ids := BlankIDs(k, 2)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设对称阵 $A={{A}}$ 与对角阵 $B={{B}}$，判断 A 与 B 是否相似及是否合同（是填1;否填0）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "similarity_congruence_pair", "lambda_min": -6, "lambda_max": 6, "max_entry": 15}},
		},
		Derived: map[string]string{
			"B":    "sc_diag_matrix(A)",
			"sim":  "is_similar(A)",
			"cong": "is_congruent(A)",
		},
		Render: map[string]string{"A": "A", "B": "B"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "sim"},
			{ID: ids[1], Expr: "cong"},
		}},
		Meta: map[string]interface{}{
			"bank_topic":  "similarity_congruence",
			"solution_zh": `**解题思路**：两个方阵相似 $\Leftrightarrow$ 它们有相同的特征值（含重数）；两个实对称矩阵合同 $\Leftrightarrow$ 它们有相同的惯性指数（正、负特征值个数分别相等）。

**步骤 1**：计算 $A={{A}}$ 的特征多项式 $\det(\lambda E-A)$，求出 $A$ 的全部特征值 $\lambda_1,\lambda_2,\lambda_3$。由于 $A$ 由 similarity_congruence_pair 生成器产生，其特征值已存储在内部，可通过计算验证。

**步骤 2**：$B={{B}}$ 为对角阵，其特征值即对角元：$\mu_1={{expr:sc_diag_comp(B,1)}}$，$\mu_2={{expr:sc_diag_comp(B,2)}}$，$\mu_3={{expr:sc_diag_comp(B,3)}}$。

**步骤 3**：判断相似性。将 $A$ 的特征值集合与 $B$ 的对角元集合 $\{\mu_1,\mu_2,\mu_3\}$ 比较（考虑重数），若完全相同则相似（填 1），否则不相似（填 0）。本题结果：{{expr:is_similar(A)}}。

**步骤 4**：判断合同性。分别统计 $A$ 和 $B$ 的正特征值个数 $n_+$ 与负特征值个数 $n_-$。对 $B$ 而言，正对角元个数即 $n_+$，负对角元个数即 $n_-$。若两者的 $(n_+,n_-)$ 相同则合同（填 1），否则不合同（填 0）。本题结果：{{expr:is_congruent(A)}}。`,
		},
	}
}

// Ch6_3: 给定二次型 f(x₁,x₂,x₃)=...，用正交变换 x=Qy 化为标准型，
// 填写 3 个标准型系数（特征值）和 3×3 正交变换矩阵 Q 各元（共 12 空）。
func buildChapter6_3() dsl.Problem {
	k := "Chapter6_3"
	ids := BlankIDs(k, 12)
	fds := make([]dsl.AnswerFieldDef, 0, 12)
	for i := 1; i <= 3; i++ {
		fds = append(fds, dsl.AnswerFieldDef{
			ID:     ids[i-1],
			Expr:   fmt.Sprintf("sym_eigenval(S,%d)", i),
			Layout: dsl.LayoutVectorComponent("lambda", i, "eigenvalues"),
			Judge: &dsl.AnswerJudgeSpec{
				Kind: "eigen_pair", EigenGroup: "E", EigenRole: "lambda",
				EigenColumn: i, MatrixVar: "S",
			},
		})
	}
	// Q 的 9 个元素（按行填；每一列对应一个特征向量）
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
			`用正交变换化二次型 $f(x_1,x_2,x_3)={{expr}}$ 为标准型，填写三个标准型系数及正交变换矩阵 $Q$ 的 9 个元素（按行自上而下、从左到右；特征值与 $Q$ 的列顺序可任意对应，每一列允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"S": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "symmetric_eigen_reverse_3x3", "lambda_min": -6, "lambda_max": 6, "max_entry": 12}},
		},
		Derived: map[string]string{"expr": "quad_expr(S)"},
		Render:  map[string]string{"expr": "expr"},
		Answer:  dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"bank_topic":  "orthogonal_transform_standard_form",
			"solution_zh": `**解题思路**：二次型 $f(x_1,x_2,x_3)={{expr}}$ 的矩阵为对称阵 $S={{S}}$。正交变换 $x=Qy$ 化标准型 $\Leftrightarrow$ 对 $S$ 做正交相似对角化，标准型系数即 $S$ 的特征值，$Q$ 的列为对应的单位正交特征向量。

**步骤 1**：写出二次型的对称矩阵 $S$。$S_{ii}$ 为 $x_i^2$ 的系数，$S_{ij}=S_{ji}$ 为 $x_ix_j$ 系数的一半。本题 $S=\begin{pmatrix}{{expr:mget(S,1,1)}}&{{expr:mget(S,1,2)}}&{{expr:mget(S,1,3)}}\\{{expr:mget(S,2,1)}}&{{expr:mget(S,2,2)}}&{{expr:mget(S,2,3)}}\\{{expr:mget(S,3,1)}}&{{expr:mget(S,3,2)}}&{{expr:mget(S,3,3)}}\end{pmatrix}$。

**步骤 2**：求 $S$ 的特征值。解特征方程 $\det(\lambda E-S)=0$，得到三个特征值 $\lambda_1={{expr:sym_eigenval(S,1)}}$，$\lambda_2={{expr:sym_eigenval(S,2)}}$，$\lambda_3={{expr:sym_eigenval(S,3)}}$。这三个值即为标准型 $f=\lambda_1 y_1^2+\lambda_2 y_2^2+\lambda_3 y_3^2$ 的系数。

**步骤 3**：对每个特征值 $\lambda_i$，解齐次线性方程组 $(\lambda_i E-S)x=0$，得到对应的特征向量。由于 $S$ 为实对称阵，不同特征值对应的特征向量自动正交。

**步骤 4**：将每个特征向量单位化（除以模长 $\|\xi_i\|=\sqrt{\xi_{i1}^2+\xi_{i2}^2+\xi_{i3}^2}$），得到单位特征向量 $\xi_1,\xi_2,\xi_3$。若存在重特征值，需对其特征空间内的向量先做 Schmidt 正交化再单位化。

**步骤 5**：以 $\xi_1,\xi_2,\xi_3$ 为列构造正交矩阵 $Q=(\xi_1,\xi_2,\xi_3)$，即 $Q=\begin{pmatrix}\xi_{11}&\xi_{21}&\xi_{31}\\\xi_{12}&\xi_{22}&\xi_{32}\\\xi_{13}&\xi_{23}&\xi_{33}\end{pmatrix}$。$Q$ 的 9 个元素按行自上而下、从左到右依次为 $Q_{11}={{expr:sym_eigenvec_comp(S,1,1)}},Q_{12}={{expr:sym_eigenvec_comp(S,2,1)}},Q_{13}={{expr:sym_eigenvec_comp(S,3,1)}},\ldots,Q_{33}={{expr:sym_eigenvec_comp(S,3,3)}}$。经正交变换 $x=Qy$ 后，$f=\lambda_1 y_1^2+\lambda_2 y_2^2+\lambda_3 y_3^2$。`,
		},
	}
}

// Ch6_4: 给定二次型 f，化为标准型，
// 填写 3 个标准型系数和 3×3 可逆变换矩阵 P 各元（共 12 空）。
// 当前用正交变换生成标准型系数（与配方法结果一致），
// 变换矩阵 P 取正交矩阵 Q（配方法可能有其他可逆矩阵，此处 Q 也是合法答案）。
func buildChapter6_4() dsl.Problem {
	k := "Chapter6_4"
	ids := BlankIDs(k, 12)
	fds := make([]dsl.AnswerFieldDef, 0, 12)
	for i := 1; i <= 3; i++ {
		fds = append(fds, dsl.AnswerFieldDef{
			ID:     ids[i-1],
			Expr:   fmt.Sprintf("sym_eigenval(S,%d)", i),
			Layout: dsl.LayoutVectorComponent("lambda", i, "eigenvalues"),
			Judge: &dsl.AnswerJudgeSpec{
				Kind: "eigen_pair", EigenGroup: "E", EigenRole: "lambda",
				EigenColumn: i, MatrixVar: "S",
			},
		})
	}
	for row := 1; row <= 3; row++ {
		for col := 1; col <= 3; col++ {
			idx := 3 + (row-1)*3 + col
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[idx-1],
				Expr:   fmt.Sprintf("sym_eigenvec_comp(S,%d,%d)", col, row),
				Layout: dsl.LayoutMatrixCell("P", row, col, 3, 3, "P"),
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
			`化二次型 $f(x_1,x_2,x_3)={{expr}}$ 为标准型，填写三个标准型系数及变换矩阵 $P$ 的 9 个元素（按行自上而下、从左到右；特征值与 $P$ 的列顺序可任意对应，每一列允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"S": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "symmetric_eigen_reverse_3x3", "lambda_min": -6, "lambda_max": 6, "max_entry": 12}},
		},
		Derived: map[string]string{"expr": "quad_expr(S)"},
		Render:  map[string]string{"expr": "expr"},
		Answer:  dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"bank_topic":  "standard_form_transform",
			"solution_zh": `**解题思路**：二次型 $f(x_1,x_2,x_3)={{expr}}$ 的标准型系数为其对称矩阵 $S={{S}}$ 的特征值，变换矩阵 $P$ 的列为对应的特征向量（需满足 $P^\top SP=\mathrm{diag}(\lambda_1,\lambda_2,\lambda_3)$）。本题采用正交变换法（也可用配方法，结果的标准型系数相同但变换矩阵可能不同）。

**步骤 1**：写出二次型的对称矩阵 $S$。$S_{ii}$ 为 $x_i^2$ 的系数，$S_{ij}=S_{ji}$ 为 $x_ix_j$ 系数的一半。本题 $S=\begin{pmatrix}{{expr:mget(S,1,1)}}&{{expr:mget(S,1,2)}}&{{expr:mget(S,1,3)}}\\{{expr:mget(S,2,1)}}&{{expr:mget(S,2,2)}}&{{expr:mget(S,2,3)}}\\{{expr:mget(S,3,1)}}&{{expr:mget(S,3,2)}}&{{expr:mget(S,3,3)}}\end{pmatrix}$。

**步骤 2**：求 $S$ 的特征值。解特征方程 $\det(\lambda E-S)=0$，得 $\lambda_1={{expr:sym_eigenval(S,1)}}$，$\lambda_2={{expr:sym_eigenval(S,2)}}$，$\lambda_3={{expr:sym_eigenval(S,3)}}$。这三个值即为标准型 $f=\lambda_1 y_1^2+\lambda_2 y_2^2+\lambda_3 y_3^2$ 的系数。

**步骤 3**：对每个 $\lambda_i$，解 $(\lambda_i E-S)x=0$ 求特征向量。分别得到属于 $\lambda_1,\lambda_2,\lambda_3$ 的特征向量。

**步骤 4**：将每个特征向量单位化，得到单位特征向量 $\xi_1,\xi_2,\xi_3$。若存在重特征值，先对特征空间内的向量做 Schmidt 正交化。

**步骤 5**：构造变换矩阵 $P=(\xi_1,\xi_2,\xi_3)$，其 9 个元素按行自上而下、从左到右依次为 $P_{11}={{expr:sym_eigenvec_comp(S,1,1)}},P_{12}={{expr:sym_eigenvec_comp(S,2,1)}},\ldots,P_{33}={{expr:sym_eigenvec_comp(S,3,3)}}$。此时 $P^\top SP=\mathrm{diag}(\lambda_1,\lambda_2,\lambda_3)$，即 $x=Py$ 将 $f$ 化为标准型 $f=\lambda_1 y_1^2+\lambda_2 y_2^2+\lambda_3 y_3^2$。`,
		},
	}
}

// Ch6_5: 给定含参数 t 的二次型，求使二次型正定的 t 的取值范围 (a < t < b)。
// 使用 Sylvester 主子式条件推导 t 的范围。
func buildChapter6_5() dsl.Problem {
	k := "Chapter6_5"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`若二次型 $f(x_1,x_2,x_3)={{expr}}$ 为正定二次型，则参数 $t$ 的取值范围是 $t>{{blank:%s}}$`,
			id),
		Variables: map[string]dsl.Variable{
			"S": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "sylvester_range", "entry_min": -5, "entry_max": 5}},
		},
		Derived: map[string]string{
			"expr":  "quad_expr_param(S)",
			"lower": "sylvester_lower(S)",
		},
		Render: map[string]string{"expr": "expr"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: id, Expr: "lower"},
		}},
		Meta: map[string]interface{}{
			"bank_topic":  "positive_definite_t_range",
			"solution_zh": `**解题思路**：实对称矩阵正定的充要条件是其各阶顺序主子式全部大于零（Sylvester 准则）。当参数 $t$ 仅出现在某一个对角元上时，正定条件只能给出 $t$ 的下界，即 $t$ 的取值范围为单侧无穷区间 $t>a$。

	**步骤 1**：由 $f(x_1,x_2,x_3)={{expr}}$ 写出对称矩阵 $S$。$S_{ii}$ 为 $x_i^2$ 的系数（含 $t$），$S_{ij}=S_{ji}$ 为 $x_ix_j$ 系数的一半。本题 $S={{S}}$，其中参数 $t$ 位于 $S_{33}$（$x_3^2$ 的系数）。

	**步骤 2**：计算各阶顺序主子式。
	- 一阶：$\Delta_1=S_{11}={{expr:mget(S,1,1)}}>0$（恒成立，与 $t$ 无关）。
	- 二阶：$\Delta_2=S_{11}S_{22}-S_{12}^2={{expr:mget(S,1,1)}}\times{{expr:mget(S,2,2)}}-{{expr:mget(S,1,2)}}^2>0$（恒成立，与 $t$ 无关）。
	- 三阶：$\Delta_3=\det(S)$。由于 $t$ 只出现在 $S_{33}$，行列式关于 $t$ 是线性的：$\Delta_3=\Delta_2\cdot t+C$，其中 $\Delta_2$ 是 $S_{33}$ 处的代数余子式（即二阶主子式），$C$ 是 $S_{33}=0$ 时的行列式值。

	**步骤 3**：令 $\Delta_3>0$，即 $\Delta_2\cdot t+C>0$。因为 $\Delta_2={{expr:mget(S,1,1)}}\times{{expr:mget(S,2,2)}}-{{expr:mget(S,1,2)}}^2>0$，所以 $t>\frac{-C}{\Delta_2}$。

	**步骤 4**：综合 $\Delta_1>0$、$\Delta_2>0$（均恒成立）和 $\Delta_3>0$（给出 $t$ 的下界），取交集得到 $t$ 的取值范围 $t>{{expr:sylvester_lower(S)}}$。`,
		},
	}
}

func buildChapter6_6() dsl.Problem {
	k := "Chapter6_6"
	ids := BlankIDs(k, 10)
	fds := make([]dsl.AnswerFieldDef, 0, 10)
	// t value (1 blank)
	fds = append(fds, dsl.AnswerFieldDef{
		ID:   ids[0],
		Expr: "param_t(S)",
	})
	// 3×3 orthogonal transformation matrix Q (9 blanks)
	for i := 1; i <= 3; i++ {
		jp := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: fmt.Sprintf("p%d", i)}
		for j := 1; j <= 3; j++ {
			idx := 1 + (i-1)*3 + j
			fds = append(fds, dsl.AnswerFieldDef{
				ID:     ids[idx-1],
				Expr:   fmt.Sprintf("param_eigenvec_comp(S,%d,%d)", j, i),
				Layout: dsl.LayoutMatrixCell("P", i, j, 3, 3, "P"),
				Judge:  jp,
			})
		}
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`若二次型 $f(x_1,x_2,x_3)={{expr}}$ 经正交变换 $x=Py$ 化为标准型 $f=%s y_1^2+%s y_2^2+%s y_3^2$，则参数 $t$={{blank:%s}}，正交变换矩阵 $P$ 的 9 个元素（按行自上而下、从左到右；每一列允许整体乘以非零整数不影响判分）：%s`,
			"{{lambda1}}", "{{lambda2}}", "{{lambda3}}", ids[0],
			joinBlankPlaceholders(ids[1:])),
		Variables: map[string]dsl.Variable{
			"S": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "param_orthogonal_diag", "lambda_min": -5, "lambda_max": 5, "max_entry": 10}},
		},
		Derived: map[string]string{
			"expr":    "quad_expr_param(S)",
			"lambda1": "param_eigenval(S,1)",
			"lambda2": "param_eigenval(S,2)",
			"lambda3": "param_eigenval(S,3)",
		},
		Render: map[string]string{"expr": "expr", "lambda1": "lambda1", "lambda2": "lambda2", "lambda3": "lambda3"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"bank_topic":  "parametric_orthogonal_transform",
			"solution_zh": `**解题思路**：二次型 $f(x_1,x_2,x_3)={{expr}}$ 的对称矩阵 $S$ 含参数 $t$。经正交变换 $x=Py$ 化为标准型 $f=\lambda_1 y_1^2+\lambda_2 y_2^2+\lambda_3 y_3^2$，其中 $\lambda_i$ 为 $S$ 的特征值。利用相似变换不改变矩阵的迹（$\mathrm{tr}(S)=\lambda_1+\lambda_2+\lambda_3$）可求出 $t$，再对确定的 $S$ 做正交对角化求 $P$。

**步骤 1**：由 $f$ 的表达式写出对称矩阵 $S$（其中含未知参数 $t$ 位于某个对角元）。本题 $S={{S}}$。

**步骤 2**：利用迹等式求 $t$。正交相似变换 $P^\top SP=\mathrm{diag}(\lambda_1,\lambda_2,\lambda_3)$ 不改变矩阵的迹，故 $\mathrm{tr}(S)=\lambda_1+\lambda_2+\lambda_3$。已知标准型系数 $\lambda_1={{lambda1}}$，$\lambda_2={{lambda2}}$，$\lambda_3={{lambda3}}$，因此 $\mathrm{tr}(S)={{lambda1}}+{{lambda2}}+{{lambda3}}$。设 $S$ 中含 $t$ 的对角元为 $S_{kk}$，则 $\mathrm{tr}(S)=S_{11}+S_{22}+S_{33}$（其中一项含 $t$），代入数值解方程得 $t={{expr:param_t(S)}}$。

**步骤 3**：将 $t={{expr:param_t(S)}}$ 代回 $S$，得到完全确定的对称矩阵。

**步骤 4**：对确定的 $S$ 求特征向量。分别解 $(\lambda_i E-S)x=0$（$i=1,2,3$），即 $(\lambda_1 E-S)x=0$、$(\lambda_2 E-S)x=0$、$(\lambda_3 E-S)x=0$，得到三个特征向量。

**步骤 5**：将每个特征向量单位化，得到单位正交特征向量 $\xi_1,\xi_2,\xi_3$。构造正交矩阵 $P=(\xi_1,\xi_2,\xi_3)$，其 9 个元素按行自上而下、从左到右依次为 $P_{11}={{expr:param_eigenvec_comp(S,1,1)}},P_{12}={{expr:param_eigenvec_comp(S,2,1)}},P_{13}={{expr:param_eigenvec_comp(S,3,1)}},\ldots,P_{33}={{expr:param_eigenvec_comp(S,3,3)}}$。`,
		},
	}
}
