package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetI18nLangListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetI18nLangListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetI18nLangListLogic {
	return &GetI18nLangListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetI18nLangListLogic) GetI18nLangList(req *types.I18nLangListReq) (resp *types.I18nLangListResp, err error) {
	out, err := l.svcCtx.Core.GetI18NLangList(l.ctx, convert.I18nLangListReq(req))
	if err != nil {
		return nil, err
	}
	return convert.I18nLangList(out), nil
}
