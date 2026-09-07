package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateI18nLangLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateI18nLangLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateI18nLangLogic {
	return &UpdateI18nLangLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateI18nLangLogic) UpdateI18nLang(req *types.UpdateI18nLangReq) (resp *types.BaseMsgResp, err error) {
	_, err = l.svcCtx.Core.UpdateI18NLang(l.ctx, convert.UpdateI18nLangReq(req))
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Result: "success"}, nil
}
