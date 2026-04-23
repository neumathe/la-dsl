package bank

import (
	"regexp"
	"sort"
	"strconv"
)

// chapterTitles la-dsl 各章在教材中的常用名称（与题库章号 1~7 对应）。
var chapterTitles = map[int]string{
	1: "行列式",
	2: "矩阵",
	3: "向量组的线性相关性",
	4: "线性方程组",
	5: "矩阵相似对角化",
	6: "二次型",
	7: "线性空间与线性变换",
}

// ChapterTitle 返回指定章号的中文名；未知章号返回空字符串。
func ChapterTitle(chapterNo int) string {
	return chapterTitles[chapterNo]
}

// chapterKeyRe 匹配逻辑题键前缀 "ChapterN_" 中的章号 N。
var chapterKeyRe = regexp.MustCompile(`^Chapter(\d+)_`)

// publishedBlockedKeys 为不希望对学生端开放的题键集合。
// 审计结论（见 la-dsl/最新审计报告.md）：S0/S1/S4 全部达标，目前可上线全量；
// 如后续发现需要下线某题，追加到此集合即可，backend/web 无需改动。
var publishedBlockedKeys = map[string]struct{}{}

// ChapterNoOf 解析逻辑题键所属章号（1~7）。
// 返回 ok=false 表示题键格式不规范。
func ChapterNoOf(key string) (int, bool) {
	m := chapterKeyRe.FindStringSubmatch(key)
	if m == nil {
		return 0, false
	}
	n, err := strconv.Atoi(m[1])
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}

// KeysByChapter 返回指定章号下 AllQuestionKeys 中的全部题键，按字典序稳定排序。
func KeysByChapter(chapterNo int) []string {
	if chapterNo <= 0 {
		return nil
	}
	out := make([]string, 0, 16)
	for _, k := range AllQuestionKeys {
		if n, ok := ChapterNoOf(k); ok && n == chapterNo {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

// PublishedKeysByChapter 在 KeysByChapter 基础上剔除 publishedBlockedKeys。
// 供 backend 在生成章节会话时使用，保证学生端仅看到已发布题。
func PublishedKeysByChapter(chapterNo int) []string {
	raw := KeysByChapter(chapterNo)
	if len(publishedBlockedKeys) == 0 {
		return raw
	}
	out := raw[:0]
	for _, k := range raw {
		if _, blocked := publishedBlockedKeys[k]; blocked {
			continue
		}
		out = append(out, k)
	}
	return out
}

// PublishedChapterNos 返回存在已发布题目的章号列表，按升序排序。
func PublishedChapterNos() []int {
	seen := map[int]struct{}{}
	for _, k := range AllQuestionKeys {
		n, ok := ChapterNoOf(k)
		if !ok {
			continue
		}
		if _, blocked := publishedBlockedKeys[k]; blocked {
			continue
		}
		seen[n] = struct{}{}
	}
	out := make([]int, 0, len(seen))
	for n := range seen {
		out = append(out, n)
	}
	sort.Ints(out)
	return out
}
