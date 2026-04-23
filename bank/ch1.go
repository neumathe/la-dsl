package bank

import (
	"fmt"

	"github.com/neumathe/la-dsl/dsl"
)

func buildChapter1_1() dsl.Problem {
	k := "Chapter1_1"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`三阶行列式 $D={{vA}}$ 等于 {{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "range", "min": -5, "max": 5}},
		},
		Derived: map[string]string{"d": "det(A)", "vA": "vmatrix_title(A)"},
		Render:  map[string]string{"vA": "vA"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "d"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：三阶行列式用按第一行展开（Laplace 展开法），每个代数余子式是一个二阶行列式。

**步骤 1**：记矩阵元素为
$$a_{11}={{expr:mget(A,1,1)}},\quad a_{12}={{expr:mget(A,1,2)}},\quad a_{13}={{expr:mget(A,1,3)}}$$
$$a_{21}={{expr:mget(A,2,1)}},\quad a_{22}={{expr:mget(A,2,2)}},\quad a_{23}={{expr:mget(A,2,3)}}$$
$$a_{31}={{expr:mget(A,3,1)}},\quad a_{32}={{expr:mget(A,3,2)}},\quad a_{33}={{expr:mget(A,3,3)}}$$

**步骤 2**：按第一行展开：
$$D = a_{11}A_{11} + a_{12}A_{12} + a_{13}A_{13}$$
其中 $A_{1j}=(-1)^{1+j}M_{1j}$，$M_{1j}$ 是删去第 1 行第 j 列后的二阶子行列式。

**步骤 3**：计算三个代数余子式：

- $A_{11}=(-1)^{1+1}\cdot\begin{vmatrix}a_{22}&a_{23}\\a_{32}&a_{33}\end{vmatrix}=(+1)\cdot\begin{vmatrix}{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\end{vmatrix}={{expr:mget(A,2,2)}}\times{{expr:mget(A,3,3)}}-{{expr:mget(A,2,3)}}\times{{expr:mget(A,3,2)}}={{expr:cofactor(A,1,1)}}$

- $A_{12}=(-1)^{1+2}\cdot\begin{vmatrix}a_{21}&a_{23}\\a_{31}&a_{33}\end{vmatrix}=(-1)\cdot\begin{vmatrix}{{expr:mget(A,2,1)}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,3)}}\end{vmatrix}=-({{expr:mget(A,2,1)}}\times{{expr:mget(A,3,3)}}-{{expr:mget(A,2,3)}}\times{{expr:mget(A,3,1)}})={{expr:cofactor(A,1,2)}}$

- $A_{13}=(-1)^{1+3}\cdot\begin{vmatrix}a_{21}&a_{22}\\a_{31}&a_{32}\end{vmatrix}=(+1)\cdot\begin{vmatrix}{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}\end{vmatrix}={{expr:mget(A,2,1)}}\times{{expr:mget(A,3,2)}}-{{expr:mget(A,2,2)}}\times{{expr:mget(A,3,1)}}={{expr:cofactor(A,1,3)}}$

**步骤 4**：代入求和：
$$D = {{expr:mget(A,1,1)}}\times({{expr:cofactor(A,1,1)}}) + {{expr:mget(A,1,2)}}\times({{expr:cofactor(A,1,2)}}) + {{expr:mget(A,1,3)}}\times({{expr:cofactor(A,1,3)}})$$

**步骤 5**：化简得 $D={{d}}$。`,
		},
	}
}

func buildChapter1_2() dsl.Problem {
	k := "Chapter1_2"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`四阶行列式 $D={{vA}}$ 等于 {{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "range", "min": -8, "max": 8}},
		},
		Derived: map[string]string{"d": "det(A)", "vA": "vmatrix_title(A)"},
		Render:  map[string]string{"vA": "vA"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "d"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：四阶行列式无 Sarrus 法则，必须按某行展开为三阶行列式再计算。这里按第一行展开。

**步骤 1**：记第一行元素为
$$a_{11}={{expr:mget(A,1,1)}},\quad a_{12}={{expr:mget(A,1,2)}},\quad a_{13}={{expr:mget(A,1,3)}},\quad a_{14}={{expr:mget(A,1,4)}}$$

**步骤 2**：按第一行展开：
$$D = a_{11}A_{11} + a_{12}A_{12} + a_{13}A_{13} + a_{14}A_{14}$$
其中 $A_{1j}=(-1)^{1+j}M_{1j}$，$M_{1j}$ 是删去第 1 行第 j 列后的三阶子行列式。

**步骤 3**：计算四个代数余子式（每个是一个三阶行列式，系统已用 Bareiss 算法计算）：

- $A_{11}=(-1)^{2}\cdot M_{11}=$ 删去第 1 行第 1 列后余下的三阶行列式 $={{expr:cofactor(A,1,1)}}$

- $A_{12}=(-1)^{3}\cdot M_{12}=-$（删去第 1 行第 2 列后余下的三阶行列式）$={{expr:cofactor(A,1,2)}}$

- $A_{13}=(-1)^{4}\cdot M_{13}=$ 删去第 1 行第 3 列后余下的三阶行列式 $={{expr:cofactor(A,1,3)}}$

- $A_{14}=(-1)^{5}\cdot M_{14}=-$（删去第 1 行第 4 列后余下的三阶行列式）$={{expr:cofactor(A,1,4)}}$

**步骤 4**：代入求和：
$$D = {{expr:mget(A,1,1)}}\times({{expr:cofactor(A,1,1)}}) + {{expr:mget(A,1,2)}}\times({{expr:cofactor(A,1,2)}}) + {{expr:mget(A,1,3)}}\times({{expr:cofactor(A,1,3)}}) + {{expr:mget(A,1,4)}}\times({{expr:cofactor(A,1,4)}})$$

**步骤 5**：化简得 $D={{d}}$。`,
		},
	}
}

func buildChapter1_3() dsl.Problem {
	k := "Chapter1_3"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`四阶行列式 $D={{vA}}$ 等于 {{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "upper_triangular", "min": -6, "max": 6}},
		},
		Derived: map[string]string{"d": "det(A)", "vA": "vmatrix_title(A)"},
		Render:  map[string]string{"vA": "vA"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "d"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：上三角矩阵（主对角线以下全为零）的行列式等于对角线元素之积。

**步骤 1**：验证矩阵为上三角矩阵：
$$A=\begin{pmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}&{{expr:mget(A,1,4)}}\\0&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}&{{expr:mget(A,2,4)}}\\0&0&{{expr:mget(A,3,3)}}&{{expr:mget(A,3,4)}}\\0&0&0&{{expr:mget(A,4,4)}}\end{pmatrix}$$
主对角线以下元素 $a_{21}={{expr:mget(A,2,1)}}, a_{31}={{expr:mget(A,3,1)}}, a_{32}={{expr:mget(A,3,2)}}, a_{41}={{expr:mget(A,4,1)}}, a_{42}={{expr:mget(A,4,2)}}, a_{43}={{expr:mget(A,4,3)}}$ 均为 0，确为上三角矩阵。

**步骤 2**：读取对角线元素：
$$a_{11}={{expr:mget(A,1,1)}},\quad a_{22}={{expr:mget(A,2,2)}},\quad a_{33}={{expr:mget(A,3,3)}},\quad a_{44}={{expr:mget(A,4,4)}}$$

**步骤 3**：上三角行列式公式：
$$\det(A) = a_{11}\cdot a_{22}\cdot a_{33}\cdot a_{44} = {{expr:mget(A,1,1)}}\times{{expr:mget(A,2,2)}}\times{{expr:mget(A,3,3)}}\times{{expr:mget(A,4,4)}}$$

**步骤 4**：逐步计算：
$${{expr:mget(A,1,1)}}\times{{expr:mget(A,2,2)}} = {{expr:mget(A,1,1)}}\times{{expr:mget(A,2,2)}}$$
$$\left({{expr:mget(A,1,1)}}\times{{expr:mget(A,2,2)}}\right)\times{{expr:mget(A,3,3)}} = {{expr:mget(A,1,1)}}\times{{expr:mget(A,2,2)}}\times{{expr:mget(A,3,3)}}$$
$$\left({{expr:mget(A,1,1)}}\times{{expr:mget(A,2,2)}}\times{{expr:mget(A,3,3)}}\right)\times{{expr:mget(A,4,4)}} = {{d}}$$

最终 $D={{d}}$。`,
		},
	}
}

func buildChapter1_4() dsl.Problem {
	k := "Chapter1_4"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`7阶行列式 $D_n={{vA}}$ 的值等于 {{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"A": {
				Kind: "matrix", Rows: 7, Cols: 7,
				Generator: map[string]interface{}{
					"rule": "equidiagonal", "diag_min": -2, "diag_max": 2, "off_min": -5, "off_max": -1,
				},
			},
		},
		Derived: map[string]string{"d": "det(A)", "vA": "equidiagonal_title(A)"},
		Render:  map[string]string{"vA": "vA"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "d"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：等对角线行列式（主对角线元素全相同、紧邻次对角线元素全相同）可用递推法。

**步骤 1**：识别矩阵结构。读取对角元和次对角元：
- 对角元（主对角线）$a = a_{11} = {{expr:mget(A,1,1)}}$（验证：$a_{22}={{expr:mget(A,2,2)}}, a_{33}={{expr:mget(A,3,3)}}$ 等均等于 ${{expr:mget(A,1,1)}}$）
- 次对角元（紧邻主对角线上下方）$b = a_{12} = {{expr:mget(A,1,2)}}$（验证：$a_{23}={{expr:mget(A,2,3)}}, a_{34}={{expr:mget(A,3,4)}}$ 等均等于 ${{expr:mget(A,1,2)}}$）
- 其余位置均为 0

**步骤 2**：$n$ 阶等对角线行列式的递推公式：
$$D_1 = a = {{expr:mget(A,1,1)}}$$
$$D_2 = a^2 - b^2 = {{expr:mget(A,1,1)}}^2 - {{expr:mget(A,1,2)}}^2 = {{expr:mget(A,1,1)}}\times{{expr:mget(A,1,1)}} - {{expr:mget(A,1,2)}}\times{{expr:mget(A,1,2)}}$$
$$D_n = a\cdot D_{n-1} - b^2\cdot D_{n-2}\quad(n\geq 3)$$

**步骤 3**：逐项递推计算：

- $D_3 = a\cdot D_2 - b^2\cdot D_1 = {{expr:mget(A,1,1)}}\cdot D_2 - {{expr:mget(A,1,2)}}^2\cdot{{expr:mget(A,1,1)}}$

- $D_4 = a\cdot D_3 - b^2\cdot D_2$

- $D_5 = a\cdot D_4 - b^2\cdot D_3$

- $D_6 = a\cdot D_5 - b^2\cdot D_4$

- $D_7 = a\cdot D_6 - b^2\cdot D_5$

**步骤 4**：经递推计算得 $D_7={{d}}$。`,
		},
	}
}

func buildChapter1_5() dsl.Problem {
	k := "Chapter1_5"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`六阶行列式 $D_6={{vA}}$ 等于 {{blank:%s}}`, id),
		Variables: map[string]dsl.Variable{
			"A": {
				Kind: "matrix", Rows: 6, Cols: 6,
				Generator: map[string]interface{}{"rule": "sparse", "values": []interface{}{-9, -8, -6, -5, -2, -1, 1, 2, 4, 5, 7, 8}, "density": 0.22},
			},
		},
		Derived: map[string]string{"d": "det(A)", "vA": "vmatrix_title(A)"},
		Render:  map[string]string{"vA": "vA"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "d"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：稀疏矩阵（大部分元素为零）的行列式计算，优先选择含零最多的行或列展开，使非零项最少，逐次降阶。

**步骤 1**：观察矩阵 $D_6={{vA}}$ 中零的分布，统计每行每列的非零元个数：

- 第 1 行非零元：$a_{11}={{expr:mget(A,1,1)}}, a_{12}={{expr:mget(A,1,2)}}, a_{13}={{expr:mget(A,1,3)}}, a_{14}={{expr:mget(A,1,4)}}, a_{15}={{expr:mget(A,1,5)}}, a_{16}={{expr:mget(A,1,6)}}$
- 第 2 行非零元：$a_{21}={{expr:mget(A,2,1)}}, a_{22}={{expr:mget(A,2,2)}}, a_{23}={{expr:mget(A,2,3)}}, a_{24}={{expr:mget(A,2,4)}}, a_{25}={{expr:mget(A,2,5)}}, a_{26}={{expr:mget(A,2,6)}}$
- 第 3 行非零元：$a_{31}={{expr:mget(A,3,1)}}, a_{32}={{expr:mget(A,3,2)}}, a_{33}={{expr:mget(A,3,3)}}, a_{34}={{expr:mget(A,3,4)}}, a_{35}={{expr:mget(A,3,5)}}, a_{36}={{expr:mget(A,3,6)}}$
- 第 4 行非零元：$a_{41}={{expr:mget(A,4,1)}}, a_{42}={{expr:mget(A,4,2)}}, a_{43}={{expr:mget(A,4,3)}}, a_{44}={{expr:mget(A,4,4)}}, a_{45}={{expr:mget(A,4,5)}}, a_{46}={{expr:mget(A,4,6)}}$
- 第 5 行非零元：$a_{51}={{expr:mget(A,5,1)}}, a_{52}={{expr:mget(A,5,2)}}, a_{53}={{expr:mget(A,5,3)}}, a_{54}={{expr:mget(A,5,4)}}, a_{55}={{expr:mget(A,5,5)}}, a_{56}={{expr:mget(A,5,6)}}$
- 第 6 行非零元：$a_{61}={{expr:mget(A,6,1)}}, a_{62}={{expr:mget(A,6,2)}}, a_{63}={{expr:mget(A,6,3)}}, a_{64}={{expr:mget(A,6,4)}}, a_{65}={{expr:mget(A,6,5)}}, a_{66}={{expr:mget(A,6,6)}}$

选取含零最多的行（或列）作为展开行，设其为第 $r$ 行。

**步骤 2**：按第 $r$ 行展开：
$$D_6 = \sum_{j=1}^{6} a_{rj}A_{rj}$$
其中只有 $a_{rj}\neq 0$ 的项需要实际计算，其余项为零。

**步骤 3**：对每个非零项对应的代数余子式 $A_{rj}$（五阶行列式），继续选取含零最多的行或列展开，将问题降为四阶行列式。

**步骤 4**：重复上述策略，逐次降阶：六阶 $\to$ 五阶 $\to$ 四阶 $\to$ 三阶 $\to$ 二阶，最终化为可计算的二阶行列式。

**步骤 5**：系统已按最优展开策略完成全部计算，最终结果 $D_6={{d}}$。`,
		},
	}
}

func buildChapter1_6() dsl.Problem {
	k := "Chapter1_6"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`行列式 $D={{vA}}$ 中代数余子式 $A_{24}={{blank:%s}}$`, id),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 4, Cols: 4, Generator: map[string]interface{}{"rule": "range", "min": -5, "max": 5}},
		},
		Derived: map[string]string{"A24": "cofactor(A,2,4)", "vA": "vmatrix_title(A)"},
		Render:  map[string]string{"vA": "vA"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "A24"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：代数余子式 $A_{ij}=(-1)^{i+j}M_{ij}$，其中 $M_{ij}$ 是删去第 i 行第 j 列后的子行列式。

**步骤 1**：原行列式 $D={{vA}}$，求 $A_{24}$（第 2 行第 4 列的代数余子式）。

**步骤 2**：符号因子：$(-1)^{2+4}=(-1)^6=+1$，故 $A_{24}=+M_{24}=M_{24}$。

**步骤 3**：删去第 2 行、第 4 列，得到三阶余子式 $M_{24}$，其元素来自原矩阵的以下位置：

第 1 行（来自原矩阵第 1 行，去掉第 4 列）：
$$a_{11}={{expr:mget(A,1,1)}},\quad a_{12}={{expr:mget(A,1,2)}},\quad a_{13}={{expr:mget(A,1,3)}}$$

第 2 行（来自原矩阵第 3 行，去掉第 4 列）：
$$a_{31}={{expr:mget(A,3,1)}},\quad a_{32}={{expr:mget(A,3,2)}},\quad a_{33}={{expr:mget(A,3,3)}}$$

第 3 行（来自原矩阵第 4 行，去掉第 4 列）：
$$a_{41}={{expr:mget(A,4,1)}},\quad a_{42}={{expr:mget(A,4,2)}},\quad a_{43}={{expr:mget(A,4,3)}}$$

即：
$$M_{24}=\begin{vmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\\{{expr:mget(A,4,1)}}&{{expr:mget(A,4,2)}}&{{expr:mget(A,4,3)}}\end{vmatrix}$$

**步骤 4**：计算三阶行列式 $M_{24}$（按第一行展开）：
$$M_{24} = {{expr:mget(A,1,1)}}\cdot\begin{vmatrix}{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\\{{expr:mget(A,4,2)}}&{{expr:mget(A,4,3)}}\end{vmatrix} - {{expr:mget(A,1,2)}}\cdot\begin{vmatrix}{{expr:mget(A,3,1)}}&{{expr:mget(A,3,3)}}\\{{expr:mget(A,4,1)}}&{{expr:mget(A,4,3)}}\end{vmatrix} + {{expr:mget(A,1,3)}}\cdot\begin{vmatrix}{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}\\{{expr:mget(A,4,1)}}&{{expr:mget(A,4,2)}}\end{vmatrix}$$

其中各二阶行列式：
- $\begin{vmatrix}{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\\{{expr:mget(A,4,2)}}&{{expr:mget(A,4,3)}}\end{vmatrix} = {{expr:mget(A,3,2)}}\times{{expr:mget(A,4,3)}} - {{expr:mget(A,3,3)}}\times{{expr:mget(A,4,2)}}$
- $\begin{vmatrix}{{expr:mget(A,3,1)}}&{{expr:mget(A,3,3)}}\\{{expr:mget(A,4,1)}}&{{expr:mget(A,4,3)}}\end{vmatrix} = {{expr:mget(A,3,1)}}\times{{expr:mget(A,4,3)}} - {{expr:mget(A,3,3)}}\times{{expr:mget(A,4,1)}}$
- $\begin{vmatrix}{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}\\{{expr:mget(A,4,1)}}&{{expr:mget(A,4,2)}}\end{vmatrix} = {{expr:mget(A,3,1)}}\times{{expr:mget(A,4,2)}} - {{expr:mget(A,3,2)}}\times{{expr:mget(A,4,1)}}$

**步骤 5**：代入化简得 $M_{24}={{expr:cofactor(A,2,4)}}$，又因符号因子为 $+1$，故
$$A_{24}={{A24}}$$`,
		},
	}
}

func buildChapter1_7() dsl.Problem {
	k := "Chapter1_7"
	ids := BlankIDs(k, 3)
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title: fmt.Sprintf(
			`用Cramer法则解方程组 $${{eq}}$$ 其解为 $x_1={{blank:%s}},\;x_2={{blank:%s}},\;x_3={{blank:%s}}$`,
			ids[0], ids[1], ids[2]),
		Variables: map[string]dsl.Variable{
			"A": {Kind: "matrix", Rows: 3, Cols: 3, Generator: map[string]interface{}{"rule": "upper_unit", "min": -4, "max": 4}},
			"x": {Kind: "vector", Size: 3, Generator: map[string]interface{}{"rule": "range", "min": -5, "max": 5}},
		},
		Derived: map[string]string{"b": "A * x", "eq": "cases_title(A,b)"},
		Render:  map[string]string{"eq": "eq"},
		Answer: dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{
			{ID: ids[0], Expr: "x[1]"}, {ID: ids[1], Expr: "x[2]"}, {ID: ids[2], Expr: "x[3]"},
		}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：Cramer 法则——当系数行列式 $D=\det(A)\neq 0$ 时，方程组有唯一解 $x_j=D_j/D$，其中 $D_j$ 是将 $A$ 的第 j 列替换为常数项 b 后的行列式。

**步骤 1**：方程组为 ${{eq}}$，系数矩阵 $A$ 和常数项 $b$ 为：
$$A=\begin{pmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\end{pmatrix},\quad b=\begin{pmatrix}{{expr:b[1]}}\\{{expr:b[2]}}\\{{expr:b[3]}}\end{pmatrix}$$

**步骤 2**：计算系数行列式 $D=\det(A)$：
$$D=\begin{vmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\end{vmatrix}$$
经计算（按第一行 Laplace 展开）：$D={{expr:det(A)}}$。因 $D\neq 0$，方程组有唯一解。

**步骤 3**：构造 $D_1$（将 $A$ 的第 1 列替换为 $b$）：
$$D_1=\begin{vmatrix}{{expr:b[1]}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}\\{{expr:b[2]}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}\\{{expr:b[3]}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\end{vmatrix}$$
经计算得 $D_1$ 的值。

**步骤 4**：构造 $D_2$（将 $A$ 的第 2 列替换为 $b$）：
$$D_2=\begin{vmatrix}{{expr:mget(A,1,1)}}&{{expr:b[1]}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,2,1)}}&{{expr:b[2]}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:b[3]}}&{{expr:mget(A,3,3)}}\end{vmatrix}$$
经计算得 $D_2$ 的值。

**步骤 5**：构造 $D_3$（将 $A$ 的第 3 列替换为 $b$）：
$$D_3=\begin{vmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:b[1]}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:b[2]}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:b[3]}}\end{vmatrix}$$
经计算得 $D_3$ 的值。

**步骤 6**：由 Cramer 法则：
$$x_1 = \frac{D_1}{D} = {{expr:x[1]}},\quad x_2 = \frac{D_2}{D} = {{expr:x[2]}},\quad x_3 = \frac{D_3}{D} = {{expr:x[3]}}$$

**步骤 7**（验算）：将解代回原方程组验证 $A\cdot x=b$：
$$\begin{pmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\end{pmatrix}\begin{pmatrix}{{expr:x[1]}}\\{{expr:x[2]}}\\{{expr:x[3]}}\end{pmatrix} = \begin{pmatrix}{{expr:b[1]}}\\{{expr:b[2]}}\\{{expr:b[3]}}\end{pmatrix}$$
等式成立，解正确。`,
		},
	}
}

func buildChapter1_8_1() dsl.Problem {
	k := "Chapter1_8_1"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`已知齐次线性方程组 $${{eq}}$$ 有非零解，则 $\lambda={{blank:%s}}$`, id),
		Variables: map[string]dsl.Variable{
			"A": {
				Kind: "matrix", Rows: 3, Cols: 3,
				Generator: map[string]interface{}{
					"rule": "lambda_linear_det_zero", "param_var": "lambda",
					"param_row": 2, "param_col": 2,
					"entry_min": -8, "entry_max": 8, "lambda_min": -12, "lambda_max": 12,
					"max_attempts": 120,
				},
			},
		},
		Derived: map[string]string{"eq": "lambda_cases_title(A)"},
		Render:  map[string]string{"eq": "eq"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "lambda"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：齐次线性方程组 $Ax=0$ 有非零解的充分必要条件是系数行列式 $\det(A)=0$。

**步骤 1**：原方程组 ${{eq}}$，系数矩阵 $A$ 为：
$$A=\begin{pmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,2,1)}}&\lambda&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\end{pmatrix}$$
其中第 $(2,2)$ 位置为参数 $\lambda$（即 $a_{22}=\lambda$），其余元素均为已知常数。

**步骤 2**：计算 $\det(A)$。按第一行展开：
$$\det(A) = a_{11}A_{11} + a_{12}A_{12} + a_{13}A_{13}$$
$$= {{expr:mget(A,1,1)}}\cdot\begin{vmatrix}\lambda&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\end{vmatrix} - {{expr:mget(A,1,2)}}\cdot\begin{vmatrix}{{expr:mget(A,2,1)}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,3)}}\end{vmatrix} + {{expr:mget(A,1,3)}}\cdot\begin{vmatrix}{{expr:mget(A,2,1)}}&\lambda\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}\end{vmatrix}$$

各二阶行列式：
- $\begin{vmatrix}\lambda&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,2)}}&{{expr:mget(A,3,3)}}\end{vmatrix} = \lambda\cdot{{expr:mget(A,3,3)}} - {{expr:mget(A,2,3)}}\cdot{{expr:mget(A,3,2)}}$
- $\begin{vmatrix}{{expr:mget(A,2,1)}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,3)}}\end{vmatrix} = {{expr:mget(A,2,1)}}\cdot{{expr:mget(A,3,3)}} - {{expr:mget(A,2,3)}}\cdot{{expr:mget(A,3,1)}}$
- $\begin{vmatrix}{{expr:mget(A,2,1)}}&\lambda\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}\end{vmatrix} = {{expr:mget(A,2,1)}}\cdot{{expr:mget(A,3,2)}} - \lambda\cdot{{expr:mget(A,3,1)}}$

**步骤 3**：将以上代入展开式，整理得 $\det(A)$ 关于 $\lambda$ 的表达式。由于 $\lambda$ 仅出现在一个位置，$\det(A)$ 是 $\lambda$ 的一次多项式：
$$\det(A) = {{expr:mget(A,1,1)}}\cdot(\lambda\cdot{{expr:mget(A,3,3)}} - {{expr:mget(A,2,3)}}\cdot{{expr:mget(A,3,2)}}) - {{expr:mget(A,1,2)}}\cdot({{expr:mget(A,2,1)}}\cdot{{expr:mget(A,3,3)}} - {{expr:mget(A,2,3)}}\cdot{{expr:mget(A,3,1)}}) + {{expr:mget(A,1,3)}}\cdot({{expr:mget(A,2,1)}}\cdot{{expr:mget(A,3,2)}} - \lambda\cdot{{expr:mget(A,3,1)}})$$

合并 $\lambda$ 项和常数项：
$$\det(A) = \lambda\cdot\left({{expr:mget(A,1,1)}}\cdot{{expr:mget(A,3,3)}} - {{expr:mget(A,1,3)}}\cdot{{expr:mget(A,3,1)}}\right) + \text{常数项}$$

**步骤 4**：令 $\det(A)=0$，解出 $\lambda$：
$$\lambda = {{expr:lambda}}$$

**步骤 5**（验算）：将 $\lambda={{expr:lambda}}$ 代回矩阵 $A$ 的第 $(2,2)$ 位置，重新计算行列式验证 $\det(A)=0$，确认方程组确有非零解。`,
		},
	}
}

func buildChapter1_8_2() dsl.Problem {
	k := "Chapter1_8_2"
	id := BlankIDs(k, 1)[0]
	return dsl.Problem{
		ID:      ProblemID(k),
		Version: "bank-v1",
		Title:   fmt.Sprintf(`已知齐次线性方程组 $${{eq}}$$ 有非零解，则 $\lambda={{blank:%s}}$`, id),
		Variables: map[string]dsl.Variable{
			"A": {
				Kind: "matrix", Rows: 3, Cols: 3,
				Generator: map[string]interface{}{
					"rule": "lambda_linear_det_zero", "param_var": "lambda",
					"param_row": 3, "param_col": 3,
					"entry_min": -8, "entry_max": 8, "lambda_min": -12, "lambda_max": 12,
					"max_attempts": 120,
				},
			},
		},
		Derived: map[string]string{"eq": "lambda_cases_title(A)"},
		Render:  map[string]string{"eq": "eq"},
		Answer:  dsl.AnswerSchema{FieldDefs: []dsl.AnswerFieldDef{{ID: id, Expr: "lambda"}}},
		Meta: map[string]interface{}{
			"solution_zh": `**解题思路**：齐次线性方程组 $Ax=0$ 有非零解的充分必要条件是系数行列式 $\det(A)=0$。

**步骤 1**：原方程组 ${{eq}}$，系数矩阵 $A$ 为：
$$A=\begin{pmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}\\{{expr:mget(A,3,1)}}&{{expr:mget(A,3,2)}}&\lambda\end{pmatrix}$$
其中第 $(3,3)$ 位置为参数 $\lambda$（即 $a_{33}=\lambda$），其余元素均为已知常数。

**步骤 2**：计算 $\det(A)$。按第三行展开：
$$\det(A) = a_{31}A_{31} + a_{32}A_{32} + a_{33}A_{33}$$
$$= {{expr:mget(A,3,1)}}\cdot A_{31} + {{expr:mget(A,3,2)}}\cdot A_{32} + \lambda\cdot A_{33}$$

各代数余子式：
- $A_{31}=(-1)^{3+1}\cdot\begin{vmatrix}{{expr:mget(A,1,2)}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,2,2)}}&{{expr:mget(A,2,3)}}\end{vmatrix} = (+1)\cdot\left({{expr:mget(A,1,2)}}\cdot{{expr:mget(A,2,3)}} - {{expr:mget(A,1,3)}}\cdot{{expr:mget(A,2,2)}}\right)$
- $A_{32}=(-1)^{3+2}\cdot\begin{vmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,3)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,3)}}\end{vmatrix} = (-1)\cdot\left({{expr:mget(A,1,1)}}\cdot{{expr:mget(A,2,3)}} - {{expr:mget(A,1,3)}}\cdot{{expr:mget(A,2,1)}}\right)$
- $A_{33}=(-1)^{3+3}\cdot\begin{vmatrix}{{expr:mget(A,1,1)}}&{{expr:mget(A,1,2)}}\\{{expr:mget(A,2,1)}}&{{expr:mget(A,2,2)}}\end{vmatrix} = (+1)\cdot\left({{expr:mget(A,1,1)}}\cdot{{expr:mget(A,2,2)}} - {{expr:mget(A,1,2)}}\cdot{{expr:mget(A,2,1)}}\right)$

**步骤 3**：将以上代入展开式，整理得 $\det(A)$ 关于 $\lambda$ 的表达式。由于 $\lambda$ 仅出现在一个位置，$\det(A)$ 是 $\lambda$ 的一次多项式：
$$\det(A) = {{expr:mget(A,3,1)}}\cdot A_{31} + {{expr:mget(A,3,2)}}\cdot A_{32} + \lambda\cdot A_{33}$$

即：
$$\det(A) = \lambda\cdot A_{33} + \left({{expr:mget(A,3,1)}}\cdot A_{31} + {{expr:mget(A,3,2)}}\cdot A_{32}\right)$$

其中 $A_{33}$ 和括号内的常数项均可由已知矩阵元素算出。

**步骤 4**：令 $\det(A)=0$，解出 $\lambda$：
$$\lambda = -\frac{{{expr:mget(A,3,1)}}\cdot A_{31} + {{expr:mget(A,3,2)}}\cdot A_{32}}{A_{33}} = {{expr:lambda}}$$

**步骤 5**（验算）：将 $\lambda={{expr:lambda}}$ 代回矩阵 $A$ 的第 $(3,3)$ 位置，重新计算行列式验证 $\det(A)=0$，确认方程组确有非零解。`,
		},
	}
}
