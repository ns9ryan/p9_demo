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

type GetI18nListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetI18nListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetI18nListLogic {
	return &GetI18nListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetI18nListLogic) GetI18NList(in *core.I18NListReq) (*core.I18NListResp, error) {
	list, total, err := l.svcCtx.Deps.ListI18ns(l.ctx, service.I18nListReq{
		PageReq:   service.PageReq{Page: int(in.GetPage()), PageSize: int(in.GetPageSize())},
		I18nGroup: in.GetI18NGroup(),
		TransKey:  in.GetTransKey(),
		Lang:      in.GetLang(),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	out := make([]*core.I18NInfo, 0, len(list))
	for _, row := range list {
		out = append(out, logic.ToI18nInfo(row))
	}
	return &core.I18NListResp{List: out, Total: total}, nil
}
