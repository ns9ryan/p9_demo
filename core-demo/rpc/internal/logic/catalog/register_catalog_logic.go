package catalog

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterCatalogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterCatalogLogic {
	return &RegisterCatalogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterCatalogLogic) RegisterCatalog(in *core.RegisterCatalogReq) (*core.Empty, error) {
	menus := make([]service.RegisterMenuReq, 0, len(in.Menus))
	for _, m := range in.Menus {
		menus = append(menus, service.RegisterMenuReq{
			Name: m.Name, Title: m.Title, MenuType: int16(m.MenuType), Path: m.Path,
			Component: m.Component, Redirect: m.Redirect, Icon: m.Icon, Permission: m.Permission,
			HideMenu: int16(m.HideMenu), Sort: int(m.Sort), Disabled: int16(m.Disabled),
			ParentName: m.ParentName,
		})
	}
	apis := make([]service.CreateAPIReq, 0, len(in.Apis))
	for _, a := range in.Apis {
		apis = append(apis, service.CreateAPIReq{
			Description: a.Description, APIGroup: a.ApiGroup, Method: a.Method, Path: a.Path,
			IsRequired: int16(a.IsRequired), ServiceName: a.ServiceName,
		})
	}
	items := make([]service.I18nItem, 0, len(in.GetI18N()))
	for _, it := range in.GetI18N() {
		items = append(items, service.I18nItem{
			I18nGroup: it.GetI18NGroup(), TransKey: it.GetTransKey(), Lang: it.GetLang(), Value: it.GetValue(),
		})
	}
	langs := make([]service.CreateI18nLangReq, 0, len(in.GetI18NLangs()))
	for _, it := range in.GetI18NLangs() {
		langs = append(langs, service.CreateI18nLangReq{
			Lang: it.GetLang(), Name: it.GetName(), Disabled: int16(it.GetDisabled()), IsDefault: int16(it.GetIsDefault()),
		})
	}
	if err := l.svcCtx.Deps.RegisterCatalog(l.ctx, menus, apis, items, langs); err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
