package i18n

import (
	"context"
	"strings"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportI18nLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExportI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportI18nLogic {
	return &ExportI18nLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExportI18nLogic) ExportI18N(in *core.ExportI18NReq) (*core.ExportI18NResp, error) {
	list, err := l.svcCtx.Deps.ExportI18ns(l.ctx, in.GetLang(), in.GetI18NCode(), in.GetI18NGroup())
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	items := make([]*core.I18NFileItem, 0, len(list))
	for _, row := range list {
		items = append(items, &core.I18NFileItem{
			I18NCode:  row.I18nCode,
			I18NGroup: row.I18nGroup,
			I18NKey:   row.I18nKey,
			I18NValue: row.I18nValue,
		})
	}
	return &core.ExportI18NResp{Lang: strings.TrimSpace(in.GetLang()), I18NItems: items}, nil
}
