package bank

import "github.com/neumathe/la-dsl/dsl"

// GenerateBankExplanation 与 GenerateBankQuestion / JudgeBankQuestion 使用相同 questionKey、seed、salt，
// 生成结构化解析（已知量、派生量、各空标准答案），供学习核对。
func GenerateBankExplanation(questionKey, seedStr, serverSalt string) (*dsl.QuestionExplanation, error) {
	p, err := BuildProblem(questionKey)
	if err != nil {
		return nil, err
	}
	return dsl.GenerateExplanation(p, seedStr, serverSalt)
}
