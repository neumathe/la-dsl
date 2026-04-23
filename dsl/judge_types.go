package dsl

// AnswerJudgeSpec 非空时启用结构化判题（需 JudgeGeneratedQuestionContext 传入 *Instance）。
// Kind 取值：
//   - "" / "scalar"：标量/有理数相等（默认）
//   - "rational_line"：同 LineGroup 的多个空按 AnswerFields 顺序组成向量，与标准答案在 Q 上成比例（允许非零标量倍）
//   - "affine_rational"：同 AffineGroup 组成向量 x，验证 MatrixVar*x == BVecVar（有理数分量）
//   - "sorted_basis_columns"：同 BasisGroup 的多个空为列下标（1..Ncols，不足填 0），验证是否构成 A 的一组极大线性无关列（张成列空间）
//   - "permutation_multiset"：同 PermGroup 的多个空为一组数值，判分与顺序无关（multiset 相等）
//   - "eigen_pair"：同 EigenGroup 的多个空共同描述若干特征对 (λ_j, v_j)，
//     用户可按任意顺序填写；判分要求：
//       1) 每列 j 的向量 v_j 非零，且 MatrixVar · v_j == λ_j · v_j
//       2) 用户给出的全部 λ 与标准答案特征值 multiset 相等
//     RefLambdaGroup 允许 vec-only 的组（如 Q 或 α 的第二分组）借用另一组的 λ 字段。
type AnswerJudgeSpec struct {
	Kind string `json:"kind,omitempty"`

	LineGroup string `json:"line_group,omitempty"`

	AffineGroup string `json:"affine_group,omitempty"`
	MatrixVar   string `json:"matrix_var,omitempty"`
	BVecVar     string `json:"bvec_var,omitempty"`

	BasisGroup string `json:"basis_group,omitempty"`
	Ncols      int    `json:"ncols,omitempty"`

	// permutation_multiset: 同组答案与顺序无关。
	PermGroup string `json:"perm_group,omitempty"`

	// eigen_pair 相关字段。
	EigenGroup     string `json:"eigen_group,omitempty"`      // 所属逻辑组标识
	EigenRole      string `json:"eigen_role,omitempty"`       // "lambda" 或 "vec"
	EigenColumn    int    `json:"eigen_column,omitempty"`     // 1-based：该字段对应第几个特征对
	EigenComponent int    `json:"eigen_component,omitempty"`  // role=vec 时：向量第几个分量（1-based）
	RefLambdaGroup string `json:"ref_lambda_group,omitempty"` // 仅 vec 组：借用另一组的 λ 字段
}
