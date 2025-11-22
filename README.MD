# 线性代数题库统一 DSL

## 目录

1. 概览（Overview）
2. 题目对象结构（Problem Object）
3. 随机变量（Variables）
4. 随机生成规则（Generators）
5. 派生变量（Derived Variables）
6. 表达式系统（Expressions）
7. 题面渲染（Rendering）
8. 答案结构（Answer Schema）
9. DEMO

---

## 1. 概览（Overview）

DSL 的主要目标：

* **统一描述题型**：用统一的数据结构表达不同题目
* **可复现生成**：通过随机种子生成题目
* **可计算**：派生量和答案由表达式系统计算
* **统一执行环境**：由一个 runtime 解释并执行 DSL

只要 runtime 支持该 DSL，就可以在不改代码的情况下扩展题库。

---

## 2. 题目对象结构（Problem Object）

每道题用一个 JSON（或 YAML）对象表示：

```json
{
  "id": 101,
  "title": "题目描述模板",
  "variables": { ... },
  "derived": { ... },
  "render": { ... },
  "answer": { ... },
  "meta": { ... }
}
```

字段说明：

| 字段            | 含义                                  |
| ------------- | ----------------------------------- |
| **id**        | 题目模板 id                           |
| **title**     | 题干渲染模板（可包含占位符）                 |
| **variables** | 原始随机变量（矩阵、向量、标量等）              |
| **derived**   | 派生变量（由 `variables` 通过表达式计算） |
| **render**    | 题干渲染时需要暴露的变量                     |
| **answer**    | 标准答案对应的表达式或字段结构                |
| **meta**      | 章节、难度等元信息                         |

---

## 3. 随机变量（Variables）

随机变量是题目中直接生成的量，例如矩阵、向量或标量。基本结构如下：

```json
"variables": {
  "A": {
    "kind": "matrix",
    "rows": 3,
    "cols": 3,
    "generator": {
      "rule": "...",
      "...": "rule-specific"
    }
  },
  "x": {
    "kind": "vector",
    "size": 3,
    "generator": {
      "rule": "range",
      "min": -5,
      "max": 5
    }
  }
}
```

### `kind` 允许的类型

| kind       | 描述     |
| ---------- | ------ |
| **scalar** | 数字     |
| **vector** | 一维数组   |
| **matrix** | 二维数组   |

---

## 4. 随机生成规则（Generators）

随机生成逻辑通过 `generator` 字段进行配置，不同的 `rule` 对应不同的策略。

### 4.1 `rule: "range"`（整型随机数）

适用于 `scalar` / `vector` / `matrix`：

```json
"generator": {
  "rule": "range",
  "min": -5,
  "max": 10
}
```

---

### 4.2 `rule: "from_set"`（从固定集合中选）

```json
"generator": {
  "rule": "from_set",
  "set": [1, 2, 3, 5, 8]
}
```

---

### 4.3 `rule: "full_rank"`（随机满秩矩阵）

生成可逆矩阵，适合坐标变换、求逆矩阵等题目：

```json
"generator": {
  "rule": "full_rank",
  "min": -5,
  "max": 5
}
```

---

### 4.4 `rule: "sparse"`（稀疏矩阵）

```json
"generator": {
  "rule": "sparse",
  "values": [-9, -6, -3, -1, 1, 2, 4],
  "density": 0.3
}
```

---

### 4.5 `rule: "integer_solution"`（保证解为整数）

适合线性方程组、向量坐标类题目。先定义 `A`、`x`，再自动生成 `b`：

```json
"generator": {
  "rule": "integer_solution",
  "A": "A",
  "x": "x"
}
```

runtime 自动计算：

> b = A · x

---

## 5. 派生变量（Derived Variables）

派生变量通过表达式由已有变量计算得到：

```json
"derived": {
  "beta": "A * x",
  "C": "A + B",
  "d": "det(A)"
}
```

表达式由表达式系统解析和执行（见第 6 节）。

---

## 6. 表达式系统（Expressions）

表达式系统提供基本运算符和常用线性代数函数。

### 6.1 运算符

| 运算符 | 含义          |
| --- | ----------- |
| +   | 加法          |
| -   | 减法          |
| *   | 标量乘法 / 矩阵乘法 |
| []  | 索引          |

### 6.2 内置函数

| 函数                  | 描述        |
| ------------------- | --------- |
| `det(A)`            | 行列式       |
| `rank(A)`           | 秩         |
| `cofactor(A, i, j)` | 代数余子式     |
| `minor(A, i, j)`    | 子式        |
| `inverse(A)`        | 逆矩阵       |
| `solve(A, b)`       | 求解 Ax = b |
| `transpose(A)`      | 转置        |
| `col(A, i)`         | 第 i 列向量   |
| `row(A, i)`         | 第 i 行向量   |
| `eigenvalues(A)`    | 特征值       |
| `eigenvectors(A)`   | 特征向量      |

可根据需要继续扩展新的函数（见第 10 节）。

---

## 7. 题面渲染（Rendering）

`render` 区域指定题干展示需要的变量，通常配合 `title` 中的占位符使用：

```json
"render": {
  "beta": "beta",
  "alpha1": "col(A,1)",
  "alpha2": "col(A,2)",
  "alpha3": "col(A,3)"
}
```

`title` 模板中可以这样引用：

```
向量 β={{beta}} 在基 {{alpha1}},{{alpha2}},{{alpha3}} 下的坐标
```

runtime 会在渲染时把占位符替换为实际数值或排版结果。

---

## 8. 答案结构（Answer Schema）

答案可以是一个整体表达式，也可以拆成多个空格。

### 8.1 单个表达式

```json
"answer": {
  "expression": "x"
}
```

### 8.2 多个空格的情况

```json
"answer": {
  "fields": [
    "x[1]",
    "x[2]",
    "x[3]"
  ]
}
```

---

## 9. DEMO：向量在一组基下的坐标

```json
{
  "id": 401,
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
      "generator": { "rule": "full_rank", "min": -5, "max": 5 }
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
  },
  "meta": {
    "chapter": "4.2",
    "difficulty": 3
  }
}
```