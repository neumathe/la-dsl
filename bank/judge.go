package bank

import "github.com/neumathe/la-dsl/dsl"

// JudgeBankQuestion 用与出题相同的 seed/salt 重算标准答案，并与用户提交的 id->答案字符串 比较。
func JudgeBankQuestion(questionKey, seedStr, serverSalt string, userAnswers map[string]string, opts *dsl.JudgeOptions) (*dsl.JudgeResult, error) {
	p, err := BuildProblem(questionKey)
	if err != nil {
		return nil, err
	}
	inst, err := dsl.InstantiateProblem(p, seedStr, serverSalt)
	if err != nil {
		return nil, err
	}
	g, err := dsl.GenerateQuestionFromInstance(p, inst)
	if err != nil {
		return nil, err
	}
	res := dsl.JudgeGeneratedQuestionContext(g, &p, inst, userAnswers, opts)
	return &res, nil
}
