package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	corei18n "oa.98ent.com/p9/core/common/i18n"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportI18nLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportI18nLogic {
	return &ImportI18nLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportI18nLogic) ImportI18n(req *types.ImportI18nReq) (resp *types.ImportI18nResp, err error) {
	file, err := convert.UnmarshalI18nFile(req.File)
	if err != nil {
		return nil, err
	}
	if err := convert.CheckImportI18nLang(req.Lang, file.Lang); err != nil {
		return nil, err
	}
	out, err := l.svcCtx.Core.ImportI18N(l.ctx, convert.ImportI18nReq(file))
	if err != nil {
		return nil, err
	}
	corei18n.InvalidateAll()
	return convert.ImportI18n(out), nil
}
