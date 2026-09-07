package i18n

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetI18nLangListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetI18nLangListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetI18nLangListLogic {
	return &GetI18nLangListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetI18nLangListLogic) GetI18NLangList(in *core.I18NLangListReq) (*core.I18NLangListResp, error) {
	list, total, err := l.svcCtx.Deps.ListI18nLangs(l.ctx, service.I18nLangListReq{
		PageReq:  service.PageReq{Page: int(in.GetPage()), PageSize: int(in.GetPageSize())},
		Lang:     in.GetLang(),
		Disabled: logic.ToInt16Ptr(in.Disabled),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	out := make([]*core.I18NLangInfo, 0, len(list))
	for _, row := range list {
		out = append(out, logic.ToI18nLangInfo(row))
	}
	return &core.I18NLangListResp{List: out, Total: total}, nil
}
