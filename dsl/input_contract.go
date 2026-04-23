package dsl

// AnswerInputConventionV1 与 JudgeGeneratedQuestion / ParseUserRational / NormalizeUserAnswer 对齐的约定标识。
const AnswerInputConventionV1 = "la-dsl.scalar_rational.v1"

// AnswerInputContractV1 供上层一次性下发或缓存的输入规范说明（JSON）。
func AnswerInputContractV1() AnswerInputContractDoc {
	return AnswerInputContractDoc{
		ID:      AnswerInputConventionV1,
		Summary: "每空为单个有理数：与标准答案在算术意义下相等即判对；题库当前无根号、无 π、无矩阵整体字符串答案。",
		AcceptedFormats: []string{
			"整数：0、-12、105",
			"分数：-3/2、1/4（使用 ASCII 斜杠 /）",
			"小数：0.25（内部转为有理数比较；极端浮点请优先用分数）",
		},
		RejectedFormats: []string{
			"根号、π、e、LaTeX（如 \\frac{}{}、\\sqrt{}）",
			"一空中填多个数、矩阵或逗号分隔坐标",
			"整数中的千分位逗号（归一化会去掉逗号，易与小数混淆，请勿使用）",
		},
		Normalization: []string{
			"首尾空白删除；中英文括号统一为 ( )",
			"Unicode 减号 U+2212、全角负号 U+FF0D 等统一为 ASCII -",
			"分数斜杠 U+2044、全角／ 统一为 /",
			"全角逗号 ， 与半角逗号 , 会删除（用于去千分位噪声，非小数点）",
		},
		JudgeImplementation: "ScalarAnswersEqual + big.Rat",
	}
}

// AnswerInputContractDoc 与前端/客户端约定的可序列化说明。
type AnswerInputContractDoc struct {
	ID                  string   `json:"id"`
	Summary             string   `json:"summary"`
	AcceptedFormats     []string `json:"accepted_formats"`
	RejectedFormats     []string `json:"rejected_formats"`
	Normalization       []string `json:"normalization"`
	JudgeImplementation string   `json:"judge_implementation,omitempty"`
}
