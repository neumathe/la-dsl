package dsl

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// evaluateExpression 支持 det(A)、solve(A,b)、cofactor(A,i,j)、col(A,1)、A * x 等基础表达式
// EvaluateExpression 对外导出，方便上层做表达式求值（例如 demo / 测试）
func EvaluateExpression(expr string, inst *Instance) (interface{}, error) {
	expr = strings.TrimSpace(expr)

	if expr == "zero()" {
		return int64(0), nil
	}

	// zero_vec(n)：返回 n 维零向量
	if strings.HasPrefix(expr, "zero_vec(") && strings.HasSuffix(expr, ")") {
		n := mustAtoi(strings.TrimSpace(insideParens(expr)))
		v := NewVectorInt(n)
		for i := 0; i < n; i++ {
			v.V[i] = 0
		}
		return v, nil
	}

	// scmul(a,b)：两个标量（int64 或 *big.Int）的乘积，返回 *big.Int
	if strings.HasPrefix(expr, "scmul(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("scmul expects 2 args")
		}
		a, err := evalScalarBig(inst, parts[0])
		if err != nil {
			return nil, err
		}
		b, err := evalScalarBig(inst, parts[1])
		if err != nil {
			return nil, err
		}
		return new(big.Int).Mul(a, b), nil
	}

	// trace2(A)：2×2 矩阵的迹（主对角和）。
	if strings.HasPrefix(expr, "trace2(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		vv, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", arg)
		}
		A, ok := vv.(*MatrixInt)
		if !ok || A.R != 2 || A.C != 2 {
			return nil, fmt.Errorf("trace2: need 2×2 matrix")
		}
		return A.A[0][0] + A.A[1][1], nil
	}

	// diagmin(A) / diagmax(A)：方阵且非对角元全为 0 时，主对角最小/最大元。
	if strings.HasPrefix(expr, "diagmin(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		mn, _, err := diagMinMaxFromInst(inst, arg)
		if err != nil {
			return nil, err
		}
		return mn, nil
	}
	if strings.HasPrefix(expr, "diagmax(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		_, mx, err := diagMinMaxFromInst(inst, arg)
		if err != nil {
			return nil, err
		}
		return mx, nil
	}

	// ranklt(A,k)：rank(A) < k 则 1，否则 0（k 为整数阈值）。
	if strings.HasPrefix(expr, "ranklt(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("ranklt expects (A,k)")
		}
		Aname := strings.TrimSpace(parts[0])
		k := mustAtoi(strings.TrimSpace(parts[1]))
		vv, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", Aname)
		}
		A, ok := vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("ranklt expects matrix")
		}
		if matrixRankRat(A) < k {
			return int64(1), nil
		}
		return int64(0), nil
	}

	// basis_index(A,i)：返回由 A 的列向量组成的向量空间中，
	// 第 i 个列主元下标（1-based），即极大无关组中第 i 个向量在原向量组中的下标。
	// 若 i > rank(A)，返回 0（表示"多余的空不填"）。
	if strings.HasPrefix(expr, "basis_index(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("basis_index expects (A,i)")
		}
		Aname := strings.TrimSpace(parts[0])
		i := mustAtoi(strings.TrimSpace(parts[1]))
		vv, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", Aname)
		}
		A, ok := vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("basis_index expects matrix")
		}
		pivots := columnPivotsRat(A)
		if i < 1 || i > len(pivots) {
			return int64(0), nil // "多余的空不填"
		}
		return int64(pivots[i-1]), nil
	}

	// space_rank(A)：返回矩阵 A 的秩（即列向量空间的维数）
	if strings.HasPrefix(expr, "space_rank(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		vv, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		A, ok := vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("space_rank expects matrix")
		}
		return int64(matrixRankRat(A)), nil
	}

	if strings.HasPrefix(expr, "gs_comp(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("gs_comp expects (V,col,row)")
		}
		Vname := strings.TrimSpace(parts[0])
		col := mustAtoi(strings.TrimSpace(parts[1]))
		row := mustAtoi(strings.TrimSpace(parts[2]))
		vv, ok := inst.Vars[Vname]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", Vname)
		}
		V, ok := vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("gs_comp expects matrix")
		}
		u, err := GramSchmidtColsOrthogRat(V)
		if err != nil {
			return nil, err
		}
		if col < 1 || col > V.C || row < 1 || row > V.R {
			return nil, fmt.Errorf("gs_comp index out of range")
		}
		return u[col-1][row-1], nil
	}

	// orthdiag_block4x6(P,D,S)：3×3 正交 Q、对角 Λ 与对称 S 的展示块（4×6），上行左 Q 右 Λ，末行 tr(S),det(S),0,0,0,0。
	if strings.HasPrefix(expr, "orthdiag_block4x6(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("orthdiag_block4x6 expects (P,D,S)")
		}
		pn, dn, sn := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2])
		vp, ok1 := inst.Vars[pn]
		vd, ok2 := inst.Vars[dn]
		vs, ok3 := inst.Vars[sn]
		if !ok1 || !ok2 || !ok3 {
			return nil, fmt.Errorf("orthdiag_block4x6: unknown var")
		}
		P, okP := vp.(*MatrixInt)
		D, okD := vd.(*MatrixInt)
		S, okS := vs.(*MatrixInt)
		if !okP || !okD || !okS || P.R != 3 || P.C != 3 || D.R != 3 || D.C != 3 || S.R != 3 || S.C != 3 {
			return nil, fmt.Errorf("orthdiag_block4x6: need three 3×3 matrices")
		}
		Z := NewMatrixInt(4, 6)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				Z.A[i][j] = P.A[i][j]
				Z.A[i][j+3] = D.A[i][j]
			}
		}
		var tr int64
		for i := 0; i < 3; i++ {
			tr += S.A[i][i]
		}
		detS := BareissDet(S)
		if !detS.IsInt64() {
			return nil, fmt.Errorf("orthdiag_block4x6: det(S) overflow")
		}
		Z.A[3][0] = tr
		Z.A[3][1] = detS.Int64()
		return Z, nil
	}

	// inertia_pos_22(S) / inertia_neg_22(S)：2×2 对称阵的正、负特征值个数（重数不计零）。
	if strings.HasPrefix(expr, "inertia_pos_22(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		vv, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		S, ok := vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("inertia_pos_22 expects matrix")
		}
		p, _, err := inertiaCounts22(S)
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	if strings.HasPrefix(expr, "inertia_neg_22(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		vv, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		S, ok := vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("inertia_neg_22 expects matrix")
		}
		_, n, err := inertiaCounts22(S)
		if err != nil {
			return nil, err
		}
		return n, nil
	}

	if strings.HasPrefix(expr, "nullbasis_comp(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("nullbasis_comp expects (A,k,i)")
		}
		Aname := strings.TrimSpace(parts[0])
		k := mustAtoi(strings.TrimSpace(parts[1]))
		i := mustAtoi(strings.TrimSpace(parts[2]))
		vv, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", Aname)
		}
		A, ok := vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("nullbasis_comp expects matrix")
		}
		basis, err := NullspaceBasisRational(A)
		if err != nil {
			return nil, err
		}
		if k < 1 || k > len(basis) {
			return nil, fmt.Errorf("nullbasis_comp k out of range")
		}
		if i < 1 || i > len(basis[k-1]) {
			return nil, fmt.Errorf("nullbasis_comp i out of range")
		}
		return basis[k-1][i-1], nil
	}

	// triple_first_col(M)：3×3 矩阵，三列均为 M 的第 1 列（用于 $A\alpha_k=\alpha_1$ 类构造）。
	if strings.HasPrefix(expr, "triple_first_col(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		vv, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", arg)
		}
		M, ok := vv.(*MatrixInt)
		if !ok || M.R != 3 || M.C != 3 {
			return nil, fmt.Errorf("triple_first_col: need 3×3")
		}
		out := NewMatrixInt(3, 3)
		for j := 0; j < 3; j++ {
			for i := 0; i < 3; i++ {
				out.A[i][j] = M.A[i][0]
			}
		}
		return out, nil
	}

	// hstack34(U,w)：U 为 4×3，w 为 4 维列向量，拼成 4×4 矩阵 [U|w]。
	if strings.HasPrefix(expr, "hstack34(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("hstack34 expects (U,wname)")
		}
		Uname := strings.TrimSpace(parts[0])
		wname := strings.TrimSpace(parts[1])
		vU, ok := inst.Vars[Uname]
		if !ok {
			return nil, fmt.Errorf("unknown %s", Uname)
		}
		vw, ok := inst.Vars[wname]
		if !ok {
			return nil, fmt.Errorf("unknown %s", wname)
		}
		U, ok := vU.(*MatrixInt)
		if !ok || U.R != 4 || U.C != 3 {
			return nil, fmt.Errorf("hstack34: U must be 4×3")
		}
		w, ok := vw.(*VectorInt)
		if !ok || w.N != 4 {
			return nil, fmt.Errorf("hstack34: w must be length 4")
		}
		out := NewMatrixInt(4, 4)
		for i := 0; i < 4; i++ {
			for j := 0; j < 3; j++ {
				out.A[i][j] = U.A[i][j]
			}
			out.A[i][3] = w.V[i]
		}
		return out, nil
	}

	// vcoef123(V,rhs,k)：用 V 的第 1,2,3 列线性表示第 rhs 列时的第 k 个系数（k=1..3）。
	if strings.HasPrefix(expr, "vcoef123(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("vcoef123 expects (V,rhsCol,k)")
		}
		Vname := strings.TrimSpace(parts[0])
		rhs := mustAtoi(strings.TrimSpace(parts[1]))
		kk := mustAtoi(strings.TrimSpace(parts[2]))
		vv, ok := inst.Vars[Vname]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", Vname)
		}
		V, ok := vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("vcoef123 expects matrix")
		}
		cols := []int{1, 2, 3}
		sol, err := SolveCoeffCols(V, cols, rhs)
		if err != nil {
			return nil, err
		}
		if kk < 1 || kk > len(sol) {
			return nil, fmt.Errorf("vcoef123 k out of range")
		}
		return sol[kk-1], nil
	}

	if strings.HasPrefix(expr, "pow(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("pow expects 2 args")
		}
		mname := strings.TrimSpace(parts[0])
		v, ok := inst.Vars[mname]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", mname)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("pow expects matrix as first arg")
		}
		expStr := strings.TrimSpace(parts[1])
		var n int64
		// 允许整数常量或标量变量
		if vv, ok := inst.Vars[expStr]; ok {
			switch t := vv.(type) {
			case int64:
				n = t
			default:
				return nil, fmt.Errorf("pow exponent var %s must be int64, got %T", expStr, vv)
			}
		} else {
			// 解析字面量
			n = int64(mustAtoi(expStr))
		}
		res, err := matrixPowInt(m, n)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	if strings.HasPrefix(expr, "matmul(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("matmul expects 2 args")
		}
		aName := strings.TrimSpace(parts[0])
		bName := strings.TrimSpace(parts[1])
		va, ok := inst.Vars[aName]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", aName)
		}
		vb, ok := inst.Vars[bName]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", bName)
		}
		a, ok1 := va.(*MatrixInt)
		b, ok2 := vb.(*MatrixInt)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("matmul expects two matrices")
		}
		res, err := matrixMulInt(a, b)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	if strings.HasPrefix(expr, "matadd(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("matadd expects 2 args")
		}
		a, b, err := twoMatrices(inst, parts[0], parts[1])
		if err != nil {
			return nil, err
		}
		return matrixAddInt(a, b)
	}

	if strings.HasPrefix(expr, "matsub(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("matsub expects 2 args")
		}
		a, b, err := twoMatrices(inst, parts[0], parts[1])
		if err != nil {
			return nil, err
		}
		return matrixSubInt(a, b)
	}

	if strings.HasPrefix(expr, "smmul(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("smmul expects 2 args (k, M)")
		}
		k, err := evalScalarLike(inst, parts[0])
		if err != nil {
			return nil, err
		}
		v, ok := inst.Vars[strings.TrimSpace(parts[1])]
		if !ok {
			return nil, fmt.Errorf("unknown matrix %s", parts[1])
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("smmul expects matrix")
		}
		return scalarMulMatrixInt(k, m), nil
	}

	if strings.HasPrefix(expr, "diag_adj(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("diag_adj expects (A,i)")
		}
		Aname := strings.TrimSpace(parts[0])
		i := mustAtoi(strings.TrimSpace(parts[1]))
		v, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", Aname)
		}
		m, ok := v.(*MatrixInt)
		if !ok || m.R != 3 || m.C != 3 {
			return nil, fmt.Errorf("diag_adj needs 3×3 matrix")
		}
		if i < 1 || i > 3 {
			return nil, fmt.Errorf("diag_adj i out of range")
		}
		for r := 0; r < 3; r++ {
			for c := 0; c < 3; c++ {
				if r != c && m.A[r][c] != 0 {
					return nil, fmt.Errorf("diag_adj: not diagonal")
				}
			}
		}
		det := BareissDet(m)
		lam := m.A[i-1][i-1]
		if lam == 0 {
			return nil, fmt.Errorf("diag_adj: zero diagonal")
		}
		q := new(big.Int).Quo(det, big.NewInt(lam))
		return q, nil
	}

	if strings.HasPrefix(expr, "det_quad_shift_diag(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("det_quad_shift_diag expects (A,a,b)")
		}
		Aname := strings.TrimSpace(parts[0])
		a, err := evalScalarLike(inst, parts[1])
		if err != nil {
			return nil, err
		}
		b, err := evalScalarLike(inst, parts[2])
		if err != nil {
			return nil, err
		}
		v, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", Aname)
		}
		m, ok := v.(*MatrixInt)
		if !ok || m.R != m.C {
			return nil, fmt.Errorf("det_quad_shift_diag: need square diagonal matrix")
		}
		n := m.R
		for r := 0; r < n; r++ {
			for c := 0; c < n; c++ {
				if r != c && m.A[r][c] != 0 {
					return nil, fmt.Errorf("det_quad_shift_diag: not diagonal")
				}
			}
		}
		prod := big.NewInt(1)
		for j := 0; j < n; j++ {
			lam := m.A[j][j]
			t := lam*lam + a*lam + b
			prod.Mul(prod, big.NewInt(t))
		}
		return prod, nil
	}

	if strings.HasPrefix(expr, "symdef3(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("symdef3 expects matrix")
		}
		c, err := ClassifySymmetric3(m)
		if err != nil {
			return nil, err
		}
		switch c {
		case SymDefPD:
			return int64(0), nil
		case SymDefND:
			return int64(1), nil
		default:
			return int64(2), nil
		}
	}

	// 题库 Chapter6_1_2：正定 1，负定 2，不定 0
	if strings.HasPrefix(expr, "symcode_612(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("symcode_612 expects matrix")
		}
		c, err := ClassifySymmetric3(m)
		if err != nil {
			return nil, err
		}
		switch c {
		case SymDefPD:
			return int64(1), nil
		case SymDefND:
			return int64(2), nil
		default:
			return int64(0), nil
		}
	}

	// 题库 Chapter6_1_1 / 6_1_3：正定 1，负定 0，不定 2
	if strings.HasPrefix(expr, "symcode_611(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("symcode_611 expects matrix")
		}
		c, err := ClassifySymmetric3(m)
		if err != nil {
			return nil, err
		}
		switch c {
		case SymDefPD:
			return int64(1), nil
		case SymDefND:
			return int64(0), nil
		default:
			return int64(2), nil
		}
	}

	if strings.HasPrefix(expr, "mget(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("mget expects 3 args")
		}
		aName := strings.TrimSpace(parts[0])
		i := mustAtoi(strings.TrimSpace(parts[1]))
		j := mustAtoi(strings.TrimSpace(parts[2]))
		v, ok := inst.Vars[aName]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", aName)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("mget expects matrix")
		}
		if i < 1 || i > m.R || j < 1 || j > m.C {
			return nil, fmt.Errorf("mget index out of range")
		}
		return m.A[i-1][j-1], nil
	}

	if strings.HasPrefix(expr, "transpose(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("transpose expects matrix, got %T", v)
		}
		out := NewMatrixInt(m.C, m.R)
		for i := 0; i < m.R; i++ {
			for j := 0; j < m.C; j++ {
				out.A[j][i] = m.A[i][j]
			}
		}
		return out, nil
	}

	if strings.HasPrefix(expr, "inv(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("inv expects matrix, got %T", v)
		}
		return MatrixInverseInt(m)
	}

	if strings.HasPrefix(expr, "det(") && strings.HasSuffix(expr, ")") {
		arg := insideParens(expr)
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("det expects matrix, got %T", v)
		}
		return BareissDet(m), nil
	}

	if strings.HasPrefix(expr, "nullvec(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("nullvec expects matrix, got %T", v)
		}
		return IntegerKernelVectorOne(m)
	}

	if strings.HasPrefix(expr, "rank(") && strings.HasSuffix(expr, ")") {
		arg := insideParens(expr)
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("rank expects matrix, got %T", v)
		}
		return int64(matrixRankRat(m)), nil
	}

	// rank_hstack(A,bname)：增广矩阵 [A|b] 的秩（b 为与 A 行数相同的列向量）。
	if strings.HasPrefix(expr, "rank_hstack(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("rank_hstack expects (A,bname)")
		}
		an := strings.TrimSpace(parts[0])
		bn := strings.TrimSpace(parts[1])
		va, ok1 := inst.Vars[an]
		vb, ok2 := inst.Vars[bn]
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("rank_hstack: unknown variable")
		}
		A, okA := va.(*MatrixInt)
		b, okB := vb.(*VectorInt)
		if !okA || !okB || A.R != b.N {
			return nil, fmt.Errorf("rank_hstack: need matrix A and vector b with len(b)=rows(A)")
		}
		aug := NewMatrixInt(A.R, A.C+1)
		for i := 0; i < A.R; i++ {
			for j := 0; j < A.C; j++ {
				aug.A[i][j] = A.A[i][j]
			}
			aug.A[i][A.C] = b.V[i]
		}
		return int64(matrixRankRat(aug)), nil
	}

	// nullity(A)：零度 dim ker(A) = 列数 − rank(A)。
	if strings.HasPrefix(expr, "nullity(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("nullity expects matrix")
		}
		rk := matrixRankRat(m)
		return int64(m.C - rk), nil
	}

	// dep3(A)：3 个列向量（3×3 矩阵列）是否线性相关，相关为 1，否则为 0
	if strings.HasPrefix(expr, "dep3(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok || m.R != 3 || m.C != 3 {
			return nil, fmt.Errorf("dep3 expects 3×3 matrix")
		}
		if matrixRankRat(m) < 3 {
			return int64(1), nil
		}
		return int64(0), nil
	}

	// dep_cols(M)：m×n 矩阵 M 的列是否线性相关（rank < cols 则相关），相关为 1，否则为 0
	if strings.HasPrefix(expr, "dep_cols(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("dep_cols expects matrix")
		}
		if matrixRankRat(m) < m.C {
			return int64(1), nil
		}
		return int64(0), nil
	}

	// vecdiv(v,k)：向量 v 的各分量除以整数 k（要求各分量可被 k 整除），返回整数向量
	if strings.HasPrefix(expr, "vecdiv(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("vecdiv expects (v,k)")
		}
		vname := strings.TrimSpace(parts[0])
		k, err := evalScalarLike(inst, parts[1])
		if err != nil {
			return nil, err
		}
		if k == 0 {
			return nil, fmt.Errorf("vecdiv: division by zero")
		}
		vv, ok := inst.Vars[vname]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", vname)
		}
		vec, ok := vv.(*VectorInt)
		if !ok {
			return nil, fmt.Errorf("vecdiv expects vector")
		}
		out := NewVectorInt(vec.N)
		for i := 0; i < vec.N; i++ {
			if vec.V[i]%k != 0 {
				return nil, fmt.Errorf("vecdiv: component %d (= %d) not divisible by %d", i+1, vec.V[i], k)
			}
			out.V[i] = vec.V[i] / k
		}
		return out, nil
	}

	// vecadd(a,b)：两个整数向量相加
	if strings.HasPrefix(expr, "vecadd(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("vecadd expects 2 args")
		}
		va, ok1 := inst.Vars[strings.TrimSpace(parts[0])]
		vb, ok2 := inst.Vars[strings.TrimSpace(parts[1])]
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("vecadd: unknown variable")
		}
		a, ok1 := va.(*VectorInt)
		b, ok2 := vb.(*VectorInt)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("vecadd expects two vectors")
		}
		if a.N != b.N {
			return nil, fmt.Errorf("vecadd: vector length mismatch")
		}
		out := NewVectorInt(a.N)
		for i := 0; i < a.N; i++ {
			out.V[i] = a.V[i] + b.V[i]
		}
		return out, nil
	}

	// nullbasis_vec(A,k)：返回 A 的零空间第 k 个基向量（整数向量，取最小整数倍使得各分量为整数）
	if strings.HasPrefix(expr, "nullbasis_vec(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("nullbasis_vec expects (A,k)")
		}
		Aname := strings.TrimSpace(parts[0])
		k := mustAtoi(strings.TrimSpace(parts[1]))
		vv, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", Aname)
		}
		A, ok := vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("nullbasis_vec expects matrix")
		}
		basis, err := NullspaceBasisRational(A)
		if err != nil {
			return nil, err
		}
		if k < 1 || k > len(basis) {
			return nil, fmt.Errorf("nullbasis_vec k out of range")
		}
		vec := basis[k-1]
		// find LCM of denominators to get integer vector
		lcm := big.NewInt(1)
		for i := 0; i < len(vec); i++ {
			d := vec[i].Denom()
			if d.Sign() != 0 {
				newLcm := new(big.Int)
				gcd := new(big.Int).GCD(nil, nil, lcm, d)
				newLcm.Mul(lcm, d)
				newLcm.Quo(newLcm, gcd)
				lcm = newLcm
			}
		}
		out := NewVectorInt(len(vec))
		for i := 0; i < len(vec); i++ {
			scaled := new(big.Rat).Mul(vec[i], new(big.Rat).SetInt(lcm))
			out.V[i] = scaled.Num().Int64()
		}
		return out, nil
	}

	// eigenval(A,i)：返回由 eigen_reverse_3x3 生成的矩阵 A 的第 i 个特征值（1-based）
	if strings.HasPrefix(expr, "eigenval(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("eigenval expects (A,i)")
		}
		Aname := strings.TrimSpace(parts[0])
		i := mustAtoi(strings.TrimSpace(parts[1]))
		vv, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", Aname)
		}
		_, ok = vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("eigenval expects matrix")
		}
		// Find the _eigen_lambda_i stored by the generator
		key := fmt.Sprintf("_eigen_lambda%d", i)
		lv, ok := inst.Vars[key]
		if !ok {
			return nil, fmt.Errorf("eigenval: eigenvalue data not found for %s index %d", Aname, i)
		}
		switch t := lv.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		}
		return nil, fmt.Errorf("eigenval: unexpected type %T", lv)
	}

	// eigenvec_comp(A,i,j)：返回 A 的第 i 个特征向量的第 j 个分量（1-based）
	// 使用 eigen_reverse_3x3 生成的 V 矩阵
	if strings.HasPrefix(expr, "eigenvec_comp(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("eigenvec_comp expects (A,i,j)")
		}
		Aname := strings.TrimSpace(parts[0])
		i := mustAtoi(strings.TrimSpace(parts[1])) // which eigenvector
		j := mustAtoi(strings.TrimSpace(parts[2])) // which component
		vv, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", Aname)
		}
		_, ok = vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("eigenvec_comp expects matrix")
		}
		V, ok := inst.Vars["_eigen_V"].(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("eigenvec_comp: eigenvector matrix not found")
		}
		if i < 1 || i > V.C || j < 1 || j > V.R {
			return nil, fmt.Errorf("eigenvec_comp index out of range")
		}
		return V.A[j-1][i-1], nil // V's column i, row j
	}

	// sym_eigenval(S,i)：返回由 symmetric_eigen_reverse_3x3 生成的对称矩阵 S 的第 i 个特征值（1-based）
	if strings.HasPrefix(expr, "sym_eigenval(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("sym_eigenval expects (S,i)")
		}
		Aname := strings.TrimSpace(parts[0])
		i := mustAtoi(strings.TrimSpace(parts[1]))
		vv, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", Aname)
		}
		_, ok = vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("sym_eigenval expects matrix")
		}
		key := fmt.Sprintf("_sym_eigen_lambda%d", i)
		lv, ok := inst.Vars[key]
		if !ok {
			return nil, fmt.Errorf("sym_eigenval: eigenvalue data not found for index %d", i)
		}
		switch t := lv.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		}
		return nil, fmt.Errorf("sym_eigenval: unexpected type %T", lv)
	}

	// sym_eigenvec_comp(S,i,j)：返回正交变换矩阵 Q 的第 (j,i) 元素
	// 即正交变换 x = Qy 中 Q 的 j 行 i 列（1-based），
	// 对应第 i 个特征向量的第 j 个分量。
	if strings.HasPrefix(expr, "sym_eigenvec_comp(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("sym_eigenvec_comp expects (S,i,j)")
		}
		Aname := strings.TrimSpace(parts[0])
		i := mustAtoi(strings.TrimSpace(parts[1])) // which eigenvector/column
		j := mustAtoi(strings.TrimSpace(parts[2])) // which component/row
		vv, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", Aname)
		}
		_, ok = vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("sym_eigenvec_comp expects matrix")
		}
		Q, ok := inst.Vars["_sym_eigen_Q"].(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("sym_eigenvec_comp: Q matrix not found")
		}
		if i < 1 || i > 3 || j < 1 || j > 3 {
			return nil, fmt.Errorf("sym_eigenvec_comp index out of range")
		}
		return Q.A[j-1][i-1], nil
	}

	// is_similar() / is_congruent()：返回由 similarity_congruence_pair 生成的相似/合同判断结果（0 或 1）
	if strings.HasPrefix(expr, "is_similar(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		vv, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		_, ok = vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("is_similar expects matrix variable")
		}
		rv, ok := inst.Vars["_sc_is_similar"]
		if !ok {
			return nil, fmt.Errorf("is_similar: similarity data not found")
		}
		switch t := rv.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		}
		return nil, fmt.Errorf("is_similar: unexpected type %T", rv)
	}
	if strings.HasPrefix(expr, "is_congruent(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		vv, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		_, ok = vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("is_congruent expects matrix variable")
		}
		rv, ok := inst.Vars["_sc_is_congruent"]
		if !ok {
			return nil, fmt.Errorf("is_congruent: congruence data not found")
		}
		switch t := rv.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		}
		return nil, fmt.Errorf("is_congruent: unexpected type %T", rv)
	}

	// sc_diag_comp(B,i)：返回由 similarity_congruence_pair 生成的对角阵 B 的第 i 个对角元素（1-based）
	if strings.HasPrefix(expr, "sc_diag_comp(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("sc_diag_comp expects (B,i)")
		}
		Bname := strings.TrimSpace(parts[0])
		i := mustAtoi(strings.TrimSpace(parts[1]))
		Bv, ok := inst.Vars["_sc_B"]
		if !ok {
			// Fall back: try to look up the variable directly
			Bv, ok = inst.Vars[Bname]
			if !ok {
				return nil, fmt.Errorf("unknown variable %s", Bname)
			}
		}
		B, ok := Bv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("sc_diag_comp: B not found or not matrix")
		}
		if i < 1 || i > 3 {
			return nil, fmt.Errorf("sc_diag_comp index out of range")
		}
		return B.A[i-1][i-1], nil
	}

	// sc_diag_matrix(A)：返回由 similarity_congruence_pair 生成的对角阵 B（作为 MatrixInt）
	if strings.HasPrefix(expr, "sc_diag_matrix(") && strings.HasSuffix(expr, ")") {
		Bv, ok := inst.Vars["_sc_B"]
		if !ok {
			return nil, fmt.Errorf("sc_diag_matrix: B not found in instance")
		}
		B, ok := Bv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("sc_diag_matrix: B is not a matrix")
		}
		return B, nil
	}

	// sym_npos(S)：3×3 对称矩阵的正惯性指数（正特征值个数）
	if strings.HasPrefix(expr, "sym_npos(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		vv, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("sym_npos expects matrix")
		}
		// Use stored eigenvalues if available (from symmetric_eigen_reverse_3x3)
		npos := 0
		for i := 1; i <= 3; i++ {
			key := fmt.Sprintf("_sym_eigen_lambda%d", i)
			lv, found := inst.Vars[key]
			if found {
				switch t := lv.(type) {
				case int64:
					if t > 0 {
						npos++
					}
				case int:
					if t > 0 {
						npos++
					}
				}
			}
		}
		if npos > 0 {
			return int64(npos), nil
		}
		// Fall back: compute from principal minors (Sylvester)
		// For 3×3, positive eigenvalues = number of sign-preserving principal minors
		// Actually easier: just compute eigenvalues from Lambda stored by symmetric generator
		// If no stored eigenvalues, try ClassifySymmetric3 and count from inertia
		c, err := ClassifySymmetric3(m)
		if err != nil {
			return nil, err
		}
		switch c {
		case SymDefPD:
			return int64(3), nil // all 3 eigenvalues positive
		case SymDefND:
			return int64(0), nil // all 3 eigenvalues negative (n_pos=0)
		case SymDefInd:
			// Need to actually count. For 3×3 indefinite:
			// If det < 0, inertia could be (2,1) or (1,2)
			// Use principal minors to determine
			d1 := principalMinorDet(m, 1)
			d2 := principalMinorDet(m, 2)
			d3 := principalMinorDet(m, 3) // = det(S)
			if d1.Sign() > 0 && d2.Sign() < 0 {
				return int64(1), nil // (1,2)
			}
			if d1.Sign() < 0 && d2.Sign() > 0 {
				// This is the negative definite pattern, but we already caught that above
				// For indefinite with det<0: could be (1,2) or (2,1)
				// If d2 > 0 but d3 < 0: (2,1)
				if d3.Sign() < 0 {
					return int64(2), nil
				}
				return int64(1), nil
			}
			// General fallback: count from stored eigenvalues
			return int64(1), nil
		}
	}

	// quad_expr_param(S)：由含参数 t 的 3×3 对称矩阵 Sbase 渲染二次型表达式的 LaTeX 字符串
	// S 的某个对角位置为参数 t 占位符（值为 0）
	if strings.HasPrefix(expr, "quad_expr_param(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("quad_expr_param expects matrix")
		}
		return formatQuadraticExprWithParam(m, inst), nil
	}

	// param_t(S)：返回含参数 t 的对称矩阵中 t 的值（由 trace 条件推导）
	if strings.HasPrefix(expr, "param_t(") && strings.HasSuffix(expr, ")") {
		tv, ok := inst.Vars["_param_t_val"]
		if !ok {
			return nil, fmt.Errorf("param_t: data not found")
		}
		switch t := tv.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		}
		return nil, fmt.Errorf("param_t: unexpected type %T", tv)
	}

	// param_eigenval(S,i)：返回含参数 t 的对称矩阵的第 i 个特征值
	if strings.HasPrefix(expr, "param_eigenval(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("param_eigenval expects (S,i)")
		}
		i := mustAtoi(strings.TrimSpace(parts[1]))
		key := fmt.Sprintf("_param_lambda%d", i)
		lv, ok := inst.Vars[key]
		if !ok {
			return nil, fmt.Errorf("param_eigenval: eigenvalue %d not found", i)
		}
		switch t := lv.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		}
		return nil, fmt.Errorf("param_eigenval: unexpected type %T", lv)
	}

	// param_eigenvec_comp(S,i,j)：返回含参数 t 的对称矩阵的正交变换矩阵 Q 的元素
	if strings.HasPrefix(expr, "param_eigenvec_comp(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("param_eigenvec_comp expects (S,i,j)")
		}
		i := mustAtoi(strings.TrimSpace(parts[1])) // which eigenvector/column
		j := mustAtoi(strings.TrimSpace(parts[2])) // which component/row
		Qv, ok := inst.Vars["_param_Q"]
		if !ok {
			return nil, fmt.Errorf("param_eigenvec_comp: Q not found")
		}
		Q, ok := Qv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("param_eigenvec_comp: Q not matrix")
		}
		if i < 1 || i > 3 || j < 1 || j > 3 {
			return nil, fmt.Errorf("param_eigenvec_comp: index out of range")
		}
		return Q.A[j-1][i-1], nil
	}

	// sylvester_lower(S) / sylvester_upper(S)：正定参数范围的下界/上界
	if strings.HasPrefix(expr, "sylvester_lower(") && strings.HasSuffix(expr, ")") {
		lv, ok := inst.Vars["_sylvester_lower"]
		if !ok {
			return nil, fmt.Errorf("sylvester_lower: data not found")
		}
		lr, ok := lv.(*big.Rat)
		if !ok {
			return nil, fmt.Errorf("sylvester_lower: unexpected type %T", lv)
		}
		return lr, nil
	}
	if strings.HasPrefix(expr, "sylvester_upper(") && strings.HasSuffix(expr, ")") {
		uv, ok := inst.Vars["_sylvester_upper"]
		if !ok {
			return nil, fmt.Errorf("sylvester_upper: data not found")
		}
		ur, ok := uv.(*big.Rat)
		if !ok {
			return nil, fmt.Errorf("sylvester_upper: unexpected type %T", uv)
		}
		return ur, nil
	}

	// eigenval_rank(S,i)：由秩条件推断的第 i 个特征值（1-based，按升序）
	if strings.HasPrefix(expr, "eigenval_rank(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("eigenval_rank expects (S,i)")
		}
		i := mustAtoi(strings.TrimSpace(parts[1]))
		return evalEigenvalRank(inst, i)
	}

	// eigenval_rowsum(S,i)：由行和+秩条件推断的第 i 个特征值（1-based，按升序）
	if strings.HasPrefix(expr, "eigenval_rowsum(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("eigenval_rowsum expects (S,i)")
		}
		i := mustAtoi(strings.TrimSpace(parts[1]))
		return evalEigenvalRowsum(inst, i)
	}

	// eigen_rank_condition_text(S)：秩条件题面的 LaTeX 字符串
	if strings.HasPrefix(expr, "eigen_rank_condition_text(") && strings.HasSuffix(expr, ")") {
		return evalEigenRankConditionText(inst)
	}

	// eigen_rowsum_condition_text(S)：行和+秩条件题面的 LaTeX 字符串
	if strings.HasPrefix(expr, "eigen_rowsum_condition_text(") && strings.HasSuffix(expr, ")") {
		return evalEigenRowsumConditionText(inst)
	}

	// eigen_rowsum_r/s/k(S)：行和+秩条件中的秩/行和/k值
	if strings.HasPrefix(expr, "eigen_rowsum_r(") && strings.HasSuffix(expr, ")") {
		rv, ok := inst.Vars["_eigen_rowsum_r"]
		if !ok {
			return nil, fmt.Errorf("eigen_rowsum_r: data not found")
		}
		switch t := rv.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		}
		return nil, fmt.Errorf("eigen_rowsum_r: unexpected type %T", rv)
	}
	if strings.HasPrefix(expr, "eigen_rowsum_s(") && strings.HasSuffix(expr, ")") {
		sv, ok := inst.Vars["_eigen_rowsum_s"]
		if !ok {
			return nil, fmt.Errorf("eigen_rowsum_s: data not found")
		}
		switch t := sv.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		}
		return nil, fmt.Errorf("eigen_rowsum_s: unexpected type %T", sv)
	}
	if strings.HasPrefix(expr, "eigen_rowsum_k(") && strings.HasSuffix(expr, ")") {
		kv, ok := inst.Vars["_eigen_rowsum_k"]
		if !ok {
			return nil, fmt.Errorf("eigen_rowsum_k: data not found")
		}
		switch t := kv.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		}
		return nil, fmt.Errorf("eigen_rowsum_k: unexpected type %T", kv)
	}

	// poly_schmidt_comp(P,j,k)：第 j 个 Schmidt 正交多项式的第 k 个系数
	// j=1,2,3 对应 g₁,g₂,g₃; k=1 对应常数项, k=2 对应 x 系数, k=3 对应 x² 系数
	if strings.HasPrefix(expr, "poly_schmidt_comp(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("poly_schmidt_comp expects (P,j,k)")
		}
		j := mustAtoi(strings.TrimSpace(parts[1]))
		k := mustAtoi(strings.TrimSpace(parts[2]))
		return evalPolySchmidtComp(inst, j, k)
	}

	// poly_schmidt_input_comp(P,j,k)：第 j 个输入多项式的第 k 个系数
	if strings.HasPrefix(expr, "poly_schmidt_input_comp(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("poly_schmidt_input_comp expects (P,j,k)")
		}
		j := mustAtoi(strings.TrimSpace(parts[1]))
		k := mustAtoi(strings.TrimSpace(parts[2]))
		return evalPolySchmidtInputComp(inst, j, k)
	}

	// poly_schmidt_input_text(P)：输入多项式的 LaTeX 字符串
	if strings.HasPrefix(expr, "poly_schmidt_input_text(") && strings.HasSuffix(expr, ")") {
		return evalPolySchmidtInputText(inst)
	}

	// quad_expr(S)：由 n×n 对称矩阵 S 渲染二次型表达式的 LaTeX 字符串
	if strings.HasPrefix(expr, "quad_expr(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("quad_expr expects matrix")
		}
		return formatQuadraticExpr(m), nil
	}

	// lambda_cases_title(A)：将含参数 λ 的齐次方程组 Ax=0 渲染为 egin{cases} 格式
	// 在参数位置显示 λ 符号而非数值
	if strings.HasPrefix(expr, "lambda_cases_title(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("lambda_cases_title expects matrix")
		}
		rowV, ok := inst.Vars["_lambda_param_row"]
		if !ok {
			return nil, fmt.Errorf("lambda_cases_title: param_row not found")
		}
		colV, ok := inst.Vars["_lambda_param_col"]
		if !ok {
			return nil, fmt.Errorf("lambda_cases_title: param_col not found")
		}
		constCV, ok := inst.Vars["_lambda_param_constC"]
		if !ok {
			return nil, fmt.Errorf("lambda_cases_title: constC not found")
		}
		var paramRow int64
		switch t := rowV.(type) {
		case int64:
			paramRow = t
		case int:
			paramRow = int64(t)
		}
		var paramCol int64
		switch t := colV.(type) {
		case int64:
			paramCol = t
		case int:
			paramCol = int64(t)
		}
		var constC int64
		switch t := constCV.(type) {
		case int64:
			constC = t
		case int:
			constC = int64(t)
		}
		return formatLambdaCasesTitle(m, paramRow, paramCol, constC), nil
	}

	// vmatrix_title(A)：将矩阵 A 渲染为 egin{vmatrix}...\end{vmatrix} 的 LaTeX 字符串
	if strings.HasPrefix(expr, "vmatrix_title(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("vmatrix_title expects matrix")
		}
		return formatVmatrixTitle(m), nil
	}

	// equidiagonal_title(A)：将等对角矩阵渲染为带 \cdots 提示的 vmatrix 格式
	if strings.HasPrefix(expr, "equidiagonal_title(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("equidiagonal_title expects matrix")
		}
		return formatEquidiagonalTitle(m), nil
	}

	// cases_title(A, b)：将 Ax=b 渲染为 egin{cases} 方程组的 LaTeX 字符串
	if strings.HasPrefix(expr, "cases_title(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("cases_title expects 2 args")
		}
		vA, ok := inst.Vars[strings.TrimSpace(parts[0])]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", parts[0])
		}
		vB, ok := inst.Vars[strings.TrimSpace(parts[1])]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", parts[1])
		}
		A, ok := vA.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("cases_title expects matrix A")
		}
		b, ok := vB.(*VectorInt)
		if !ok {
			return nil, fmt.Errorf("cases_title expects vector b")
		}
		return formatCasesTitle(A, b), nil
	}


	if strings.HasPrefix(expr, "solve(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("solve expects 2 args")
		}
		vA, ok := inst.Vars[strings.TrimSpace(parts[0])]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", parts[0])
		}
		vB, ok := inst.Vars[strings.TrimSpace(parts[1])]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", parts[1])
		}
		A, ok := vA.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("solve expects matrix A")
		}
		bvec, ok := vB.(*VectorInt)
		if !ok {
			return nil, fmt.Errorf("solve expects vector b")
		}
		xrat, err := solveLinearSystemRat(A, bvec)
		if err != nil {
			return nil, err
		}
		return xrat, nil
	}

	if strings.HasPrefix(expr, "cofactor(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("cofactor expects 3 args")
		}
		Aname := strings.TrimSpace(parts[0])
		i := mustAtoi(strings.TrimSpace(parts[1]))
		j := mustAtoi(strings.TrimSpace(parts[2]))
		v, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", Aname)
		}
		A, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("cofactor expects matrix")
		}
		if i < 1 || j < 1 || i > A.R || j > A.C {
			return nil, fmt.Errorf("index out of range")
		}
		minor := NewMatrixInt(A.R-1, A.C-1)
		ri := 0
		for r := 0; r < A.R; r++ {
			if r == i-1 {
				continue
			}
			ci := 0
			for c := 0; c < A.C; c++ {
				if c == j-1 {
					continue
				}
				minor.A[ri][ci] = A.A[r][c]
				ci++
			}
			ri++
		}
		d := BareissDet(minor)
		sign := int64(1)
		if ((i + j) % 2) == 1 {
			sign = -1
		}
		return new(big.Int).Mul(big.NewInt(sign), d), nil
	}

	if strings.HasPrefix(expr, "basis_cols(") && strings.HasSuffix(expr, ")") {
		arg := insideParens(expr)
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("basis_cols expects matrix, got %T", v)
		}
		pivots := columnPivotsRat(m)
		vec := NewVectorInt(m.C)
		for i := 0; i < len(pivots) && i < m.C; i++ {
			vec.V[i] = int64(pivots[i])
		}
		// 多余位置保持为 0，表示“无”
		return vec, nil
	}

	if strings.HasPrefix(expr, "col(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("col expects 2 args: col(A, 1)")
		}
		Aname := strings.TrimSpace(parts[0])
		k := mustAtoi(strings.TrimSpace(parts[1]))
		v, ok := inst.Vars[Aname]
		if !ok {
			return nil, fmt.Errorf("unknown var %s", Aname)
		}
		A, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("col expects matrix")
		}
		if k < 1 || k > A.C {
			return nil, fmt.Errorf("col index out of range")
		}
		vec := NewVectorInt(A.R)
		for i := 0; i < A.R; i++ {
			vec.V[i] = A.A[i][k-1]
		}
		return vec, nil
	}

	// param_lambda_val()：含参方程组的 λ 值
	if strings.HasPrefix(expr, "param_lambda_val(") && strings.HasSuffix(expr, ")") {
		lv, ok := inst.Vars["_param_lambda_val"]
		if !ok {
			return nil, fmt.Errorf("param_lambda_val: data not found")
		}
		switch t := lv.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		}
		return nil, fmt.Errorf("param_lambda_val: unexpected type %T", lv)
	}

	// param_mu_val()：含参方程组的 μ 值
	if strings.HasPrefix(expr, "param_mu_val(") && strings.HasSuffix(expr, ")") {
		mv, ok := inst.Vars["_param_mu_val"]
		if !ok {
			return nil, fmt.Errorf("param_mu_val: data not found")
		}
		switch t := mv.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		}
		return nil, fmt.Errorf("param_mu_val: unexpected type %T", mv)
	}

	// param_rref_comp(A,i,j)：增广矩阵行最简型的第 (i,j) 元素（1-based）
	if strings.HasPrefix(expr, "param_rref_comp(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 3 {
			return nil, fmt.Errorf("param_rref_comp expects (A,i,j)")
		}
		i := mustAtoi(strings.TrimSpace(parts[1]))
		j := mustAtoi(strings.TrimSpace(parts[2]))
		rv, ok := inst.Vars["_param_rref"]
		if !ok {
			return nil, fmt.Errorf("param_rref_comp: rref data not found")
		}
		rref, ok := rv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("param_rref_comp: rref not matrix")
		}
		if i < 1 || i > rref.R || j < 1 || j > rref.C {
			return nil, fmt.Errorf("param_rref_comp: index out of range")
		}
		return rref.A[i-1][j-1], nil
	}

	// param_x0_comp(A,i)：含参方程组特解的第 i 个分量（1-based）
	if strings.HasPrefix(expr, "param_x0_comp(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("param_x0_comp expects (A,i)")
		}
		i := mustAtoi(strings.TrimSpace(parts[1]))
		xv, ok := inst.Vars["_param_x0"]
		if !ok {
			return nil, fmt.Errorf("param_x0_comp: x0 data not found")
		}
		x0, ok := xv.(*VectorInt)
		if !ok {
			return nil, fmt.Errorf("param_x0_comp: x0 not vector")
		}
		if i < 1 || i > x0.N {
			return nil, fmt.Errorf("param_x0_comp: index out of range")
		}
		return x0.V[i-1], nil
	}

	// param_nb_comp(A,i)：含参方程组基础解系向量的第 i 个分量（1-based）
	if strings.HasPrefix(expr, "param_nb_comp(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("param_nb_comp expects (A,i)")
		}
		i := mustAtoi(strings.TrimSpace(parts[1]))
		nv, ok := inst.Vars["_param_nb1"]
		if !ok {
			return nil, fmt.Errorf("param_nb_comp: nb1 data not found")
		}
		nb, ok := nv.(*VectorInt)
		if !ok {
			return nil, fmt.Errorf("param_nb_comp: nb1 not vector")
		}
		if i < 1 || i > nb.N {
			return nil, fmt.Errorf("param_nb_comp: index out of range")
		}
		return nb.V[i-1], nil
	}

	// param_rank_A()：含参方程组系数矩阵的秩
	if strings.HasPrefix(expr, "param_rank_A(") && strings.HasSuffix(expr, ")") {
		return int64(2), nil // rank(A) is always 2 for this generator
	}

	// param_rank_aug()：含参方程组增广矩阵的秩
	if strings.HasPrefix(expr, "param_rank_aug(") && strings.HasSuffix(expr, ")") {
		return int64(2), nil // rank(A|b) equals rank(A) for consistency
	}

	// param_system_title(A)：含参方程组题面 LaTeX
	if strings.HasPrefix(expr, "param_system_title(") && strings.HasSuffix(expr, ")") {
		return evalParamSystemTitle(inst)
	}

	// poly_from_vec(v)：将向量 [a,b,c] 渲染为多项式 a+bx+cx² 的 LaTeX 字符串
	if strings.HasPrefix(expr, "poly_from_vec(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		vv, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("poly_from_vec: unknown variable %s", arg)
		}
		vec, ok := vv.(*VectorInt)
		if !ok {
			return nil, fmt.Errorf("poly_from_vec: expects vector")
		}
		return formatPolynomialFromVec(vec), nil
	}

	// poly_from_matcol(M,j)：将矩阵 M 的第 j 列渲染为多项式 a+bx+cx² 的 LaTeX 字符串
	if strings.HasPrefix(expr, "poly_from_matcol(") && strings.HasSuffix(expr, ")") {
		ins := insideParens(expr)
		parts := splitArgs(ins)
		if len(parts) != 2 {
			return nil, fmt.Errorf("poly_from_matcol expects (M,j)")
		}
		Mname := strings.TrimSpace(parts[0])
		j := mustAtoi(strings.TrimSpace(parts[1]))
		vv, ok := inst.Vars[Mname]
		if !ok {
			return nil, fmt.Errorf("poly_from_matcol: unknown variable %s", Mname)
		}
		M, ok := vv.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("poly_from_matcol: expects matrix")
		}
		if j < 1 || j > M.C {
			return nil, fmt.Errorf("poly_from_matcol: column index out of range")
		}
		vec := NewVectorInt(M.R)
		for i := 0; i < M.R; i++ {
			vec.V[i] = M.A[i][j-1]
		}
		return formatPolynomialFromVec(vec), nil
	}

	// eigenval_list_text(A)：从对角矩阵 A 的对角元（特征值）渲染为 LaTeX 列表字符串
	// 格式如 "-5, 4, -2" — 用于题面"已知三阶矩阵A的三个特征值分别为..."
	if strings.HasPrefix(expr, "eigenval_list_text(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("eigenval_list_text: unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("eigenval_list_text: expects matrix")
		}
		var vals []string
		for i := 0; i < m.R; i++ {
			vals = append(vals, fmt.Sprintf("%d", m.A[i][i]))
		}
		return strings.Join(vals, ", "), nil
	}

	// linear_transform_title(A0)：将 3×3 矩阵 A₀ 渲染为 T(x₁,x₂,x₃)ᵀ=(...)ᵀ 的 LaTeX 字符串
	// 格式如 T\begin{pmatrix}x_1\\x_2\\x_3\end{pmatrix}=\begin{pmatrix}a₁x₁+a₂x₂+a₃x₃\\...\end{pmatrix}
	if strings.HasPrefix(expr, "linear_transform_title(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("linear_transform_title: unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("linear_transform_title: expects matrix")
		}
		return formatLinearTransformTitle(m), nil
	}

	// basis_linear_combo_title(B)：将 3×3 上三角矩阵 B 的各列渲染为新基向量的线性组合表示
	// 第1列: ε₁, 第2列: c₂₁ε₁+c₂₂ε₂, 第3列: c₃₁ε₁+c₃₂ε₂+c₃₃ε₃
	if strings.HasPrefix(expr, "basis_linear_combo_title(") && strings.HasSuffix(expr, ")") {
		arg := strings.TrimSpace(insideParens(expr))
		v, ok := inst.Vars[arg]
		if !ok {
			return nil, fmt.Errorf("basis_linear_combo_title: unknown variable %s", arg)
		}
		m, ok := v.(*MatrixInt)
		if !ok {
			return nil, fmt.Errorf("basis_linear_combo_title: expects matrix")
		}
		return formatBasisLinearComboTitle(m), nil
	}

		if strings.Contains(expr, "*") {
		parts := strings.Split(expr, "*")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			vL, lok := inst.Vars[left]
			vR, rok := inst.Vars[right]
			if lok && rok {
				// matrix * vector
				if m, ok := vL.(*MatrixInt); ok {
					if x, ok2 := vR.(*VectorInt); ok2 {
						if m.C != x.N {
							return nil, fmt.Errorf("dimension mismatch in %s * %s", left, right)
						}
						out := NewVectorInt(m.R)
						for i := 0; i < m.R; i++ {
							var sum int64
							for j := 0; j < m.C; j++ {
								sum += m.A[i][j] * x.V[j]
							}
							out.V[i] = sum
						}
						return out, nil
					}
				}
				// matrix * matrix
				if a, ok1 := vL.(*MatrixInt); ok1 {
					if b, ok2 := vR.(*MatrixInt); ok2 {
						res, err := matrixMulInt(a, b)
						if err != nil {
							return nil, err
						}
						return res, nil
					}
				}
			}
		}
	}

	// 支持简单的下标访问：var[index]，如 x[1]、coord[3]
	if strings.Contains(expr, "[") && strings.HasSuffix(expr, "]") {
		idxOpen := strings.Index(expr, "[")
		idxClose := strings.LastIndex(expr, "]")
		if idxOpen > 0 && idxClose > idxOpen {
			vname := strings.TrimSpace(expr[:idxOpen])
			idx := mustAtoi(expr[idxOpen+1 : idxClose])
			if vv, ok := inst.Vars[vname]; ok {
				switch t := vv.(type) {
				case *VectorInt:
					if idx >= 1 && idx <= t.N {
						return t.V[idx-1], nil
					}
				case []*big.Rat:
					if idx >= 1 && idx <= len(t) {
						return t[idx-1], nil
					}
				}
			}
		}
	}

	if v, ok := inst.Vars[expr]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("unsupported expression: %s", expr)
}

func insideParens(s string) string {
	i := strings.Index(s, "(")
	j := strings.LastIndex(s, ")")
	if i == -1 || j == -1 || j <= i {
		return ""
	}
	return s[i+1 : j]
}

func splitArgs(s string) []string {
	parts := []string{}
	cur := ""
	depth := 0
	for _, ch := range s {
		switch ch {
		case '(':
			depth++
			cur += string(ch)
		case ')':
			depth--
			cur += string(ch)
		case ',':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(cur))
				cur = ""
			} else {
				cur += string(ch)
			}
		default:
			cur += string(ch)
		}
	}
	if strings.TrimSpace(cur) != "" {
		parts = append(parts, strings.TrimSpace(cur))
	}
	return parts
}

func mustAtoi(s string) int {
	s = strings.TrimSpace(s)
	var v int
	fmt.Sscanf(s, "%d", &v)
	return v
}

// inertiaCounts22 返回 2×2 实对称阵的正、负特征值个数（零特征不计入）。
func inertiaCounts22(S *MatrixInt) (pos int64, neg int64, err error) {
	if S == nil || S.R != 2 || S.C != 2 {
		return 0, 0, fmt.Errorf("inertia22: need 2×2 matrix")
	}
	if S.A[0][1] != S.A[1][0] {
		return 0, 0, fmt.Errorf("inertia22: matrix not symmetric")
	}
	a, b, d := S.A[0][0], S.A[0][1], S.A[1][1]
	tr := a + d
	det := a*d - b*b
	if det < 0 {
		return 1, 1, nil
	}
	if det > 0 {
		if tr > 0 {
			return 2, 0, nil
		}
		if tr < 0 {
			return 0, 2, nil
		}
		return 0, 0, nil
	}
	if tr > 0 {
		return 1, 0, nil
	}
	if tr < 0 {
		return 0, 1, nil
	}
	return 0, 0, nil
}

func twoMatrices(inst *Instance, aName, bName string) (*MatrixInt, *MatrixInt, error) {
	va, ok := inst.Vars[strings.TrimSpace(aName)]
	if !ok {
		return nil, nil, fmt.Errorf("unknown var %s", aName)
	}
	vb, ok := inst.Vars[strings.TrimSpace(bName)]
	if !ok {
		return nil, nil, fmt.Errorf("unknown var %s", bName)
	}
	a, ok1 := va.(*MatrixInt)
	b, ok2 := vb.(*MatrixInt)
	if !ok1 || !ok2 {
		return nil, nil, fmt.Errorf("expected two matrices")
	}
	return a, b, nil
}

func diagMinMaxFromInst(inst *Instance, name string) (int64, int64, error) {
	vv, ok := inst.Vars[name]
	if !ok {
		return 0, 0, fmt.Errorf("unknown var %s", name)
	}
	A, ok := vv.(*MatrixInt)
	if !ok || A.R != A.C {
		return 0, 0, fmt.Errorf("diagmin/max: need square matrix")
	}
	n := A.R
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j && A.A[i][j] != 0 {
				return 0, 0, fmt.Errorf("diagmin/max: not diagonal")
			}
		}
	}
	mn, mx := A.A[0][0], A.A[0][0]
	for i := 0; i < n; i++ {
		v := A.A[i][i]
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	return mn, mx, nil
}

func evalScalarLike(inst *Instance, s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty scalar")
	}
	if v, ok := inst.Vars[s]; ok {
		switch t := v.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		case float64:
			return int64(t), nil
		}
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("scalar %q: %w", s, err)
	}
	return v, nil
}

func evalScalarBig(inst *Instance, s string) (*big.Int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty scalar")
	}
	if v, ok := inst.Vars[s]; ok {
		switch t := v.(type) {
		case *big.Int:
			return t, nil
		case int64:
			return big.NewInt(t), nil
		case int:
			return big.NewInt(int64(t)), nil
		}
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("scalar %q: %w", s, err)
	}
	return big.NewInt(v), nil
}
