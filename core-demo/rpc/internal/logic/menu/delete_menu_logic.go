package menu

import (
	"context"

	"oa.98ent.com/p9/core/common/utils"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMenuLogic) DeleteMenu(in *core.IDsReq) (*core.Empty, error) {
	if err := l.svcCtx.Deps.DeleteMenus(l.ctx, utils.ResolveIDs(in.Id, in.Ids)); err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
