package menu

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMenuLogic) UpdateMenu(in *core.UpdateMenuReq) (*core.Empty, error) {
	err := l.svcCtx.Deps.UpdateMenu(l.ctx, service.UpdateMenuReq{
		ID: in.Id, ParentID: in.ParentId, MenuType: logic.ToInt16Ptr(in.MenuType), Path: in.Path, Name: in.Name,
		Component: in.Component, Redirect: in.Redirect, Title: in.Title, Icon: in.Icon, Permission: in.Permission,
		HideMenu: logic.ToInt16Ptr(in.HideMenu), Sort: logic.ToIntPtr(in.Sort), Disabled: logic.ToInt16Ptr(in.Disabled),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
