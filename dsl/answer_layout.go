package dsl

// 与前端约定的填空空间布局（JSON）；判分仍只依赖 field id + 字符串值，与 layout 无关。
const AnswerLayoutSchema = "la-dsl.answer_layout.v1"

const (
	LayoutKindMatrixCell       = "matrix_cell"        // 矩阵第 (row,col) 格，1-based，与 mget 一致
	LayoutKindVectorComponent  = "vector_component"   // 向量第 index 个分量，1-based，与 v[k] 一致
)

// AnswerFieldLayout 描述单个填空在版式中的语义位置，供 App/Web 渲染矩阵格/向量分量格。
// 前端可按 kind+matrix/vector+行列维渲染输入框，field_id 仍用于提交答案 map。
type AnswerFieldLayout struct {
	Schema string `json:"schema,omitempty"` // AnswerLayoutSchema

	Kind string `json:"kind"` // LayoutKind*

	// matrix_cell：与 DSL 变量 matrix 的 (row,col) 对齐；rows/cols 为整块矩阵形状（用于画空表）。
	Matrix string `json:"matrix,omitempty"`
	Row    int    `json:"row,omitempty"`
	Col    int    `json:"col,omitempty"`
	Rows   int    `json:"rows,omitempty"`
	Cols   int    `json:"cols,omitempty"`

	// 同一题多块矩阵时区分版块（展示用，可为纯文本或简短 LaTeX 片段）。
	GroupLabel string `json:"group_label,omitempty"`

	// vector_component
	Vector string `json:"vector,omitempty"`
	Index  int    `json:"index,omitempty"`
}

// LayoutMatrixCell 构造矩阵格布局元数据。
func LayoutMatrixCell(matrix string, row, col, rows, cols int, groupLabel string) *AnswerFieldLayout {
	return &AnswerFieldLayout{
		Schema:     AnswerLayoutSchema,
		Kind:       LayoutKindMatrixCell,
		Matrix:     matrix,
		Row:        row,
		Col:        col,
		Rows:       rows,
		Cols:       cols,
		GroupLabel: groupLabel,
	}
}

// LayoutVectorComponent 构造向量分量格布局。
func LayoutVectorComponent(vector string, index int, groupLabel string) *AnswerFieldLayout {
	return &AnswerFieldLayout{
		Schema:     AnswerLayoutSchema,
		Kind:       LayoutKindVectorComponent,
		Vector:     vector,
		Index:      index,
		GroupLabel: groupLabel,
	}
}
