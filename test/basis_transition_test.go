package test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/neumathe/la-dsl/dsl"
)

// TestBasisTransitionAndCoord
// 对应题型：
// 向量空间 R[x]_3 中，由基 1,x,x^2 到基 -1, -2+x, -1-2x+2x^2 的过渡矩阵，
// 以及 p(x)=a+bx+cx^2 在该基下的坐标。
// 这里：
// - 新基固定为题目给出的三个多项式；
// - p(x) 的系数 (a,b,c) 随机生成（整数）；
// - 使用 DSL 的 derived + answer.fields 一次性给出 3x3 过渡矩阵和 3 维坐标，共 12 个空。
func TestBasisTransitionAndCoord(t *testing.T) {
	// 构造固定的基变换矩阵 B：列为新基在旧基 {1,x,x^2} 下的坐标
	// b1(x) = -1        -> [-1,  0, 0]^T
	// b2(x) = -2 + 1 x  -> [-2,  1, 0]^T
	// b3(x) = -1 -2x +2x^2 -> [-1, -2, 2]^T
	B := dsl.NewMatrixInt(3, 3)
	B.A[0][0], B.A[1][0], B.A[2][0] = -1, 0, 0
	B.A[0][1], B.A[1][1], B.A[2][1] = -2, 1, 0
	B.A[0][2], B.A[1][2], B.A[2][2] = -1, -2, 2

	// 标准基下的 e1,e2,e3
	e1 := dsl.NewVectorInt(3)
	e1.V[0], e1.V[1], e1.V[2] = 1, 0, 0
	e2 := dsl.NewVectorInt(3)
	e2.V[0], e2.V[1], e2.V[2] = 0, 1, 0
	e3 := dsl.NewVectorInt(3)
	e3.V[0], e3.V[1], e3.V[2] = 0, 0, 1

	// 随机生成 p(x) = a + b x + c x^2 的系数（简单整数）
	seed := time.Now().UnixNano()
	a0 := int64(seed%5) - 2 // 简单取几种值，示例用，防止太大
	a1 := int64((seed/7)%5) - 2
	a2 := int64((seed/11)%5) - 2
	if a0 == 0 && a1 == 0 && a2 == 0 {
		a0 = 1
	}
	pVec := dsl.NewVectorInt(3)
	pVec.V[0], pVec.V[1], pVec.V[2] = a0, a1, a2

	prob := dsl.Problem{
		ID:      70001,
		Version: "v1",
		// 此处 Title 留空，测试中手动格式化题干，便于直接对照原题样式。
		Title: "",
		Variables: map[string]dsl.Variable{
			"B": {
				Kind:  "matrix",
				Rows:  3,
				Cols:  3,
				Fixed: B,
			},
			"e1": {
				Kind:  "vector",
				Size:  3,
				Fixed: e1,
			},
			"e2": {
				Kind:  "vector",
				Size:  3,
				Fixed: e2,
			},
			"e3": {
				Kind:  "vector",
				Size:  3,
				Fixed: e3,
			},
			"p": {
				Kind:  "vector",
				Size:  3,
				Fixed: pVec,
			},
		},
		Derived: map[string]string{
			// 过渡矩阵 P_{new <- old} = B^{-1} 的三列
			"Pcol1": "solve(B,e1)",
			"Pcol2": "solve(B,e2)",
			"Pcol3": "solve(B,e3)",
			// p 在新基下的坐标：solve(B,p)
			"coord": "solve(B,p)",
		},
		Answer: dsl.AnswerSchema{
			Fields: []string{
				"Pcol1[1]", "Pcol1[2]", "Pcol1[3]",
				"Pcol2[1]", "Pcol2[2]", "Pcol2[3]",
				"Pcol3[1]", "Pcol3[2]", "Pcol3[3]",
				"coord[1]", "coord[2]", "coord[3]",
			},
		},
	}

	seedStr := fmt.Sprintf("basis-transition-%d", seed)
	inst, err := dsl.InstantiateProblem(prob, seedStr, "test-salt")
	if err != nil {
		t.Fatalf("InstantiateProblem error: %v", err)
	}

	// 1) 校验过渡矩阵三列：使用 EvaluateExpression 再算一遍 solve(B, ei)
	checkSolveVec := func(exprName, expr string) {
		val, err := dsl.EvaluateExpression(expr, inst)
		if err != nil {
			t.Fatalf("EvaluateExpression %s error: %v", expr, err)
		}
		vExpr, ok := val.([]*big.Rat)
		if !ok {
			t.Fatalf("%s expr result is not []*big.Rat, got %T", expr, val)
		}
		raw, ok := inst.Vars[exprName]
		if !ok {
			t.Fatalf("%s not found in Vars", exprName)
		}
		vDer, ok := raw.([]*big.Rat)
		if !ok {
			t.Fatalf("%s is not []*big.Rat, got %T", exprName, raw)
		}
		if len(vExpr) != len(vDer) {
			t.Fatalf("%s length mismatch", exprName)
		}
		for i := range vExpr {
			if vExpr[i].Cmp(vDer[i]) != 0 {
				t.Fatalf("%s[%d] mismatch: %s vs %s", exprName, i, vExpr[i].RatString(), vDer[i].RatString())
			}
		}
	}

	checkSolveVec("Pcol1", "solve(B,e1)")
	checkSolveVec("Pcol2", "solve(B,e2)")
	checkSolveVec("Pcol3", "solve(B,e3)")
	checkSolveVec("coord", "solve(B,p)")

	// 2) 用 ExtractAnswer 拿到 12 个空的答案
	ans, err := dsl.ExtractAnswer(prob, inst)
	if err != nil {
		t.Fatalf("ExtractAnswer error: %v", err)
	}
	ansSlice, ok := ans.([]interface{})
	if !ok {
		t.Fatalf("answer is not []interface{}, got %T", ans)
	}
	if len(ansSlice) != 12 {
		t.Fatalf("answer length != 12, got %d", len(ansSlice))
	}

	// 3) 打印出题目和答案示例（LaTeX）
	qText := formatBasisTransitionQuestion(a0, a1, a2)
	Ptex := formatTransitionMatrixLatex(inst)
	coordTex := formatCoordLatex(inst)

	t.Logf("示例题目（LaTeX）：%s", qText)
	t.Logf("过渡矩阵为：%s", Ptex)
	t.Logf("坐标为：%s\\(^T.\\)", coordTex)
	t.Logf("示例答案字段（按照 9 个矩阵元素 + 3 个坐标）：%v", ansSlice)
}

// 按原题样式格式化题干文字（仅关键数据 a,b,c 随机）
func formatBasisTransitionQuestion(a0, a1, a2 int64) string {
	// 新基保持固定：-1, -2+x, -1-2x+2x^2
	// 格式化 p(x)，处理正负号
	var pTerms []string

	// 常数项
	if a0 != 0 {
		pTerms = append(pTerms, fmt.Sprintf("%d", a0))
	}

	// x 项
	if a1 != 0 {
		if a1 == 1 {
			if len(pTerms) > 0 {
				pTerms = append(pTerms, "+x")
			} else {
				pTerms = append(pTerms, "x")
			}
		} else if a1 == -1 {
			pTerms = append(pTerms, "-x")
		} else if a1 > 0 {
			if len(pTerms) > 0 {
				pTerms = append(pTerms, fmt.Sprintf("+%dx", a1))
			} else {
				pTerms = append(pTerms, fmt.Sprintf("%dx", a1))
			}
		} else {
			pTerms = append(pTerms, fmt.Sprintf("%dx", a1))
		}
	}

	// x^2 项
	if a2 != 0 {
		if a2 == 1 {
			if len(pTerms) > 0 {
				pTerms = append(pTerms, "+x^2")
			} else {
				pTerms = append(pTerms, "x^2")
			}
		} else if a2 == -1 {
			pTerms = append(pTerms, "-x^2")
		} else if a2 > 0 {
			if len(pTerms) > 0 {
				pTerms = append(pTerms, fmt.Sprintf("+%dx^2", a2))
			} else {
				pTerms = append(pTerms, fmt.Sprintf("%dx^2", a2))
			}
		} else {
			pTerms = append(pTerms, fmt.Sprintf("%dx^2", a2))
		}
	}

	pStr := strings.Join(pTerms, "")
	if pStr == "" {
		pStr = "0"
	}

	var b strings.Builder
	b.WriteString("向量空间\\(R[x]_3\\)中，由基1，x，\\(x^2\\)到基-1，-2+1x，-1-2x+2\\(x^2\\)的过渡矩阵为")
	b.WriteString("，p(x)=")
	b.WriteString(pStr)
	b.WriteString(" 在基-1，-2+1x，-1-2x+2\\(x^2\\)下的坐标为")
	return b.String()
}

// 从 inst 中取出 Pcol1~3，格式化为 3x3 LaTeX 矩阵
func formatTransitionMatrixLatex(inst *dsl.Instance) string {
	get := func(name string) []*big.Rat {
		raw := inst.Vars[name]
		v, ok := raw.([]*big.Rat)
		if !ok {
			return nil
		}
		return v
	}
	c1 := get("Pcol1")
	c2 := get("Pcol2")
	c3 := get("Pcol3")
	if c1 == nil || c2 == nil || c3 == nil {
		return ""
	}
	var rows []string
	for i := 0; i < 3; i++ {
		row := fmt.Sprintf("%s&%s&%s",
			c1[i].RatString(),
			c2[i].RatString(),
			c3[i].RatString(),
		)
		rows = append(rows, row)
	}
	return `\begin{bmatrix}` + strings.Join(rows, `\\`) + `\end{bmatrix}`
}

// 从 inst 中取出 coord，格式化为 1x3 LaTeX 向量 (a,b,c)
func formatCoordLatex(inst *dsl.Instance) string {
	raw := inst.Vars["coord"]
	v, ok := raw.([]*big.Rat)
	if !ok || len(v) < 3 {
		return ""
	}
	return fmt.Sprintf("(%s,%s,%s)", v[0].RatString(), v[1].RatString(), v[2].RatString())
}
