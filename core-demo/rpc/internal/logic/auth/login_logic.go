package auth

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *core.LoginReq) (*core.LoginResp, error) {
	res, err := l.svcCtx.Deps.Login(l.ctx, service.LoginReq{
		Username: in.Username, Password: in.Password, OperatorCode: in.OperatorCode, ClientIP: in.ClientIp, UserAgent: in.UserAgent,
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return logic.ToLoginResp(res), nil
}
