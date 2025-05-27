package types


// 测试用例结构
type Test struct {
	TestId      int    `json:"test_id"` 
	InputOutput string `json:"input_output"` // 输入输出数据
}

// 测试用例切片类型
type Tests []Test

// 判题数据
type JudgementData struct {
	Language      string `json:"language"`   // 编程语言名称: python | java | c | cpp | go
	QuestionId    int `json:"question_id"` // 当前题目对应的 ID
	JudgeTemplate string `json:"judge_template"` // 该编程语言对应的判题模板代码
	SolutionCode  string `json:"solution_code"` // 用户编写的代码
	Tests         Tests `json:"tests"` // 测试用例列表
	TimeLimit     int `json:"time_limit"` // 时间限制（单位：ms）
	MemoryLimit   float32 `json:"memory_limit"` // 内存限制（单位：MB）
	UserId        string `json:"user_id"` // 用户ID
}

// 单个测试用例的结果
type Result struct {
	TestId         int    `json:"test_id"`
	Status         int    `json:"status"` // 代码运行异常时为-1
	Result         string `json:"result"`  // 判题结果（OK | ErrorMessage）
	Answer         any `json:"answer"` // 解题函数的输出（字符串）
	Criterion      any `json:"criterion"` // 正确的输出（字符串）
	TimeConsumed   int  `json:"time_consumed"`  // 运行该测试用例消耗的时间
	MemoryConsumed float32  `json:"memory_consumed"`  // 运行该测试用例消耗的内存
	IsSuccess      bool   `json:"is_success"` // 是否通过该测试用例
}

// 判题结果切片类型
type Results []Result
