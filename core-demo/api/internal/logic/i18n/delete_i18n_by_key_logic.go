package i18n

import (
	"context"

	corei18n "oa.98ent.com/p9/common/i18n"
	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteI18nByKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteI18nByKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteI18nByKeyLogic {
	return &DeleteI18nByKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteI18nByKeyLogic) DeleteI18nByKey(req *types.DeleteI18nByKeyReq) (resp *types.BaseMsgResp, err error) {
	_, err = l.svcCtx.Core.DeleteI18NByKey(l.ctx, convert.DeleteI18nByKeyReq(req))
	if err != nil {
		return nil, err
	}
	corei18n.InvalidateAll()
	return &types.BaseMsgResp{Result: "success"}, nil
}
