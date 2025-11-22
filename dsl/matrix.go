package dsl

// MatrixInt 使用 int64 存储矩阵元素，生成阶段使用
type MatrixInt struct {
	R int
	C int
	A [][]int64
}

func NewMatrixInt(r, c int) *MatrixInt {
	m := &MatrixInt{R: r, C: c}
	m.A = make([][]int64, r)
	for i := 0; i < r; i++ {
		m.A[i] = make([]int64, c)
	}
	return m
}

// VectorInt 使用 int64 存储向量元素
type VectorInt struct {
	N int
	V []int64
}

func NewVectorInt(n int) *VectorInt {
	v := &VectorInt{N: n}
	v.V = make([]int64, n)
	return v
}
