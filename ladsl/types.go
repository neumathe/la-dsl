package ladsl

import "github.com/neumathe/la-dsl/dsl"

// BlankInfo 单个填空位（下发客户端时不含表达式与答案）。
type BlankInfo struct {
	ID     string               `json:"id"`
	Order  int                  `json:"order"` // 1 起，与表单顺序一致
	Layout *dsl.AnswerFieldLayout `json:"layout,omitempty"`
}

// BlankDescriptor 某题键下的静态空位布局（与具体随机种子无关）。
type BlankDescriptor struct {
	QuestionKey         string      `json:"question_key"`
	ProblemID           int64       `json:"problem_id"`
	Version             string      `json:"version,omitempty"`
	InputConventionID   string      `json:"input_convention_id,omitempty"`
	BlankCount          int         `json:"blank_count"`
	Blanks              []BlankInfo `json:"blanks"`
}

// QuestionPublic 一次随机实例的对外题面（供学生端展示与收题）。
type QuestionPublic struct {
	QuestionKey         string               `json:"question_key"`
	ProblemID           int64                `json:"problem_id"`
	Version             string               `json:"version,omitempty"`
	InputConventionID   string               `json:"input_convention_id,omitempty"`
	FieldHints          []dsl.FieldInputHint `json:"field_hints,omitempty"`
	Seed                string               `json:"seed"`
	Title               string               `json:"title"`
	Blanks              []BlankInfo          `json:"blanks"`
}

// QuestionServerBundle 服务端一次生成：对外题面 + dsl 标准答案（Private 勿下发给学生端）。
type QuestionServerBundle struct {
	Public  *QuestionPublic          `json:"public"`
	Private *dsl.GeneratedQuestion   `json:"private"`
}
