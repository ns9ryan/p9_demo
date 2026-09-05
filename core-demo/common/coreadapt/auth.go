package coreadapt

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/middleware"
	"oa.98ent.com/p9/core/rpc/coreclient"
)

// auth 认证客户端
type auth struct {
	cli coreclient.Core
}

// Auth 创建认证客户端
func Auth(cli coreclient.Core) middleware.Client {
	return &auth{cli: cli}
}

// CheckToken 检查令牌
func (a *auth) CheckToken(ctx context.Context, token string) (*ctxdata.Claims, error) {
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

// Enforce 强制执行检查权限
func (a *auth) Enforce(ctx context.Context, claims *ctxdata.Claims, path, method string) (bool, error) {
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
