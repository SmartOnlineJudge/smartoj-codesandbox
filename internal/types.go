package internal


// 测试用例结构
type Test struct {
	TestId      int    `json:"test_id"` 
	InputOutput string `json:"input_output"` // 输入输出数据
}

// 判题数据
type JudgementData struct {
	Language      string `json:"language"`   // 编程语言名称: python | java | c | cpp | go
	QuestionId    int `json:"question_id"` // 当前题目对应的 ID
	Template      string `json:"judge_template"` // 该编程语言对应的判题模板代码
	Code          string `json:"solution_code"` // 用户编写的代码
	Tests         []Test `json:"tests"` // 测试用例列表
	TimeLimit     int `json:"time_limit"` // 时间限制（单位：ms）
	MemoryLimit   float32 `json:"memory_limit"` // 内存限制（单位：MB）
	UserId        string `json:"user_id"` // 用户ID
}

// 单个测试用例的结果
type Result struct {
	TestId         int    `json:"test_id"`
	Status         int    `json:"status"` // 代码运行异常时为-1
	Result         string `json:"result"`  // 判题结果（OK | ErrorMessage）
	Solution       string `json:"solution"` // 解题函数的输出（字符串）
	Criterion      string `json:"criterion"` // 正确的输出（字符串）
	TimeConsumed   int  `json:"time_consumed"`  // 运行该测试用例消耗的时间
	MemoryConsumed int  `json:"memory_consumed"`  // 运行该测试用例消耗的内存
	IsSuccess      bool   `json:"is_success"` // 是否通过该测试用例
}

// 判题响应结构
type SandboxResponse struct {
	Code    int       `json:"code"`    // 200表示判题成功（不代表通过测试用例），422表示输入数据有问题
	Message string    `json:"message"` // OK | 其他消息
	Results []Result  `json:"results"` // 测试用例结果列表
}

// 语言配置
type LanguageConfig struct {
	CompileCmd   []string // 编译命令（如果需要）
	RunCmd       []string // 运行命令
	FileExt      string   // 文件扩展名
}

// 支持的语言配置
var SupportedLanguages = map[string]LanguageConfig{
	"python": {
		CompileCmd:   nil, // Python不需要编译
		RunCmd:       []string{"python3"}, // 直接使用系统安装的python3
		FileExt:      ".py",
	},
	"cpp": {
		CompileCmd:   []string{"g++", "-std=c++17", "-O2"}, // 使用系统安装的g++
		RunCmd:       []string{"./a.out"}, // 编译后的可执行文件
		FileExt:      ".cpp",
	},
	// 其他语言的支持
}
