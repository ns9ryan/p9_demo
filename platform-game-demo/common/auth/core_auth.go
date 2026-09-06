package auth

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/rpc/coreclient"
)

// CoreAuth 实现 core 服务的鉴权客户端接口
type CoreAuth struct {
	cli coreclient.Core
}

// NewCoreAuth 创建新的鉴权客户端
func NewCoreAuth(cli coreclient.Core) *CoreAuth {
	return &CoreAuth{cli: cli}
}

// CheckToken 验证 token 并返回声明信息
// 参数: ctx 上下文, token 访问令牌
// 返回: 声明信息, 错误信息
func (a *CoreAuth) CheckToken(ctx context.Context, token string) (*ctxdata.Claims, error) {
	resp, err := a.cli.CheckToken(ctx, &coreclient.CheckTokenReq{AccessToken: token})
	if err != nil {
		return nil, err
	}

	codes := resp.RoleCodes
	if codes == nil {
		codes = []string{}
	}

	return &ctxdata.Claims{
		UserID:       resp.UserId,
		UserCode:     resp.UserCode,
		Username:     resp.Username,
		OperatorID:   resp.OperatorId,
		OperatorCode: resp.OperatorCode,
		RoleCodes:    codes,
		Salt:         resp.Salt,
		ExpiresAt:    resp.ExpiresAt,
		IsPlatform:   resp.IsPlatform,
		TokenType:    resp.TokenType,
	}, nil
}

// Enforce 检查用户是否有权限执行指定的操作
// 参数: ctx 上下文, claims 声明信息, path 访问路径, method HTTP 方法
// 返回: 是否有权限, 错误信息
func (a *CoreAuth) Enforce(ctx context.Context, claims *ctxdata.Claims, path, method string) (bool, error) {
	req := &coreclient.EnforceReq{Path: path, Method: method}
	if claims != nil {
		req.RoleCodes = claims.RoleCodes
		req.OperatorId = claims.OperatorID
	}

	resp, err := a.cli.Enforce(ctx, req)
	if err != nil {
		return false, err
	}

	return resp.Allowed, nil
}
