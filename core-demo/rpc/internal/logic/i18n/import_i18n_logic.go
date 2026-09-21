package i18n

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportI18nLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewImportI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportI18nLogic {
	return &ImportI18nLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ImportI18nLogic) ImportI18N(in *core.ImportI18NReq) (*core.ImportI18NResp, error) {
	items := make([]service.I18nFileItem, 0, len(in.GetI18NItems()))
	for _, row := range in.GetI18NItems() {
		if row == nil {
			continue
		}
		items = append(items, service.I18nFileItem{
			I18nCode:  row.GetI18NCode(),
			I18nGroup: row.GetI18NGroup(),
			I18nKey:   row.GetI18NKey(),
			I18nValue: row.GetI18NValue(),
		})
	}
	out, err := l.svcCtx.Deps.ImportI18ns(l.ctx, in.GetLang(), items)
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.ImportI18NResp{Created: out.Created, Updated: out.Updated, Skipped: out.Skipped}, nil
}
