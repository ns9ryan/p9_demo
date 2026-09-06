package logic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/logger"
)

// TestAuthCheckLogic 测试鉴权检查逻辑
type TestAuthCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewTestAuthCheckLogic 创建新的测试鉴权检查逻辑
func NewTestAuthCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestAuthCheckLogic {
	return &TestAuthCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// TestAuthCheck 执行测试鉴权检查
// 此接口不需要认证，用于指导用户如何获取和测试 token
func (l *TestAuthCheckLogic) TestAuthCheck(req *types.TestAuthCheckReq) (*types.TestAuthCheckResp, error) {
	logger.Infof("[TestAuthCheck] 收到测试鉴权检查请求")

	resp := &types.TestAuthCheckResp{
		Message: "Platform-Game 鉴权测试接口已就绪",
		Instructions: `
步骤 1: 登录获取 Token
- 在 Core 服务中登录，获取 access_token
- API 端点: POST http://core-api:8000/auth/login
- 请求体:
  {
    "username": "admin",
    "password": "password"
  }
- 响应包含 access_token 和 refresh_token

步骤 2: 验证 Token
- 使用获取的 access_token 调用本接口
- API 端点: GET http://platform-game-api:8001/test/auth-with-token
- 请求头: Authorization: Bearer {access_token}
- 成功响应会显示您的用户信息和权限

步骤 3: 查看权限
- 如果 TokenType == "preview"，则只有读权限
- 如果 IsPlatform == true，则为平台管理员
- RoleCodes 包含您的所有角色编码
`,
	}

	logger.Infof("[TestAuthCheck] 返回测试说明: %+v", resp)
	return resp, nil
}

// TestAuthWithTokenLogic 测试 token 有效性逻辑
type TestAuthWithTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewTestAuthWithTokenLogic 创建新的测试 token 有效性逻辑
func NewTestAuthWithTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestAuthWithTokenLogic {
	return &TestAuthWithTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// TestAuthWithToken 执行测试 token 有效性
// 此接口需要有效的 JWT token（通过 JWT 中间件验证）
func (l *TestAuthWithTokenLogic) TestAuthWithToken() (*types.TestAuthWithTokenResp, error) {
	// 从上下文中获取声明信息
	claims := ctxdata.ClaimsFromCtx(l.ctx)

	if claims == nil {
		logger.Warnf("[TestAuthWithToken] 声明信息为空，可能 JWT 中间件未正确配置")
		return nil, fmt.Errorf("no claims found in context, JWT middleware may not be configured")
	}

	logger.Infof("[TestAuthWithToken] 成功验证 Token:")
	logger.Infof("  - UserID: %d", claims.UserID)
	logger.Infof("  - Username: %s", claims.Username)
	logger.Infof("  - OperatorCode: %s", claims.OperatorCode)
	logger.Infof("  - RoleCodes: %v", claims.RoleCodes)
	logger.Infof("  - IsPlatform: %v", claims.IsPlatform)
	logger.Infof("  - TokenType: %s", claims.TokenType)

	resp := &types.TestAuthWithTokenResp{
		UserID:       claims.UserID,
		UserCode:     claims.UserCode,
		Username:     claims.Username,
		OperatorID:   claims.OperatorID,
		OperatorCode: claims.OperatorCode,
		RoleCodes:    claims.RoleCodes,
		TokenType:    claims.TokenType,
		IsPlatform:   claims.IsPlatform,
		ExpiresAt:    claims.ExpiresAt,
		Message:      "✓ Token 有效，鉴权成功！",
	}

	return resp, nil
}
