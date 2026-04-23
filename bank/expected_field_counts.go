package bank

// ExpectedAnswerFieldCount 与 la-dsl/html 中 blankcount / calculationcount 按题号顺序一致，
// 用于单元测试断言 GenerateBankQuestion 返回的空数。
var ExpectedAnswerFieldCount = map[string]int{
	// 1.html blankcount
	"Chapter1_8_1": 1, "Chapter1_4": 1, "Chapter1_3": 1, "Chapter1_1": 1, "Chapter1_2": 1,
	"Chapter1_6": 1, "Chapter1_5": 1, "Chapter1_7": 3, "Chapter1_8_2": 1,
	// 4804... Chapter 3
	"Chapter3_5": 1, "Chapter3_8": 1, "Chapter3_4": 1, "Chapter3_3": 9, "Chapter3_10": 1,
	"Chapter3_1": 4, "Chapter3_2": 4, "Chapter3_6": 3, "Chapter3_7": 4, "Chapter3_9": 1,
	"Chapter3_11": 24,
	// 5370... Chapter 2
	"Chapter2_4_1": 9, "Chapter2_7_2": 9, "Chapter2_7_3": 9, "Chapter2_4_2": 16, "Chapter2_4_3": 16,
	"Chapter2_6": 1, "Chapter2_7_1": 9,
	"Chapter2_5_2": 21, "Chapter2_1": 16, "Chapter2_3": 8, "Chapter2_2": 8,
	// 5dddb... Chapter 7
	"Chapter7_7": 4, "Chapter7_4": 3, "Chapter7_5_1": 5, "Chapter7_2": 5, "Chapter7_3": 4,
	"Chapter7_1": 1, "Chapter7_5_2": 9, "Chapter7_8": 9, "Chapter7_6": 9, "Chapter7_5_3": 5,
	"Chapter7_10": 12, "Chapter7_9": 9,
	// 6cf42... Chapter 6
	"Chapter6_1_2": 10, "Chapter6_1_3": 10, "Chapter6_1_1": 10, "Chapter6_5": 1,
	"Chapter6_4": 12, "Chapter6_2": 2, "Chapter6_6": 10, "Chapter6_3": 12,
	// 81c04... Chapter 4
	"Chapter4_8": 12, "Chapter4_3_1": 20, "Chapter4_1": 5, "Chapter4_4": 1, "Chapter4_2": 3,
	"Chapter4_3_3": 20, "Chapter4_3_2": 20, "Chapter4_5_1": 12, "Chapter4_5_2": 12, "Chapter4_7": 8,
	"Chapter4_6": 22,
	// dd926... Chapter 5
	"Chapter5_1": 12, "Chapter5_3": 4, "Chapter5_5": 5, "Chapter5_2": 12, "Chapter5_4": 4,
	"Chapter5_8": 24, "Chapter5_6": 9, "Chapter5_7": 12,
}
