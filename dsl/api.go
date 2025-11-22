package dsl

// 面向上层的封装 API：
// 上层仅需传入 Problem + seedStr + serverSalt，
// 即可得到渲染好的题目文本和带空 ID 的答案列表。

// AnswerField 是对单个答案字段的封装，方便前端展示与后端判分。
type AnswerField struct {
	ID    string      `json:"id"`
	Expr  string      `json:"expr"`
	Value interface{} `json:"value"`
}

// GeneratedQuestion 是一次完整的题目生成结果。
type GeneratedQuestion struct {
	Title        string        `json:"title"`
	AnswerFields []AnswerField `json:"answer_fields"`
}

// GenerateQuestion 从 DSL Problem 和种子生成一道题目：
// - 按 Problem 的定义实例化变量与派生量
// - 渲染题干 Title（如果有）
// - 计算所有答案字段，绑定稳定的空 ID
//
// 注意：同一 Problem + seedStr + serverSalt，多次调用结果完全一致。
func GenerateQuestion(p Problem, seedStr, serverSalt string) (*GeneratedQuestion, error) {
	inst, err := InstantiateProblem(p, seedStr, serverSalt)
	if err != nil {
		return nil, err
	}
	title, err := RenderTitle(p, inst)
	if err != nil {
		return nil, err
	}
	fields, err := ExtractAnswerWithMeta(p, inst)
	if err != nil {
		return nil, err
	}
	return &GeneratedQuestion{
		Title:        title,
		AnswerFields: fields,
	}, nil
}
