package types

// TestAuthCheckReq 测试鉴权检查请求
type TestAuthCheckReq struct {
}

// TestAuthCheckResp 测试鉴权检查响应
type TestAuthCheckResp struct {
	// 测试消息
	Message string `json:"message"`
	// 如何获取 token 的说明
	Instructions string `json:"instructions"`
}

// TestAuthWithTokenResp 测试 token 有效性的响应
type TestAuthWithTokenResp struct {
	// 用户 ID
	UserID int64 `json:"user_id"`
	// 用户编码
	UserCode string `json:"user_code"`
	// 用户名
	Username string `json:"username"`
	// 操作员 ID
	OperatorID int64 `json:"operator_id"`
	// 操作员编码
	OperatorCode string `json:"operator_code"`
	// 角色编码列表
	RoleCodes []string `json:"role_codes"`
	// Token 类型
	TokenType string `json:"token_type"`
	// 是否为平台用户
	IsPlatform bool `json:"is_platform"`
	// 过期时间
	ExpiresAt int64 `json:"expires_at"`
	// 验证消息
	Message string `json:"message"`
}
