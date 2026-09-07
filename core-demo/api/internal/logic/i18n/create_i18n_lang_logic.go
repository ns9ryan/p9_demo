package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateI18nLangLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateI18nLangLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateI18nLangLogic {
	return &CreateI18nLangLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateI18nLangLogic) CreateI18nLang(req *types.CreateI18nLangReq) (resp *types.I18nLangInfo, err error) {
	out, err := l.svcCtx.Core.CreateI18NLang(l.ctx, convert.CreateI18nLangReq(req))
	if err != nil {
		return nil, err
	}
	return convert.I18nLangInfo(out), nil
}
