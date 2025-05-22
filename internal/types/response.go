package types

// 判题响应结构
type SandboxResponse struct {
	// 200表示判题成功（不代表通过测试用例），422表示输入数据有问题
	Code    int       `json:"code"`
	// OK | 其他消息
	Message string    `json:"message"`
	 // 测试用例结果列表
	Results *Results  `json:"results"`
}
