# 线性代数 DSL 题库系统

用 JSON/YAML 描述线性代数题目，通过统一 runtime 生成随机题目和答案。

## 核心特性

- **声明式题目定义**：用 DSL 描述题型，不写代码
- **确定性随机**：相同种子生成相同题目，便于复现和测试
- **答案可计算**：表达式系统自动计算派生量和答案
- **填空题支持**：题干包含 `{{blank:ID}}`，答案自动绑定 ID

## 快速开始

```go
import "github.com/neumathe/la-dsl/dsl"

// 从 JSON 解析题目定义
var prob dsl.Problem
json.Unmarshal(data, &prob)

// 生成题目
q, err := dsl.GenerateQuestion(prob, "user123", "server-salt")

// q.Title: 渲染后的题干（含 {{blank:ID}} 占位符）
// q.AnswerFields: [{ID, Expr, Value}, ...]
```

## 目录

1. [题目结构](#题目结构)
2. [随机变量](#随机变量)
3. [生成规则](#生成规则)
4. [派生变量](#派生变量)
5. [表达式系统](#表达式系统)
6. [题面渲染](#题面渲染)
7. [答案结构](#答案结构)
8. [完整示例](#完整示例)

---

## 题目结构

每道题用一个 JSON 对象表示：

```json
{
  "id": 101,
  "version": "v1",
  "title": "向量 β={{beta}} 在基 {{alpha1}}, {{alpha2}}, {{alpha3}} 下的坐标为 {{blank:c1}} {{blank:c2}} {{blank:c3}}",
  "variables": { ... },
  "derived": { ... },
  "answer": { ... }
}
```

| 字段          | 说明                                    |
| ----------- | ------------------------------------- |
| `id`        | 题目 ID                                 |
| `version`   | 版本号（可选）                               |
| `title`     | 题干模板，支持 `{{var}}` 和 `{{blank:ID}}` 占位符 |
| `variables` | 随机变量定义                                |
| `derived`   | 派生变量（由表达式计算）                          |
| `answer`    | 答案定义                                  |
| `meta`      | 元信息（可选）                               |

---

## 随机变量

变量分三类：`scalar`（标量）、`vector`（向量）、`matrix`（矩阵）。

可以用 `fixed` 指定固定值，或用 `generator` 随机生成：

```json
"variables": {
  "A": {
    "kind": "matrix",
    "rows": 3,
    "cols": 3,
    "generator": { "rule": "full_rank", "min": -5, "max": 5 }
  },
  "x": {
    "kind": "vector",
    "size": 3,
    "fixed": [1, 2, 3]
  }
}
```

---

## 生成规则

### `range` - 整数范围

适用于标量/向量/矩阵的随机整数生成：

```json
"generator": { "rule": "range", "min": -5, "max": 5 }
```

### `from_set` - 从集合选取

从指定集合中随机选择：

```json
"generator": { "rule": "from_set", "set": [1, 2, 3, 5, 8] }
```

### `full_rank` - 满秩矩阵

生成可逆矩阵（行列式非零）：

```json
"generator": { "rule": "full_rank", "min": -5, "max": 5 }
```

### `sparse` - 稀疏矩阵

指定密度和值域：

```json
"generator": {
  "rule": "sparse",
  "values": [-9, -6, -3, -1, 1, 2, 4],
  "density": 0.3
}
```

### `integer_solution` - 整数解保证

先随机 `A` 和 `x`，自动计算 `b = A·x`：

```json
"b": {
  "kind": "vector",
  "size": 3,
  "generator": { "rule": "integer_solution", "A": "A", "x": "x" }
}
```

### `lambda_linear_det_zero` - 参数矩阵

生成带参数 λ 的矩阵，保证 `det(A) = 0` 且 λ 为整数：

```json
"A": {
  "kind": "matrix",
  "rows": 3,
  "cols": 3,
  "generator": {
    "rule": "lambda_linear_det_zero",
    "min": -10, "max": 10,
    "lambda_min": -20, "lambda_max": 20,
    "lambda_pos": [2, 2],
    "lambda_var": "lambda"
  }
}
```

---

## 派生变量

用表达式从已有变量计算新变量：

```json
"derived": {
  "beta": "A * x",
  "d": "det(A)",
  "coord": "solve(A, b)"
}
```

---

## 表达式系统

### 已支持的函数

| 函数                  | 说明              |
| ------------------- | --------------- |
| `det(A)`            | 行列式             |
| `rank(A)`           | 秩               |
| `solve(A, b)`       | 求解线性方程组 Ax = b  |
| `cofactor(A, i, j)` | 代数余子式（i,j 从 1 开始） |
| `col(A, i)`         | 第 i 列向量         |
| `A * x`             | 矩阵向量乘法          |
| `basis_cols(A)`     | 返回主元列索引向量       |

### 数组访问

```json
"field_defs": [
  {"id": "x1", "expr": "x[1]"},
  {"id": "x2", "expr": "x[2]"}
]
```

索引从 1 开始，支持访问向量和 `solve()` 返回的有理数数组

---

## 题面渲染

`title` 支持两种占位符：

- `{{var}}`：替换为变量值（LaTeX 格式）
- `{{blank:ID}}`：填空位置

```json
{
  "title": "行列式 D={{A}} 的代数余子式 A24 = {{blank:ans}}",
  "variables": {
    "A": { "kind": "matrix", "rows": 4, "cols": 4, ... }
  },
  "derived": {
    "a24": "cofactor(A, 2, 4)"
  },
  "answer": {
    "field_defs": [
      {"id": "ans", "expr": "a24"}
    ]
  }
}
```

渲染后的 `title` 会将 `{{A}}` 替换为 LaTeX 矩阵，`{{blank:ans}}` 保留给前端插入输入框。

---

## 答案结构

### 方式一：单个表达式

```json
"answer": {
  "expression": "det(A)"
}
```

自动分配 ID `field_1`。

### 方式二：字段列表

```json
"answer": {
  "fields": ["x[1]", "x[2]", "x[3]"]
}
```

自动分配 ID `field_1`, `field_2`, `field_3`。

### 方式三：显式 ID 绑定（推荐）

```json
"answer": {
  "field_defs": [
    {"id": "lambda", "expr": "lambda"},
    {"id": "x1", "expr": "coord[1]"},
    {"id": "x2", "expr": "coord[2]"}
  ]
}
```

ID 与题干中的 `{{blank:ID}}` 对应。

---

## 完整示例

求解带参数的齐次线性方程组：

```json
{
  "id": 1001,
  "version": "v1",
  "title": "已知齐次线性方程组有非零解，求 λ = {{blank:lambda}}",
  "variables": {
    "A": {
      "kind": "matrix",
      "rows": 3,
      "cols": 3,
      "generator": {
        "rule": "lambda_linear_det_zero",
        "min": -10,
        "max": 10,
        "lambda_min": -20,
        "lambda_max": 20,
        "lambda_pos": [2, 2],
        "lambda_var": "lambda"
      }
    },
    "lambda": { "kind": "scalar", "fixed": 0 }
  },
  "answer": {
    "field_defs": [
      {"id": "lambda", "expr": "lambda"}
    ]
  }
}
```

调用：

```go
q, _ := dsl.GenerateQuestion(prob, "seed123", "salt")
// q.Title: "已知齐次线性方程组有非零解，求 λ = {{blank:lambda}}"
// q.AnswerFields[0]: {ID: "lambda", Expr: "lambda", Value: 5}
```