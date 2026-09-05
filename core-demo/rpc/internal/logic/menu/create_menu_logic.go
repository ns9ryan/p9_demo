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

type CreateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMenuLogic) CreateMenu(in *core.CreateMenuReq) (*core.MenuInfo, error) {
	row, err := l.svcCtx.Deps.CreateMenu(l.ctx, service.CreateMenuReq{
		ParentID: in.ParentId, MenuType: int16(in.MenuType), Path: in.Path, Name: in.Name, Component: in.Component,
		Redirect: in.Redirect, Title: in.Title, Icon: in.Icon, Permission: in.Permission,
		HideMenu: int16(in.HideMenu), Sort: int(in.Sort), Disabled: int16(in.Disabled),
	})
	if err != nil {
		l.Errorf("create menu failed: %v", err)
		return nil, xerr.RpcErr(err)
	}
	return logic.ToMenuInfo(*row), nil
}
