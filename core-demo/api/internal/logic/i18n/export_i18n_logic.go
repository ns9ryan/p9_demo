package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportI18nLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportI18nLogic {
	return &ExportI18nLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportI18nLogic) ExportI18n(req *types.ExportI18nReq) (resp *types.ExportI18nResp, err error) {
	out, err := l.svcCtx.Core.ExportI18N(l.ctx, convert.ExportI18nReq(req))
	if err != nil {
		return nil, err
	}
	return convert.ExportI18n(out), nil
}
