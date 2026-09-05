package user

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPermLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserPermLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPermLogic {
	return &GetUserPermLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserPermLogic) GetUserPerm(in *core.Empty) (*core.PermResp, error) {
	codes, err := l.svcCtx.Deps.PermCodes(l.ctx, ctxdata.ClaimsFromCtx(l.ctx))
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	if codes == nil {
		codes = []string{}
	}
	return &core.PermResp{Permissions: codes}, nil
}
