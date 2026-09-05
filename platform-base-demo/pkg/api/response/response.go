package response

import "context"

// SuccessResponse 成功响应
type SuccessResponse struct {
	Code int    `json:"code"`           // 状态码
	Msg  string `json:"msg"`            // 提示信息
	Data any    `json:"data,omitempty"` // 响应数据
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code  int        `json:"code"`            // 状态码
	Msg   string     `json:"msg"`             // 提示信息
	Debug *DebugInfo `json:"debug,omitempty"` // 调试信息
}

// DebugInfo 调试信息
type DebugInfo struct {
	Source string `json:"source"`        // 错误来源
	Rpc    string `json:"rpc,omitempty"` // RPC方法
	Error  string `json:"error"`         // 原始错误
}

// Ok 创建成功响应
func Ok(_ context.Context, data any) any {
	return &SuccessResponse{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
}
