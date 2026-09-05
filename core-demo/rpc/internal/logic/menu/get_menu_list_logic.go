package menu

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuListLogic {
	return &GetMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuListLogic) GetMenuList(in *core.Empty) (*core.MenuListResp, error) {
	list, err := l.svcCtx.Deps.ListMenus(l.ctx)
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	out := make([]*core.MenuInfo, 0, len(list))
	for _, m := range list {
		out = append(out, logic.ToMenuInfo(m))
	}
	return &core.MenuListResp{List: out}, nil
}
