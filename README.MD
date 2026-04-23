# 线性代数 DSL 题库系统

用 JSON/YAML 描述线性代数题目，通过统一 runtime 生成随机题目和答案。

### 环境与依赖

- **Go**：`1.23+`（见根目录 `go.mod`）
- **模块路径**：`github.com/neumathe/la-dsl`（子路径包：`dsl`、`bank`、`ladsl`）

克隆后在本目录执行：

```bash
go test ./...
go run .
```

`go run .` 会跑 `main.go` 中的小规模演示（实例化、渲染、答案抽取、行列式与 `solve` 示例）。

### 代码分层

| 目录 / 包 | 职责 |
| --------- | ---- |
| `dsl/` | 核心：Problem 解析、随机实例化、表达式求值、题干渲染、判分/解析相关类型与逻辑 |
| `bank/` | 正式题库：`ch1.go` … `ch7.go` 中 `buildChapter*()` 返回 `dsl.Problem`；`registry.go` 将 **逻辑题键** 映射到 builder |
| `bank/catalog.go` | `AllQuestionKeys`：全量已发布题键列表（顺序与历史 HTML 题单对齐，改序需谨慎） |
| `bank/expected_field_counts.go` | 每题 `GenerateBankQuestion` 期望的答案字段数，供 `bank/all_questions_test.go` 等断言 |
| `bank/chapter_index.go` | 章号解析、`ChapterTitle`、**`publishedBlockedKeys`**：临时下线某题键时只改此处即可 |
| `ladsl/` | 给 **backend** 用的门面：`Service` 统一出题、判题、解析、空位描述；内部调用 `bank` + `dsl` |
| `test/`、`qa/`、`bank/*_test.go` | 回归与专项测试 |

上层服务（如 monorepo 内 `backend`）通常只依赖 **`ladsl`** 与 **`bank`**，不直接拼装 JSON 题面。

### 新增或修改一道 bank 题

1. 在对应 `bank/chN.go` 中实现 `buildChapterN_xxx() dsl.Problem`（题面里的 `{{blank:...}}` 与 `Answer.FieldDefs` 的 `ID` 必须一致；复杂题可参考同章已有题的 `Meta["solution_zh"]`、判题 `Judge` 等写法）。
2. 在 `bank/registry.go` 的 `builders` 中注册：`"ChapterN_xxx": buildChapterN_xxx`。
3. 在 `bank/catalog.go` 的 `AllQuestionKeys` 中按章加入该键（保持与同章其它键的分组与展示顺序约定）。
4. 在 `bank/expected_field_counts.go` 中为该键写入 `GenerateBankQuestion(..., seed, salt)` 返回的 `len(AnswerFields)`（与 HTML 空位数一致时以题单为准）。
5. 运行 `go test ./...`，确保 `bank` 全量用例与相关 `test/` 通过。

**逻辑题键命名**：`Chapter{章号}_{小节编号...}`，须能被 `bank/chapter_index.go` 中的正则解析出章号（1～7）。

**临时不下线代码、只对学生端隐藏**：把题键加入 `publishedBlockedKeys`（见 `chapter_index.go` 注释），不必改 backend / Web 路由。

### 合并前自检清单

- [ ] `go test ./...` 通过  
- [ ] 新题已注册 `registry` + `catalog` + `expected_field_counts`  
- [ ] 若动了表达式或判题语义，检查 `dsl` 侧单测与 `bank/judge_test.go` 等是否需增补  
---

## 核心特性

- **声明式题目定义**：用 DSL 描述题型，bank 中也可用 Go 结构体直写 `dsl.Problem`
- **确定性随机**：相同 `seedStr` + `serverSalt` 生成相同实例，便于复现和测试
- **答案可计算**：表达式系统自动计算派生量和答案
- **填空题支持**：题干包含 `{{blank:ID}}`，答案通过 `field_defs` 绑定 ID

## 快速开始

```go
import "github.com/neumathe/la-dsl/dsl"

// 从 JSON 解析题目定义
var prob dsl.Problem
json.Unmarshal(data, &prob)

// 生成题目（推荐统一入口）
q, err := dsl.GenerateQuestion(prob, "user123", "server-salt")

// q.Title: 渲染后的题干（含 {{blank:ID}} 占位符）
// q.AnswerFields: [{ID, Expr, Value}, ...]
```

底层如需分步调试，可使用 `InstantiateProblem` → `RenderInst` / `ExtractAnswer`（见 `main.go`）。

## 目录（DSL 参考）

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

每道题用一个 JSON 对象表示（Go 中为 `dsl.Problem`）：

```json
{
  "id": 101,
  "version": "v1",
  "title": "向量 β={{beta}} 在基 {{alpha1}}, {{alpha2}}, {{alpha3}} 下的坐标为 {{blank:c1}} {{blank:c2}} {{blank:c3}}",
  "variables": { },
  "derived": { },
  "render": { },
  "answer": { }
}
```

| 字段          | 说明 |
| ------------- | ---- |
| `id`          | 题目 ID |
| `version`     | 版本号（可选） |
| `title`       | 题干模板，支持 `{{var}}` 和 `{{blank:ID}}` 占位符 |
| `variables`   | 随机变量定义 |
| `derived`     | 派生变量（由表达式计算） |
| `render`      | 题面中 `{{var}}` 与内部变量名的映射（表达式字符串）；未在 `title` 出现的可省略 |
| `answer`      | 答案定义 |
| `meta`        | 元信息（可选），如解析文案 `solution_zh` |

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
