package auth

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTokenLogic {
	return &CheckTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckTokenLogic) CheckToken(in *core.CheckTokenReq) (*core.CheckTokenResp, error) {
	c, err := l.svcCtx.Deps.CheckToken(l.ctx, in.AccessToken)
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.CheckTokenResp{
		UserId: c.UserID, UserCode: c.UserCode, Username: c.Username, OperatorId: c.OperatorID,
		RoleCodes: c.RoleCodes, Salt: c.Salt, ExpiresAt: c.ExpiresAt, OperatorCode: c.OperatorCode,
		IsPlatform: c.IsPlatform, TokenType: c.TokenType,
	}, nil
}
