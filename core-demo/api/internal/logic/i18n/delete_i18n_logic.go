package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	corei18n "oa.98ent.com/p9/core/common/i18n"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteI18nLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteI18nLogic {
	return &DeleteI18nLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteI18nLogic) DeleteI18n(req *types.IDsReq) (resp *types.BaseMsgResp, err error) {
	_, err = l.svcCtx.Core.DeleteI18N(l.ctx, convert.IDsReq(req))
	if err != nil {
		return nil, err
	}
	corei18n.InvalidateAll()
	return &types.BaseMsgResp{Result: "success"}, nil
}
