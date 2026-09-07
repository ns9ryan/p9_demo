package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteI18nLangLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteI18nLangLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteI18nLangLogic {
	return &DeleteI18nLangLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteI18nLangLogic) DeleteI18nLang(req *types.IDsReq) (resp *types.BaseMsgResp, err error) {
	_, err = l.svcCtx.Core.DeleteI18NLang(l.ctx, convert.IDsReq(req))
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Result: "success"}, nil
}
