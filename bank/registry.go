package bank

import (
	"fmt"

	"github.com/neumathe/la-dsl/dsl"
)

var builders = map[string]func() dsl.Problem{
	"Chapter1_8_1": buildChapter1_8_1, "Chapter1_4": buildChapter1_4, "Chapter1_3": buildChapter1_3,
	"Chapter1_1": buildChapter1_1, "Chapter1_2": buildChapter1_2, "Chapter1_6": buildChapter1_6,
	"Chapter1_5": buildChapter1_5, "Chapter1_7": buildChapter1_7, "Chapter1_8_2": buildChapter1_8_2,

	"Chapter2_4_1": buildChapter2_4_1, "Chapter2_7_2": buildChapter2_7_2, "Chapter2_7_3": buildChapter2_7_3,
	"Chapter2_4_2": buildChapter2_4_2, "Chapter2_4_3": buildChapter2_4_3, "Chapter2_6": buildChapter2_6,
	"Chapter2_7_1": buildChapter2_7_1,
	"Chapter2_5_2": buildChapter2_5_2, "Chapter2_1": buildChapter2_1, "Chapter2_3": buildChapter2_3, "Chapter2_2": buildChapter2_2,

	"Chapter3_5": buildChapter3_5, "Chapter3_8": buildChapter3_8, "Chapter3_4": buildChapter3_4,
	"Chapter3_3": buildChapter3_3, "Chapter3_10": buildChapter3_10, "Chapter3_1": buildChapter3_1,
	"Chapter3_2": buildChapter3_2, "Chapter3_6": buildChapter3_6, "Chapter3_7": buildChapter3_7,
	"Chapter3_9": buildChapter3_9, "Chapter3_11": buildChapter3_11,

	"Chapter4_8": buildChapter4_8, "Chapter4_3_1": buildChapter4_3_1, "Chapter4_1": buildChapter4_1,
	"Chapter4_4": buildChapter4_4, "Chapter4_2": buildChapter4_2, "Chapter4_3_3": buildChapter4_3_3,
	"Chapter4_3_2": buildChapter4_3_2, "Chapter4_5_1": buildChapter4_5_1, "Chapter4_5_2": buildChapter4_5_2,
	"Chapter4_7": buildChapter4_7, "Chapter4_6": buildChapter4_6,

	"Chapter5_1": buildChapter5_1, "Chapter5_3": buildChapter5_3, "Chapter5_5": buildChapter5_5,
	"Chapter5_2": buildChapter5_2, "Chapter5_4": buildChapter5_4,
	"Chapter5_8": buildChapter5_8, "Chapter5_6": buildChapter5_6, "Chapter5_7": buildChapter5_7,

	"Chapter6_1_2": buildChapter6_1_2, "Chapter6_1_3": buildChapter6_1_3, "Chapter6_1_1": buildChapter6_1_1,
	"Chapter6_5": buildChapter6_5, "Chapter6_4": buildChapter6_4, "Chapter6_2": buildChapter6_2,
	"Chapter6_6": buildChapter6_6, "Chapter6_3": buildChapter6_3,

	"Chapter7_7": buildChapter7_7, "Chapter7_4": buildChapter7_4, "Chapter7_5_1": buildChapter7_5_1,
	"Chapter7_2": buildChapter7_2, "Chapter7_3": buildChapter7_3, "Chapter7_1": buildChapter7_1,
	"Chapter7_5_2": buildChapter7_5_2, "Chapter7_8": buildChapter7_8, "Chapter7_6": buildChapter7_6,
	"Chapter7_5_3": buildChapter7_5_3, "Chapter7_10": buildChapter7_10, "Chapter7_9": buildChapter7_9,
}

// BuildProblem 返回与题库逻辑题号对应的 DSL 题目（填空 id 与 HTML 一致）。
func BuildProblem(questionKey string) (dsl.Problem, error) {
	fn, ok := builders[questionKey]
	if !ok {
		return dsl.Problem{}, fmt.Errorf("bank: unknown question key %q", questionKey)
	}
	return fn(), nil
}

// GenerateBankQuestion 按题库键 + 种子生成题目与答案。
func GenerateBankQuestion(questionKey, seedStr, serverSalt string) (*dsl.GeneratedQuestion, error) {
	p, err := BuildProblem(questionKey)
	if err != nil {
		return nil, err
	}
	return dsl.GenerateQuestion(p, seedStr, serverSalt)
}
