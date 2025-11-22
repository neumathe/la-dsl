package dsl

// Problem 定义了一道题目的 DSL 结构
type Problem struct {
	ID        int64                  `json:"id"`
	Title     string                 `json:"title"`
	Variables map[string]Variable    `json:"variables"`
	Derived   map[string]string      `json:"derived"` // expression -> var
	Render    map[string]string      `json:"render"`
	Answer    AnswerSchema           `json:"answer"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
	Version   string                 `json:"version,omitempty"`
}

// Variable 描述一个随机变量（标量 / 向量 / 矩阵）
type Variable struct {
	Kind      string                 `json:"kind"` // "matrix","vector","scalar"
	Rows      int                    `json:"rows,omitempty"`
	Cols      int                    `json:"cols,omitempty"`
	Size      int                    `json:"size,omitempty"`
	Generator map[string]interface{} `json:"generator"`
	// Fixed 为可选的固定值（不走生成器）
	Fixed interface{} `json:"fixed,omitempty"`
}

// AnswerSchema 定义答案的 DSL 结构
type AnswerSchema struct {
	Expression string   `json:"expression,omitempty"`
	Fields     []string `json:"fields,omitempty"`
	// 可选的增强版字段定义，支持为每个答案表达式绑定一个稳定的 ID
	// （用于前端填空 ID 与后端判分的对应）。
	FieldDefs []AnswerFieldDef `json:"field_defs,omitempty"`
}

// AnswerFieldDef 为一个答案字段的 DSL 定义
type AnswerFieldDef struct {
	ID   string `json:"id,omitempty"`
	Expr string `json:"expr,omitempty"`
}

// Instance 表示一次题目实例化的结果
type Instance struct {
	ProblemID int64
	Seed      string
	Vars      map[string]interface{} // variable name -> value (matrix/vector/scalar)
	Derived   map[string]interface{} // derived vars
}
