package user

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeOwnPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeOwnPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeOwnPasswordLogic {
	return &ChangeOwnPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeOwnPasswordLogic) ChangeOwnPassword(in *core.SelfPasswordReq) (*core.Empty, error) {
	err := l.svcCtx.Deps.ChangeOwnPassword(l.ctx, ctxdata.ClaimsFromCtx(l.ctx), in.OldPassword, in.Password)
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
