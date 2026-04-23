package bank

import (
	"fmt"

	"github.com/neumathe/la-dsl/dsl"
)

func buildChapter3_5() dsl.Problem {
	k := "Chapter3_5"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $\alpha_1,\alpha_2,\alpha_3$ 为矩阵 $A={{A}}$ 的三列，若它们线性相关则填 1，否则填 0：{{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "rank3_coin", "min": -5, "max": 5}},
		},
		Derived: map[string]string{"flag": "dep3(A)"},
		Render:  map[string]string{"A": "A"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "flag"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：三个三维向量线性相关当且仅当它们构成的行列式为零。

**步骤 1**：将 $\alpha_1,\alpha_2,\alpha_3$ 作为列构成矩阵 $A=\begin{bmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\end{bmatrix}$。

**步骤 2**：计算行列式 $\det(A)={{expr:det(A)}}$。

**步骤 3**：若 $\det(A)=0$ 则线性相关（填 1），否则线性无关（填 0）。

**答案**：填 {{expr:dep3(A)}}（0=线性无关，1=线性相关）。`,
		},
	}
}

func buildChapter3_4() dsl.Problem {
	k := "Chapter3_4"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $\alpha_1,\alpha_2,\alpha_3$ 为 $M={{M}}$ 的三列（$\alpha_i\in\mathbb{R}^4$），线性相关填 1，否则填 0：{{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"M": {Kind: "matrix", Rows: 4, Cols: 3, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
		},
		Derived: map[string]string{"flag": "dep_cols(M)"},
		Render:  map[string]string{"M": "M"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "flag"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：三个四维向量线性相关当且仅当它们构成的 $4\times 3$ 矩阵的秩小于 3（列数）。

**步骤 1**：$\alpha_1,\alpha_2,\alpha_3$ 为 $M=\begin{bmatrix}{{expr:mget(M,1,1)}}&{{expr:mget(M,1,2)}}&{{expr:mget(M,1,3)}}\\{{expr:mget(M,2,1)}}&{{expr:mget(M,2,2)}}&{{expr:mget(M,2,3)}}\\{{expr:mget(M,3,1)}}&{{expr:mget(M,3,2)}}&{{expr:mget(M,3,3)}}\\{{expr:mget(M,4,1)}}&{{expr:mget(M,4,2)}}&{{expr:mget(M,4,3)}}\end{bmatrix}$ 的三列。

**步骤 2**：对 $M$ 做初等行变换化为行阶梯形，求得其秩 $\mathrm{rank}(M)={{expr:rank(M)}}$。

**步骤 3**：若 $\mathrm{rank}(M)<3$ 则线性相关（填 1），否则线性无关（填 0）。

**答案**：填 {{expr:dep_cols(M)}}（0=线性无关，1=线性相关）。`,
		},
	}
}

func buildChapter3_8() dsl.Problem {
	k := "Chapter3_8"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $A={{A}}$ 为四阶方阵，填写 $\mathrm{rank}(A)$：{{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
		},
		Derived: map[string]string{"r": "rank(A)"},
		Render:  map[string]string{"A": "A"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "r"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：矩阵的秩等于其行阶梯形中非零行的个数。

**步骤 1**：给定四阶方阵 $A=\begin{bmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}&{{expr:mget(A,1,4)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}&{{expr:mget(A,2,4)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}&{{expr:mget(A,3,4)}}\\{{expr:mget(A,4,1)}}&{{expr:mget(A,4,2)}}&{{expr:mget(A,4,3)}}&{{expr:mget(A,4,4)}}\end{bmatrix}$。

**步骤 2**：对 $A$ 做初等行变换化为行阶梯形。非零行的个数即为矩阵的秩。

**步骤 3**：$\mathrm{rank}(A)={{r}}$。`,
		},
	}
}

func buildChapter3_9() dsl.Problem {
	k := "Chapter3_9"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $A={{A}}$，填写矩阵的秩 $\mathrm{rank}(A)$：{{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "range", "min": -5, "max": 5}},
		},
		Derived: map[string]string{"r": "rank(A)"},
		Render:  map[string]string{"A": "A"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "r"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：矩阵的秩等于行阶梯形中非零行的个数。

**步骤 1**：$A=\begin{bmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}&{{expr:mget(A,1,4)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}&{{expr:mget(A,2,4)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}&{{expr:mget(A,3,4)}}\\{{expr:mget(A,4,1)}}&{{expr:mget(A,4,2)}}&{{expr:mget(A,4,3)}}&{{expr:mget(A,4,4)}}\end{bmatrix}$。

**步骤 2**：通过初等行变换将 $A$ 化为行阶梯形，统计非零行个数。

**步骤 3**：$\mathrm{rank}(A)={{r}}$。`,
		},
	}
}

func buildChapter3_10() dsl.Problem {
	k := "Chapter3_10"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $A=(\alpha_1,\alpha_2,\alpha_3)= {{A}}$，$B=(2\alpha_1+\alpha_2+3\alpha_3,\;-2\alpha_2+6\alpha_3,\;5\alpha_2+5\alpha_3)$，已知 $|A|={{detA}}$，填写 $|B|$：{{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "full_rank", "min": -3, "max": 3}},
			"C": {Kind: "matrix", Rows: 3, Cols: 3, Fixed: [][]interface{}{{2, 1, 3}, {0, -2, 6}, {0, 5, 5}}},
		},
		Derived: map[string]string{
			"detA": "det(A)",
			"detC": "det(C)",
			"detB": "scmul(detA,detC)",
		},
		Render:  map[string]string{"A": "A", "detA": "detA"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "detB"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：$B$ 的各列是 $A$ 各列的线性组合，即 $B=AC$，故 $|B|=|A|\cdot|C|$。

**步骤 1**：写出组合矩阵 $C$。因为 $B$ 的第 1 列为 $2\alpha_1+\alpha_2+3\alpha_3$，第 2 列为 $0\alpha_1-2\alpha_2+6\alpha_3$，第 3 列为 $0\alpha_1+5\alpha_2+5\alpha_3$，所以
$$C=\begin{bmatrix}2&0&0\\1&-2&5\\3&6&5\end{bmatrix}$$
（注意：组合系数作为列排列，即 $B=(\alpha_1,\alpha_2,\alpha_3)C=AC$。）

**步骤 2**：计算 $\det(C)={{expr:det(C)}}$。

**步骤 3**：已知 $|A|={{detA}}$，利用行列式乘法公式：
$$|B|=|A|\cdot|C|={{detA}}\times{{expr:det(C)}}={{detB}}$$

**答案**：$|B|={{detB}}$。`,
		},
	}
}

func buildChapter3_1() dsl.Problem {
	k := "Chapter3_1"
	ids := BlankIDs(k, 4)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $\alpha_1,\alpha_2,\alpha_3$ 为 $M={{M}}$ 的三列（$\alpha_i$ 取 $M$ 的第 $i$ 列），求 $-2\alpha_1-2\alpha_2-3\alpha_3$ 的坐标分量：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"M": {Kind: "matrix", Rows: 4, Cols: 3, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
			"c": {Kind: "vector", Size: 3, Fixed: []interface{}{-2, -2, -3}},
		},
		Derived: map[string]string{"v": "M * c"},
		Render:  map[string]string{"M": "M"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "v[1]"}, {ID: ids[1], Expr: "v[2]"}, {ID: ids[2], Expr: "v[3]"}, {ID: ids[3], Expr: "v[4]"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：向量的线性组合 $-2\alpha_1-2\alpha_2-3\alpha_3$ 等于矩阵 $M$ 乘以系数向量 $(-2,-2,-3)^T$。

**步骤 1**：$\alpha_1,\alpha_2,\alpha_3$ 为 $M$ 的三列：
$$\alpha_1=\begin{bmatrix}{{expr:mget(M,1,1)}}\\{{expr:mget(M,2,1)}}\\{{expr:mget(M,3,1)}}\\{{expr:mget(M,4,1)}}\end{bmatrix},\quad \alpha_2=\begin{bmatrix}{{expr:mget(M,1,2)}}\\{{expr:mget(M,2,2)}}\\{{expr:mget(M,3,2)}}\\{{expr:mget(M,4,2)}}\end{bmatrix},\quad \alpha_3=\begin{bmatrix}{{expr:mget(M,1,3)}}\\{{expr:mget(M,2,3)}}\\{{expr:mget(M,3,3)}}\\{{expr:mget(M,4,3)}}\end{bmatrix}$$

**步骤 2**：系数向量 $v=(-2,-2,-3)^T={{c}}$。

**步骤 3**：计算 $Mv$：
- 第 1 分量：$(-2)\times{{expr:mget(M,1,1)}}+(-2)\times{{expr:mget(M,1,2)}}+(-3)\times{{expr:mget(M,1,3)}}$
- 第 2 分量：$(-2)\times{{expr:mget(M,2,1)}}+(-2)\times{{expr:mget(M,2,2)}}+(-3)\times{{expr:mget(M,2,3)}}$
- 第 3 分量：$(-2)\times{{expr:mget(M,3,1)}}+(-2)\times{{expr:mget(M,3,2)}}+(-3)\times{{expr:mget(M,3,3)}}$
- 第 4 分量：$(-2)\times{{expr:mget(M,4,1)}}+(-2)\times{{expr:mget(M,4,2)}}+(-3)\times{{expr:mget(M,4,3)}}$

**步骤 4**：结果为 $v={{v}}$。

**答案**：四个分量依次为 {{Chapter3_1_1}}, {{Chapter3_1_2}}, {{Chapter3_1_3}}, {{Chapter3_1_4}}。`,
		},
	}
}

func buildChapter3_2() dsl.Problem {
	k := "Chapter3_2"
	ids := BlankIDs(k, 4)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $\alpha_1,\alpha_2,\alpha_3$ 为 $M={{M}}$ 的三列（$\alpha_i$ 取 $M$ 的第 $i$ 列），求满足 $4(\alpha-\alpha_1)+3(\alpha+\alpha_2)=4\alpha+4\alpha_3$ 的 $\alpha$（填写各分量）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"M": {Kind: "matrix", Rows: 4, Cols: 3, Generator: map[string]interface{}{"rule": "from_set", "set": []interface{}{-9, -6, -3, 0, 3, 6, 9}}},
			"c": {Kind: "vector", Size: 3, Fixed: []interface{}{4, -3, 4}},
		},
		Derived: map[string]string{"v3": "M * c", "alpha": "vecdiv(v3,3)"},
		Render:  map[string]string{"M": "M"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "alpha[1]"}, {ID: ids[1], Expr: "alpha[2]"}, {ID: ids[2], Expr: "alpha[3]"}, {ID: ids[3], Expr: "alpha[4]"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：从向量方程化简出 $\alpha$ 的表达式，再用矩阵乘法计算。

**步骤 1**：化简方程 $4(\alpha-\alpha_1)+3(\alpha+\alpha_2)=4\alpha+4\alpha_3$：

展开得 $4\alpha-4\alpha_1+3\alpha+3\alpha_2=4\alpha+4\alpha_3$

合并得 $7\alpha-4\alpha_1+3\alpha_2=4\alpha+4\alpha_3$

移项得 $3\alpha=4\alpha_1-3\alpha_2+4\alpha_3$

故 $\alpha=\dfrac{1}{3}(4\alpha_1-3\alpha_2+4\alpha_3)$。

**步骤 2**：$\alpha_1,\alpha_2,\alpha_3$ 为 $M=\begin{bmatrix}{{expr:mget(M,1,1)}}&{{expr:mget(M,1,2)}}&{{expr:mget(M,1,3)}}\\{{expr:mget(M,2,1)}}&{{expr:mget(M,2,2)}}&{{expr:mget(M,2,3)}}\\{{expr:mget(M,3,1)}}&{{expr:mget(M,3,2)}}&{{expr:mget(M,3,3)}}\\{{expr:mget(M,4,1)}}&{{expr:mget(M,4,2)}}&{{expr:mget(M,4,3)}}\end{bmatrix}$ 的三列。

先计算 $M\cdot(4,-3,4)^T$：
- 第 1 分量：$4\times{{expr:mget(M,1,1)}}+(-3)\times{{expr:mget(M,1,2)}}+4\times{{expr:mget(M,1,3)}}$
- 第 2 分量：$4\times{{expr:mget(M,2,1)}}+(-3)\times{{expr:mget(M,2,2)}}+4\times{{expr:mget(M,2,3)}}$
- 第 3 分量：$4\times{{expr:mget(M,3,1)}}+(-3)\times{{expr:mget(M,3,2)}}+4\times{{expr:mget(M,3,3)}}$
- 第 4 分量：$4\times{{expr:mget(M,4,1)}}+(-3)\times{{expr:mget(M,4,2)}}+4\times{{expr:mget(M,4,3)}}$

**步骤 3**：得 $M\cdot(4,-3,4)^T={{v3}}$，故 $\alpha=\dfrac{1}{3}{{v3}}={{alpha}}$。

**答案**：四个分量依次为 {{Chapter3_2_1}}, {{Chapter3_2_2}}, {{Chapter3_2_3}}, {{Chapter3_2_4}}。`,
		},
	}
}

func buildChapter3_6() dsl.Problem {
	k := "Chapter3_6"
	ids := BlankIDs(k, 3)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $\alpha_1,\alpha_2,\alpha_3$ 为 $M={{M}}$ 的三列（$\alpha_i$ 取 $M$ 的第 $i$ 列），$\beta={{b}}$，将 $\beta$ 表示为 $\alpha_1,\alpha_2,\alpha_3$ 的线性组合，填写系数：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"M": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "upper_unit", "min": -4, "max": 4}},
			"b": {Kind: "vector", Size: 3, Generator: map[string]interface{}{"rule": "range", "min": -5, "max": 5}},
		},
		Derived: map[string]string{"c": "solve(M,b)"},
		Render:  map[string]string{"M": "M", "b": "b"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "c[1]"}, {ID: ids[1], Expr: "c[2]"}, {ID: ids[2], Expr: "c[3]"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：将 $\beta$ 表示为 $\alpha_1,\alpha_2,\alpha_3$ 的线性组合 $\beta=k_1\alpha_1+k_2\alpha_2+k_3\alpha_3$，等价于解线性方程组 $Mx=\beta$，其中 $M$ 的列为 $\alpha_1,\alpha_2,\alpha_3$。

**步骤 1**：$M=\begin{bmatrix}{{expr:mget(M,1,1)}}&{{expr:mget(M,1,2)}}&{{expr:mget(M,1,3)}}\\{{expr:mget(M,2,1)}}&{{expr:mget(M,2,2)}}&{{expr:mget(M,2,3)}}\\{{expr:mget(M,3,1)}}&{{expr:mget(M,3,2)}}&{{expr:mget(M,3,3)}}\end{bmatrix}$，$\beta={{b}}$。

**步骤 2**：解方程组 $Mx=\beta$，即求 $x=M^{-1}\beta$。

**步骤 3**：解得 $x={{c}}$。

**答案**：$k_1={{Chapter3_6_1}},\;k_2={{Chapter3_6_2}},\;k_3={{Chapter3_6_3}}$，即 $\beta={{Chapter3_6_1}}\alpha_1+{{Chapter3_6_2}}\alpha_2+{{Chapter3_6_3}}\alpha_3$。`,
		},
	}
}

func buildChapter3_7() dsl.Problem {
	k := "Chapter3_7"
	ids := BlankIDs(k, 4)
	j := &dsl.AnswerJudgeSpec{Kind: "sorted_basis_columns", BasisGroup: "bc", MatrixVar: "A", Ncols: 4}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $A={{A}}$ 的四列 $\alpha_1,\alpha_2,\alpha_3,\alpha_4$，写出一个极大线性无关列组下标（1–4，从小到大，不足补 0）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
		},
		Derived: map[string]string{"piv": "basis_cols(A)"},
		Render:  map[string]string{"A": "A"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "piv[1]", Judge: j},
			{ID: ids[1], Expr: "piv[2]", Judge: j},
			{ID: ids[2], Expr: "piv[3]", Judge: j},
			{ID: ids[3], Expr: "piv[4]", Judge: j},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：对 $A$ 做初等行变换化为行阶梯形，主元（pivot）所在的列即为极大线性无关列组。

**步骤 1**：$A=\begin{bmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}&{{expr:mget(A,1,4)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}&{{expr:mget(A,2,4)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}&{{expr:mget(A,3,4)}}\\{{expr:mget(A,4,1)}}&{{expr:mget(A,4,2)}}&{{expr:mget(A,4,3)}}&{{expr:mget(A,4,4)}}\end{bmatrix}$。

**步骤 2**：对 $A$ 做初等行变换化为行阶梯形，$\mathrm{rank}(A)={{expr:rank(A)}}$。

**步骤 3**：行阶梯形中主元所在列的下标即为极大线性无关列组的下标。

**答案**：四个空依次填写 {{Chapter3_7_1}}, {{Chapter3_7_2}}, {{Chapter3_7_3}}, {{Chapter3_7_4}}（下标为 0 的表示该位置无主元）。`,
		},
	}
}

func buildChapter3_3() dsl.Problem {
	k := "Chapter3_3"
	ids := BlankIDs(k, 9)
	fds := make([]dsl.AnswerFieldDef, 0, 9)
	for col := 1; col <= 3; col++ {
		jb := &dsl.AnswerJudgeSpec{Kind: "rational_line", LineGroup: fmt.Sprintf("b%d", col)}
		for row := 1; row <= 3; row++ {
			idx := (col-1)*3 + row - 1
			fds = append(fds, dsl.AnswerFieldDef{
				ID:    ids[idx],
				Expr:  fmt.Sprintf("gs_comp(V,%d,%d)", col, row),
				Judge: jb,
			})
		}
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $\alpha_1,\alpha_2,\alpha_3$ 为 $V={{V}}$ 的三列（$\alpha_i$ 取 $V$ 的第 $i$ 列，线性无关），对其做 Gram–Schmidt 正交化（不单位化），依次填写 $\beta_1,\beta_2,\beta_3$ 各分量（每个向量允许整体乘以非零整数不影响判分）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"V": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "full_rank", "min": -5, "max": 5}},
		},
		Render: map[string]string{"V": "V"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：Gram-Schmidt 正交化公式：
$$\beta_1=\alpha_1$$
$$\beta_2=\alpha_2-\frac{\langle\alpha_2,\beta_1\rangle}{\langle\beta_1,\beta_1\rangle}\beta_1$$
$$\beta_3=\alpha_3-\frac{\langle\alpha_3,\beta_1\rangle}{\langle\beta_1,\beta_1\rangle}\beta_1-\frac{\langle\alpha_3,\beta_2\rangle}{\langle\beta_2,\beta_2\rangle}\beta_2$$
其中 $\langle u,v\rangle=u_1v_1+u_2v_2+u_3v_3$ 为标准欧氏内积。

**步骤 1**：$\alpha_1,\alpha_2,\alpha_3$ 为 $V=\begin{bmatrix}{{expr:mget(V,1,1)}}&{{expr:mget(V,1,2)}}&{{expr:mget(V,1,3)}}\\{{expr:mget(V,2,1)}}&{{expr:mget(V,2,2)}}&{{expr:mget(V,2,3)}}\\{{expr:mget(V,3,1)}}&{{expr:mget(V,3,2)}}&{{expr:mget(V,3,3)}}\end{bmatrix}$ 的三列：
$$\alpha_1=\begin{bmatrix}{{expr:mget(V,1,1)}}\\{{expr:mget(V,2,1)}}\\{{expr:mget(V,3,1)}}\end{bmatrix},\quad \alpha_2=\begin{bmatrix}{{expr:mget(V,1,2)}}\\{{expr:mget(V,2,2)}}\\{{expr:mget(V,3,2)}}\end{bmatrix},\quad \alpha_3=\begin{bmatrix}{{expr:mget(V,1,3)}}\\{{expr:mget(V,2,3)}}\\{{expr:mget(V,3,3)}}\end{bmatrix}$$

**步骤 2**：$\beta_1=\alpha_1=\begin{bmatrix}{{expr:mget(V,1,1)}}\\{{expr:mget(V,2,1)}}\\{{expr:mget(V,3,1)}}\end{bmatrix}$。

**步骤 3**：计算 $\beta_2=\alpha_2-\dfrac{\langle\alpha_2,\beta_1\rangle}{\langle\beta_1,\beta_1\rangle}\beta_1$：
- $\langle\alpha_2,\beta_1\rangle = {{expr:mget(V,1,2)}}\times{{expr:mget(V,1,1)}}+{{expr:mget(V,2,2)}}\times{{expr:mget(V,2,1)}}+{{expr:mget(V,3,2)}}\times{{expr:mget(V,3,1)}}$
- $\langle\beta_1,\beta_1\rangle = {{expr:mget(V,1,1)}}^2+{{expr:mget(V,2,1)}}^2+{{expr:mget(V,3,1)}}^2$
- $\beta_2=\begin{bmatrix}{{expr:gs_comp(V,2,1)}}\\{{expr:gs_comp(V,2,2)}}\\{{expr:gs_comp(V,2,3)}}\end{bmatrix}$。

**步骤 4**：计算 $\beta_3=\alpha_3-\dfrac{\langle\alpha_3,\beta_1\rangle}{\langle\beta_1,\beta_1\rangle}\beta_1-\dfrac{\langle\alpha_3,\beta_2\rangle}{\langle\beta_2,\beta_2\rangle}\beta_2$：
- $\langle\alpha_3,\beta_1\rangle = {{expr:mget(V,1,3)}}\times{{expr:mget(V,1,1)}}+{{expr:mget(V,2,3)}}\times{{expr:mget(V,2,1)}}+{{expr:mget(V,3,3)}}\times{{expr:mget(V,3,1)}}$
- $\langle\alpha_3,\beta_2\rangle$ 为 $\alpha_3$ 与 $\beta_2$ 的标准内积
- $\beta_3=\begin{bmatrix}{{expr:gs_comp(V,3,1)}}\\{{expr:gs_comp(V,3,2)}}\\{{expr:gs_comp(V,3,3)}}\end{bmatrix}$。

**答案**：
$\beta_1=\begin{bmatrix}{{expr:gs_comp(V,1,1)}}\\{{expr:gs_comp(V,1,2)}}\\{{expr:gs_comp(V,1,3)}}\end{bmatrix}$,
$\beta_2=\begin{bmatrix}{{expr:gs_comp(V,2,1)}}\\{{expr:gs_comp(V,2,2)}}\\{{expr:gs_comp(V,2,3)}}\end{bmatrix}$,
$\beta_3=\begin{bmatrix}{{expr:gs_comp(V,3,1)}}\\{{expr:gs_comp(V,3,2)}}\\{{expr:gs_comp(V,3,3)}}\end{bmatrix}$。`,
		},
	}
}

func buildChapter3_11() dsl.Problem {
	k := "Chapter3_11"
	ids := BlankIDs(k, 24)
	jb := &dsl.AnswerJudgeSpec{Kind: "sorted_basis_columns", BasisGroup: "bc311", MatrixVar: "V", Ncols: 4}
	fds := []dsl.AnswerFieldDef{
		{ID: ids[0], Expr: "rank(V)"},
		{ID: ids[1], Expr: "ranklt(V,4)"},
		{ID: ids[2], Expr: "piv[1]", Judge: jb},
		{ID: ids[3], Expr: "piv[2]", Judge: jb},
		{ID: ids[4], Expr: "piv[3]", Judge: jb},
		{ID: ids[5], Expr: "piv[4]", Judge: jb},
		{ID: ids[6], Expr: "lhs4"},
		{ID: ids[7], Expr: "vcoef123(V,4,1)"},
		{ID: ids[8], Expr: "one1"},
		{ID: ids[9], Expr: "vcoef123(V,4,2)"},
		{ID: ids[10], Expr: "two2"},
		{ID: ids[11], Expr: "vcoef123(V,4,3)"},
		{ID: ids[12], Expr: "thr3"},
		{ID: ids[13], Expr: "zero()"},
		{ID: ids[14], Expr: "zero()"},
	}
	for i := 15; i < 24; i++ {
		fds = append(fds, dsl.AnswerFieldDef{ID: ids[i], Expr: "zero()"})
	}
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`设 $V={{V}}$ 的四列为 $\alpha_1,\ldots,\alpha_4\in\mathbb{R}^4$，且 $\mathrm{rank}(V)=3$、$\alpha_4$ 可由 $\alpha_1,\alpha_2,\alpha_3$ 线性表示。填写秩、相关性、极大无关组列标及 $\alpha_4$ 的表示系数（其余空填 0）：%s`,
			joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"V":     {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "rank334_last_dep", "min": -4, "max": 4}},
			"lhs4":  {Kind: "scalar", Fixed: 4},
			"one1":  {Kind: "scalar", Fixed: 1},
			"two2":  {Kind: "scalar", Fixed: 2},
			"thr3":  {Kind: "scalar", Fixed: 3},
		},
		Derived: map[string]string{"piv": "basis_cols(V)"},
		Render:  map[string]string{"V": "V"},
		Answer:  dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：本题综合运用矩阵的秩、向量组的线性相关性、极大线性无关组、线性表示四个核心概念。

**步骤 1 — 求秩**：$V=\begin{bmatrix}{{expr:mget(V,1,1)}}&{{expr:mget(V,1,2)}}&{{expr:mget(V,1,3)}}&{{expr:mget(V,1,4)}}\\{{expr:mget(V,2,1)}}&{{expr:mget(V,2,2)}}&{{expr:mget(V,2,3)}}&{{expr:mget(V,2,4)}}\\{{expr:mget(V,3,1)}}&{{expr:mget(V,3,2)}}&{{expr:mget(V,3,3)}}&{{expr:mget(V,3,4)}}\\{{expr:mget(V,4,1)}}&{{expr:mget(V,4,2)}}&{{expr:mget(V,4,3)}}&{{expr:mget(V,4,4)}}\end{bmatrix}$，对 $V$ 做初等行变换化为行阶梯形，得 $\mathrm{rank}(V)={{expr:rank(V)}}$。

**步骤 2 — 判定相关性**：因为 $\mathrm{rank}(V)={{expr:rank(V)}}<4$（列数），所以 $\alpha_1,\alpha_2,\alpha_3,\alpha_4$ 线性相关（填 {{Chapter3_11_2}}，即 1=相关）。

**步骤 3 — 极大无关组**：对 $V$ 做行变换找主元列，极大无关组列下标为 {{Chapter3_11_3}}, {{Chapter3_11_4}}, {{Chapter3_11_5}}, {{Chapter3_11_6}}（下标 0 表示该位置无对应主元）。

**步骤 4 — 线性表示**：$\alpha_4$ 可由极大无关组线性表示，设 $\alpha_4=k_1\alpha_i+k_2\alpha_j+k_3\alpha_k$（$i,j,k$ 为极大无关组的列下标），表示系数见第 8、10、12 空。

**答案汇总**：
- 第 1 空（秩）：{{Chapter3_11_1}}
- 第 2 空（相关性）：{{Chapter3_11_2}}
- 第 3-6 空（极大无关组下标）：{{Chapter3_11_3}}, {{Chapter3_11_4}}, {{Chapter3_11_5}}, {{Chapter3_11_6}}
- 第 7-12 空（$\alpha_4$ 表示式）：第 {{Chapter3_11_7}} 列、系数 {{Chapter3_11_8}}，第 {{Chapter3_11_9}} 列、系数 {{Chapter3_11_10}}，第 {{Chapter3_11_11}} 列、系数 {{Chapter3_11_12}}
- 第 13-24 空：其余位置填 0。`,
		},
	}
}
