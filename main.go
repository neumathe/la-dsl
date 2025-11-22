package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/neumathe/la-dsl/dsl"
)

func main() {
	// 向量在一组基下坐标的示例
	probJSON := `
{
  "id": 401,
  "version": "v1",
  "title": "向量 β={{beta}} 在基 {{alpha1}}, {{alpha2}}, {{alpha3}} 下的坐标为 ____",
  "variables": {
    "x": {
      "kind": "vector",
      "size": 3,
      "generator": { "rule": "range", "min": -5, "max": 5 }
    },
    "A": {
      "kind": "matrix",
      "rows": 3,
      "cols": 3,
      "generator": { "rule": "full_rank", "min": -6, "max": 6 }
    }
  },
  "derived": {
    "beta": "A * x"
  },
  "render": {
    "beta": "beta",
    "alpha1": "col(A,1)",
    "alpha2": "col(A,2)",
    "alpha3": "col(A,3)"
  },
  "answer": {
    "fields": ["x[1]", "x[2]", "x[3]"]
  }
}
`
	var p dsl.Problem
	if err := json.Unmarshal([]byte(probJSON), &p); err != nil {
		log.Fatalf("parse prob json: %v", err)
	}

	seedStr := "exam2025:user42" // example seed
	serverSalt := "server-secret-please-store-in-vault"

	inst, err := dsl.InstantiateProblem(p, seedStr, serverSalt)
	if err != nil {
		log.Fatalf("instantiate error: %v", err)
	}

	// Render for frontend
	rendered, err := dsl.RenderInst(p, inst)
	if err != nil {
		log.Fatalf("render error: %v", err)
	}
	fmt.Println("Rendered (for frontend):")
	printRendered(rendered)

	// Extract answer
	ans, err := dsl.ExtractAnswer(p, inst)
	if err != nil {
		log.Fatalf("extract answer error: %v", err)
	}
	fmt.Println("\nAnswer fields (canonical):")
	switch t := ans.(type) {
	case []interface{}:
		fmt.Println(t)
	default:
		fmt.Printf("%#v\n", t)
	}

	fmt.Println("\nDemo: 4x4 circulant-like matrix and its det using Bareiss")
	circ := &dsl.MatrixInt{R: 4, C: 4, A: [][]int64{
		{5, 8, 4, 8},
		{8, 4, 8, 5},
		{4, 8, 5, 8},
		{8, 5, 8, 4},
	}}
	d := dsl.BareissDet(circ)
	fmt.Printf("matrix: %+v\n", circ.A)
	fmt.Printf("det = %s\n", d.String())

	// Demonstrating cofactor evaluation using expr evaluator
	// Build instance with var M
	inst2 := &dsl.Instance{
		ProblemID: 999,
		Seed:      "demo",
		Vars:      map[string]interface{}{"M": circ},
		Derived:   map[string]interface{}{},
	}
	cof, err := dsl.EvaluateExpression("cofactor(M, 2, 4)", inst2)
	if err != nil {
		fmt.Println("cofactor err:", err)
	} else {
		fmt.Println("cofactor(M,2,4) =", cof)
	}

	// Solve demo: small integer-solution generator example
	fmt.Println("\nDemo: integer_solution generation (A,x -> b then solve)")
	// build A and x
	rand.Seed(time.Now().UnixNano())
	A := &dsl.MatrixInt{R: 3, C: 3}
	A.A = [][]int64{
		{3, 4, 2},
		{-3, -6, 2},
		{6, 12, -1},
	}
	x := &dsl.VectorInt{N: 3, V: []int64{2, -1, 3}}
	b := dsl.NewVectorInt(3)
	for i := 0; i < 3; i++ {
		var s int64 = 0
		for j := 0; j < 3; j++ {
			s += A.A[i][j] * x.V[j]
		}
		b.V[i] = s
	}
	inst3 := &dsl.Instance{Vars: map[string]interface{}{"A": A, "b": b}}
	sol, err := dsl.EvaluateExpression("solve(A,b)", inst3)
	if err != nil {
		fmt.Println("solve err:", err)
	} else {
		fmt.Println("A =", A.A)
		fmt.Println("b =", b.V)
		fmt.Println("solve(A,b) =", ratsToString(sol.([]*big.Rat)))
	}
}

// helper: print renderable objects
func printRendered(m map[string]interface{}) {
	for k, v := range m {
		fmt.Printf("  %s: ", k)
		switch t := v.(type) {
		case *dsl.VectorInt:
			fmt.Printf("%v\n", t.V)
		case *dsl.MatrixInt:
			fmt.Printf("%v\n", t.A)
		case *big.Int:
			fmt.Printf("%s\n", t.String())
		default:
			fmt.Printf("%#v\n", v)
		}
	}
}

func ratsToString(v []*big.Rat) string {
	parts := []string{}
	for _, r := range v {
		parts = append(parts, r.RatString())
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
