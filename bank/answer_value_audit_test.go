package bank

import (
	"fmt"
	"testing"

	"github.com/neumathe/la-dsl/dsl"
)

// 题库所有题的标准答案值类型须落在有理数判分路径支持的集合内；不得出现矩阵整体、字符串等。
func TestAllBankAnswerValueKinds(t *testing.T) {
	seed := "audit-answer-value-kind"
	salt := "audit-salt"
	allowed := map[string]struct{}{
		"int":      {},
		"int64":    {},
		"*big.Int": {},
		"*big.Rat": {},
	}
	for _, key := range AllQuestionKeys {
		t.Run(key, func(t *testing.T) {
			g, err := GenerateBankQuestion(key, seed+":"+key, salt)
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range g.AnswerFields {
				typ := fmt.Sprintf("%T", f.Value)
				if _, ok := allowed[typ]; !ok {
					t.Fatalf("unsupported answer value type %s for field %s", typ, f.ID)
				}
				k := dsl.ClassifyScalarAnswerValue(f.Value)
				if k == "unsupported" {
					t.Fatalf("classify unsupported for field %s type %s", f.ID, typ)
				}
				_, err := dsl.ParseUserRational(dsl.ValueToCanonicalString(f.Value))
				if err != nil {
					t.Fatalf("canonical not parseable as user rational for %s: %v", f.ID, err)
				}
			}
		})
	}
}
