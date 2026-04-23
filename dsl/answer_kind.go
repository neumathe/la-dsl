package dsl

import "math/big"

// 标准答案值的粗分类，便于上层选择键盘/占位提示（判分仍统一按有理数比较）。
const (
	FieldAnswerKindInteger  = "integer"  // int / int64 / 整型 big.Rat
	FieldAnswerKindBigInt   = "bigint"   // *big.Int（行列式等，输入仍为十进制整数字符串）
	FieldAnswerKindRational = "rational" // 非整 *big.Rat，建议展示「可填分数」
)

// FieldInputHint 与单次出题实例的标准答案形态对应。
type FieldInputHint struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

// ClassifyScalarAnswerValue 根据标准答案 Go 类型与值分类。
func ClassifyScalarAnswerValue(v interface{}) string {
	switch t := v.(type) {
	case *big.Rat:
		if t.IsInt() {
			return FieldAnswerKindInteger
		}
		return FieldAnswerKindRational
	case *big.Int:
		return FieldAnswerKindBigInt
	case int:
		return FieldAnswerKindInteger
	case int64:
		return FieldAnswerKindInteger
	case float64:
		return FieldAnswerKindRational
	default:
		return "unsupported"
	}
}

// FieldInputHintsFromGenerated 由出题结果生成每空的输入提示种类。
func FieldInputHintsFromGenerated(g *GeneratedQuestion) []FieldInputHint {
	if g == nil {
		return nil
	}
	out := make([]FieldInputHint, len(g.AnswerFields))
	for i, f := range g.AnswerFields {
		out[i] = FieldInputHint{ID: f.ID, Kind: ClassifyScalarAnswerValue(f.Value)}
	}
	return out
}
