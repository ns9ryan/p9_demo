package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReorderI18nLangLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReorderI18nLangLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReorderI18nLangLogic {
	return &ReorderI18nLangLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReorderI18nLangLogic) ReorderI18nLang(req *types.ReorderI18nLangReq) (resp *types.BaseMsgResp, err error) {
	_, err = l.svcCtx.Core.ReorderI18NLang(l.ctx, convert.ReorderI18nLangReq(req))
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Result: "success"}, nil
}
