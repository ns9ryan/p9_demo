package menu

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuByRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuByRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuByRoleLogic {
	return &GetMenuByRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuByRoleLogic) GetMenuByRole(in *core.Empty) (*core.MenuTreeResp, error) {
	tree, err := l.svcCtx.Deps.MenuTreeByRole(l.ctx, ctxdata.ClaimsFromCtx(l.ctx))
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	list := make([]*core.MenuNode, 0, len(tree))
	for _, n := range tree {
		list = append(list, logic.ToMenuNode(n))
	}
	return &core.MenuTreeResp{List: list}, nil
}
