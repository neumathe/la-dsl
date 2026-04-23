package ladsl

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/neumathe/la-dsl/bank"
	"github.com/neumathe/la-dsl/dsl"
)

// Service 题库统一入口：同一 serverSalt 下出题、判题、解析与 seed 规则与 bank 包一致。
type Service struct {
	serverSalt string
}

// NewService 创建服务；serverSalt 建议为服务端机密常量（与 Judge/Explain 使用同值）。
func NewService(serverSalt string) *Service {
	return &Service{serverSalt: serverSalt}
}

// QuestionKeys 返回题库中全部逻辑题键（可随机抽题）。
func QuestionKeys() []string {
	return append([]string(nil), bank.AllQuestionKeys...)
}

// ValidQuestionKey 判断 key 是否可 BuildProblem。
func ValidQuestionKey(key string) bool {
	_, err := bank.BuildProblem(key)
	return err == nil
}

// RandomSeed 生成可用作 RollQuestion 的随机种子字符串（服务端抽题）。
func RandomSeed() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// DescribeQuestion 返回该题静态空位布局（不实例化随机数，无题面文本）。
func DescribeQuestion(questionKey string) (*BlankDescriptor, error) {
	p, err := bank.BuildProblem(questionKey)
	if err != nil {
		return nil, err
	}
	blanks, err := blanksFromProblem(p)
	if err != nil {
		return nil, err
	}
	return &BlankDescriptor{
		QuestionKey:       questionKey,
		ProblemID:         p.ID,
		Version:           p.Version,
		InputConventionID: dsl.AnswerInputConventionV1,
		BlankCount:        len(blanks),
		Blanks:            blanks,
	}, nil
}

// RollQuestion 生成一题随机实例的对外数据（题面 + 空位 id，不含答案与表达式）。
func (s *Service) RollQuestion(questionKey, seed string) (*QuestionPublic, error) {
	p, err := bank.BuildProblem(questionKey)
	if err != nil {
		return nil, err
	}
	g, err := bank.GenerateBankQuestion(questionKey, seed, s.serverSalt)
	if err != nil {
		return nil, err
	}
	return publicFromGenerated(questionKey, p, seed, g), nil
}

// RollQuestionServer 一次生成：对外题面 + 完整标准答案（Private 仅服务端使用）。
func (s *Service) RollQuestionServer(questionKey, seed string) (*QuestionServerBundle, error) {
	p, err := bank.BuildProblem(questionKey)
	if err != nil {
		return nil, err
	}
	g, err := bank.GenerateBankQuestion(questionKey, seed, s.serverSalt)
	if err != nil {
		return nil, err
	}
	return &QuestionServerBundle{
		Public:  publicFromGenerated(questionKey, p, seed, g),
		Private: g,
	}, nil
}

// Judge 根据与用户出题时相同的 key、seed、serverSalt 重算标准答案并判分。
func (s *Service) Judge(questionKey, seed string, userAnswers map[string]string, opts *dsl.JudgeOptions) (*dsl.JudgeResult, error) {
	return bank.JudgeBankQuestion(questionKey, seed, s.serverSalt, userAnswers, opts)
}

// Explain 生成结构化解析（与 Judge 同源数据）。
func (s *Service) Explain(questionKey, seed string) (*dsl.QuestionExplanation, error) {
	return bank.GenerateBankExplanation(questionKey, seed, s.serverSalt)
}

// InputContractV1 返回与判分一致的输入约定文档，上层可缓存或单独接口下发。
func InputContractV1() dsl.AnswerInputContractDoc {
	return dsl.AnswerInputContractV1()
}

func publicFromGenerated(key string, p dsl.Problem, seed string, g *dsl.GeneratedQuestion) *QuestionPublic {
	blanks := make([]BlankInfo, len(g.AnswerFields))
	for i, f := range g.AnswerFields {
		blanks[i] = BlankInfo{ID: f.ID, Order: i + 1, Layout: f.Layout}
	}
	return &QuestionPublic{
		QuestionKey:       key,
		ProblemID:         p.ID,
		Version:           p.Version,
		InputConventionID: dsl.AnswerInputConventionV1,
		FieldHints:          dsl.FieldInputHintsFromGenerated(g),
		Seed:                seed,
		Title:               g.Title,
		Blanks:              blanks,
	}
}

// blanksFromProblem 与 dsl.ExtractAnswerWithMeta 对 FieldDefs / Expression / Fields 的编号规则一致。
func blanksFromProblem(p dsl.Problem) ([]BlankInfo, error) {
	if p.Answer.Expression != "" {
		id := "ans"
		if len(p.Answer.FieldDefs) > 0 && p.Answer.FieldDefs[0].ID != "" {
			id = p.Answer.FieldDefs[0].ID
		}
		var ly *dsl.AnswerFieldLayout
		if len(p.Answer.FieldDefs) > 0 {
			ly = p.Answer.FieldDefs[0].Layout
		}
		return []BlankInfo{{ID: id, Order: 1, Layout: ly}}, nil
	}
	if len(p.Answer.FieldDefs) > 0 {
		out := make([]BlankInfo, 0, len(p.Answer.FieldDefs))
		order := 0
		for i, fd := range p.Answer.FieldDefs {
			if strings.TrimSpace(fd.Expr) == "" {
				continue
			}
			order++
			id := fd.ID
			if id == "" {
				id = fmt.Sprintf("field_%d", i+1)
			}
			out = append(out, BlankInfo{ID: id, Order: order, Layout: fd.Layout})
		}
		if len(out) == 0 {
			return nil, fmt.Errorf("ladsl: no answer field defs in problem %d", p.ID)
		}
		return out, nil
	}
	if len(p.Answer.Fields) > 0 {
		out := make([]BlankInfo, len(p.Answer.Fields))
		for i := range p.Answer.Fields {
			out[i] = BlankInfo{ID: fmt.Sprintf("field_%d", i+1), Order: i + 1}
		}
		return out, nil
	}
	return nil, fmt.Errorf("ladsl: empty answer schema for problem %d", p.ID)
}
