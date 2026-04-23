package bank

import (
	"fmt"

	"github.com/neumathe/la-dsl/dsl"
)

func invTitlePlaceholders(key string, n int) string {
	ids := BlankIDs(key, n*n)
	s := ""
	for _, id := range ids {
		s += fmt.Sprintf(` ${{blank:%s}}$`, id)
	}
	return s
}

func buildChapter2_4_1() dsl.Problem {
	k := "Chapter2_4_1"
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设矩阵 $A={{A}}$，则 $A^{-1}$ 为：%s`, invTitlePlaceholders(k, 3)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "lower_unit", "min": -4, "max": 4}},
		},
		Derived: map[string]string{"Inv": "inv(A)"},
		Render:  map[string]string{"A": "A"},
		Answer:  dsl.AnswerSchema{FieldDefs: MatrixFieldDefsSquare(k, 3, "Inv", "A^{-1}")},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：对三阶下三角单位矩阵 $A$，用初等行变换法 $(A|E)\to(E|A^{-1})$ 求逆矩阵。

**步骤 1**：写出矩阵 $A=\begin{bmatrix}a_{11}&a_{12}&a_{13}\\a_{21}&a_{22}&a_{23}\\a_{31}&a_{32}&a_{33}\end{bmatrix}=\begin{bmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\end{bmatrix}$，构造增广矩阵 $(A|E)$。

**步骤 2**：因 $A$ 为下三角单位矩阵（对角元全为 1），变换过程为从第 2 行开始逐行消去下三角非零元，再将右半部分对应的初等变换结果保留。具体地：
- 用第 1 行消第 2、3 行第 1 列：第 2 行加 $-a_{21}=-{{expr:mget(A,2,1)}}$ 倍第 1 行，第 3 行加 $-a_{31}=-{{expr:mget(A,3,1)}}$ 倍第 1 行；
- 用第 2 行消第 3 行第 2 列：第 3 行加 $-a_{32}'$ 倍第 2 行（$a_{32}'$ 为第 1 步消去后的值）。

**步骤 3**：左半部分化为 $E$ 后，右半部分即为 $A^{-1}={{Inv}}$。

**步骤 4**：验证 $AA^{-1}={{expr:matmul(A,Inv)}}=I_3$（单位矩阵），确认求逆正确。`,
		},
	}
}

func buildChapter2_4_2() dsl.Problem {
	k := "Chapter2_4_2"
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $A={{A}}$，则 $A^{-1}$ 为：%s`, invTitlePlaceholders(k, 4)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "lower_unit", "min": -3, "max": 3}},
		},
		Derived: map[string]string{"Inv": "inv(A)"},
		Render:  map[string]string{"A": "A"},
		Answer:  dsl.AnswerSchema{FieldDefs: MatrixFieldDefsSquare(k, 4, "Inv", "A^{-1}")},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：四阶下三角单位矩阵求逆——初等行变换法 $(A|E)\to(E|A^{-1})$。

**步骤 1**：写出矩阵 $A=\begin{bmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}&{{expr:mget(A,1,4)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}&{{expr:mget(A,2,4)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}&{{expr:mget(A,3,4)}}\\{{expr:mget(A,4,1)}}&{{expr:mget(A,4,2)}}&{{expr:mget(A,4,3)}}&{{expr:mget(A,4,4)}}\end{bmatrix}$，构造 $4\times 8$ 增广矩阵 $(A|E)$。

**步骤 2**：由于 $A$ 为下三角单位矩阵（对角元全为 1，上三角全为 0），消元过程自下而上：
- 先消第 4 列下三角：第 4 行已是主元行，无需消；
- 消第 3 列：第 4 行加 $-a_{43}=-{{expr:mget(A,4,3)}}$ 倍第 3 行，消去 $(4,3)$ 元；
- 消第 2 列：第 4 行加适当倍数消 $(4,2)$ 元，第 3 行加适当倍数消 $(3,2)$ 元；
- 消第 1 列：第 4、3、2 行分别加适当倍数消去第 1 列非对角元。

**步骤 3**：左半化为 $I_4$ 后，右半即为 $A^{-1}={{Inv}}$。

**步骤 4**：验证 $AA^{-1}={{expr:matmul(A,Inv)}}=I_4$。`,
		},
	}
}

func buildChapter2_4_3() dsl.Problem {
	k := "Chapter2_4_3"
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $A={{A}}$，则 $A^{-1}$ 为：%s`, invTitlePlaceholders(k, 4)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "lower_unit", "min": -4, "max": 4}},
		},
		Derived: map[string]string{"Inv": "inv(A)"},
		Render:  map[string]string{"A": "A"},
		Answer:  dsl.AnswerSchema{FieldDefs: MatrixFieldDefsSquare(k, 4, "Inv", "A^{-1}")},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：四阶下三角单位矩阵求逆——初等行变换法 $(A|E)\to(E|A^{-1})$。

**步骤 1**：写出矩阵 $A=\begin{bmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}&{{expr:mget(A,1,4)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}&{{expr:mget(A,2,4)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}&{{expr:mget(A,3,4)}}\\{{expr:mget(A,4,1)}}&{{expr:mget(A,4,2)}}&{{expr:mget(A,4,3)}}&{{expr:mget(A,4,4)}}\end{bmatrix}$，构造增广矩阵 $(A|E)$。

**步骤 2**：$A$ 为下三角单位矩阵，逐列从下到上消去非对角元：
- 第 4 行作为主元行，消去第 3、2、1 行对应列的元素；
- 依次处理第 3、2 列，各主元行消去上方行对应元素。

**步骤 3**：左半化为 $I_4$ 后右半部分即为 $A^{-1}={{Inv}}$。

**步骤 4**：验证 $AA^{-1}={{expr:matmul(A,Inv)}}=I_4$，确认结果正确。`,
		},
	}
}

func buildChapter2_7_1() dsl.Problem {
	k := "Chapter2_7_1"
	return elemPAequalsC(k,
		"E", [][]interface{}{{1, 0, 0}, {0, 1, 1}, {0, 0, 1}},
		"S", [][]interface{}{{0, 1, 0}, {1, 0, 0}, {0, 0, 1}},
		fmt.Sprintf(`设三阶矩阵 $A={{A}}$，交换 $A$ 的第 2 行和第 1 行得到矩阵 $B$，将矩阵 $B$ 的第 3 行的 1 倍加到第 2 行得到矩阵 $C$，则满足 $PA=C$ 的矩阵 $P$：%s`, invTitlePlaceholders(k, 3)),
		map[string]interface{}{
			"solution_zh": `**解题思路**：每步行变换对应一个初等矩阵，$PA=C$ 中 $P$ 是各步初等矩阵之积（后做变换在左）。

**步骤 1**：第一步"交换第 1 行和第 2 行"对应初等矩阵 $S=\begin{bmatrix}0&1&0\\1&0&0\\0&0&1\end{bmatrix}$，即 $B=SA$。

**步骤 2**：第二步"将第 3 行的 1 倍加到第 2 行"对应初等矩阵 $E=\begin{bmatrix}1&0&0\\0&1&1\\0&0&1\end{bmatrix}$，即 $C=EB=ESA$。

**步骤 3**：因此 $P=E\cdot S={{P}}$（先做的变换 $S$ 在右侧）。

**步骤 4**：验证 $PA={{expr:matmul(P,A)}}={{C}}$，与题中 $C$ 一致，求 $P$ 正确。`,
		},
	)
}

func buildChapter2_7_2() dsl.Problem {
	k := "Chapter2_7_2"
	return elemPAequalsC(k,
		"D", [][]interface{}{{3, 0, 0}, {0, 1, 0}, {0, 0, 1}},
		"S", [][]interface{}{{1, 0, 0}, {0, 0, 1}, {0, 1, 0}},
		fmt.Sprintf(`设三阶矩阵 $A={{A}}$，交换 $A$ 的第 2 行和第 3 行得到矩阵 $B$，将矩阵 $B$ 的第 1 行乘以数 3 得到矩阵 $C$，则满足 $PA=C$ 的矩阵 $P$：%s`, invTitlePlaceholders(k, 3)),
		map[string]interface{}{
			"solution_zh": `**解题思路**：每步行变换对应一个初等矩阵，$P$ 是各步初等矩阵之积。

**步骤 1**："交换第 2 行和第 3 行"对应初等矩阵 $S=\begin{bmatrix}1&0&0\\0&0&1\\0&1&0\end{bmatrix}$，即 $B=SA$。

**步骤 2**："将第 1 行乘以 3"对应初等矩阵 $D=\begin{bmatrix}3&0&0\\0&1&0\\0&0&1\end{bmatrix}$，即 $C=DB=DSA$。

**步骤 3**：因此 $P=D\cdot S={{P}}$（先做 $S$ 在右侧）。

**步骤 4**：验证 $PA={{expr:matmul(P,A)}}={{C}}$，与题中 $C$ 一致。`,
		},
	)
}

func buildChapter2_7_3() dsl.Problem {
	k := "Chapter2_7_3"
	return elemPAequalsC(k,
		"E2", [][]interface{}{{1, 0, 0}, {0, -5, 0}, {0, 0, 1}},
		"E1", [][]interface{}{{1, 0, 0}, {0, 1, 0}, {4, 0, 1}},
		fmt.Sprintf(`设三阶矩阵 $A={{A}}$，将矩阵 $A$ 的第 1 行的 4 倍加到第 3 行得到矩阵 $B$，将矩阵 $B$ 的第 2 行乘以数 $-5$ 得到矩阵 $C$，则满足 $PA=C$ 的矩阵 $P$：%s`, invTitlePlaceholders(k, 3)),
		map[string]interface{}{
			"solution_zh": `**解题思路**：每步行变换对应一个初等矩阵，$P$ 是各步初等矩阵之积。

**步骤 1**："将第 1 行的 4 倍加到第 3 行"对应初等矩阵 $E_1=\begin{bmatrix}1&0&0\\0&1&0\\4&0&1\end{bmatrix}$，即 $B=E_1A$。

**步骤 2**："将第 2 行乘以 $-5$"对应初等矩阵 $E_2=\begin{bmatrix}1&0&0\\0&-5&0\\0&0&1\end{bmatrix}$，即 $C=E_2B=E_2E_1A$。

**步骤 3**：因此 $P=E_2\cdot E_1={{P}}$（先做 $E_1$ 在右侧）。

**步骤 4**：验证 $PA={{expr:matmul(P,A)}}={{C}}$，与题中 $C$ 一致。`,
		},
	)
}

func elemPAequalsC(key, leftName string, left [][]interface{}, rightName string, right [][]interface{}, title string, meta map[string]interface{}) dsl.Problem {
	return dsl.Problem{
		ID:      ProblemID(key),
		Version: "bank-v1",
		Title:   title,
		Variables: map[string]dsl.Variable{
			leftName:  {Kind: "matrix", Rows: 3, Cols: 3, Fixed: left},
			rightName: {Kind: "matrix", Rows: 3, Cols: 3, Fixed: right},
			"A":       {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "full_rank", "min": -5, "max": 5}},
		},
		Derived: map[string]string{
			"P": fmt.Sprintf("matmul(%s,%s)", leftName, rightName),
			"C": "matmul(P,A)",
		},
		Render: map[string]string{"A": "A"},
		Answer: dsl.AnswerSchema{FieldDefs: MatrixFieldDefsSquare(key, 3, "P", "P")},
		Meta:   meta,
	}
}

func buildChapter2_6() dsl.Problem {
	k := "Chapter2_6"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $A={{A}}$，$B$ 满足 $AB=3A-B$，填写 $\det(B)$：{{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"U": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "upper_unit", "min": -3, "max": 3}},
			"L": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "lower_unit", "min": -3, "max": 3}},
			"I": {Kind: "matrix", Rows: 3, Cols: 3, Fixed: [][]interface{}{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}},
		},
		Derived: map[string]string{
			"ApI":    "matmul(U,L)",
			"A":      "matsub(ApI,I)",
			"invApI": "inv(ApI)",
			"3A":     "smmul(3,A)",
			"B":      "matmul(invApI,3A)",
			"detB":   "det(B)",
		},
		Render: map[string]string{"A": "A"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "detB"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：将矩阵方程 $AB=3A-B$ 化为 $(A+I)B=3A$，左乘 $(A+I)^{-1}$ 解出 $B$，再取行列式。

**步骤 1**：$AB=3A-B$ 移项得 $AB+B=3A$，提取 $B$：$(A+I)B=3A$。

**步骤 2**：计算 $A+I={{ApI}}$，验证 $\det(A+I)={{expr:det(ApI)}}\neq 0$，故 $(A+I)$ 可逆。

**步骤 3**：$(A+I)^{-1}={{invApI}}$，$3A={{3A}}$。

**步骤 4**：$B=(A+I)^{-1}\cdot 3A={{expr:matmul(invApI,3A)}}={{B}}$。

**步骤 5**：$\det(B)={{expr:det(B)}}={{detB}}$，即为所求。`,
		},
	}
}

func buildChapter2_2() dsl.Problem {
	k := "Chapter2_2"
	ids := BlankIDs(k, 8)
	fds := append(
		MatrixFieldDefsRectIDs(ids, 0, 2, 2, "A2", "A^2"),
		MatrixFieldDefsRectIDs(ids, 4, 2, 2, "A5", "A^5")...,
	)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $A={{A}}$，求 $A^2$ 与 $A^5$：%s`, joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 2, Cols: 2, Generator: map[string]interface{}{"rule": "full_rank", "min": -3, "max": 3}},
		},
		Derived: map[string]string{"A2": "pow(A,2)", "A5": "pow(A,5)"},
		Render:  map[string]string{"A": "A"},
		Answer:  dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：矩阵幂通过递推乘法计算：$A^2=A\cdot A$，$A^n=A^{n-1}\cdot A$。

**步骤 1**：$A=\begin{bmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}\end{bmatrix}$。计算 $A^2=A\cdot A$：
- $(A^2)_{11}=a_{11}^2+a_{12}a_{21}={{expr:mget(A,1,1)}}\times{{expr:mget(A,1,1)}}+{{expr:mget(A,1,2)}}\times{{expr:mget(A,2,1)}}={{expr:mget(A2,1,1)}}$；
- $(A^2)_{12}=a_{11}a_{12}+a_{12}a_{22}={{expr:mget(A,1,1)}}\times{{expr:mget(A,1,2)}}+{{expr:mget(A,1,2)}}\times{{expr:mget(A,2,2)}}={{expr:mget(A2,1,2)}}$；
- $(A^2)_{21}=a_{21}a_{11}+a_{22}a_{21}={{expr:mget(A,2,1)}}\times{{expr:mget(A,1,1)}}+{{expr:mget(A,2,2)}}\times{{expr:mget(A,2,1)}}={{expr:mget(A2,2,1)}}$；
- $(A^2)_{22}=a_{21}a_{12}+a_{22}^2={{expr:mget(A,2,1)}}\times{{expr:mget(A,1,2)}}+{{expr:mget(A,2,2)}}\times{{expr:mget(A,2,2)}}={{expr:mget(A2,2,2)}}$。
故 $A^2={{A2}}$。

**步骤 2**：递推计算 $A^3=A^2\cdot A={{expr:pow(A,3)}}$，$A^4=A^3\cdot A={{expr:pow(A,4)}}$，$A^5=A^4\cdot A={{A5}}$。

**步骤 3**：最终 $A^2={{A2}}$，$A^5={{A5}}$。`,
		},
	}
}

func buildChapter2_3() dsl.Problem {
	k := "Chapter2_3"
	ids := BlankIDs(k, 8)
	fds := append(
		MatrixFieldDefsRectIDs(ids, 0, 2, 2, "A8", `(P^{-1}AP)^8`),
		MatrixFieldDefsRectIDs(ids, 4, 2, 2, "A7", "A^7")...,
	)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $A={{A}}$，$P={{P}}$，求 $(P^{-1}AP)^8$ 与 $A^7$：%s`, joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A":  {Kind: "matrix", Rows: 2, Cols: 2, Generator: map[string]interface{}{"rule": "diagonalizable_2x2", "lambda_min": -4, "lambda_max": 4, "entry_min": -3, "entry_max": 3}},
			"PU": {Kind: "matrix", Rows: 2, Cols: 2, Generator: map[string]interface{}{"rule": "upper_unit", "min": -4, "max": 4}},
			"PL": {Kind: "matrix", Rows: 2, Cols: 2, Generator: map[string]interface{}{"rule": "lower_unit", "min": -4, "max": 4}},
		},
		Derived: map[string]string{
			"A7":         "pow(A,7)",
			"A8":         "pow(A,8)",
			"P":          "matmul(PU,PL)",
			"InvP":       "inv(P)",
			"A8P":        "matmul(A8,P)",
			"P_A8P":      "matmul(P,A8P)",
			"InvP_A8P":   "matmul(InvP,A8P)",
			"P_inv_A8_P": "InvP_A8P",
			"P_A8_P_inv": "P_A8P",
		},
		Render:  map[string]string{"A": "A", "P": "P"},
		Answer:  dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：利用相似矩阵性质 $(P^{-1}AP)^n=P^{-1}A^nP$，先算 $A^n$ 再做相似变换。

**步骤 1**：先递推计算 $A$ 的各次幂。$A=\begin{bmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}\end{bmatrix}$，$P=\begin{bmatrix}{{expr:mget(P,1,1)}}&{{expr:mget(P,1,2)}}\\{{expr:mget(P,2,1)}}&{{expr:mget(P,2,2)}}\end{bmatrix}$。

**步骤 2**：逐次计算 $A^2,A^3,\ldots,A^8={{A8}}$。

**步骤 3**：求 $P^{-1}={{InvP}}$，然后 $(P^{-1}AP)^8=P^{-1}A^8P={{P_inv_A8_P}}$。

**步骤 4**：递推计算 $A^7={{A7}}$（$A^7=A^8\cdot A^{-1}$ 或从 $A^6$ 递推）。

**步骤 5**：$(P^{-1}AP)^8={{P_inv_A8_P}}$，$A^7={{A7}}$。`,
		},
	}
}

func buildChapter2_1() dsl.Problem {
	k := "Chapter2_1"
	ids := BlankIDs(k, 16)
	fds := append(
		MatrixFieldDefsRectIDs(ids, 0, 4, 2, "AB", "AB"),
		MatrixFieldDefsRectIDs(ids, 8, 4, 2, "U", "AB+AC^T")...,
	)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $A={{A}}$，$B={{B}}$，$C={{C}}$，求 $AB$ 与 $AB+AC^T$：%s`, joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 3, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
			"B": {Kind: "matrix", Rows: 3, Cols: 2, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
			"C": {Kind: "matrix", Rows: 2, Cols: 3, Generator: map[string]interface{}{"rule": "range", "min": -4, "max": 4}},
		},
		Derived: map[string]string{
			"AB": "matmul(A,B)",
			"tC": "transpose(C)",
			"AC": "matmul(A,tC)",
			"U":  "matadd(AB,AC)",
		},
		Render: map[string]string{"A": "A", "B": "B", "C": "C"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：矩阵乘法按行·列内积计算——$(AB)_{ij}=\sum_k a_{ik}b_{kj}$；加法对应位置相加；转置互换行和列。

**步骤 1**：$A$（$4\times 3$）乘 $B$（$3\times 2$）得 $AB$（$4\times 2$）。以第 1 行为例：
- $(AB)_{11}=a_{11}b_{11}+a_{12}b_{21}+a_{13}b_{31}={{expr:mget(A,1,1)}}\times{{expr:mget(B,1,1)}}+{{expr:mget(A,1,2)}}\times{{expr:mget(B,2,1)}}+{{expr:mget(A,1,3)}}\times{{expr:mget(B,3,1)}}={{expr:mget(AB,1,1)}}$；
- $(AB)_{12}=a_{11}b_{12}+a_{12}b_{22}+a_{13}b_{32}={{expr:mget(A,1,1)}}\times{{expr:mget(B,1,2)}}+{{expr:mget(A,1,2)}}\times{{expr:mget(B,2,2)}}+{{expr:mget(A,1,3)}}\times{{expr:mget(B,3,2)}}={{expr:mget(AB,1,2)}}$。
其余行类似计算，得到 $AB={{AB}}$。

**步骤 2**：$C^T={{tC}}$（$C$ 的行列互换，变为 $3\times 2$）。$A$（$4\times 3$）乘 $C^T$（$3\times 2$）得 $AC^T={{AC}}$（$4\times 2$）。

**步骤 3**：$AB+AC^T$ 为两 $4\times 2$ 矩阵对应位置相加：$(AB+AC^T)_{ij}=(AB)_{ij}+(AC^T)_{ij}$，结果为 $AB+AC^T={{U}}$。`,
		},
	}
}

func buildChapter2_5_2() dsl.Problem {
	k := "Chapter2_5_2"
	ids := BlankIDs(k, 21)
	fds := append(
		append(
			MatrixFieldDefsRectIDs(ids, 0, 3, 3, "Inv", "A^{-1}"),
			MatrixFieldDefsRectIDs(ids, 9, 3, 2, "D", "C-B")...,
		),
		MatrixFieldDefsRectIDs(ids, 15, 3, 2, "X", "X")...,
	)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`设 $A={{A}}$，$B={{B}}$，$C={{C}}$，求 $A^{-1}$、$C-B$ 及 $AX+B=C$ 的解 $X$：%s`, joinBlankPlaceholders(ids)),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "lower_unit", "min": -5, "max": 5}},
			"B": {Kind: "matrix", Rows: 3, Cols: 2, Generator: map[string]interface{}{"rule": "range", "min": -6, "max": 6}},
			"C": {Kind: "matrix", Rows: 3, Cols: 2, Generator: map[string]interface{}{"rule": "range", "min": -6, "max": 6}},
		},
		Derived: map[string]string{
			"Inv": "inv(A)",
			"D":   "matsub(C,B)",
			"X":   "matmul(Inv,D)",
		},
		Render: map[string]string{"A": "A", "B": "B", "C": "C"},
		Answer: dsl.AnswerSchema{FieldDefs: fds},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：矩阵方程 $AX+B=C$ 化为 $AX=C-B$，左乘 $A^{-1}$ 得 $X=A^{-1}(C-B)$。

**步骤 1**：$AX+B=C$ 移项得 $AX=C-B$。计算 $C-B$（对应位置相减）：
- $(C-B)_{ij}=c_{ij}-b_{ij}$，例如 $(C-B)_{11}={{expr:mget(C,1,1)}}-({{expr:mget(B,1,1)}})={{expr:mget(D,1,1)}}$。
结果为 $C-B={{D}}$（$3\times 2$ 矩阵）。

**步骤 2**：求 $A^{-1}$。$A=\begin{bmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\end{bmatrix}$ 为下三角单位矩阵，用初等行变换法 $(A|I)\to(I|A^{-1})$ 得 $A^{-1}={{Inv}}$。

**步骤 3**：$X=A^{-1}(C-B)={{expr:matmul(Inv,D)}}={{X}}$（$3\times 3$ 乘 $3\times 2$ 得 $3\times 2$）。`,
		},
	}
}

func joinBlankPlaceholders(ids []string) string {
	s := ""
	for _, id := range ids {
		s += fmt.Sprintf(` ${{blank:%s}}$`, id)
	}
	return s
}
